// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

package GoCrab

import (
	"errors"
	"fmt"
	"github.com/CloudWise-OpenSource/GoCrab/Core/config"
	"github.com/CloudWise-OpenSource/GoCrab/Core/logs"
	"github.com/CloudWise-OpenSource/GoCrab/Helpers"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

var (
	CrabApp             *App // GoCrab application
	AppName             string
	AppPath             string
	workPath            string
	AppConfigPath       string
	EnableHttpListen    bool
	HttpAddr            string
	HttpPort            int
	ListenTCP4          bool
	EnableHttpTLS       bool
	HttpsPort           int
	HttpCertFile        string
	HttpKeyFile         string
	RecoverPanic        bool // flag of auto recover panic
	AutoRender          bool // flag of render template automatically
	AppConfig           *GoCrabAppConfig
	RunMode             string // run mode, "dev" or "prod"
	UseFcgi             bool
	UseStdIo            bool
	MaxMemory           int64
	EnableGzip          bool // flag of enable gzip
	DirectoryIndex      bool // flag of display directory index. default is false.
	HttpServerTimeOut   int64
	ErrorsShow          bool   // flag of show errors in page. if true, show error and trace info in page rendered with error template.
	CopyRequestBody     bool   // flag of copy raw request body in context.
	GoCrabServerName    string // GoCrab server name exported in response header.
	AppConfigProvider   string // config provider
	RouterCaseSensitive bool   // router case sensitive default is true
	AccessLogs          bool   // print access logs, default is false
	StaticDir           map[string]string
	ViewsPath           string
	ResourcePath        string
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
	// create GoCrab application
	CrabApp = NewApp()

	workPath, _ = os.Getwd()
	workPath, _ = filepath.Abs(workPath)
	// initialize default configurations
	AppPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	AppConfigPath = filepath.Join(AppPath, "conf", "app.ini")

	if workPath != AppPath {
		if Helpers.FileExists(AppConfigPath) {
			os.Chdir(AppPath)
		} else {
			AppConfigPath = filepath.Join(workPath, "conf", "app.ini")
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

	DirectoryIndex = true

	// set this to 0.0.0.0 to make this app available to externally
	EnableHttpListen = true //default enable http Listen

	HttpAddr = ""
	HttpPort = 26789

	HttpsPort = 10443

	AppName = "GoCrab"

	RunMode = RUNMODE_PROD

	AutoRender = true

	RecoverPanic = true

	UseFcgi = false
	UseStdIo = false

	CopyRequestBody = true

	MaxMemory = 1 << 26 //64MB

	EnableGzip = false

	HttpServerTimeOut = 0

	ErrorsShow = true

	GoCrabServerName = "GoCrab/" + VERSION

	RouterCaseSensitive = true

	runtime.GOMAXPROCS(runtime.NumCPU())

	StaticDir = make(map[string]string)

	ResourcePath = "resource"
	ViewsPath = "views"

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
	AppConfig, err = newAppConfig(AppConfigProvider, AppConfigPath)
	if err != nil {
		return err
	}
	envRunMode := os.Getenv("GOCRAB_RUNMODE")
	// set the runmode first
	if envRunMode != "" {
		RunMode = envRunMode
	} else if runmode := AppConfig.String("RunMode"); runmode != "" {
		RunMode = runmode
	}

	HttpAddr = AppConfig.String("HttpAddr")

	if v, err := AppConfig.Int("HttpPort"); err == nil {
		HttpPort = v
	}

	if v, err := AppConfig.Bool("ListenTCP4"); err == nil {
		ListenTCP4 = v
	}

	if v, err := AppConfig.Bool("EnableHttpListen"); err == nil {
		EnableHttpListen = v
	}

	if maxmemory, err := AppConfig.Int64("MaxMemory"); err == nil {
		MaxMemory = maxmemory
	}

	if appname := AppConfig.String("AppName"); appname != "" {
		AppName = appname
	}

	if autorender, err := AppConfig.Bool("AutoRender"); err == nil {
		AutoRender = autorender
	}

	if autorecover, err := AppConfig.Bool("RecoverPanic"); err == nil {
		RecoverPanic = autorecover
	}

	if usefcgi, err := AppConfig.Bool("UseFcgi"); err == nil {
		UseFcgi = usefcgi
	}

	if enablegzip, err := AppConfig.Bool("EnableGzip"); err == nil {
		EnableGzip = enablegzip
	}

	if directoryindex, err := AppConfig.Bool("DirectoryIndex"); err == nil {
		DirectoryIndex = directoryindex
	}

	if timeout, err := AppConfig.Int64("HttpServerTimeOut"); err == nil {
		HttpServerTimeOut = timeout
	}

	if errorsshow, err := AppConfig.Bool("ErrorsShow"); err == nil {
		ErrorsShow = errorsshow
	}

	if copyrequestbody, err := AppConfig.Bool("CopyRequestBody"); err == nil {
		CopyRequestBody = copyrequestbody
	}

	if httptls, err := AppConfig.Bool("EnableHttpTLS"); err == nil {
		EnableHttpTLS = httptls
	}

	if httpsport, err := AppConfig.Int("HttpsPort"); err == nil {
		HttpsPort = httpsport
	}

	if certfile := AppConfig.String("HttpCertFile"); certfile != "" {
		HttpCertFile = certfile
	}

	if keyfile := AppConfig.String("HttpKeyFile"); keyfile != "" {
		HttpKeyFile = keyfile
	}

	if serverName := AppConfig.String("GoCrabServerName"); serverName != "" {
		GoCrabServerName = serverName
	}

	if casesensitive, err := AppConfig.Bool("RouterCaseSensitive"); err == nil {
		RouterCaseSensitive = casesensitive
	}

	if sd := AppConfig.String("StaticDir"); sd != "" {
		for k := range StaticDir {
			delete(StaticDir, k)
		}
		sds := strings.Split(sd, ",")

		for _, v := range sds {
			StaticDir["/"+v] = v
		}
	}

	if views := AppConfig.String("ViewsPath"); views != "" {
		ViewsPath = views
	}
	StaticDir["/"+ViewsPath] = ViewsPath
	StaticDir["/"+ResourcePath] = ResourcePath

	return nil
}

func Config(returnType, key string, defaultVal interface{}) (value interface{}, err error) {
	switch returnType {
	case "String":
		value = AppConfig.String(key)
	case "Bool":
		value, err = AppConfig.Bool(key)
	case "Int":
		value, err = AppConfig.Int(key)
	case "Int64":
		value, err = AppConfig.Int64(key)
	case "Float":
		value, err = AppConfig.Float(key)
	case "DIY":
		value, err = AppConfig.DIY(key)
	default:
		err = errors.New("Config keys must be of type String, Bool, Int, Int64, Float, or DIY!")
	}

	if err != nil {
		if reflect.TypeOf(returnType) != reflect.TypeOf(defaultVal) {
			err = errors.New("defaultVal type does not match returnType!")
		} else {
			value, err = defaultVal, nil
		}
	} else if reflect.TypeOf(value).Kind() == reflect.String {
		if value == "" {
			if reflect.TypeOf(defaultVal).Kind() != reflect.String {
				err = errors.New("defaultVal type must be a String if the returnType is a String")
			} else {
				value = defaultVal.(string)
			}
		}
	}

	return
}
