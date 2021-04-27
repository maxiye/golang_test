package test

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"
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

func TestSliceCopy(t *testing.T) {
	a := []int{1, 2, 3}
	b := append(a, 4)
	t.Logf("a:%d cap(a):%d; b:%d cap(b):%d", a, cap(a), b, cap(b))
	// // append后cap发生改变时，a的操作不会影响到b
	c := append(a, 5)
	t.Logf("a:%d cap(a):%d; b:%d cap(b):%d; c:%d cap(c):%d", a, cap(a), b, cap(b), c, cap(c))
	d := append(b, 6)
	e := append(b, 7)
	// cap不改变，影响到d
	t.Log(b, d, e)
}

func TestSliceRefrence(t *testing.T) {
	a := make([]int, 1, 300)
	b := append(a, 4)
	a[0] = 0
	t.Log(a, b)
	b[1] = 3
	t.Log(a, b)
	c := append(a, 5)
	t.Log(a, b, c)
	d := append(a, 6)
	// cap未发生改变时，append都是在原数组上操作的，会影响到所有使用统一原始数组的slice
	t.Log(cap(a), cap(b), cap(c), cap(d))
	t.Log(a, b, c, d)
}

func TestSliceSort(t *testing.T) {
	slice1 := []string{"aa", "我", "撒旦撒旦", "1"}
	sort.Strings(slice1)
	t.Log(slice1)
	t.Log(sort.StringsAreSorted(slice1))
	sort.Sort(sort.Reverse(sort.StringSlice(slice1)))
	t.Log(slice1, sort.StringsAreSorted(slice1)) // false
}

func TestManualSort(t *testing.T) {
	slice1 := []string{"aabb", "aabbcc", "我是", "撒旦撒旦", "2", "1"}
	less := func(i, j int) bool {
		return len(slice1[i]) > len(slice1[j])
	}
	sort.Slice(slice1, less)
	t.Log(slice1)
	sort.SliceStable(slice1, less)
	t.Log(slice1)
}

func ChunkDo(i interface{}, chunkNum int, do func(interface{}) error) error {
	if chunkNum <= 0 {
		panic("err chunkNum")
	}
	if reflect.TypeOf(i).Kind() == reflect.Slice {
		src := reflect.ValueOf(i)
		srcLen := src.Len()
		if srcLen == 0 {
			return nil
		}
		if chunkNum >= srcLen {
			return do(i)
		}
		for i := 0; i < srcLen; i += chunkNum {
			end := i + chunkNum
			if end > srcLen {
				end = srcLen
			}
			chunkSlick := src.Slice(i, end)
			if err := do(chunkSlick.Interface()); err != nil {
				return err
			}
		}
	}
	return nil
}

func TestChunk(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 13, 14}
	ChunkDo(a, 9, func(i interface{}) error {
		t.Log(i)
		return nil
	})
}

func TestIntGen(t *testing.T) {
	x := 234800
	for i := 0; i < 211; i++ {
		fmt.Println(x)
		x++
	}
}

func TestIntFormat(t *testing.T) {
	fmt.Printf("%07d", 99999991)
}

func TestShuffle(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 13, 14}
	randNew := rand.New(rand.NewSource(time.Now().UnixNano()))
	sort.Slice(a, func(i, j int) bool {
		return randNew.Intn(1000) > 500
	})
	t.Log(a)
	randNew = rand.New(rand.NewSource(time.Now().UnixNano()))
	sort.Slice(a, func(i, j int) bool {
		return randNew.Intn(1000) > 500
	})
	t.Log(a)
	randNew = rand.New(rand.NewSource(time.Now().UnixNano()))
	sort.Slice(a, func(i, j int) bool {
		return randNew.Intn(1000) > 500
	})
	t.Log(a)
}

func TestPerm(t *testing.T) {
	a := rand.Perm(10)
	fmt.Println(a)
	a = rand.Perm(10)
	fmt.Println(a)
}
