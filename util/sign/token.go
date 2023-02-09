package sign

import (
	"crypto/rand"
	"errors"
)

func GenerateToken(length int, chars []byte) (string, error) {
	if length == 0 {
		return "", nil
	}
	clen := len(chars)
	if clen < 2 || clen > 256 {
		return "", errors.New("wrong charset length for GenerateToken")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			return "", errors.New("Error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue // Skip this number to avoid modulo bias.
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b), nil
			}
		}
	}
}
