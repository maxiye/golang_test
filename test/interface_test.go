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

	dog2 := Dog{}
	dog2.Eat()
}

// 自由组合接口，必须实现子接口的所有方法
func TestComplexInterface(t *testing.T) {
	var l Living
	l = new(Something)
	l.Eat()
	l.Grow()
	t.Log(l.Name())
}

type people interface {
	Speak(string) string
}

type student struct{}

// 仅 *people可以
func (stu *student) Speak(think string) (talk string) {
	if think == "bitch" {
		talk = "You are a good boy"
	} else {
		talk = "hi"
	}
	return
}

type teacher struct{}

// people & *people都可以
func (t teacher) Speak(word string) (out string) {
	return word
}

// people & *people都可以
func (t teacher) Speak2(word string) (out string) {
	fmt.Println(t.Speak(word))
	return word + "2"
}

type suTeacher struct {
	teacher
	Su bool
}

func (st suTeacher) Speak(word string) string {
	return "su" + word
}

func TestImpl(t *testing.T) {
	//var peo people = student{}// Cannot use 'student{}' (type student) as type people Type does not implement 'people' as 'Speak' method has a pointer receiver
	var peo people = &student{}
	think := "bitch"
	fmt.Println(peo.Speak(think))

	var p2 people = teacher{}
	var p3 people = new(teacher)
	println(p2.Speak("aa"), p3.Speak("bb"))
}

func TestMix(t *testing.T) {
	sut := suTeacher{}
	t.Log(sut.Speak("aa"))
	t.Log(sut.Speak2("aa"))
}
