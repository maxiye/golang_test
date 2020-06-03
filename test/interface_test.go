package test

import (
	"fmt"
	"testing"
)

type Animal interface {
	Eat()
	Name() string
}

type Plant interface {
	Grow()
}

type Living interface {
	Animal
	Plant
}

type Dog struct {
}

func (d *Dog) Eat() {
	fmt.Println("dog eat")
}

func (d *Dog) Name() string {
	return "dog"
}

func (d *Dog) Run() {
	fmt.Println("dog run")
}

type Cat struct {
	NickName string
}

func (c *Cat) Eat() {
	fmt.Println("dog eat")
}

func (c *Cat) Name() string {
	return c.NickName
}

func getAnimalName(a Animal) {
	fmt.Println("animal:", a.Name())
}

type Something struct {
}

func (s *Something) Eat() {
	fmt.Println("something eat")
}

func (s *Something) Name() string {
	return "something"
}

func (s *Something) Grow() {
	fmt.Println("something grow...")
}

func TestInterface(t *testing.T) {
	var ani Animal
	ani = new(Dog) // 必须实现接口中所有的方法才可以
	ani.Eat()
	//ani.Run() //无此方法
	cat := new(Cat)
	cat.NickName = "cat"
	//dog := Dog{}//Cannot use 'dog' (type Dog) as type Animal Type does not implement 'Animal' as 'Eat' method has a pointer receiver
	dog := &Dog{}
	getAnimalName(cat)
	getAnimalName(dog)
}

// 自由组合接口，必须实现子接口的所有方法
func TestComplexInterface(t *testing.T) {
	var l Living
	l = new(Something)
	l.Eat()
	l.Grow()
	t.Log(l.Name())
}
