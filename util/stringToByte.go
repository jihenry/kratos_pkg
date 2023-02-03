package util

import (
	"reflect"
	"unsafe"
)

type StringPointer interface {
	StringToByte(string) []byte
	ByteToString([]byte) string
}
type strPointer struct {
}

var (
	_ StringPointer = (*strPointer)(nil)
)

func NewPointer() StringPointer {
	return &strPointer{}
}

func (s *strPointer) StringToByte(str string) []byte {
	header := (*reflect.StringHeader)(unsafe.Pointer(&str))
	newHeader := reflect.SliceHeader{
		Data: header.Data,
		Len:  header.Len,
		Cap:  header.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&newHeader))
}

func (s *strPointer) ByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
