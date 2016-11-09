// Copyright 2015 GoCrab Author neeke@php.net All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Usage:
//
// import "github.com/CloudWise-OpenSource/GoCrab/Libs/httplib"
//
//	b := httplib.Post("http://GoCrab.me/")
//	b.Param("username","Neeke")
//	b.Param("password","123456")
//	b.PostFile("uploadfile1", "httplib.pdf")
//	b.PostFile("uploadfile2", "httplib.txt")
//	str, err := b.String()
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println(str)
//
package httplib

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var compressNull = ""
var compressGzip = "gzip"

var defaultSetting = GoCrabHttpSettings{compressNull, false, "GoCrabServer", 60 * time.Second, 60 * time.Second, nil, nil, nil, false}
var defaultCookieJar http.CookieJar
var settingMutex sync.Mutex

// createDefaultCookie creates a global cookiejar to store cookies.
func createDefaultCookie() {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultCookieJar, _ = cookiejar.New(nil)
}

// Overwrite default settings
func SetDefaultSetting(setting GoCrabHttpSettings) {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultSetting = setting
	if defaultSetting.ConnectTimeout == 0 {
		defaultSetting.ConnectTimeout = 60 * time.Second
	}
	if defaultSetting.ReadWriteTimeout == 0 {
		defaultSetting.ReadWriteTimeout = 60 * time.Second
	}
}

