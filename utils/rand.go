package utils

import (
	"math/rand"
	"time"
)

var (
	allChars    = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lowerChars  = []byte("abcdefghijklmnopqrstuvwxyz")
	upperChars  = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numberChars = []byte("0123456789")
)

func randChars(l int, chs []byte) string {
	rand.Seed(time.Now().UnixNano())
	var bs []byte
	for i := 0; i < l; i++ {
		bs = append(bs, chs[rand.Intn(len(chs))])
	}
	return string(bs)
}

// RandString 生成随机字符串
func RandString(l int) string {
	return randChars(l, allChars)
}

// RandNumberString 生成随机数字字符串
func RandNumberString(l int) string {
	return randChars(l, numberChars)
}

// RandString 生成小写字母随机字符串
func RandLowerString(l int) string {
	return randChars(l, lowerChars)
}

// RandUpperString 生成大写字母随机字符串
func RandUpperString(l int) string {
	return randChars(l, upperChars)
}
