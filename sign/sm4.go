package sign

import (
	"errors"
	"strconv"
)

const (
	BlockSize = 16
	KeySize   = 16
)

var sBox = [256]rune{
	0xd6, 0x90, 0xe9, 0xfe, 0xcc, 0xe1, 0x3d, 0xb7,
	0x16, 0xb6, 0x14, 0xc2, 0x28, 0xfb, 0x2c, 0x05,
	0x2b, 0x67, 0x9a, 0x76, 0x2a, 0xbe, 0x04, 0xc3,
	0xaa, 0x44, 0x13, 0x26, 0x49, 0x86, 0x06, 0x99,
	0x9c, 0x42, 0x50, 0xf4, 0x91, 0xef, 0x98, 0x7a,
	0x33, 0x54, 0x0b, 0x43, 0xed, 0xcf, 0xac, 0x62,
	0xe4, 0xb3, 0x1c, 0xa9, 0xc9, 0x08, 0xe8, 0x95,
	0x80, 0xdf, 0x94, 0xfa, 0x75, 0x8f, 0x3f, 0xa6,
	0x47, 0x07, 0xa7, 0xfc, 0xf3, 0x73, 0x17, 0xba,
	0x83, 0x59, 0x3c, 0x19, 0xe6, 0x85, 0x4f, 0xa8,
	0x68, 0x6b, 0x81, 0xb2, 0x71, 0x64, 0xda, 0x8b,
	0xf8, 0xeb, 0x0f, 0x4b, 0x70, 0x56, 0x9d, 0x35,
	0x1e, 0x24, 0x0e, 0x5e, 0x63, 0x58, 0xd1, 0xa2,
	0x25, 0x22, 0x7c, 0x3b, 0x01, 0x21, 0x78, 0x87,
	0xd4, 0x00, 0x46, 0x57, 0x9f, 0xd3, 0x27, 0x52,
	0x4c, 0x36, 0x02, 0xe7, 0xa0, 0xc4, 0xc8, 0x9e,
	0xea, 0xbf, 0x8a, 0xd2, 0x40, 0xc7, 0x38, 0xb5,
	0xa3, 0xf7, 0xf2, 0xce, 0xf9, 0x61, 0x15, 0xa1,
	0xe0, 0xae, 0x5d, 0xa4, 0x9b, 0x34, 0x1a, 0x55,
	0xad, 0x93, 0x32, 0x30, 0xf5, 0x8c, 0xb1, 0xe3,
	0x1d, 0xf6, 0xe2, 0x2e, 0x82, 0x66, 0xca, 0x60,
	0xc0, 0x29, 0x23, 0xab, 0x0d, 0x53, 0x4e, 0x6f,
	0xd5, 0xdb, 0x37, 0x45, 0xde, 0xfd, 0x8e, 0x2f,
	0x03, 0xff, 0x6a, 0x72, 0x6d, 0x6c, 0x5b, 0x51,
	0x8d, 0x1b, 0xaf, 0x92, 0xbb, 0xdd, 0xbc, 0x7f,
	0x11, 0xd9, 0x5c, 0x41, 0x1f, 0x10, 0x5a, 0xd8,
	0x0a, 0xc1, 0x31, 0x88, 0xa5, 0xcd, 0x7b, 0xbd,
	0x2d, 0x74, 0xd0, 0x12, 0xb8, 0xe5, 0xb4, 0xb0,
	0x89, 0x69, 0x97, 0x4a, 0x0c, 0x96, 0x77, 0x7e,
	0x65, 0xb9, 0xf1, 0x09, 0xc5, 0x6e, 0xc6, 0x84,
	0x18, 0xf0, 0x7d, 0xec, 0x3a, 0xdc, 0x4d, 0x20,
	0x79, 0xee, 0x5f, 0x3e, 0xd7, 0xcb, 0x39, 0x48,
}

