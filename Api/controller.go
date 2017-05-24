// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

package GoCrab

import (
	"errors"
	"github.com/CloudWise-OpenSource/GoCrab/Core/context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

//commonly used mime-types
const (
	applicationJson = "application/json"
	applicationXml  = "application/xml"
	textXml         = "text/xml"
)

var (
	// custom error when user stop request handler manually.
	USERSTOPRUN                                            = errors.New("User stop run")
	GlobalControllerRouter map[string][]ControllerComments = make(map[string][]ControllerComments) //pkgpath+controller:comments
)

// store the comment for the controller method
type ControllerComments struct {
	Method           string
	Router           string
	AllowHTTPMethods []string
	Params           []map[string]string
}

// Controller defines some basic http request handler operations, such as
// http context, template and view, session.
type Controller struct {
	Ctx            *context.Context
	Data           map[interface{}]interface{}
	controllerName string
	actionName     string
	TplNames       string
	Layout         string
	LayoutSections map[string]string // the key is the section name and the value is the template name
	TplExt         string
	gotofunc       string
	AppController  interface{}
	EnableRender   bool
	methodMapping  map[string]func() //method:routertree
}

// ControllerInterface is an interface to uniform all controller handler.
type ControllerInterface interface {
	Init(ct *context.Context, controllerName, actionName string, app interface{})
	Prepare()
	Get()
	Post()
	Delete()
	Put()
	Head()
	Patch()
	Options()
	Finish()
	Render() error
	HandlerFunc(fn string) bool
	URLMapping()
}

// Init generates default values of controller operations.
func (c *Controller) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
	c.Layout = ""
	c.TplNames = ""
	c.controllerName = controllerName
	c.actionName = actionName
	c.Ctx = ctx
	c.TplExt = "tpl"
	c.AppController = app
	c.EnableRender = true
	c.Data = ctx.Input.Data
	c.methodMapping = make(map[string]func())
}

// Prepare runs after Init before request function execution.
func (c *Controller) Prepare() {

}

// Finish runs after request function execution.
func (c *Controller) Finish() {

}

// Get adds a request function to handle GET request.
func (c *Controller) Get() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

// Post adds a request function to handle POST request.
func (c *Controller) Post() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

// Delete adds a request function to handle DELETE request.
func (c *Controller) Delete() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

// Put adds a request function to handle PUT request.
func (c *Controller) Put() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

// Head adds a request function to handle HEAD request.
func (c *Controller) Head() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

// Patch adds a request function to handle PATCH request.
func (c *Controller) Patch() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

// Options adds a request function to handle OPTIONS request.
func (c *Controller) Options() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

// call function fn
func (c *Controller) HandlerFunc(fnname string) bool {
	if v, ok := c.methodMapping[fnname]; ok {
		v()
		return true
	} else {
		return false
	}
}

// URLMapping register the internal Controller router.
func (c *Controller) URLMapping() {
}

func (c *Controller) Mapping(method string, fn func()) {
	c.methodMapping[method] = fn
}

// Render sends the response with rendered template bytes as text/html type.
func (c *Controller) Render() error {
	if !c.EnableRender {
		return nil
	}
	rb, err := c.RenderBytes()

	if err != nil {
		return err
	} else {
		c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
		c.Ctx.Output.Body(rb)
	}
	return nil
}

// RenderString returns the rendered template string. Do not send out response.
func (c *Controller) RenderString() (string, error) {
	b, e := c.RenderBytes()
	return string(b), e
}

