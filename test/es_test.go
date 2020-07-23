package test

import (
	"github.com/olivere/elastic"
	"log"
	"os"
	"testing"
)

func TestEs(t *testing.T) {
	errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
	es, err := elastic.NewClient(
		elastic.SetErrorLog(errorlog),
		elastic.SetURL("http://172.17.12.40:9200/"),
		elastic.SetSniff(false))
	t.Log(es, err)
	//t.Log(es.ElasticsearchVersion("http://172.17.12.40:9200"))
}
