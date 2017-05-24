// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

// Usage:
//
//	import "github.com/CloudWise-OpenSource/GoCrab/Core/context"
//
//	ctx := context.Context{Request:req,ResponseWriter:rw}
//
package context

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Http request context struct including GoCrabInput, GoCrabOutput, http.Request and http.ResponseWriter.
// GoCrabInput and GoCrabOutput provides some api to operate request and response more easily.
type Context struct {
	Input          *GoCrabInput
	Output         *GoCrabOutput
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

// Redirect does redirection to localurl with http header status code.
// It sends http response header directly.
func (ctx *Context) Redirect(status int, localurl string) {
	ctx.Output.Header("Location", localurl)
	ctx.ResponseWriter.WriteHeader(status)
}

// Abort stops this request.
// if GoCrab.ErrorMaps exists, panic body.
func (ctx *Context) Abort(status int, body string) {
	ctx.ResponseWriter.WriteHeader(status)
	panic(body)
}

// Write string to response body.
// it sends response body.
func (ctx *Context) WriteString(content string) {
	ctx.ResponseWriter.Write([]byte(content))
}

// Get cookie from request by a given key.
// It's alias of GoCrabInput.Cookie.
func (ctx *Context) GetCookie(key string) string {
	return ctx.Input.Cookie(key)
}

// Set cookie for response.
// It's alias of GoCrabOutput.Cookie.
func (ctx *Context) SetCookie(name string, value string, others ...interface{}) {
	ctx.Output.Cookie(name, value, others...)
}

// Get secure cookie from request by a given key.
func (ctx *Context) GetSecureCookie(Secret, key string) (string, bool) {
	val := ctx.Input.Cookie(key)
	if val == "" {
		return "", false
	}

	parts := strings.SplitN(val, "|", 3)

	if len(parts) != 3 {
		return "", false
	}

	vs := parts[0]
	timestamp := parts[1]
	sig := parts[2]

	h := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)

	if fmt.Sprintf("%02x", h.Sum(nil)) != sig {
		return "", false
	}
	res, _ := base64.URLEncoding.DecodeString(vs)
	return string(res), true
}

// Set Secure cookie for response.
func (ctx *Context) SetSecureCookie(Secret, name, value string, others ...interface{}) {
	vs := base64.URLEncoding.EncodeToString([]byte(value))
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)
	sig := fmt.Sprintf("%02x", h.Sum(nil))
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")
	ctx.Output.Cookie(name, cookie, others...)
}
