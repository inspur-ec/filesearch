package es

import (
	"bytes"
	"context"
	"encoding/json"
	"filesearch/conf"
	"filesearch/office"
	"fmt"
	"github.com/elastic/go-elasticsearch/v6"
	"log"
	"os/exec"
	"strings"
	"time"
)

type FileSearch struct {
	Filetype string `json:"filetype"`
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type SearchResult struct {
	Hits Hit `json:"hits"`
}
type Hit struct {
	Total Total   `json:"total"`
	Score float64 `json:"max_score"`
	Hits  []Hits  `json:"hits"`
}
type Total struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}
type Hits struct {
	Index     string    `json:"_index"`
	Type1     string    `json:"_type"`
	Id        string    `json:"_id"`
	Score     float64   `json:"_score"`
	Source    Source    `json:"_source"`
	Highlight Highlight `json:"highlight"`
}
type Source struct {
	Filetype string `json:"filetype"`
	Filename string `json:"filename"`
	Content  string `json:"content"`
}
type Highlight struct {
	Content []string `json:"content"`
}

var client *elasticsearch.Client

func init() {
	addresses := conf.Conf.EsServer.Address
	config := elasticsearch.Config{
		Addresses: addresses,
		Username:  "",
		Password:  "",
		CloudID:   "",
		APIKey:    "",
	}
	// new client
	var err error
	client, err = elasticsearch.NewClient(config)
	failOnError(err, "Error creating the client")
	//如果es未启动，调用启动命令
	_, err = client.Info()
	if err != nil {
		cmd := exec.Command(conf.Conf.EsServer.StartCmd)
		if err := cmd.Start(); err != nil { // 运行命令
			log.Fatal(err)
		}
		//不确定多长时间能启动，暂时休眠5s
		time.Sleep(time.Second * 10)
	}
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func Index(indexname, indextype string, v interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		return err
	}
	res, err := client.Index(indexname, &buf, client.Index.WithDocumentType(indextype))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func FormatResult(res string) SearchResult {
	res = strings.TrimLeft(res, "[200 OK]")
	var result SearchResult
	b := []byte(res)
	json.Unmarshal(b, &result)
	log.Printf("result is [%v]", result)
	return result
}

func SearchALL(content string, filetype string, color  string) string {
	start := time.Now()
	// info
	res, err := client.Info()
	failOnError(err, "Error getting response")
	//fmt.Println(res.String())
	// search - highlight
	var buf bytes.Buffer
	query := map[string]interface{}{
		"_source": []string{"filename"}, //不需要返回值时可以用false
		"_type":filetype,  //该类型不对
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"content": content,
			},
		},
		"highlight": map[string]interface{}{
			"pre_tags":  []string{"<font color='" + color + "'>"},
			"post_tags": []string{"</font>"},
			"fields": map[string]interface{}{
				"content": map[string]interface{}{},
			},
		},
	}
	querys,_:=json.Marshal(query)
	log.Printf("query :%s",querys)
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		failOnError(err, "Error encoding query")
	}
	// Perform the search request.
	res, err = client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex("localindex"),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
		client.Search.WithSize(1000),
		client.Search.WithSearchType("dfs_query_then_fetch"),
	)
	if err != nil {
		failOnError(err, "Error getting response")
	}
	defer res.Body.Close()
	res1 := res.String()
	log.Printf("content[%s]\n", res1)
	log.Printf("查询消耗时间为[%d]s\n", time.Since(start)/1e9)
	//log.Printf(res1)
	return res1
}