// return *GoCrabHttpRequest with specific method
func newGoCrabRequest(url, method string) *GoCrabHttpRequest {
	var resp http.Response
	req := http.Request{
		Method:     method,
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	return &GoCrabHttpRequest{url, &req, map[string]string{}, map[string]string{}, defaultSetting, &resp, nil}
}

// Get returns *GoCrabHttpRequest with GET method.
func Get(url string) *GoCrabHttpRequest {
	return newGoCrabRequest(url, "GET")
}

// Post returns *GoCrabHttpRequest with POST method.
func Post(url string) *GoCrabHttpRequest {
	return newGoCrabRequest(url, "POST")
}

// Put returns *GoCrabHttpRequest with PUT method.
func Put(url string) *GoCrabHttpRequest {
	return newGoCrabRequest(url, "PUT")
}

// Delete returns *GoCrabHttpRequest DELETE method.
func Delete(url string) *GoCrabHttpRequest {
	return newGoCrabRequest(url, "DELETE")
}

// Head returns *GoCrabHttpRequest with HEAD method.
func Head(url string) *GoCrabHttpRequest {
	return newGoCrabRequest(url, "HEAD")
}

// GoCrabHttpSettings
type GoCrabHttpSettings struct {
	Compress         string
	ShowDebug        bool
	UserAgent        string
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
	TlsClientConfig  *tls.Config
	Proxy            func(*http.Request) (*url.URL, error)
	Transport        http.RoundTripper
	EnableCookie     bool
}

// GoCrabHttpRequest provides more useful methods for requesting one url than http.Request.
type GoCrabHttpRequest struct {
	url     string
	req     *http.Request
	params  map[string]string
	files   map[string]string
	setting GoCrabHttpSettings
	resp    *http.Response
	body    []byte
}

// Change request settings
func (b *GoCrabHttpRequest) Setting(setting GoCrabHttpSettings) *GoCrabHttpRequest {
	b.setting = setting
	return b
}

// SetBasicAuth sets the request's Authorization header to use HTTP Basic Authentication with the provided username and password.
func (b *GoCrabHttpRequest) SetBasicAuth(username, password string) *GoCrabHttpRequest {
	b.req.SetBasicAuth(username, password)
	return b
}

// SetEnableCookie sets enable/disable cookiejar
func (b *GoCrabHttpRequest) SetEnableCookie(enable bool) *GoCrabHttpRequest {
	b.setting.EnableCookie = enable
	return b
}

// SetUserAgent sets User-Agent header field
func (b *GoCrabHttpRequest) SetUserAgent(useragent string) *GoCrabHttpRequest {
	b.setting.UserAgent = useragent
	return b
}

// Debug sets show debug or not when executing request.
func (b *GoCrabHttpRequest) Debug(isdebug bool) *GoCrabHttpRequest {
	b.setting.ShowDebug = isdebug
	return b
}

// SetTimeout sets connect time out and read-write time out for GoCrabRequest.
func (b *GoCrabHttpRequest) SetTimeout(connectTimeout, readWriteTimeout time.Duration) *GoCrabHttpRequest {
	b.setting.ConnectTimeout = connectTimeout
	b.setting.ReadWriteTimeout = readWriteTimeout
	return b
}

// SetTLSClientConfig sets tls connection configurations if visiting https url.
func (b *GoCrabHttpRequest) SetTLSClientConfig(config *tls.Config) *GoCrabHttpRequest {
	b.setting.TlsClientConfig = config
	return b
}

// Header add header item string in request.
func (b *GoCrabHttpRequest) Header(key, value string) *GoCrabHttpRequest {
	b.req.Header.Set(key, value)
	return b
}

// Set the protocol version for incoming requests.
// Client requests always use HTTP/1.1.
func (b *GoCrabHttpRequest) SetProtocolVersion(vers string) *GoCrabHttpRequest {
	if len(vers) == 0 {
		vers = "HTTP/1.1"
	}

	major, minor, ok := http.ParseHTTPVersion(vers)
	if ok {
		b.req.Proto = vers
		b.req.ProtoMajor = major
		b.req.ProtoMinor = minor
	}

	return b
}

// SetCookie add cookie into request.
func (b *GoCrabHttpRequest) SetCookie(cookie *http.Cookie) *GoCrabHttpRequest {
	b.req.Header.Add("Cookie", cookie.String())
	return b
}

// Set transport to
func (b *GoCrabHttpRequest) SetTransport(transport http.RoundTripper) *GoCrabHttpRequest {
	b.setting.Transport = transport
	return b
}

// Set http proxy
// example:
//
//	func(req *http.Request) (*url.URL, error) {
// 		u, _ := url.ParseRequestURI("http://127.0.0.1:8118")
// 		return u, nil
// 	}
func (b *GoCrabHttpRequest) SetProxy(proxy func(*http.Request) (*url.URL, error)) *GoCrabHttpRequest {
	b.setting.Proxy = proxy
	return b
}

func (b *GoCrabHttpRequest) SetCompress(compress string) *GoCrabHttpRequest {
	b.setting.Compress = compress
	return b
}

// Param adds query param in to request.
// params build query string as ?key1=value1&key2=value2...
func (b *GoCrabHttpRequest) Param(key, value string) *GoCrabHttpRequest {
	b.params[key] = value
	return b
}

func (b *GoCrabHttpRequest) PostFile(formname, filename string) *GoCrabHttpRequest {
	b.files[formname] = filename
	return b
}

// Body adds request raw body.
// it supports string and []byte.
func (b *GoCrabHttpRequest) Body(data interface{}) *GoCrabHttpRequest {
	switch t := data.(type) {
	case string:
		bf := bytes.NewBufferString(t)
		b.req.Body = ioutil.NopCloser(bf)
		b.req.ContentLength = int64(len(t))
	case []byte:
		bf := bytes.NewBuffer(t)
		b.req.Body = ioutil.NopCloser(bf)
		b.req.ContentLength = int64(len(t))
	}
	return b
}

func (b *GoCrabHttpRequest) buildUrl(paramBody string) {
	// build GET url with query string
	if b.req.Method == "GET" && len(paramBody) > 0 {
		if strings.Index(b.url, "?") != -1 {
			b.url += "&" + paramBody
		} else {
			b.url = b.url + "?" + paramBody
		}
		return
	}

	// build POST url and body
	if b.req.Method == "POST" && b.req.Body == nil {
		// with files
		if len(b.files) > 0 {
			pr, pw := io.Pipe()
			bodyWriter := multipart.NewWriter(pw)
			go func() {
				for formname, filename := range b.files {
					fileWriter, err := bodyWriter.CreateFormFile(formname, filename)
					if err != nil {
						log.Fatal(err)
					}
					fh, err := os.Open(filename)
					if err != nil {
						log.Fatal(err)
					}
					//iocopy
					_, err = io.Copy(fileWriter, fh)
					fh.Close()
					if err != nil {
						log.Fatal(err)
					}
				}
				for k, v := range b.params {
					bodyWriter.WriteField(k, v)
				}
				bodyWriter.Close()
				pw.Close()
			}()
			b.Header("Content-Type", bodyWriter.FormDataContentType())
			b.req.Body = ioutil.NopCloser(pr)
			return
		}

		// with params
		if len(paramBody) > 0 {

			if b.setting.Compress == compressGzip {
				var input = []byte(paramBody)
				var buf bytes.Buffer
				compressor, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
				if err != nil {

					b.Header("Content-Type", "application/x-www-form-urlencoded")
					b.Body(paramBody)

					return
				}
				compressor.Write(input)
				compressor.Close()

				b.Header("Content-Type", "text/plain;charset=utf-8")
				b.Header("Content-Encoding", "gzip")
				b.Header("Transfer-Encoding", "chunked")

				b.Body(buf.Bytes())

			} else {
				b.Header("Content-Type", "application/x-www-form-urlencoded")
				b.Body(paramBody)
			}
		}
	}
}

func (b *GoCrabHttpRequest) getResponse() (*http.Response, error) {
	if b.resp.StatusCode != 0 {
		return b.resp, nil
	}
	var paramBody string
	if len(b.params) > 0 {
		var buf bytes.Buffer
		for k, v := range b.params {
			buf.WriteString(url.QueryEscape(k))
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
			buf.WriteByte('&')
		}
		paramBody = buf.String()
		paramBody = paramBody[0 : len(paramBody)-1]
	}

	b.buildUrl(paramBody)
	url, err := url.Parse(b.url)
	if err != nil {
		return nil, err
	}

	b.req.URL = url

	trans := b.setting.Transport

	if trans == nil {
		// create default transport
		trans = &http.Transport{
			TLSClientConfig: b.setting.TlsClientConfig,
			Proxy:           b.setting.Proxy,
			Dial:            TimeoutDialer(b.setting.ConnectTimeout, b.setting.ReadWriteTimeout),
		}
	} else {
		// if b.transport is *http.Transport then set the settings.
		if t, ok := trans.(*http.Transport); ok {
			if t.TLSClientConfig == nil {
				t.TLSClientConfig = b.setting.TlsClientConfig
			}
			if t.Proxy == nil {
				t.Proxy = b.setting.Proxy
			}
			if t.Dial == nil {
				t.Dial = TimeoutDialer(b.setting.ConnectTimeout, b.setting.ReadWriteTimeout)
			}
		}
	}

	var jar http.CookieJar = nil
	if b.setting.EnableCookie {
		if defaultCookieJar == nil {
			createDefaultCookie()
		}
		jar = defaultCookieJar
	}

	client := &http.Client{
		Transport: trans,
		Jar:       jar,
	}

	if b.setting.UserAgent != "" && b.req.Header.Get("User-Agent") == "" {
		b.req.Header.Set("User-Agent", b.setting.UserAgent)
	}

	if b.setting.ShowDebug {
		dump, err := httputil.DumpRequest(b.req, true)
		if err != nil {
			println(err.Error())
		}
		println(string(dump))
	}

	resp, doErr := client.Do(b.req)
	if doErr != nil {
		return nil, doErr
	}
	b.resp = resp
	return resp, nil
}

// String returns the body string in response.
// it calls Response inner.
func (b *GoCrabHttpRequest) String() (string, error) {
	data, err := b.Bytes()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Bytes returns the body []byte in response.
// it calls Response inner.
func (b *GoCrabHttpRequest) Bytes() ([]byte, error) {
	if b.body != nil {
		return b.body, nil
	}
	resp, err := b.getResponse()

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode == 204 {
			return []byte(""), nil
		}
		return nil, errors.New("Response Code is not 200")
	}

	if resp.Body == nil {
		return nil, errors.New("Response is null")
	}
	b.body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b.body, nil
}

//只获取响应status
func (b *GoCrabHttpRequest) Status() (int, error) {
	if b.body != nil {
		return 0, nil
	}

	respObj, respErr := b.getResponse()

	if respErr != nil {
		return 0, respErr
	}

	defer respObj.Body.Close()

	if respObj.StatusCode == 200 {
		return respObj.StatusCode, nil
	}

	return 0, nil
}

// ToFile saves the body data in response to one file.
// it calls Response inner.
func (b *GoCrabHttpRequest) ToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := b.getResponse()
	if err != nil {
		return err
	}
	if resp.Body == nil {
		return nil
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// ToJson returns the map that marshals from the body bytes as json in response .
// it calls Response inner.
func (b *GoCrabHttpRequest) ToJson(v interface{}) error {
	data, err := b.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// ToXml returns the map that marshals from the body bytes as xml in response .
// it calls Response inner.
func (b *GoCrabHttpRequest) ToXml(v interface{}) error {
	data, err := b.Bytes()
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}

// Response executes request client gets response mannually.
func (b *GoCrabHttpRequest) Response() (*http.Response, error) {
	return b.getResponse()
}

// TimeoutDialer returns functions of connection dialer with timeout settings for http.Transport Dial field.
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}