var cK = [32]int{
	0x00070e15, 0x1c232a31, 0x383f464d, 0x545b6269,
	0x70777e85, 0x8c939aa1, 0xa8afb6bd, 0xc4cbd2d9,
	0xe0e7eef5, 0xfc030a11, 0x181f262d, 0x343b4249,
	0x50575e65, 0x6c737a81, 0x888f969d, 0xa4abb2b9,
	0xc0c7ced5, 0xdce3eaf1, 0xf8ff060d, 0x141b2229,
	0x30373e45, 0x4c535a61, 0x686f767d, 0x848b9299,
	0xa0a7aeb5, 0xbcc3cad1, 0xd8dfe6ed, 0xf4fb0209,
	0x10171e25, 0x2c333a41, 0x484f565d, 0x646b7279,
}

var fK = [4]int{0xa3b1bac6, 0x56aa3350, 0x677d9197, 0xb27022dc}

type KeySizeError int

func (k KeySizeError) Error() string {
	return "sm4: invalid key size " + strconv.Itoa(int(k))
}

type sm4Cipher struct {
	enc []int32
	dec []int32
}

type Block interface {
	BlockSize() int
	Encrypt(dst []byte, src []byte)
	Decrypt(dst []byte, src []byte)
}

func SWAP(sk []int32, i int) {
	t := sk[i]
	sk[int32(i)] = sk[int32(31-i)]
	sk[int32(31-i)] = t
}

func NewCipher(key []byte, check bool) (Block, error) {
	n := len(key)
	if n != KeySize {
		return nil, KeySizeError(n)
	}
	c := new(sm4Cipher)
	c.enc = encKey(key)
	if !check {
		for i := 0; i < 16; i++ {
			SWAP(c.enc, i)
		}
	}
	return c, nil
}

func (c *sm4Cipher) BlockSize() int {
	return BlockSize
}

func (c *sm4Cipher) Encrypt(dst []byte, src []byte) {
	if len(src) < BlockSize {
		panic("sm4: input not full block")
	}
	if len(dst) < BlockSize {
		panic("sm4: output not full block")
	}
	processBlock(c.enc, src, dst)
}

func (c *sm4Cipher) Decrypt(dst []byte, src []byte) {
	if len(src) < BlockSize {
		panic("sm4: input not full block")
	}
	if len(dst) < BlockSize {
		panic("sm4: output not full block")
	}
	processBlock(c.enc, src, dst)
}

func sm4Sbox(inch byte) byte {
	i := inch & 0xFF
	return byte(int8(sBox[rune(i)]))
}

func sm4CalciRK(ka int32) int32 {
	var (
		bb, rk int32
	)
	a, b := make([]byte, 4), make([]byte, 4)
	putUlongBe(ka, a, 0)
	b[0] = sm4Sbox(a[0])
	b[1] = sm4Sbox(a[1])
	b[2] = sm4Sbox(a[2])
	b[3] = sm4Sbox(a[3])
	bb = RuneInt32(b, 0)
	rk = bb ^ ROTL(int(bb), 13) ^ ROTL(int(bb), 23)
	return rk
}

func ROTL(x, n int) int32 {
	return int32(SHL(x, n) | x>>(32-n))
}

func SHL(x, n int) int {
	return x & 0xFFFFFFFF << n
}

func encKey(key []byte) []int32 {
	var mK = make([]int32, 4)
	mK[0] = Int32(key, 0)
	mK[1] = Int32(key, 4)
	mK[2] = Int32(key, 8)
	mK[3] = Int32(key, 12)
	var k = make([]int32, 36)
	k[0] = mK[0] ^ int32(fK[0])
	k[1] = mK[1] ^ int32(fK[1])
	k[2] = mK[2] ^ int32(fK[2])
	k[3] = mK[3] ^ int32(fK[3])

	var rk = make([]int32, 32)
	for i := 0; i < 32; i++ {
		kkk := k[i+1] ^ k[i+2] ^ k[i+3] ^ int32(cK[i])
		k[i+4] = int32(k[i] ^ sm4CalciRK(kkk))
		rk[i] = k[i+4]
	}
	return rk
}

