// Copyright 2016 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

package config

import (
	"os"
	"testing"
)

var inicontext = `
;comment one
#comment two
appname = GoCrabApi
httpport = 8080
mysqlport = 3600
PI = 3.1415976
runmode = "dev"
autorender = false
copyrequestbody = true
[demo]
key1="asta"
key2 = "xie"
CaseInsensitive = true
peers = one;two;three
`

func TestIni(t *testing.T) {
	f, err := os.Create("testini.conf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(inicontext)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove("testini.conf")
	iniconf, err := NewConfig("ini", "testini.conf")
	if err != nil {
		t.Fatal(err)
	}
	if iniconf.String("appname") != "GoCrabApi" {
		t.Fatal("appname not equal to GoCrabApi")
	}
	if port, err := iniconf.Int("httpport"); err != nil || port != 8080 {
		t.Error(port)
		t.Fatal(err)
	}
	if port, err := iniconf.Int64("mysqlport"); err != nil || port != 3600 {
		t.Error(port)
		t.Fatal(err)
	}
	if pi, err := iniconf.Float("PI"); err != nil || pi != 3.1415976 {
		t.Error(pi)
		t.Fatal(err)
	}
	if iniconf.String("runmode") != "dev" {
		t.Fatal("runmode not equal to dev")
	}
	if v, err := iniconf.Bool("autorender"); err != nil || v != false {
		t.Error(v)
		t.Fatal(err)
	}
	if v, err := iniconf.Bool("copyrequestbody"); err != nil || v != true {
		t.Error(v)
		t.Fatal(err)
	}
	if err = iniconf.Set("name", "Neeke"); err != nil {
		t.Fatal(err)
	}
	if iniconf.String("name") != "Neeke" {
		t.Fatal("get name error")
	}
	if iniconf.String("demo::key1") != "asta" {
		t.Fatal("get demo.key1 error")
	}
	if iniconf.String("demo::key2") != "xie" {
		t.Fatal("get demo.key2 error")
	}
	if v, err := iniconf.Bool("demo::caseinsensitive"); err != nil || v != true {
		t.Fatal("get demo.caseinsensitive error")
	}

	if data := iniconf.Strings("demo::peers"); len(data) != 3 {
		t.Fatal("get strings error", data)
	} else if data[0] != "one" {
		t.Fatal("get first params error not equat to one")
	}

}
