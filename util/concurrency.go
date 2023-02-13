package util

import (
	"context"
	"runtime/debug"

	zaplog "gitlab.yeahka.com/gaas/pkg/log"
)

func Go(c context.Context, f func(context.Context)) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zaplog.FromContext(c).Errorf("recover:%s %v", err, string(debug.Stack()))
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
				zaplog.FromContext(c).Errorf("recover:%s %v", err, string(debug.Stack()))
				return
			}
		}()
		f(c, args...)
	}()
}
