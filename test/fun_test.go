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
	/*TimerFunc(func(i ...int) int {//Cannot call non-function TimerFunc(func(i ...int) int
		res := 1
		for _, v := range i {
			res *= v
		}
		return res
	})(2, 3, 4)*/
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
		return function(n...), nil //...
	}, nil
}

func TestFunArgs(t *testing.T) {
	a := [5]int{0, 1}
	funcWithArray(a)
	t.Log(a)
	a2 := [5]int{0, 1}
	funcWithArrayP(&a2)
	t.Log(a2)
	sl := make([]int, 5, 10)
	funcWithSlice(sl)
	t.Log(sl)
	sl2 := make([]int, 5, 10)
	// 原分片改变
	funcWithSliceP(&sl2)
	t.Log(sl2)
}

// 传值，不改变原数组
func funcWithArray(arr [5]int) {
	arr[0] = 22222
}

// 传指针，改变原数组
func funcWithArrayP(arr *[5]int) {
	arr[0] = 22222
}

func funcWithSlice(slice []int) {
	//slice = []int{0, 1}//不改动原数据
	slice[0] = 22222
}

//传指针的引用，赋值操作直接改变原指针地址
func funcWithSliceP(slc *[]int) {
	*slc = []int{22222}
}

// 生成器，闭包
func intGenerator() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

func TestGen(t *testing.T) {
	intGetter := intGenerator()
	t.Log(intGetter())
	t.Log(intGetter())
	t.Log(intGetter())
}

// 先进后出
func TestDefer(t *testing.T) {
	defer func() {
		t.Log(1111)
	}()
	defer func() {
		t.Log(2222)
	}()
}

func TestDefer2(t *testing.T) {
	deferfun := func() {
		defer func() {
			println("aaa")
		}()
		defer func() {
			println("bbb")
		}()
	}
	defer println("111")
	deferfun()
	defer println("222")
	println("ccc")
	defer panic("fff")
	panic("eee")
}

func TestDefer3(t *testing.T) {
	calc := func(index string, a, b int) int {
		ret := a + b
		fmt.Println(index, a, b, ret)
		return ret
	}
	a := 1
	b := 2
	defer calc("1", a, calc("10", a, b))
	a = 0
	defer func() { println(a, b) }()     // 引用外部变量，
	defer calc("2", a, calc("20", a, b)) // defer时，参数已确定
	b = 1
}
