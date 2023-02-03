package log

import (
	"testing"

	"github.com/google/uuid"
)

func TestInitZapLogger(t *testing.T) {
	InitZapLogger(ZapLoggerConf{
		Level:       "info",
		FileName:    "game",
		FilePath:    "./mysql",
		MaxSize:     100,
		MaxBackups:  30,
		MaxAge:      30,
		Compress:    true,
		ShowConsole: true,
	})

	t.Log("zaplog:", sugarLogger, "token:", uuid.New().String())
}

/**
goos: darwin
goarch: amd64
pkg: gitlab.yeahka.com/gaas/pkg/zaplog
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkInitZapLogger-12         171686              5961 ns/op            2954 B/op         52 allocs/op
PASS
ok      gitlab.yeahka.com/gaas/pkg/zaplog     1.424s
*/
func BenchmarkInitZapLogger(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		InitZapLogger(ZapLoggerConf{
			Level:       "info",
			FileName:    "game",
			FilePath:    "./mysql",
			MaxSize:     100,
			MaxBackups:  30,
			MaxAge:      30,
			Compress:    true,
			ShowConsole: true,
		})
	}
}
