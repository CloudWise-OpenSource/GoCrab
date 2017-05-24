// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

// Usage:
//
//	import "github.com/CloudWise-OpenSource/GoCrab/Core/crtype"
//  crtype.SetDesKey("NeekeGao")
//
//	encodeResult := crtype.DesEncode("aaabbbccc")
//  decodeResult := crtype.DesDecode(encodeResult)
//
package crtype

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"strings"
)

const SIGN_PLUS = "+"
const SIGN_PLUS_REPLACE = "**"

const SIGN_EQUATE = "="
const SIGN_EQUATE_REPLACE = "!!"

var (
	DesKey string
	desKey []byte
)

func SetDesKey(key string) bool {
	DesKey = key
	desKey = []byte(key)

	return true
}

func GetDesKey() string {
	return DesKey
}

func DesEncode(str string) string {
	result, err := desEncrypt([]byte(str), desKey)
	if err != nil {
		panic(err)
	}
	sresult := base64.StdEncoding.EncodeToString(result)
	sresult = strings.Replace(sresult, SIGN_PLUS, SIGN_PLUS_REPLACE, -1)
	sresult = strings.Replace(sresult, SIGN_EQUATE, SIGN_EQUATE_REPLACE, -1)

	return sresult
}

func DesDecode(str string) string {
	sresult := strings.Replace(str, SIGN_PLUS_REPLACE, SIGN_PLUS, -1)
	sresult = strings.Replace(sresult, SIGN_EQUATE_REPLACE, SIGN_EQUATE, -1)
	result, _ := base64.StdEncoding.DecodeString(sresult)

	origData, err := desDecrypt(result, desKey)
	if err != nil {
		panic(err)
	}

	return string(origData)
}

func desEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = pKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)

	return crypted, nil
}

func desDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))

	blockMode.CryptBlocks(origData, crypted)
	origData = pKCS5UnPadding(origData)

	return origData, nil
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])

	return origData[:(length - unpadding)]
}
