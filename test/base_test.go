package test

import (
	"fmt"
	"os"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestT(t *testing.T) {
	t.Log(234821 % 10)
}

func TestAssign(t *testing.T) {
	var (
		a = 1
		b = 2
	)
	a, b = b, a // 交换
	t.Log(a, b)
	const (
		c = iota + 1
		d
		e
		f
	)
	const g = 5
	t.Log(c, d, e, f, g)
	const (
		h = 1 << iota
		i
		j
		k = 3 << iota
		l
		m
	)
	t.Log(h, i, j, k, l, m)
	const (
		n = 2
		o
		p
	)
	// 不可用连续赋值
	var (
		q = 3
		r int
		s uint32
	)
	t.Log(n, o, p, q, r, s)
	const (
		u string = "aa"
		v
		w = 2
		x
	)
	t.Log(u, v, w, x)
}

func BenchmarkGoquery(b *testing.B) {
	b.ResetTimer()
	var res map[string]interface{}
	_ = jsoniter.UnmarshalFromString("{\"a\":777}", &res)
	b.StopTimer()
}

func TestSwtichIf(t *testing.T) {
	i := 5
	i++
	switch {
	case i < 5:
		t.Log("i < 5")
	case i >= 5:
		t.Log("i >= 5")
	}
}

func TestNilPointer(t *testing.T) {
	a := func(i *int) {
		*i = 0
	}
	i1 := 4
	a(&i1)
	t.Log(i1)
	var i2 int
	a(&i2)
	//a(nil)// panic: runtime error: invalid memory address or nil pointer dereference
	t.Log(i2)
	var slice1 []int // 0 0
	t.Log(len(slice1), cap(slice1))
	lenObj := func(i []int) int {
		return len(i) // 不报错 。。。。。
	}
	t.Log(lenObj(nil))
}

func TestMult(t *testing.T) {
	fmt.Println(9 * 0.9)
	num := 99
	fmt.Println(float64(num) * 0.95)
}

func TestMkdir(t *testing.T) {
	t.Log(os.Mkdir("/tmp/aaa", os.ModePerm))
}
