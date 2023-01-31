package uuid

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

type base57 struct {
	alphabet alphabet
}

func (b base57) Encode(u uuid.UUID) string {
	var num big.Int
	num.SetString(strings.Replace(u.String(), "-", "", 4), 16)

	length := math.Ceil(math.Log(math.Pow(2, 128)) / math.Log(float64(b.alphabet.Length())))

	return b.numToString(&num, int(length))
}

func (b *base57) numToString(number *big.Int, padToLen int) string {
	var (
		out   []rune
		digit *big.Int
	)

	alphaLen := big.NewInt(b.alphabet.Length())

	zero := new(big.Int)
	for number.Cmp(zero) > 0 {
		number, digit = new(big.Int).DivMod(number, alphaLen, new(big.Int))
		out = append(out, b.alphabet.char[digit.Int64()])
	}
	if padToLen > 0 {
		remainder := math.Max(float64(padToLen-len(out)), 0)
		out = append(out, []rune(strings.Repeat(string(b.alphabet.char[0]), int(remainder)))...)
	}
	reverse(out)

	return string(out)
}

func (b base57) Decode(s string) (uuid.UUID, error) {
	str, err := b.stringToNum(s)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(str)
}

func (b *base57) stringToNum(s string) (string, error) {
	n := big.NewInt(0)
	for _, char := range s {
		n.Mul(n, big.NewInt(b.alphabet.Length()))
		index, err := b.alphabet.Index(char)
		if err != nil {
			return "", err
		}
		n.Add(n, big.NewInt(index))
	}
	if n.BitLen() > 128 {
		return "", fmt.Errorf("number is out of range (need a 128-bit value)")
	}
	return fmt.Sprintf("%032x", n), nil
}

func reverse(a []rune) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}
