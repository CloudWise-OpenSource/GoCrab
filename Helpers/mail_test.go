// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

package Helpers

import "testing"

func TestMail(t *testing.T) {
	config := `{"username":"Neeke@gmail.com","password":"Neeke","host":"smtp.gmail.com","port":587}`
	mail := NewEMail(config)
	if mail.Username != "Neeke@gmail.com" {
		t.Fatal("email parse get username error")
	}
	if mail.Password != "Neeke" {
		t.Fatal("email parse get password error")
	}
	if mail.Host != "smtp.gmail.com" {
		t.Fatal("email parse get host error")
	}
	if mail.Port != 587 {
		t.Fatal("email parse get port error")
	}
	mail.To = []string{"xiemengjun@gmail.com"}
	mail.From = "Neeke@gmail.com"
	mail.Subject = "hi, just from GoCrab!"
	mail.Text = "Text Body is, of course, supported!"
	mail.HTML = "<h1>Fancy Html is supported, too!</h1>"
	mail.AttachFile("/Users/Neeke/github/github.com/CloudWise-OpenSource/GoCrab/Core.go")
	mail.Send()
}
