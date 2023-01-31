package uuid

import (
	"fmt"
	"sort"
	"strings"
)

const (
	DefaultAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

type alphabet struct {
	char [57]rune
	len  int64
}

func newAlphabet(s string) alphabet {
	strs := debupe(strings.Split(s, ""))
	if len(strs) != 57 {
		panic("encoding is not 57-bytes long")
	}
	//排序
	sort.Strings(strs)
	a := alphabet{
		len: int64(len(strs)),
	}
	//深度拷贝
	for i, char := range strings.Join(strs, "") {
		a.char[i] = char
	}
	return a
}

func (a *alphabet) Length() int64 {
	return a.len
}
func (a *alphabet) Index(t rune) (int64, error) {
	for i, char := range a.char {
		if char == t {
			return int64(i), nil
		}
	}
	return 0, fmt.Errorf("element '%v' is not part of the alphabet", t)
}

func debupe(s []string) []string {
	out := make([]string, 0, 0)
	m := make(map[string]bool)
	for _, char := range s {
		if _, ok := m[char]; !ok {
			m[char] = true
			out = append(out, char)
		}
	}
	return out
}
