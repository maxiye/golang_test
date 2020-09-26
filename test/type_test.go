package test

import (
	"go/types"
	"math"
	"testing"
	"unsafe"
)

type myFun func(n ...int) int
type aa int32

func TestType(t *testing.T) {
	var a int64 = 1 << 39
	var b int32
	var c aa
	b = int32(a) //不支持隐式转换
	c = aa(b)    //甚至不支持别名隐式转换
	t.Log(a, b, c, math.MaxFloat32, math.MaxUint32)
	var sum myFun = func(n ...int) int {
		s := 0
		for _, v := range n {
			s += v
		}
		return s
	}
	t.Log(sum(0, 1, 2, 3, 4))
}

func TestPoint(t *testing.T) {
	var a aa = 1
	b := &a
	//b++
	t.Log(a, b, *b)
	t.Logf("%T %T %T", a, b, *b)
	var d uintptr = 0xc000054290
	t.Log(d, *(*int)(unsafe.Pointer(d)))
	var s string //空字符串
	t.Log("*" + s + "*")
}

func TestOp(t *testing.T) {
	a := 7
	b := 2
	t.Log(a&^b, ^b, a&(^b))
	t.Logf("二进制%b %b %b", a&^b, ^b, a&(^b))
}

func TestAssert(t *testing.T) {
	a := 1
	b := "aa"
	c := 'a'
	var d types.Nil
	var e []int
	f := func() int {
		return 1
	}
	t.Log(printType(a), printType(b), printType(c), printType(d), printType(e), printType(f))
}

func printType(t interface{}) string {
	switch t.(type) {
	case int:
		return "int"
	case string:
		return "string"
	case func() int: // 精准匹配参数和返回值
		return "func"
	case rune:
		return "rune"
	case types.Nil:
		return "nil"
	default:
		return "don't know"
	}
}
