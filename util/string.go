package util

import "math/rand"

func RandString(n int) string {
	byteArr := make([]byte, n)
	for i := 0; i < n; i++ {
		str := rand.Intn(26) + 65
		byteArr[i] = byte(str)
	}
	return string(byteArr)
}
