package util

import (
	"math/rand"
	"regexp"
)

func RandString(n int) string {
	byteArr := make([]byte, n)
	for i := 0; i < n; i++ {
		str := rand.Intn(26) + 65
		byteArr[i] = byte(str)
	}
	return string(byteArr)
}

func IsMobile(phone string) bool {
	regRuler := "^1[345789]{1}\\d{9}$"
	reg := regexp.MustCompile(regRuler)
	return reg.MatchString(phone)
}
