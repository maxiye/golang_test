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

func TestCmp(t *testing.T) {
	m1 := new(Man)
	m2 := new(Man)
	//m3 := Man{}
	//m4 := new(ManActor)
	//m5 := ManActor{}
	//m6 := ManActor{}
	t.Log(m1 == m2) // false，指针指向的变量类型相同可毕竟
	//t.Log(m1 == m3)// 类型不同不可比较
	//t.Log(m4 == m5)// 实例与结构体不可比较
	//t.Log(m2 == m4)// 指针类型不同不可比较
	//t.Log(*m1 == *m2) //Invalid operation: *m1 == *m2 (operator == is not defined on Man)
	//t.Log(m3 == m5) // Invalid operation: m3 == m5 (mismatched types Man and ManActor)
	//t.Log(m5 == m6) // Invalid operation: m5 == m6 (operator == is not defined on ManActor)

	type student struct {
		Name string
		List chan string
	}
	s1 := student{Name: "a"}
	s2 := student{}
	t.Log(s1 == s2) // 仅包含可比较的数据类型的结构是可比较

	c1 := make(chan string)
	c2 := make(chan string)
	t.Log(c1 == c2) //chan可比较

	var a interface{} = "a"
	var b interface{} = 2
	t.Log(a == b) // interface 可比较

	//c := 0
	//d := "c"
	//t.Log(c == d) // Invalid operation: c == d (mismatched types int and string)
}

func TestUintSub(t *testing.T) {
	var a uint = 2
	var b uint = 4
	t.Log(a - b) // 18446744073709551614
}

func TestAlign(t *testing.T) {
	type exp struct {
		B bool
		I int16
		S string
	}
	exp1 := exp{
		B: false,
		I: 0,
		S: "",
	}
	aliagn := unsafe.Alignof(exp1)
	t.Log(aliagn, unsafe.Sizeof(exp1))
	t.Log(unsafe.Sizeof(exp1.B), unsafe.Offsetof(exp1.B), &exp1.B)
	t.Log(unsafe.Sizeof(exp1.I), unsafe.Offsetof(exp1.I), &exp1.I)
	t.Log(unsafe.Sizeof(exp1.S), unsafe.Offsetof(exp1.S), &exp1.S)

	type exp2 struct {
		B bool
		S string
		I int16
	}
	exp21 := exp2{
		I: 0,
		B: false,
		S: "",
	}
	aliagn2 := unsafe.Alignof(exp21)
	t.Log(aliagn2, unsafe.Sizeof(exp21))
	t.Log(unsafe.Sizeof(exp21.B), unsafe.Offsetof(exp21.B), &exp21.B)
	t.Log(unsafe.Sizeof(exp21.S), unsafe.Offsetof(exp21.S), &exp21.S)
	t.Log(unsafe.Sizeof(exp21.I), unsafe.Offsetof(exp21.I), &exp21.I)
}
