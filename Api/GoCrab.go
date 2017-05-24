// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

package GoCrab

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type hookfunc func() error //hook function to run
var hooks []hookfunc       //hook function slice to store the hookfunc

// Router adds a patterned controller handler to CrabApp.
// it's an alias method of App.Router.
// usage:
//  simple router
//  GoCrab.Router("/admin", &admin.UserController{})
//  GoCrab.Router("/admin/index", &admin.ArticleController{})
//
//  regex router
//
//  GoCrab.Router("/api/:id([0-9]+)", &controllers.RController{})
//
//  custom rules
//  GoCrab.Router("/api/list",&RestController{},"*:ListFood")
//  GoCrab.Router("/api/create",&RestController{},"post:CreateFood")
//  GoCrab.Router("/api/update",&RestController{},"put:UpdateFood")
//  GoCrab.Router("/api/delete",&RestController{},"delete:DeleteFood")
func Router(rootpath string, c ControllerInterface, mappingMethods ...string) *App {
	CrabApp.Handlers.Add(rootpath, c, mappingMethods...)
	return CrabApp
}

// Router add list from
// usage:
// GoCrab.Include(&BankAccount{}, &OrderController{},&RefundController{},&ReceiptController{})
// type BankAccount struct{
//   GoCrab.Controller
// }
//
// register the function
// func (b *BankAccount)Mapping(){
//  b.Mapping("ShowAccount" , b.ShowAccount)
//  b.Mapping("ModifyAccount", b.ModifyAccount)
//}
//
// //@router /account/:id  [get]
// func (b *BankAccount) ShowAccount(){
//    //logic
// }
//
//
// //@router /account/:id  [post]
// func (b *BankAccount) ModifyAccount(){
//    //logic
// }
//
// the comments @router url methodlist
// url support all the function Router's pattern
// methodlist [get post head put delete options *]
func Include(cList ...ControllerInterface) *App {
	CrabApp.Handlers.Include(cList...)
	return CrabApp
}

// RESTRouter adds a restful controller handler to CrabApp.
// its' controller implements GoCrab.ControllerInterface and
// defines a param "pattern/:objectId" to visit each resource.
func RESTRouter(rootpath string, c ControllerInterface) *App {
	Router(rootpath, c)
	Router(path.Join(rootpath, ":objectId"), c)
	return CrabApp
}

// AutoRouter adds defined controller handler to CrabApp.
// it's same to App.AutoRouter.
// if GoCrab.AddAuto(&MainContorlller{}) and MainController has methods List and Page,
// visit the url /main/list to exec List function or /main/page to exec Page function.
func AutoRouter(c ControllerInterface) *App {
	CrabApp.Handlers.AddAuto(c)
	return CrabApp
}

// AutoPrefix adds controller handler to CrabApp with prefix.
// it's same to App.AutoRouterWithPrefix.
// if GoCrab.AutoPrefix("/admin",&MainContorlller{}) and MainController has methods List and Page,
// visit the url /admin/main/list to exec List function or /admin/main/page to exec Page function.
func AutoPrefix(prefix string, c ControllerInterface) *App {
	CrabApp.Handlers.AddAutoPrefix(prefix, c)
	return CrabApp
}

// register router for Get method
// usage:
//    GoCrab.Get("/", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Get(rootpath string, f FilterFunc) *App {
	CrabApp.Handlers.Get(rootpath, f)
	return CrabApp
}

// register router for Post method
// usage:
//    GoCrab.Post("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Post(rootpath string, f FilterFunc) *App {
	CrabApp.Handlers.Post(rootpath, f)
	return CrabApp
}

// register router for Delete method
// usage:
//    GoCrab.Delete("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Delete(rootpath string, f FilterFunc) *App {
	CrabApp.Handlers.Delete(rootpath, f)
	return CrabApp
}

// register router for Put method
// usage:
//    GoCrab.Put("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Put(rootpath string, f FilterFunc) *App {
	CrabApp.Handlers.Put(rootpath, f)
	return CrabApp
}

// register router for Head method
// usage:
//    GoCrab.Head("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Head(rootpath string, f FilterFunc) *App {
	CrabApp.Handlers.Head(rootpath, f)
	return CrabApp
}

// register router for Options method
// usage:
//    GoCrab.Options("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Options(rootpath string, f FilterFunc) *App {
	CrabApp.Handlers.Options(rootpath, f)
	return CrabApp
}

// register router for Patch method
// usage:
//    GoCrab.Patch("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Patch(rootpath string, f FilterFunc) *App {
	CrabApp.Handlers.Patch(rootpath, f)
	return CrabApp
}

// register router for all method
// usage:
//    GoCrab.Any("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Any(rootpath string, f FilterFunc) *App {
	CrabApp.Handlers.Any(rootpath, f)
	return CrabApp
}

// register router for own Handler
// usage:
//    GoCrab.Handler("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Handler(rootpath string, h http.Handler, options ...interface{}) *App {
	CrabApp.Handlers.Handler(rootpath, h, options...)
	return CrabApp
}

// InsertFilter adds a FilterFunc with pattern condition and action constant.
// The pos means action constant including
// GoCrab.BeforeStatic, GoCrab.BeforeRouter, GoCrab.BeforeExec, GoCrab.AfterExec and GoCrab.FinishRouter.
// The bool params is for setting the returnOnOutput value (false allows multiple filters to execute)
func InsertFilter(pattern string, pos int, filter FilterFunc, params ...bool) *App {
	CrabApp.Handlers.InsertFilter(pattern, pos, filter, params...)
	return CrabApp
}

// The hookfunc will run in GoCrab.Run()
// such as middlerware start, buildtemplate, admin start
func AddAPPStartHook(hf hookfunc) {
	hooks = append(hooks, hf)
}

// Run GoCrab application.
// GoCrab.Run() default run on HttpPort
// GoCrab.Run(":8089")
// GoCrab.Run("127.0.0.1:8089")
func Run(params ...string) {
	if len(params) > 0 && params[0] != "" {
		strs := strings.Split(params[0], ":")
		if len(strs) > 0 && strs[0] != "" {
			HttpAddr = strs[0]
		}
		if len(strs) > 1 && strs[1] != "" {
			HttpPort, _ = strconv.Atoi(strs[1])
		}
	}
	initBeforeHttpRun()

	CrabApp.Run()
}

func initBeforeHttpRun() {
	// if AppConfigPath not In the conf/app.conf reParse config
	if AppConfigPath != filepath.Join(AppPath, "conf", "app.conf") {
		err := ParseConfig()
		if err != nil && AppConfigPath != filepath.Join(workPath, "conf", "app.conf") {
			// configuration is critical to app, panic here if parse failed
			panic(err)
		}
	}

	//init mime
	//AddAPPStartHook(initMime)

	// do hooks function
	for _, hk := range hooks {
		err := hk()
		if err != nil {
			panic(err)
		}
	}

	registerDefaultErrorHandler()
}

// this function is for test package init
func TestGoCrabInit(apppath string) {
	AppPath = apppath
	RunMode = "test"
	AppConfigPath = filepath.Join(AppPath, "conf", "app.conf")
	err := ParseConfig()
	if err != nil && !os.IsNotExist(err) {
		// for init if doesn't have app.conf will not panic
		Info(err)
	}
	os.Chdir(AppPath)
	initBeforeHttpRun()
}

func init() {
	hooks = make([]hookfunc, 0)
}
