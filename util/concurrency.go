package util

import (
	"context"
	"runtime/debug"

	"github.com/go-kratos/kratos/v2/log"
)

func Go(c context.Context, f func(context.Context)) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("recover:%s %v", err, string(debug.Stack()))
				return
			}
		}()
		f(c)
	}()

}

func GoWithParam(c context.Context, f func(context.Context, ...interface{}), args ...interface{}) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("recover:%s\n%v", err, string(debug.Stack()))
				return
			}
		}()
		f(c, args...)
	}()
}
