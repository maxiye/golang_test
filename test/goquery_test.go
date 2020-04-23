package test

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"testing"
)

func TestGoQuery(t *testing.T) {
	if res, err := http.Get("https://baidu.com"); err == nil && res != nil {
		t.Log(res)
		doc, err2 :=goquery.NewDocumentFromReader(res.Body)
		if err2 == nil {
			t.Log(doc.Find("a").Size())
		}
	}
}

func BenchGoqueryB(b *testing.B) {
	b.ResetTimer()
	res, _ := http.Get("xq.2345.com")
	b.Log(res)
	b.StopTimer()
}
