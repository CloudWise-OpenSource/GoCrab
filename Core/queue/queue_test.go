// Copyright 2016 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

package queue

import (
	"fmt"
	"testing"
)

//func Fibonacci(n int64) int64 {
//	if n < 2 {
//		return n
//	}
//	return Fibonacci(n-1) + Fibonacci(n-2)
//}

//func TestFibonacci(t *testing.T) {
//	r := Fibonacci(10)
//	if r != 55 {
//		t.Errorf("Fibonacci(10) failed. Got %d, expected 55.", r)
//	}
//}

type Data struct {
	ObjID      string
	ApiEnum    string
	RoutingKey string
	Content    string
	Method     string
	Direct     int
}

var QTest Queue

func TestQueue(t *testing.T) {
	QTest, err := InitQueue(10)

	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	var Ob Data
	Ob.Content = "sss"
	Ob.Direct = 44

	QTest.Push(Ob, 1)
	QTest.Push("asdf")
	QTest.Push("asdf", 1)
	QTest.Push("asdf", 1, 2, 3)

	iLen := QTest.Len(1)
	fmt.Println(iLen)

	iAllLen := QTest.Lens()
	fmt.Println(iAllLen)

	a, _ := QTest.Pop(1)

	C := a.Value.(Data)

	fmt.Println(C.Content)
	fmt.Println(C.Direct)
	fmt.Println(C.Method)

	//	var bb interface{}
	//	bb = "dddddd"

	//	fmt.Println(bb.(string))

	//	switch btype := a.(type) {
	//	case string:
	//		fmt.Println("is string")
	//		fmt.Println(btype)
	//		break
	//	}
	//	fmt.Println(bb.(type))
}
