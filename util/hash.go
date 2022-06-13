package util

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/creachadair/cityhash"
)

func Hash(key string) int {
	return int(cityhash.Hash32([]byte(key)))
}

func Md5(data string) string {
	md := md5.New()
	md.Write([]byte(data))
	return hex.EncodeToString(md.Sum([]byte("")))
}
