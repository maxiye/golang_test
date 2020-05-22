package test

import (
	jsoniter "github.com/json-iterator/go"
	"os"
	"testing"
)

func TestArray(t *testing.T) {
	// 加‘...’生成数组，不加则是slice
	a := [...]int{0, 1, 2, 3, 4, 5, 6}
	for idx, v := range a {
		t.Log(idx, v)
	}
	c := a[1:3]
	d := a[:2]
	c[0] = 999
	d[0] = 888
	t.Log(a, c, d) //切片也相同
}

func TestSlice(t *testing.T) {
	var b []int
	b = append(b, 1)
	b = append(b, 2)
	t.Log(len(b), cap(b))
	b = append(b, 3)
	t.Log(len(b), cap(b))

	a := []int{1, 2, 3, 4, 5}
	t.Log(len(a), cap(a))
	a = append(a, 6)
	t.Log(len(a), cap(a))

	c := make([]byte, 3, 6)
	//t.Log(c[3])//index out of range
	t.Log(len(c), cap(c))
}

func TestSliceAppend(t *testing.T) {
	b := []int{3, 4, 5}
	c := append(b, 9)
	t.Log(b, c)
	t.Logf("b:%p c:%p", b, c)
	// ？？？
	bytes := append([]byte("hello "), "world你好"...)
	//bytes = append(bytes, 2, 3, 4, "aaaa"...)//Invalid use of ..., corresponding parameter is non-variadic
	t.Log(bytes)
	cs := append([]rune("cd"), 101, 'f')
	cs2 := append([]rune{'a', 98}, cs...)
	t.Log(cs, cs2)
	for k, v := range cs2 {
		t.Log(k, v, string(v))
	}
}

func TestJson(t *testing.T) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	t.Log(json.NewEncoder(os.Stdout).Encode(map[string]interface{}{"aa": "bb", "bb": 1, "cc": [2]int{1, 2}}))
}
