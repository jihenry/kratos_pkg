package util

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/go-kratos/kratos/v2/log"
)

type SHA1withRSA struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (s *SHA1withRSA) SetPriKey(pkey []byte) {
	block, _ := pem.Decode(pkey)
	if block == nil {
		log.Error("pem decode private key err")
		return
	}
	private, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Errorf("ParsePKCS8PrivateKey err: %v", err)
		return
	}
	s.privateKey = private.(*rsa.PrivateKey)
}

func (s *SHA1withRSA) SetPublicKey(pKey []byte) {
	block, _ := pem.Decode(pKey)
	if block == nil {
		log.Error("pem decode public key err")
		return
	}
	public, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Errorf("parse publicKey err: %v", err)
		return
	}
	s.publicKey = public.(*rsa.PublicKey)
}

func (s *SHA1withRSA) Sign(data []byte) (string, error) {
	hashed := sha1Hash(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA1, hashed)
	if err != nil {
		log.Infof("Error from signing: %v", err)
		return "", err
	}
	log.Infof("signature: %x", signature)
	return base64.StdEncoding.EncodeToString(signature), nil
}

func (s *SHA1withRSA) CheckSign(sign string, data []byte) bool {
	hashed := sha1Hash(data)
	sigBytes, _ := base64.StdEncoding.DecodeString(sign)
	err := rsa.VerifyPKCS1v15(s.publicKey, crypto.SHA1, hashed[:], sigBytes)
	if err != nil {
		log.Errorf("verify sign err: %v", err)
		return false
	}
	log.Info("check sign true")
	return true
}

func (s *SHA1withRSA) Encrypt(data []byte) (string, error) {
	inputLen := len(data)
	// 分段加密单次需要减掉padding长度，PKCS1为11.
	offSet, once := 0, s.publicKey.Size()-11
	buffer := bytes.Buffer{}
	for offSet < inputLen {
		endIndex := offSet + once
		if endIndex > inputLen {
			endIndex = inputLen
		}
		byteOnce, err := rsa.EncryptPKCS1v15(rand.Reader, s.publicKey, data[offSet:endIndex])
		if err != nil {
			log.Errorf("encrypt err: %v", err)
			return "", err
		}
		buffer.Write(byteOnce)
		offSet = endIndex
	}
	res := base64.StdEncoding.EncodeToString(buffer.Bytes())
	log.Info("encrypt: ", res)
	return res, nil

}

func (s *SHA1withRSA) Decrypt(data []byte) (string, error) {
	inputLen := len(data)
	offSet, keySize := 0, s.privateKey.Size()
	buffer := bytes.Buffer{}
	for offSet < inputLen {
		endIndex := offSet + keySize
		if endIndex > inputLen {
			endIndex = inputLen
		}
		byteOnce, err := rsa.DecryptPKCS1v15(rand.Reader, s.privateKey, data[offSet:endIndex])
		if err != nil {
			log.Errorf("decrypt err: %v", err)
			return "", err
		}
		buffer.Write(byteOnce)
		offSet = endIndex
	}
	res := buffer.String()
	log.Info("decrypt: ", res)
	return res, nil
}

func sha1Hash(data []byte) []byte {
	h := crypto.Hash.New(crypto.SHA1)
	h.Write(data)
	return h.Sum(nil)
}
