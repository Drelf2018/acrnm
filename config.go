package acrnm

import (
	"flag"

	"github.com/Drelf2018/initial"
	"github.com/Drelf2020/utils"
	"gopkg.in/yaml.v2"
)

type XPath struct {
	List     string `default:"tbody > .m-product-table__row"`
	Name     string `default:".m-product-table__title_cell span"`
	Price    string `default:".m-product-table__price_cell span"`
	Variants string `default:".m-product-table__variant_cell span[class=o-item-row]"`
	Color    string `default:"div > span"`
	Size     string `default:"div + span"`
}

type Config struct {
	Interval float64 `default:"10"`
	Url      string  `default:"https://acrnm.com/?sort=default&filter=txt"`
	XPath    XPath   `default:"initial.Default"`
	Chrome   Chrome  `default:"initial.Default"`
}

// 获取命令行参数
func FilePath() (filepath string) {
	flag.StringVar(&filepath, "config", "config.yml", "配置文件路径")
	flag.Parse()
	return
}

// 读取配置文件
func ReadConfig() *Config {
	var conf Config
	s := utils.ReadFile(FilePath())
	err := yaml.Unmarshal([]byte(s), &conf)
	if utils.LogErr(err) {
		panic(err)
	}
	return initial.Default(&conf)
}
