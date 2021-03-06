package test

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
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

func TestAtoI(t *testing.T) {
	var a interface{}
	a = ""
	pkgId, _ := strconv.Atoi(a.(string))
	assert.Equal(t, 0, pkgId)
}

func TestRangeString(t *testing.T) {
	for i, c := range "23cd我搜索" {
		// i为字符的起始字节位置，【我】时，i=4，【搜】时，i=7
		t.Log(i, c)
	}
}

func TestSplit(t *testing.T) {
	s := ""
	sArr := strings.Split(s, ",")
	t.Log(len(sArr), sArr) // len = 1
	s2 := "aaa"
	t.Log(strings.Split(s2, ","))
}

func TestTimeFormat(t *testing.T) {
	form := "20060102150405-.000.000000"
	t.Log(time.Now().Format(form))
	t.Log(time.Now().Format("2006-01-02-15:04:05:1233"))
}

func TestJsonSlice(t *testing.T) {
	var a []int
	var b = []int{}
	s1, _ := jsoniter.MarshalToString(a)
	s2, _ := jsoniter.MarshalToString(b)
	t.Log(s1, s2)
}

func TestStringIndex(t *testing.T) {
	src := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	t.Log(string(src[1%26]))
	t.Log(string(src[2%26]))
	t.Log(string(src[26%26]))
	t.Log(string(src[25%26]))
	t.Log(string(src[27%26]))
}

func TestJsonStr(t *testing.T) {
	s, err := jsoniter.Marshal("aaaa")
	fmt.Println(string(s), err)
}

func TestCalcDS(tt *testing.T) {
	salt := "w9p2p72p9octwd7lj1oa913hncq1k4td"
	t := time.Now().Unix()
	r := "J6h4er"
	secret := fmt.Sprintf("salt=%s&t=%d&r=%s", salt, t, r)
	s := fmt.Sprintf("%x", md5.Sum([]byte(secret)))
	fmt.Printf("%d,%s,%s", t, r, s)
	fmt.Println()
}
