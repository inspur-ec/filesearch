package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type BaseConfig struct {
	Paths    Paths    `yaml:"paths"`
	EsServer EsServer `yaml:"es"`
	Monitor  Monitor  `yaml:"monitor"`
	Log      Log      `yaml:"log"`
}

//解析yml文件
type Paths struct {
	Include []string `yaml:"include,flow"`
	Exclude []string `yaml:"exclude,flow"`
}

type EsServer struct {
	Address  []string `yaml:"address,flow"`
	StartCmd string   `yaml:"startcmd"`
}

type Monitor struct {
	Enable     bool   `yaml:"enable"`
	Outpath    string `yaml:"outpath"`
	Filescount int    `yaml:"filescount"`
	Interval   int    `yaml:"interval"`
}
type Log struct {
	Path string `yaml:"path"`
}

var Conf BaseConfig

func init() {

	//定义默认值
	Conf = BaseConfig{Monitor: Monitor{Enable: false, Filescount: 30, Interval: 10}}
	yamlFile, err := ioutil.ReadFile("config.yaml")
	fmt.Printf("[%v]", yamlFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		panic(err)
	}
	return

}

func GetMonitor() Monitor {
	return Conf.Monitor
}
func GetEsServer() EsServer {
	return Conf.EsServer
}
func GetLog() Log {
	return Conf.Log
}
func GetPath() Paths {
	return Conf.Paths
}
