package main

import (
	"encoding/json"
	"fmt"
	"filesearch/conf"
	"filesearch/gui"
	"github.com/olivere/elastic"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"runtime/pprof"
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
		"number_of_shards": 1,
		"number_of_replicas": 0
	}
}`

func init() {
	logFile, err := os.OpenFile(conf.Conf.Log.Path+"/filesearch.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
}
func main() {
	log.Printf("start filesearch ...")
	b, _ := json.Marshal(conf.Conf)
	log.Printf("config is [%s]\n", string(b))
	if conf.Conf.Monitor.Enable {
		//w := bufio.NewWriter(file)
		go func() {
			for i := 0; ; i++ {
				if i > conf.Conf.Monitor.Filescount {
					i = 0
				}
				file, _ := os.Create(conf.Conf.Monitor.Outpath + "cpu." + strconv.Itoa(i) + ".pprof")
				file1, _ := os.Create(conf.Conf.Monitor.Outpath + "mem." + strconv.Itoa(i) + ".pprof")
				pprof.StartCPUProfile(file)
				time.Sleep(time.Duration(conf.Conf.Monitor.Interval) * time.Second)
				pprof.WriteHeapProfile(file1)
				pprof.StopCPUProfile()
				file1.Close()
				file.Close()

			}
		}()
	}
	//监控好像很慢，暂时去掉
	//	monitor.MonitorFiles(conf.GetPath().Include)
	gui.OpenSearchWindows() //打开主界面
	//gui.Createlinklabel()  //测试用
	//es.Delete()
	//InitIndex()
	//ss:=es.Search("content","*中台*","green")
	//es.FormatResult(ss)
	//Create()
	//Query()
	//QueryString("技术")
	//HighLightQuery()
	//TestCreate()
	//rest.InsertIndex("中台")
	//rest.Queryhighlight("localindex","中台")
	//es.Index()
	//es.Search()
}

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth string, outpath string, suffix ...string) (files map[string][]string, err error) {
	files = make(map[string][]string)
	for i, _ := range suffix {
		suffix[i] = strings.ToUpper(suffix[i])
	} //忽略后缀匹配的大小写
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		//if err != nil { //忽略错误
		// return err
		//}
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if !strings.HasPrefix(fi.Name(), "~") {
			for j, _ := range suffix {
				if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix[j]) {
					files[suffix[j]] = append(files[suffix[j]], filename)
				}
			}
		}
		return nil
	})
	return files, err
}
