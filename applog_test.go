package goapplog

import (
	"github.com/olivere/elastic/v7"
	"testing"
)

func TestNewLogger(t *testing.T) {
	c, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	l := NewLogger(Elasticsearch(c), IP("127.0.0.1"))
	l.Infoln("这是测试数据 Test data")
}
