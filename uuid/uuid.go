package uuid

import (
	"strings"

	"github.com/google/uuid"
)

var (
	defaultEncoder = &base57{newAlphabet(DefaultAlphabet)}
)

type EnCoding interface {
	Encode(uuid.UUID) string
	Decode(string) (uuid.UUID, error)
}

func New() string {
	return defaultEncoder.Encode(uuid.New())
}

func NewWithNamespace(name string) string {
	var u uuid.UUID
	switch {
	case name == "":
		u = uuid.New()
	case strings.HasPrefix(strings.ToLower(name), "http://"), strings.HasPrefix(strings.ToLower(name), "https://"):
		u = uuid.NewSHA1(uuid.NameSpaceURL, []byte(name))
	default:
		u = uuid.NewSHA1(uuid.NameSpaceDNS, []byte(name))
	}
	return defaultEncoder.Encode(u)
}

func NewWithAlphabet(abc string) string {
	return base57{newAlphabet(abc)}.Encode(uuid.New())
}
