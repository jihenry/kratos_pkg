package log

import (
	"io"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func createFileWriter(dir, fileName string, maxAge int) (io.Writer, error) {
	month := time.Now().Format("200601")
	filePath := filepath.Join(dir, month)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.MkdirAll(filePath, 0755)
		if err != nil {
			return nil, err
		}
	}
	return rotatelogs.New(
		filepath.Join(filePath, fileName+".%Y%m%d"+".log"),
		rotatelogs.WithLinkName(filepath.Join(dir, fileName+".log")),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(maxAge)),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
}
