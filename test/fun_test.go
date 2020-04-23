package test

import (
	"fmt"
	"testing"
	"time"
)

func TestFun(t *testing.T) {
	func1 := func(i ...int) int {
		sum := 0
		for _, v := range i {
			sum += v
		}
		return sum
	}
	func2, err := TimerFunc(func1)
	//t.Logf("%T", TimerFunc(func1))//multiple-value TimerFunc() in single-value context
	t.Logf("%T %T", func2, err)
	t.Log(func2(1, 2, 3, 4))
}

func TimerFunc(function func(i ...int) int) (func(i ...int) (int, error), error) {
	return func(n ...int) (int, error) {
		start := time.Now()
		defer func() {
			fmt.Println("time used(s): ", time.Since(start).Seconds())
		}()
		time.Sleep(time.Second * 2)
		return function(n...),nil//...
	},nil
}