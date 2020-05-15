package test

import "testing"

type Girl struct {
	Name string
	Age  byte
}

type Man struct {
	Name string
	Age  byte
	Wifi []Girl
}

func TestObj(t *testing.T) {
	g := Girl{"AAA", 17}
	g.Name = "Ali"
	g2 := Girl{
		Name: "Bli",
		Age:  30, //nice ,
	}
	g3 := new(Girl) // 引用/指针
	g3.Name = "Cli"
	t.Log(g, g2, g3)
	//doObj(g3)//Cannot use 'g3' (type *Girl) as type Girl
	doObj(g)
	t.Log(g)
	doObjP(&g)
	t.Log(g)
	doObjP(g3)
	t.Log(g3)
}

func doObj(girl Girl) {
	girl.Age++
}

func doObjP(girl *Girl) {
	girl.Age++
}
