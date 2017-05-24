// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

package logs

import (
	"testing"
)

func TestConn(t *testing.T) {
	log := NewLogger(1000)
	log.SetLogger("conn", `{"net":"tcp","addr":":7020"}`)
	log.Informational("informational")
}
