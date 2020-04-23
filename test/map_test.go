package test

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	m1 := map[string]int{"aa":2}
	t.Log(m1, m1["bb"], len(m1))
	m2 := make(map[int]string, 20)
	m2[2] = "hahaha"
	// 没有的key是空字符串
	t.Log(m2, reflect.TypeOf(m2[233]),len(m2))
	if v, ok := m2[999];ok {
		t.Log(v)
	} else {
		t.Log("no exist")
	}
}

func TestMapFor(t *testing.T) {
	m := map[string]interface{}{"a": 2, "b": "aa"}
	for k, v:= range m {
		t.Log(k, v)
	}
}

func TestMapFun(t *testing.T) {
	m := map[string]func(...interface{})interface{}{}
	m["sum"] = func(i ...interface{}) interface {}{
		sum := 0
		for _, a := range i {
			switch a.(type) {// 只能用于switch
				case int:
					sum += a.(int)
				default:
					t.Log("type err: ", reflect.TypeOf(a))
			}
		}
		return sum
	}
	t.Log(m["sum"](1, 2, 3, 2, 1))
	t.Log(m["sum"]("2", 1))
}

func TestMapForSet(t *testing.T) {
	set := map[string]bool{}
	set["bb"] = true
	if v, ok := set["aa"];ok {
		t.Log(v)
	} else {
		t.Log("not exist")
	}
	t.Log(len(set))
	delete(set, "bb")
}