func Search(content, color string) string {
	start := time.Now()
	// info
	res, err := client.Info()
	failOnError(err, "Error getting response")
	//fmt.Println(res.String())
	// search - highlight
	var buf bytes.Buffer
	query := map[string]interface{}{
		"_source": []string{"filename"}, //不需要返回值时可以用false
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"content": content,
			},
		},
		"highlight": map[string]interface{}{
			"pre_tags":  []string{"<font color='" + color + "'>"},
			"post_tags": []string{"</font>"},
			"fields": map[string]interface{}{
				"content": map[string]interface{}{},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		failOnError(err, "Error encoding query")
	}
	// Perform the search request.
	res, err = client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex("localindex"),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
		client.Search.WithSize(1000),
		client.Search.WithSearchType("dfs_query_then_fetch"),
	)
	if err != nil {
		failOnError(err, "Error getting response")
	}
	defer res.Body.Close()
	res1 := res.String()
	log.Printf("content[%s]\n", res1)
	log.Printf("查询消耗时间为[%d]s\n", time.Since(start)/1e9)
	//log.Printf(res1)
	return res1
}
func DeleteByQuery(q string) {

	// DeleteByQuery deletes documents matching the provided query
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"wildcard": q,
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		failOnError(err, "Error encoding query")
	}
	index := []string{"localindex"}
	res, err := client.DeleteByQuery(index, &buf)
	if err != nil {
		failOnError(err, "Error delete by query response")
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}
func Delete(indexname string) {

	// Delete removes a document from the index
	res, err := client.Indices.Delete([]string{indexname})
	if err != nil {
		failOnError(err, "Error delete by id response")
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

func Create(indexname string, filename string, suf string) error {
	suffix := strings.ToLower(suf)
	c := ""
	var err error
	if "docx" == suffix {
		c, err = office.ReadXml(office.ReadWord(filename))
	} else if "xlsx" == suffix {
		c, err = office.ReadExcel(filename)
	} else if "pdf" == suffix {
		c, err = office.ReadPdf(filename)
	}
	if err != nil {
		return err
	}
	f := FileSearch{Filetype: suffix, Filename: filename, Content: c}
	err = Index(indexname, "doc", f)
	return err
}

func Get() {

	res, err := client.Get("demo", "esd")
	if err != nil {
		failOnError(err, "Error get response")
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}
func Update() {

	// Update updates a document with a script or partial document.
	var buf bytes.Buffer
	doc := map[string]interface{}{
		"doc": map[string]interface{}{
			"title":   "更新你看到外面的世界是什么样的？",
			"content": "更新外面的世界真的很精彩",
		},
	}
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		failOnError(err, "Error encoding doc")
	}
	res, err := client.Update("demo", "esd", &buf, client.Update.WithDocumentType("doc"))
	if err != nil {
		failOnError(err, "Error Update response")
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

func UpdateByQuery() {

	// UpdateByQuery performs an update on every document in the index without changing the source,
	// for example to pick up a mapping change.
	index := []string{"demo"}
	var buf bytes.Buffer
	doc := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "外面",
			},
		},
		// 根据搜索条件更新title
		/*
		   "script": map[string]interface{}{
		       "source": "ctx._source['title']='更新你看到外面的世界是什么样的？'",
		   },
		*/
		// 根据搜索条件更新title、content
		/*
		   "script": map[string]interface{}{
		       "source": "ctx._source=params",
		       "params": map[string]interface{}{
		           "title": "外面的世界真的很精彩",
		           "content": "你看到外面的世界是什么样的？",
		       },
		       "lang": "painless",
		   },
		*/
		// 根据搜索条件更新title、content
		"script": map[string]interface{}{
			"source": "ctx._source.title=params.title;ctx._source.content=params.content;",
			"params": map[string]interface{}{
				"title":   "看看外面的世界真的很精彩",
				"content": "他们和你看到外面的世界是什么样的？",
			},
			"lang": "painless",
		},
	}
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		failOnError(err, "Error encoding doc")
	}
	res, err := client.UpdateByQuery(
		index,
		client.UpdateByQuery.WithDocumentType("doc"),
		client.UpdateByQuery.WithBody(&buf),
		client.UpdateByQuery.WithContext(context.Background()),
		client.UpdateByQuery.WithPretty(),
	)
	if err != nil {
		failOnError(err, "Error Update response")
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}
