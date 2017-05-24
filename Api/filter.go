// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

package GoCrab

import "github.com/CloudWise-OpenSource/GoCrab/Core/context"

// FilterFunc defines filter function type.
type FilterFunc func(*context.Context)

// FilterRouter defines filter operation before controller handler execution.
// it can match patterned url and do filter function when action arrives.
type FilterRouter struct {
	filterFunc     FilterFunc
	tree           *Tree
	pattern        string
	returnOnOutput bool
}

// ValidRouter check current request is valid for this filter.
// if matched, returns parsed params in this request by defined filter router pattern.
func (f *FilterRouter) ValidRouter(router string) (bool, map[string]string) {
	isok, params := f.tree.Match(router)
	if isok == nil {
		return false, nil
	}
	if isok, ok := isok.(bool); ok {
		return isok, params
	} else {
		return false, nil
	}
}
