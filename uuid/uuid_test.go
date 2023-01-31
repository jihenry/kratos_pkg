package uuid

import (
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/require"
)

func TestGeneration(t *testing.T) {
	tests := []string{
		"",
		"http://www.baidu.com/",
		"https://www.baidu.com/",
		"baidu.com/",
	}
	should := require.New(t)
	for _, test := range tests {
		u := NewWithNamespace(test)
		if len(u) < 20 || len(u) > 24 {
			t.Errorf("expected %q to be in range[20,24],got %d", u, len(u))
		}
		should.Equal(u, "")
	}
}
func TestDecoding(t *testing.T) {
	id := uuid.New()
	t.Log(id)
	xId := defaultEncoder.Encode(id)
	t.Log(xId)
	uuid, _ := defaultEncoder.Decode(xId)
	t.Log(uuid)

	uu := NewWithNamespace("10.30.107.196")
	ip, _ := defaultEncoder.Decode(uu)
	uuids, _ := defaultEncoder.Decode(ip.String())
	t.Log(uuids)
}

/*
goos: darwin
goarch: amd64
pkg: gitlab.yeahka.com/gaas/pkg/uuid
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkUUID
BenchmarkUUID-12    	  221862	      5376 ns/op
PASS
*/
func BenchmarkUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}

/**
goos: darwin
goarch: amd64
pkg: gitlab.yeahka.com/gaas/pkg/uuid
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkEncode
BenchmarkEncode-12    	  245908	      4704 ns/op
PASS
*/
func BenchmarkEncode(b *testing.B) {
	u := uuid.New()
	for i := 0; i < b.N; i++ {
		defaultEncoder.Encode(u)
	}
}

/**
goos: darwin
goarch: amd64
pkg: gitlab.yeahka.com/gaas/pkg/uuid
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkDecoding
BenchmarkDecoding-12    	 1301613	       923.9 ns/op
PASS
*/
func BenchmarkDecoding(b *testing.B) {
	for i := 0; i < b.N; i++ {
		defaultEncoder.Decode("b485a394-9cab-4395-bb91-d4bd7fd0f2ec")
	}
}