// RenderBytes returns the bytes of rendered template string. Do not send out response.
func (c *Controller) RenderBytes() ([]byte, error) {
	//if the controller has set layout, then first get the tplname's content set the content to the layout

	return nil, nil

	//if c.Layout != "" {
	//	if c.TplNames == "" {
	//		c.TplNames = strings.ToLower(c.controllerName) + "/" + strings.ToLower(c.actionName) + "." + c.TplExt
	//	}
	//	newbytes := bytes.NewBufferString("")
	//	if _, ok := Templates[c.TplNames]; !ok {
	//		panic("can't find templatefile in the path:" + c.TplNames)
	//	}
	//	err := Templates[c.TplNames].ExecuteTemplate(newbytes, c.TplNames, c.Data)
	//	if err != nil {
	//		Trace("template Execute err:", err)
	//		return nil, err
	//	}
	//	tplcontent, _ := ioutil.ReadAll(newbytes)
	//	c.Data["LayoutContent"] = template.HTML(string(tplcontent))

	//	if c.LayoutSections != nil {
	//		for sectionName, sectionTpl := range c.LayoutSections {
	//			if sectionTpl == "" {
	//				c.Data[sectionName] = ""
	//				continue
	//			}

	//			sectionBytes := bytes.NewBufferString("")
	//			err = Templates[sectionTpl].ExecuteTemplate(sectionBytes, sectionTpl, c.Data)
	//			if err != nil {
	//				Trace("template Execute err:", err)
	//				return nil, err
	//			}
	//			sectionContent, _ := ioutil.ReadAll(sectionBytes)
	//			c.Data[sectionName] = template.HTML(string(sectionContent))
	//		}
	//	}

	//	ibytes := bytes.NewBufferString("")
	//	err = Templates[c.Layout].ExecuteTemplate(ibytes, c.Layout, c.Data)
	//	if err != nil {
	//		Trace("template Execute err:", err)
	//		return nil, err
	//	}
	//	icontent, _ := ioutil.ReadAll(ibytes)
	//	return icontent, nil
	//} else {
	//	if c.TplNames == "" {
	//		c.TplNames = strings.ToLower(c.controllerName) + "/" + strings.ToLower(c.actionName) + "." + c.TplExt
	//	}
	//	ibytes := bytes.NewBufferString("")
	//	if _, ok := Templates[c.TplNames]; !ok {
	//		panic("can't find templatefile in the path:" + c.TplNames)
	//	}
	//	err := Templates[c.TplNames].ExecuteTemplate(ibytes, c.TplNames, c.Data)
	//	if err != nil {
	//		Trace("template Execute err:", err)
	//		return nil, err
	//	}
	//	icontent, _ := ioutil.ReadAll(ibytes)
	//	return icontent, nil
	//}
}

// Redirect sends the redirection response to url with status code.
func (c *Controller) Redirect(url string, code int) {
	c.Ctx.Redirect(code, url)
}

// Aborts stops controller handler and show the error data if code is defined in ErrorMap or code string.
func (c *Controller) Abort(code string) {
	status, err := strconv.Atoi(code)
	if err != nil {
		status = 200
	}
	c.CustomAbort(status, code)
}

// CustomAbort stops controller handler and show the error data, it's similar Aborts, but support status code and body.
func (c *Controller) CustomAbort(status int, body string) {
	c.Ctx.ResponseWriter.WriteHeader(status)
	// first panic from ErrorMaps, is is user defined error functions.
	if _, ok := ErrorMaps[body]; ok {
		panic(body)
	}
	// last panic user string
	c.Ctx.ResponseWriter.Write([]byte(body))
	panic(USERSTOPRUN)
}

// StopRun makes panic of USERSTOPRUN error and go to recover function if defined.
func (c *Controller) StopRun() {
	panic(USERSTOPRUN)
}

// UrlFor does another controller handler in this request function.
// it goes to this controller method if endpoint is not clear.
func (c *Controller) UrlFor(endpoint string, values ...interface{}) string {

	return ""
}

// ServeJson sends a json response with encoding charset.
func (c *Controller) ServeJson(encoding ...bool) {
	var hasIndent bool
	var hasencoding bool
	if RunMode == RUNMODE_PROD {
		hasIndent = false
	} else {
		hasIndent = true
	}
	if len(encoding) > 0 && encoding[0] == true {
		hasencoding = true
	}
	c.Ctx.Output.Json(c.Data["json"], hasIndent, hasencoding)
}

//ServeString sends a string response.
func (c *Controller) ServeString(data []byte) {
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body(data)
}

// ServeJsonp sends a jsonp response.
func (c *Controller) ServeJsonp() {
	var hasIndent bool
	if RunMode == RUNMODE_PROD {
		hasIndent = false
	} else {
		hasIndent = true
	}
	c.Ctx.Output.Jsonp(c.Data["jsonp"], hasIndent)
}

// ServeXml sends xml response.
func (c *Controller) ServeXml() {
	var hasIndent bool
	if RunMode == RUNMODE_PROD {
		hasIndent = false
	} else {
		hasIndent = true
	}
	c.Ctx.Output.Xml(c.Data["xml"], hasIndent)
}

// ServeFormatted serve Xml OR Json, depending on the value of the Accept header
func (c *Controller) ServeFormatted() {
	accept := c.Ctx.Input.Header("Accept")
	switch accept {
	case applicationJson:
		c.ServeJson()
	case applicationXml, textXml:
		c.ServeXml()
	default:
		c.ServeJson()
	}
}

// Input returns the input data map from POST or PUT request body and query string.
func (c *Controller) Input() url.Values {
	if c.Ctx.Request.Form == nil {
		c.Ctx.Request.ParseForm()
	}
	return c.Ctx.Request.Form
}

// ParseForm maps input data map to obj struct.
//func (c *Controller) ParseForm(obj interface{}) error {
//return ParseForm(c.Input(), obj)
//}

// GetString returns the input value by key string or the default value while it's present and input is blank
func (c *Controller) GetString(key string, def ...string) string {
	var defv string
	if len(def) > 0 {
		defv = def[0]
	}

	if v := c.Ctx.Input.Query(key); v != "" {
		return v
	} else {
		return defv
	}
}

