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
	Gf   *[]Girl
	Bf   []*Girl
}

func (m *Man) Sleep(g Girl) *Man {
	// 必须赋值，直接改变不生效
	//m.Wifi = append(m.Wifi, g)
	m.Wifi[2] = g
	return new(Man)
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

func TestObj2(t *testing.T) {
	g1 := Girl{"A", 20}
	g2 := Girl{Name: "B", Age: 22}
	g3 := new(Girl)
	g3.Name = "C"
	girls := []Girl{g1, g2}
	//m := Man{"M", 40, []Girl{g1, g2, g3}}//Cannot use 'g3' (type *Girl) as type Girl
	m := Man{Name: "M", Age: 40, Wifi: girls}
	t.Log(m)
	m.Gf = &girls
	m.Bf = []*Girl{g3}
	t.Log(m)
	m.Sleep(g1)
	t.Log(m)
}

func doObj(girl Girl) {
	girl.Age++
}

func doObjP(girl *Girl) {
	girl.Age++
}