func Int32(b []byte, i int) int32 {
	return int32(b[i]&0xff)<<24 | int32(b[i+1]&0xff)<<16 | int32(b[i+2]&0xff)<<8 | int32(int(b[i+3]&0xff)&0xffffffff)
}

func RuneInt32(b []byte, i int) int32 {
	return int32(b[i]&0xff)<<24 | int32(b[i+1]&0xff)<<16 | int32(b[i+2]&0xff)<<8 | int32(int(b[i+3]&0xff)&0xffffffff)
}

func processBlock(rk []int32, in []byte, out []byte) {
	var x = make([]int32, 36)
	x[0] = Int32(in, 0)
	x[1] = Int32(in, 4)
	x[2] = Int32(in, 8)
	x[3] = Int32(in, 12)
	for i := 0; i < 32; i++ {
		x[i+4] = sm4F(x[i], x[i+1], x[i+2], x[i+3], rk[i])
	}
	putUlongBe(x[35], out, 0)
	putUlongBe(x[34], out, 4)
	putUlongBe(x[33], out, 8)
	putUlongBe(x[32], out, 12)

}

func putUlongBe(n int32, b []byte, i int) {
	kk := int8(0xFF & (n >> 24))
	b[i] = byte(kk)
	b[i+1] = byte(int8(0xFF & (n >> 16)))
	b[i+2] = byte(int8(0xFF & (n >> 8)))
	b[i+3] = byte(int8(0xFF & n))
}

func sm4Lt(ka int32) int32 {
	var (
		c  int32
		bb int32
	)
	a := make([]byte, 4)
	b := make([]byte, 4)
	putUlongBe(ka, a, 0)
	b[0] = sm4Sbox(a[0])
	b[1] = sm4Sbox(a[1])
	b[2] = sm4Sbox(a[2])
	b[3] = sm4Sbox(a[3])
	bb = RuneInt32(b, 0)
	c = bb ^ ROTL(int(bb), 2) ^ ROTL(int(bb), 10) ^ ROTL(int(bb), 18) ^ ROTL(int(bb), 24)
	return c
}

func sm4F(x0, x1, x2, x3, rk int32) int32 {
	return x0 ^ sm4Lt(x1^x2^x3^rk)
}

// ECBEncrypt 加密
func ECBEncrypt(key, plainText []byte) (cipherText []byte, err error) {
	plainTextLen := len(plainText)
	if plainTextLen%BlockSize != 0 {
		return nil, errors.New("input not full blocks")
	}
	c, err := NewCipher(key, true)
	if err != nil {
		return nil, err
	}
	cipherText = make([]byte, plainTextLen)
	for i := 0; i < plainTextLen; i += BlockSize {
		c.Encrypt(cipherText[i:i+BlockSize], plainText[i:i+BlockSize])
	}
	return cipherText, nil
}

// ECBDecrypt 解密
func ECBDecrypt(key []byte, cipherText []byte) (plainText []byte, err error) {
	cipherTextLen := len(cipherText)
	if cipherTextLen%BlockSize != 0 {
		return nil, errors.New("input not full blocks")
	}
	c, err := NewCipher(key, false)
	if err != nil {
		return nil, err
	}
	plainText = make([]byte, cipherTextLen)
	for i := 0; i < cipherTextLen; i += BlockSize {
		c.Decrypt(plainText[i:i+BlockSize], cipherText[i:i+BlockSize])
	}
	return plainText, nil
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	ret := make([]byte, len(src)+padding)
	for k := range src {
		ret[k] = src[k]
	}
	for i := 0; i < padding; i++ {
		ret[len(src)+i] = byte(padding)
	}
	return ret
}
