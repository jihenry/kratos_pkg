package log

import (
	"io"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func createFileWriter(dir, fileName string, maxAge int) (io.Writer, error) {
	_ = os.MkdirAll(dir, 0755)
	return rotatelogs.New(
		filepath.Join(dir, fileName+".%Y%m%d"+".log"),
		rotatelogs.WithLinkName(filepath.Join(dir, fileName+".log")),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(maxAge)),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
}
