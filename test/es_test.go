package test

import (
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"testing"
)

func TestEs(t *testing.T) {
	errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
	pwd := "2" + "3"
	pwd += "4" + "5" + ".com"
	es, err := elastic.NewClient(
		elastic.SetErrorLog(errorlog),
		elastic.SetURL("http://elastic:"+pwd+"@172.17.210.70:9200/"),
		elastic.SetSniff(false))
	if err == nil {
		indexes, _ := es.IndexNames()
		t.Log(indexes)
	} else {
		t.Log(err)
	}
}