// GetStrings returns the input string slice by key string or the default value while it's present and input is blank
// it's designed for multi-value input field such as checkbox(input[type=checkbox]), multi-selection.
func (c *Controller) GetStrings(key string, def ...[]string) []string {
	var defv []string
	if len(def) > 0 {
		defv = def[0]
	}

	f := c.Input()
	if f == nil {
		return defv
	}

	vs := f[key]
	if len(vs) > 0 {
		return vs
	} else {
		return defv
	}
}

// GetInt returns input as an int or the default value while it's present and input is blank
func (c *Controller) GetInt(key string, def ...int) (int, error) {
	var defv int
	if len(def) > 0 {
		defv = def[0]
	}

	if strv := c.Ctx.Input.Query(key); strv != "" {
		return strconv.Atoi(strv)
	} else {
		return defv, nil
	}
}

// GetInt8 return input as an int8 or the default value while it's present and input is blank
func (c *Controller) GetInt8(key string, def ...int8) (int8, error) {
	var defv int8
	if len(def) > 0 {
		defv = def[0]
	}

	if strv := c.Ctx.Input.Query(key); strv != "" {
		i64, err := strconv.ParseInt(strv, 10, 8)
		i8 := int8(i64)
		return i8, err
	} else {
		return defv, nil
	}
}

// GetInt16 returns input as an int16 or the default value while it's present and input is blank
func (c *Controller) GetInt16(key string, def ...int16) (int16, error) {
	var defv int16
	if len(def) > 0 {
		defv = def[0]
	}

	if strv := c.Ctx.Input.Query(key); strv != "" {
		i64, err := strconv.ParseInt(strv, 10, 16)
		i16 := int16(i64)

		return i16, err
	} else {
		return defv, nil
	}
}

// GetInt32 returns input as an int32 or the default value while it's present and input is blank
func (c *Controller) GetInt32(key string, def ...int32) (int32, error) {
	var defv int32
	if len(def) > 0 {
		defv = def[0]
	}

	if strv := c.Ctx.Input.Query(key); strv != "" {
		i64, err := strconv.ParseInt(c.Ctx.Input.Query(key), 10, 32)
		i32 := int32(i64)
		return i32, err
	} else {
		return defv, nil
	}
}

// GetInt64 returns input value as int64 or the default value while it's present and input is blank.
func (c *Controller) GetInt64(key string, def ...int64) (int64, error) {
	var defv int64
	if len(def) > 0 {
		defv = def[0]
	}

	if strv := c.Ctx.Input.Query(key); strv != "" {
		return strconv.ParseInt(strv, 10, 64)
	} else {
		return defv, nil
	}
}

// GetBool returns input value as bool or the default value while it's present and input is blank.
func (c *Controller) GetBool(key string, def ...bool) (bool, error) {
	var defv bool
	if len(def) > 0 {
		defv = def[0]
	}

	if strv := c.Ctx.Input.Query(key); strv != "" {
		return strconv.ParseBool(strv)
	} else {
		return defv, nil
	}
}

// GetFloat returns input value as float64 or the default value while it's present and input is blank.
func (c *Controller) GetFloat(key string, def ...float64) (float64, error) {
	var defv float64
	if len(def) > 0 {
		defv = def[0]
	}

	if strv := c.Ctx.Input.Query(key); strv != "" {
		return strconv.ParseFloat(c.Ctx.Input.Query(key), 64)
	} else {
		return defv, nil
	}
}

// GetFile returns the file data in file upload field named as key.
// it returns the first one of multi-uploaded files.
func (c *Controller) GetFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Ctx.Request.FormFile(key)
}

// SaveToFile saves uploaded file to new path.
// it only operates the first one of mutil-upload form file field.
func (c *Controller) SaveToFile(fromfile, tofile string) error {
	file, _, err := c.Ctx.Request.FormFile(fromfile)
	if err != nil {
		return err
	}
	defer file.Close()
	f, err := os.OpenFile(tofile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, file)
	return nil
}

// IsAjax returns this request is ajax or not.
func (c *Controller) IsAjax() bool {
	return c.Ctx.Input.IsAjax()
}

// GetSecureCookie returns decoded cookie value from encoded browser cookie values.
func (c *Controller) GetSecureCookie(Secret, key string) (string, bool) {
	return c.Ctx.GetSecureCookie(Secret, key)
}

// SetSecureCookie puts value into cookie after encoded the value.
func (c *Controller) SetSecureCookie(Secret, name, value string, others ...interface{}) {
	c.Ctx.SetSecureCookie(Secret, name, value, others...)
}

// GetControllerAndAction gets the executing controller name and action name.
func (c *Controller) GetControllerAndAction() (controllerName, actionName string) {
	return c.controllerName, c.actionName
}
