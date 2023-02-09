package util

import (
	"context"
	"fmt"
	"testing"
	"time"
)


func func1(c context.Context, args ...interface{}) {
	for _, arg := range args {
		fmt.Println("func1", arg)
	}
}

func TestGoWithParam(t *testing.T) {
	ctx := context.TODO()
	GoWithParam(ctx, func1, 1, "op")
}

func func2(c context.Context) {
	fmt.Println("func2")
	panic("here")
}

func TestGo1(t *testing.T) {
	ctx := context.TODO()
	Go(ctx, func2)
	time.Sleep(time.Second)
}

type st struct {
	i int
}

func (s *st) func3(c context.Context) {
	fmt.Println("func3", s.i)
}

func TestGo2(t *testing.T) {
	ctx := context.TODO()
	s := st{4}
	Go(ctx, s.func3)
	time.Sleep(time.Second)
}
