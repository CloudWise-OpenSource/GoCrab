// Copyright 2016 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

package context

import (
	"fmt"
	"net/http"
	"testing"
)

func TestParse(t *testing.T) {
	r, _ := http.NewRequest("GET", "/?id=123&isok=true&ft=1.2&ol[0]=1&ol[1]=2&ul[]=str&ul[]=array&user.Name=Neeke", nil)
	GoCrabInput := NewInput(r)
	GoCrabInput.ParseFormOrMulitForm(1 << 20)

	var id int
	err := GoCrabInput.Bind(&id, "id")
	if id != 123 || err != nil {
		t.Fatal("id should has int value")
	}
	fmt.Println(id)

	var isok bool
	err = GoCrabInput.Bind(&isok, "isok")
	if !isok || err != nil {
		t.Fatal("isok should be true")
	}
	fmt.Println(isok)

	var float float64
	err = GoCrabInput.Bind(&float, "ft")
	if float != 1.2 || err != nil {
		t.Fatal("float should be equal to 1.2")
	}
	fmt.Println(float)

	ol := make([]int, 0, 2)
	err = GoCrabInput.Bind(&ol, "ol")
	if len(ol) != 2 || err != nil || ol[0] != 1 || ol[1] != 2 {
		t.Fatal("ol should has two elements")
	}
	fmt.Println(ol)

	ul := make([]string, 0, 2)
	err = GoCrabInput.Bind(&ul, "ul")
	if len(ul) != 2 || err != nil || ul[0] != "str" || ul[1] != "array" {
		t.Fatal("ul should has two elements")
	}
	fmt.Println(ul)

	type User struct {
		Name string
	}
	user := User{}
	err = GoCrabInput.Bind(&user, "user")
	if err != nil || user.Name != "Neeke" {
		t.Fatal("user should has name")
	}
	fmt.Println(user)
}

func TestSubDomain(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://www.example.com/?id=123&isok=true&ft=1.2&ol[0]=1&ol[1]=2&ul[]=str&ul[]=array&user.Name=Neeke", nil)
	GoCrabInput := NewInput(r)

	subdomain := GoCrabInput.SubDomains()
	if subdomain != "www" {
		t.Fatal("Subdomain parse error, got" + subdomain)
	}

	r, _ = http.NewRequest("GET", "http://localhost/", nil)
	GoCrabInput.Request = r
	if GoCrabInput.SubDomains() != "" {
		t.Fatal("Subdomain parse error, should be empty, got " + GoCrabInput.SubDomains())
	}

	r, _ = http.NewRequest("GET", "http://aa.bb.example.com/", nil)
	GoCrabInput.Request = r
	if GoCrabInput.SubDomains() != "aa.bb" {
		t.Fatal("Subdomain parse error, got " + GoCrabInput.SubDomains())
	}

	/* TODO Fix this
	r, _ = http.NewRequest("GET", "http://127.0.0.1/", nil)
	GoCrabInput.Request = r
	if GoCrabInput.SubDomains() != "" {
		t.Fatal("Subdomain parse error, got " + GoCrabInput.SubDomains())
	}
	*/

	r, _ = http.NewRequest("GET", "http://example.com/", nil)
	GoCrabInput.Request = r
	if GoCrabInput.SubDomains() != "" {
		t.Fatal("Subdomain parse error, got " + GoCrabInput.SubDomains())
	}

	r, _ = http.NewRequest("GET", "http://aa.bb.cc.dd.example.com/", nil)
	GoCrabInput.Request = r
	if GoCrabInput.SubDomains() != "aa.bb.cc.dd" {
		t.Fatal("Subdomain parse error, got " + GoCrabInput.SubDomains())
	}
}
