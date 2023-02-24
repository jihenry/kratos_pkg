package util

import (
	"runtime"

	"github.com/go-kratos/kratos/v2/log"
)

func Recover(moduleName string) {
	if err := recover(); err != nil {
		buf := make([]byte, 64<<10)
		n := runtime.Stack(buf, false)
		buf = buf[:n]
		log.Errorf("moduleName:%s err:%s:\n%s\n", moduleName, err, buf)
	}
}
