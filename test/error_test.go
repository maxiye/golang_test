package test

import (
	"errors"
	"fmt"
	"testing"
)

var AError error = errors.New("A error")

func getAError() error {
	return AError
}

func init() {
	fmt.Println("init1")
}

func init() {
	fmt.Println("init2")
}

func init() {
	fmt.Println("init3")
}

func TestError(t *testing.T) {
	e := getAError()
	t.Log(errors.Is(e, AError))
	t.Error(AError.Error())
}

func TestPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic: ", err)
		}
	}()
	//os.Exit(1)// 不执行defer
	panic("啊啊啊啊啊")
}
