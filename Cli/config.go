// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

package GoCrab

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/CloudWise-OpenSource/GoCrab/Core/config"
	"github.com/CloudWise-OpenSource/GoCrab/Core/logs"
	"github.com/CloudWise-OpenSource/GoCrab/Helpers"
)

var (
	CrabApp             *App // GoCrab application
	AppName             string
	AppPath             string
	UseCLI              bool
	workPath            string
	AppConfigPath       string
	RecoverPanic        bool // flag of auto recover panic
	AppConfig           *GoCrabAppConfig
	RunMode             string // run mode, "dev" or "prod"
	MaxMemory           int64
	AppConfigProvider   string // config provider
	RouterCaseSensitive bool   // router case sensitive default is true
)

type GoCrabAppConfig struct {
	innerConfig config.ConfigContainer
}

func newAppConfig(AppConfigProvider, AppConfigPath string) (*GoCrabAppConfig, error) {
	ac, err := config.NewConfig(AppConfigProvider, AppConfigPath)
	if err != nil {
		return nil, err
	}
	rac := &GoCrabAppConfig{ac}
	return rac, nil
}

func (b *GoCrabAppConfig) Set(key, val string) error {
	return b.innerConfig.Set(key, val)
}

func (b *GoCrabAppConfig) String(key string) string {
	v := b.innerConfig.String(RunMode + "::" + key)
	if v == "" {
		return b.innerConfig.String(key)
	}
	return v
}

func (b *GoCrabAppConfig) Strings(key string) []string {
	v := b.innerConfig.Strings(RunMode + "::" + key)
	if v[0] == "" {
		return b.innerConfig.Strings(key)
	}
	return v
}

func (b *GoCrabAppConfig) Int(key string) (int, error) {
	v, err := b.innerConfig.Int(RunMode + "::" + key)
	if err != nil {
		return b.innerConfig.Int(key)
	}
	return v, nil
}

func (b *GoCrabAppConfig) Int64(key string) (int64, error) {
	v, err := b.innerConfig.Int64(RunMode + "::" + key)
	if err != nil {
		return b.innerConfig.Int64(key)
	}
	return v, nil
}

func (b *GoCrabAppConfig) Bool(key string) (bool, error) {
	v, err := b.innerConfig.Bool(RunMode + "::" + key)
	if err != nil {
		return b.innerConfig.Bool(key)
	}
	return v, nil
}

func (b *GoCrabAppConfig) Float(key string) (float64, error) {
	v, err := b.innerConfig.Float(RunMode + "::" + key)
	if err != nil {
		return b.innerConfig.Float(key)
	}
	return v, nil
}

func (b *GoCrabAppConfig) DefaultString(key string, defaultval string) string {
	return b.innerConfig.DefaultString(key, defaultval)
}

func (b *GoCrabAppConfig) DefaultStrings(key string, defaultval []string) []string {
	return b.innerConfig.DefaultStrings(key, defaultval)
}

func (b *GoCrabAppConfig) DefaultInt(key string, defaultval int) int {
	return b.innerConfig.DefaultInt(key, defaultval)
}

func (b *GoCrabAppConfig) DefaultInt64(key string, defaultval int64) int64 {
	return b.innerConfig.DefaultInt64(key, defaultval)
}

func (b *GoCrabAppConfig) DefaultBool(key string, defaultval bool) bool {
	return b.innerConfig.DefaultBool(key, defaultval)
}

func (b *GoCrabAppConfig) DefaultFloat(key string, defaultval float64) float64 {
	return b.innerConfig.DefaultFloat(key, defaultval)
}

func (b *GoCrabAppConfig) DIY(key string) (interface{}, error) {
	return b.innerConfig.DIY(key)
}

func (b *GoCrabAppConfig) GetSection(section string) (map[string]string, error) {
	return b.innerConfig.GetSection(section)
}

func (b *GoCrabAppConfig) SaveConfigFile(filename string) error {
	return b.innerConfig.SaveConfigFile(filename)
}

func init() {
	workPath, _ = os.Getwd()
	workPath, _ = filepath.Abs(workPath)
	// initialize default configurations
	AppPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	AppConfigPath = filepath.Join(AppPath, "conf", "app.ini")

	if workPath != AppPath {
		if Helpers.FileExists(AppConfigPath) {
			os.Chdir(AppPath)
		} else {
			AppConfigPath = filepath.Join(workPath, "conf", "app.inin")
		}
	}
	AppConfigProvider = "ini"

	if !Helpers.FileExists(AppConfigPath) {
		AppConfigPath = filepath.Join(AppPath, "conf", "app.conf")
		if workPath != AppPath {
			if Helpers.FileExists(AppConfigPath) {
				os.Chdir(AppPath)
			} else {
				AppConfigPath = filepath.Join(workPath, "conf", "app.conf")
			}
		}
		AppConfigProvider = "json"
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	// init Logger
	Logger = logs.NewLogger(10000)
	err := Logger.SetLogger("console", "")
	if err != nil {
		fmt.Println("init console log error:", err)
	}
	SetLogFuncCall(true)

	err = ParseConfig()
	if err != nil && os.IsNotExist(err) {
		// for init if doesn't have app.conf will not panic
		ac := config.NewFakeConfig()
		AppConfig = &GoCrabAppConfig{ac}
		Warning(err)
	}
}

// ParseConfig parsed default config file.
// now only support ini, next will support json.
func ParseConfig() (err error) {

	// create GoCrab application
	CrabApp = NewApp()

	AppName = "GoCrab"

	RunMode = RUNMODE_PROD

	MaxMemory = 1 << 26 //64MB

	RouterCaseSensitive = true

	envRunMode := os.Getenv("GOCRAB_RUNMODE")
	// set the runmode first

	AppConfig, err = newAppConfig(AppConfigProvider, AppConfigPath)
	if err != nil {
		Error("config is error, but it's still running by default config")
		Error("config path -> ", AppConfigPath)
		Error("Error -> ", err)
		return err
	}

	if envRunMode != "" {
		RunMode = envRunMode
	} else if runmode := AppConfig.String("RunMode"); runmode != "" {
		RunMode = runmode
	}

	if maxmemory, err := AppConfig.Int64("MaxMemory"); err == nil {
		MaxMemory = maxmemory
	}

	if appname := AppConfig.String("AppName"); appname != "" {
		AppName = appname
	}

	return nil
}
