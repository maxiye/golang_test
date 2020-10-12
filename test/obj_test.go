package test

import (
	"fmt"
	"sync"
	"testing"
)

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

type Actor struct {
	Sing       string
	Jump       int
	Rap        []string
	Basketball uintptr
}

// 和接口一样组合
type ManActor struct {
	Man // 可以用man的方法
	Actor
}

func (m *Man) Sleep(g Girl) *Man {
	// 必须赋值，直接改变不生效
	fmt.Println(m, " slept ", g)
	return new(Man)
}

func TestObjAsign(t *testing.T) {
	var g Girl
	var g2 Girl
	var g3 Girl
	const (
		A byte = 1 << iota
		B byte = 1 << iota
		C byte = 1 << iota
	)
	g.Age = A
	g2.Age = B
	g3.Age = C
	t.Log(g, g2, g3)
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

func TestComplexObj(t *testing.T) {
	manActor := new(ManActor)
	manActor.Name = "Xukun Cai"
	manActor.Age = 22
	manActor.Gf = &[]Girl{}
	manActor.Sleep(Girl{"aa", 20})
	t.Log(manActor)
}

var single *Actor
var once sync.Once // 必须放到程序外边

func getSinleActor() *Actor {
	once.Do(func() {
		fmt.Println("once do")
		single = new(Actor)
	})
	return single
}

func getSingleActor2() {
	once.Do(func() {
		fmt.Println("once do2")
	})
}

func TestSingleOnce(t *testing.T) {
	actor := getSinleActor()
	actor2 := getSinleActor()
	getSingleActor2() // 不执行，once只能用一次
	t.Logf("%p %p", actor, actor2)
}
