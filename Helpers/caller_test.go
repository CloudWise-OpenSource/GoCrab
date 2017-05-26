// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

package Helpers

import (
	"strings"
	"testing"
)

func TestGetFuncName(t *testing.T) {
	name := GetFuncName(TestGetFuncName)
	t.Log(name)
	if !strings.HasSuffix(name, ".TestGetFuncName") {
		t.Error("get func name error")
	}
}
