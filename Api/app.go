// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

package GoCrab

import (
	"fmt"
	"github.com/CloudWise-OpenSource/GoCrab/Helpers"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"time"
)

// App defines GoCrab application with a new PatternServeMux.
type App struct {
	Handlers *ControllerRegistor
	Server   *http.Server
}

// NewApp returns a new GoCrab application.
func NewApp() *App {
	cr := NewControllerRegister()
	app := &App{Handlers: cr, Server: &http.Server{}}
	return app
}

// Run GoCrab application.
func (app *App) Run() {
	addr := HttpAddr

	if HttpPort != 0 {
		addr = fmt.Sprintf("%s:%d", HttpAddr, HttpPort)
	}

	var (
		err error
		l   net.Listener
	)
	endRunning := make(chan bool, 1)

	if UseFcgi {
		if UseStdIo {
			err = fcgi.Serve(nil, app.Handlers) // standard I/O
			if err == nil {
				Logger.Info("Use FCGI via standard I/O")
			} else {
				Logger.Info("Cannot use FCGI via standard I/O", err)
			}
		} else {
			if HttpPort == 0 {
				// remove the Socket file before start
				if Helpers.FileExists(addr) {
					os.Remove(addr)
				}
				l, err = net.Listen("unix", addr)
			} else {
				l, err = net.Listen("tcp", addr)
			}
			if err != nil {
				Logger.Critical("Listen: ", err)
			}
			err = fcgi.Serve(l, app.Handlers)
		}
	} else {
		app.Server.Addr = addr
		app.Server.Handler = app.Handlers
		app.Server.ReadTimeout = time.Duration(HttpServerTimeOut) * time.Second
		app.Server.WriteTimeout = time.Duration(HttpServerTimeOut) * time.Second

		if EnableHttpTLS {
			go func() {
				time.Sleep(20 * time.Microsecond)
				if HttpsPort != 0 {
					app.Server.Addr = fmt.Sprintf("%s:%d", HttpAddr, HttpsPort)
				}
				Logger.Info("https server Running on %s", app.Server.Addr)
				err := app.Server.ListenAndServeTLS(HttpCertFile, HttpKeyFile)
				if err != nil {
					Logger.Critical("ListenAndServeTLS: ", err)
					time.Sleep(100 * time.Microsecond)
					endRunning <- true
				}
			}()
		}

		if EnableHttpListen {
			go func() {
				app.Server.Addr = addr
				Logger.Info("http server Running on %s", app.Server.Addr)
				if ListenTCP4 && HttpAddr == "" {
					ln, err := net.Listen("tcp4", app.Server.Addr)
					if err != nil {
						Logger.Critical("ListenAndServe: ", err)
						time.Sleep(100 * time.Microsecond)
						endRunning <- true
						return
					}
					err = app.Server.Serve(ln)
					if err != nil {
						Logger.Critical("ListenAndServe: ", err)
						time.Sleep(100 * time.Microsecond)
						endRunning <- true
						return
					}
				} else {
					err := app.Server.ListenAndServe()
					if err != nil {
						Logger.Critical("ListenAndServe: ", err)
						time.Sleep(100 * time.Microsecond)
						endRunning <- true
					}
				}
			}()
		}
	}

	<-endRunning
}
