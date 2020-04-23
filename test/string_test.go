package test

import (
	"strconv"
	"testing"
)

func TestString(t *testing.T) {
	s := "哈下啊"
	//s[0] = 229;//不可变
	for i := 0; i < len(s); i++ {
		t.Log(strconv.FormatInt(int64(s[i]), 16))
	}
	for _, c := range s {
		//t.Log(c, string(c))
		t.Logf("%[1]c %[1]x %s", c, string(c))
	}
	runeArr := []rune(s)
	t.Log(runeArr[0], runeArr)
	// 必须为utf8的三个byte，不能为rune（Unicode）的2个
	s2 := "\xe5\x93\x88\xe4\xb8\x8b\xe5\x95\x8a"//哈下啊
	t.Log(s2, s2[3:6])
}
