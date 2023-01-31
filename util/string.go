package util

import (
	"math/rand"
	"regexp"
)

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits  = 6
	defaultRandLen = 8
	letterIdxMask  = 1<<letterIdxBits - 1
	letterIdxMax   = 63 / letterIdxBits
)

func Randn(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func RandString(n int) string {
	byteArr := make([]byte, n)
	for i := 0; i < n; i++ {
		str := rand.Intn(26) + 65
		byteArr[i] = byte(str)
	}
	return string(byteArr)
}

func IsMobile(phone string) bool {
	regRuler := "^1[0-9]{1}\\d{9}$"
	reg := regexp.MustCompile(regRuler)
	return reg.MatchString(phone)
}
