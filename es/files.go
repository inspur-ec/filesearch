package es

import (
	"filesearch/conf"
	"filesearch/utils"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const Idx_name string = "localindex"

func InitIndex() (success int, errors int, errs string) {
	urls := conf.Conf.EsServer.Address[0]
	res, err := http.Get(urls + "/" + Idx_name)
	if err != nil {
		return 0, 0, fmt.Sprintf("%v", err)
	}
	//如果索引不存在，则设置索引分词器为IK；自动创建一个索引
	if res.StatusCode == 404 {
		body := `{"settings": {"index" : {"analysis.analyzer.default.type":"ik_max_word"}}}`
		req, _ := http.NewRequest("PUT", urls+"/"+Idx_name, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		resp, _ := http.DefaultClient.Do(req)
		if err != nil {
			return 0, 0, fmt.Sprintf("%v", err)
		}
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Sprintf("%s", b)
	}
	return AddIndex(conf.Conf.Paths.Include, conf.Conf.Paths.Exclude)
}
func AddIndex(inpaths, outpaths []string) (success int, errors int, errs string) {
	success = 0
	errors = 0
	//并发设置为8个
	c := make(chan int, 8)
	for _, inpath := range inpaths {
		files := utils.ListAllFileByName(inpath, outpaths, "docx", "xlsx", "pdf")

		for suf, f := range files {
			//es.Index("localindex1","doc",f)
			log.Printf("开始创建[%s]的索引，共计[%d]条", suf, len(f))
			for _, s := range f {
				//wg.Add(1)
				c <- 1
				go func(s string, suf string, errs *string) {
					//创建索引
					err := Create(Idx_name, s, suf)
					if err != nil {
						*errs = *errs + fmt.Sprintf("file [%s] is error :[%v]\n", s, err)
						errors++
					} else {
						success++
					}
					<-c
				}(s, suf, &errs)
			}
		}
	}
	//wg.Wait()
	log.Printf("索引创建完毕,创建[%d]个文件索引，失败[%d]个文件，发生错误：\n[%s]", success, errors, errs)
	return
}
