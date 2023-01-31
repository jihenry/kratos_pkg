package uuid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReverse(t *testing.T) {
	a := []rune("abc123")
	reverse(a)
	should := require.New(t)
	should.Equal(string(a), "321cba1")
}
