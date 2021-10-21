package esolivere

import (
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"reflect"
	"testing"
	"time"
)

type Tweet struct {
	User     string                `json:"user"`
	Age      int                   `json:"age"`
	Message  string                `json:"message"`
	Retweets int                   `json:"retweets"`
	Image    string                `json:"image,omitempty"`
	Created  time.Time             `json:"created,omitempty"`
	Tags     []string              `json:"tags,omitempty"`
	Location string                `json:"location,omitempty"`
	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"`
}

var mapping = `{
	"settings":{
		"number_of_shards": 3,
		"number_of_replicas": 1
	},
	"mappings":{
		"doc":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"age":{
					"type": "integer"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"image":{
					"type":"keyword"
				},
				"created":{
					"type":"date"
				},
				"tags":{
					"type":"keyword"
				},
				"location":{
					"type":"geo_point"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}`

func TestPingNode(t *testing.T) {
	PingNode()
}

func TestIndexExists(t *testing.T) {
	result := IndexExists("car_source", "test")
	fmt.Println("all index exists: ", result)
}

func TestDeleteIndex(t *testing.T) {
	result := DelIndex("localindex1")
	fmt.Println("all index deleted: ", result)
}

func TestCreateIndex(t *testing.T) {
	result := CreateIndex("twitter", mapping)
	fmt.Println("mapping created: ", result)
}

func TestBatch(t *testing.T) {
	tweet1 := Tweet{User: "Jame1",Age: 23, Message: "Take One", Retweets: 1, Created: time.Now()}
	tweet2 := Tweet{User: "中国",Age: 32, Message: "Take Two", Retweets: 0, Created: time.Now()}
	tweet3 := Tweet{User: "Jame3",Age: 32, Message: "Take Three", Retweets: 0, Created: time.Now()}
	Batch("twitter", "_doc", tweet1, tweet2, tweet3)
}

func TestGetDoc(t *testing.T) {
	var tweet Tweet
	data := GetDoc("twitter", "1")
	if err := json.Unmarshal(data, &tweet); err == nil {
		fmt.Printf("data: %v\n", tweet)
	}
}

func TestTermQuery(t *testing.T) {
	var tweet Tweet
	result := TermQuery("twitter", "doc", "user", "Jame1")
	//获得数据, 方法一
	for _, item := range result.Each(reflect.TypeOf(tweet)) {
		if t, ok := item.(Tweet); ok {
			fmt.Printf("tweet : %v\n", t)
		}
	}
	//获得数据, 方法二
	fmt.Println("num of raws: ", result.Hits.TotalHits)
	if result.Hits.TotalHits.Value > 0 {
		for _, hit := range result.Hits.Hits {
			err := json.Unmarshal(*hit.Source, &tweet)
			if err != nil {
				fmt.Printf("source convert json failed, err: %v\n", err)
			}
			fmt.Printf("data: %v\n", tweet)
		}
	}
}

func TestSearch(t *testing.T) {
	result := Search("twitter", "doc")
	var tweet Tweet
	for _, item := range result.Each(reflect.TypeOf(tweet)) {
		if t, ok := item.(Tweet); ok {
			fmt.Printf("tweet : %v\n", t)
		}
	}
}

func TestAggsSearch(t *testing.T) {
	AggsSearch("twitter", "doc")
}

func TestWildcardQuery(t *testing.T) {
	q := elastic.NewWildcardQuery("filename", "技术方案部分（信息中心提供*")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	fmt.Printf("result:[%s]\n",got)

}


func HighLightQuery(){
	result:=es.HighLightQuery("localindex1","doc","content","*资源*")
	var o FileSearch
	for _,h:=range result.Hits.Hits{
		fmt.Printf("highLight [%s]:[%s]\n","content",h.Highlight["content"])
	}
	for _, item := range result.Each(reflect.TypeOf(o)) {
		if t, ok := item.(FileSearch); ok {
			t.Content=""
			fmt.Printf("tweet : %v\n", t)
		}
	}
}