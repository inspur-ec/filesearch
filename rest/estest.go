package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Queryhighlight(ind string, str ...string) {
	urls := "http://127.0.0.1:9200/" + ind + "/_search?pretty"
	jsons := `{"query":{ "wildcard":{ "content":"*` + str[0] + `*" } }}`
	contentType := "application/json;charset=utf-8"
	//javaJsonParam, err := json.Marshal(jsons)
	fmt.Printf("jsons:[%s]", jsons)
	req, err := http.NewRequest("GET", urls, strings.NewReader(jsons))
	req.Header.Set("Content-Type", contentType)
	client := &http.Client{}
	res, err := client.Do(req)
	//res,err:=http.Get(urls)

	if err != nil {
		fmt.Printf("error is occur [%v]", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}

func InsertIndex(str ...string) {
	urls := "http://127.0.0.1:9200/localindex/doc/1?pretty"
	jsons := `{
   "content":"` + str[0] + `"
}`
	contentType := "application/json;charset=utf-8"
	//javaJsonParam, err := json.Marshal(jsons)
	fmt.Printf("jsons:[%s]", jsons)
	req, err := http.NewRequest("PUT", urls, strings.NewReader(jsons))
	req.Header.Set("Content-Type", contentType)
	client := &http.Client{}
	res, err := client.Do(req)
	//res,err:=http.Get(urls)

	if err != nil {
		fmt.Printf("error is occur [%v]", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}
