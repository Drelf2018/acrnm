package acrnm

import (
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var config = GetConfig()

type Config struct {
	Interval float64
	Url      string
	XPath    struct {
		List     string
		Name     string
		Price    string
		Variants string
		Color    string
		Size     string
	}
}

// 获取命令行参数
func FilePath() (filepath string) {
	flag.StringVar(&filepath, "config", "config.yml", "配置文件路径")
	flag.Parse()
	return
}

// 读取配置文件
func GetConfig() (conf Config) {
	yamlFile, err := ioutil.ReadFile(FilePath())
	if !checkErr(err) {
		panic("配置文件读取失败")
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if !checkErr(err) {
		panic("配置文件解析失败")
	}
	return
}
