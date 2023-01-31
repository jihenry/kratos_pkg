package sm4

import (
	"encoding/base64"
	"errors"
	"fmt"
)

type ECB interface {
	ECBDecrypt(key []byte) (plainText []byte, err error)
	ECBEncrypt(key []byte) (cipherText []byte, err error)
	SetCipherText(cipherText string) ECB
	PKCS5Padding(src []byte) ECB
}

type ecbSM4 struct {
	plainText  []byte
	cipherText []byte
	blockSize  int
}

var (
	_ ECB = (*ecbSM4)(nil)
)

func NewECBSm4(blockSize int) ECB {
	sm4 := new(ecbSM4)
	sm4.blockSize = blockSize
	return sm4
}

func (s *ecbSM4) PKCS5Padding(src []byte) ECB {
	padding := s.blockSize - len(src)%s.blockSize
	ret := make([]byte, len(src)+padding)
	for k := range src {
		ret[k] = src[k]
	}
	for i := 0; i < padding; i++ {
		ret[len(src)+i] = byte(padding)
	}
	s.plainText = ret
	return s
}

func (s *ecbSM4) SetCipherText(cipherText string) ECB {
	if text, err := base64.StdEncoding.DecodeString(cipherText); err == nil {
		s.cipherText = text
	}
	return s
}

func (s *ecbSM4) ECBDecrypt(key []byte) (plainText []byte, err error) {
	cipherTextLen := len(s.cipherText)
	if cipherTextLen == 0 {
		return nil, fmt.Errorf("input not full blocks:%d", cipherTextLen)
	}
	if cipherTextLen%s.blockSize != 0 {
		return nil, errors.New("input not full blocks")
	}
	c, err := NewSm4Cipher(key, false)
	if err != nil {
		return nil, err
	}
	plainText = make([]byte, cipherTextLen)
	for i := 0; i < cipherTextLen; i += s.blockSize {
		if err = c.Encrypt(plainText[i:i+s.blockSize], s.cipherText[i:i+s.blockSize]); err != nil {
			return
		}
	}
	return
}

func (s *ecbSM4) ECBEncrypt(key []byte) (cipherText []byte, err error) {
	cipherTextLen := len(s.plainText)
	if cipherTextLen == 0 {
		return nil, fmt.Errorf("input not full blocks:%d", cipherTextLen)
	}
	if cipherTextLen%s.blockSize != 0 {
		return cipherText, errors.New("input not full blocks")
	}
	c, err := NewSm4Cipher(key, true)
	if err != nil {
		return cipherText, err
	}
	cipherText = make([]byte, cipherTextLen)
	for i := 0; i < cipherTextLen; i += s.blockSize {
		if err = c.Encrypt(cipherText[i:i+s.blockSize], s.plainText[i:i+s.blockSize]); err != nil {
			return
		}
	}
	return
}
