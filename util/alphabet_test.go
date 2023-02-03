package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDedupe(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{
			"01010101010101", "02",
		},
		{
			"abcabcfoo", "abcfo",
		},
	}
	should := require.New(t)
	for _, test := range tests {
		in := strings.Join(debupe(strings.Split(test.in, "")), "")
		should.Equal(in, test.out)
	}
}

func TestAlphabet_Index(t *testing.T) {
	str := newAlphabet(DefaultAlphabet)
	should := require.New(t)
	idx, err := str.Index(rune('z'))
	should.Nil(err)
	should.Equal(idx, int64(56))
}

func TestAlphabet_IndexZero(t *testing.T) {
	str := newAlphabet(DefaultAlphabet)
	should := require.New(t)
	idx, err := str.Index(rune('2'))
	should.Nil(err)
	should.Equal(idx, int64(0))
}
