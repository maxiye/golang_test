package test

import (
	"reflect"
	"testing"
)

func TestTypeValue(t *testing.T) {
	map1 := map[string]interface{}{}
	map1["a"] = 1
	map1["b"] = "b"
	map1["c"] = [...]int{1, 2, 3}
	t.Log(reflect.TypeOf(map1["a"]))
	t.Log(reflect.TypeOf(map1["c"]))
	t.Log(reflect.ValueOf(map1["c"]).Type())
	t.Log(reflect.TypeOf(map1["d"]))
	t.Log(reflect.ValueOf(map1["d"]))
	tp := reflect.TypeOf(map1).Kind()
	t.Log(tp == reflect.Map)
}

func TestTypeField(t *testing.T) {
	type b struct {
		b string
	}
	type aaa struct {
		aa string
		_  struct{}
		b
		bb string
	}
	ao := aaa{}
	tp := reflect.TypeOf(ao)
	t.Log(tp.NumField())
	for i := 0; i < tp.NumField(); i++ {
		t.Log(tp.Field(i))
	}
	t.Log(ao)
}

func TestReflectObj(t *testing.T) {
	man := &Man{
		Name: "aa",
		Age:  20,
	}
	t.Log(reflect.ValueOf(man))
	t.Log(reflect.ValueOf(*man).FieldByName("Name"))
	t.Log(reflect.ValueOf(man).MethodByName("Sleep")) // 没有*
	reflect.ValueOf(man).MethodByName("Sleep").Call([]reflect.Value{reflect.ValueOf(Girl{})})
	//t.Log(reflect.ValueOf(*man).FieldByIndex([]int{0, 1}))//必须数组
}

/**
&是取地址符号, 用在对象变量前取到对象的地址,23行
*可以表示一个变量是指针类型(r是一个指针变量)，39行
*也可以表示指针类型变量所指向的存储单元 ,也就是这个地址所指向的值，28行
*/
func TestReflectObj2(t *testing.T) {
	man := Man{
		Name: "aa",
		Age:  20,
	}
	t.Log(reflect.ValueOf(man))
	t.Log(reflect.ValueOf(man).FieldByName("Name"))
	t.Log(reflect.ValueOf(&man).MethodByName("Sleep")) // 有&
	reflect.ValueOf(&man).MethodByName("Sleep").Call([]reflect.Value{reflect.ValueOf(Girl{})})
	//t.Log(reflect.ValueOf(*man).FieldByIndex([]int{0, 1}))//必须数组
}

func TestReflectTag(t *testing.T) {
	modList := ModList{}
	//t.Log(reflect.TypeOf(&modList).FieldByName("aa"))//指针类型没有field：FieldByName of non-struct type *test.ModList
	objType, _ := reflect.TypeOf(modList).FieldByName("ModulePath")
	t.Log(objType.Tag.Get("json"))
}

func TestReflectDeepEqual(t *testing.T) {
	man1 := &Man{}
	man2 := Man{}
	man3 := Man{Name: "man"}
	t.Log(reflect.DeepEqual(man1, man2))
	t.Log(reflect.DeepEqual(man1, &man2))
	t.Log(reflect.DeepEqual(*man1, man2))
	t.Log(reflect.DeepEqual(man2, man3))
}

func TestReflectElem(t *testing.T) {
	man := Man{Name: "Hello"}
	//t.Log(reflect.ValueOf(man).Elem())//panic: reflect: call of reflect.Value.Elem on struct Value
	manElem := reflect.ValueOf(&man).Elem()
	t.Log(manElem)
	manElem.FieldByName("Name").Set(reflect.ValueOf("World"))
	t.Log(manElem, man)
}
