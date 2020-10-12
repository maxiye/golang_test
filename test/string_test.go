package test

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
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
	s2 := "\xe5\x93\x88\xe4\xb8\x8b\xe5\x95\x8a" //哈下啊
	t.Log(s2, s2[3:6])
}

func TestInt(t *testing.T) {
	t.Log(float32(12121) / 1024)
}

func TestJson(t *testing.T) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	t.Log(json.NewEncoder(os.Stdout).Encode(map[string]interface{}{"aa": "bb", "bb": 1, "cc": [2]int{1, 2}}))
	jsonStr, err := json.MarshalToString(&map[string]interface{}{"aa": "bb", "bb": 1, "cc": [2]int{1, 2}})
	if err == nil {
		t.Log(jsonStr)
	}
	var jsonObj map[string]interface{}
	if err := json.UnmarshalFromString(jsonStr, &jsonObj); err == nil {
		t.Log(jsonObj)
	}
}

func TestGJson(t *testing.T) {
	jsonStr := `
{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}
`
	json := gjson.Parse(jsonStr)
	t.Log(json.Get("friends"))
	val := gjson.Get(jsonStr, "name.first")
	children := gjson.Get(jsonStr, "children")
	gf := gjson.Get(jsonStr, "gf")
	t.Log(val, children)
	t.Log(gf)
}

type ModList struct {
	ModulePath    string  `json:"module_path,omitempty"`
	DownloadCount float64 `json:"download_count,omitempty"`
}

func TestGoQuery2(t *testing.T) {
	if res, err := http.Get("https://goproxy.cn/stats/trends/latest"); err == nil && res != nil {
		bytes, _ := ioutil.ReadAll(res.Body)
		defer func() {
			_ = res.Body.Close()
		}()
		var jsonObj []*ModList
		err = jsoniter.Unmarshal(bytes, &jsonObj)
		t.Log(len(jsonObj), jsonObj[len(jsonObj)-1])
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("%d Kb\n", m.Alloc/1024)
	} else {
		t.Log(err)
	}
}
