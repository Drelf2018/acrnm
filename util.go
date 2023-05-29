package acrnm

import (
	"io"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

// 基础输出格式
var log = &logrus.Logger{
	Out: os.Stderr,
	Formatter: &nested.Formatter{
		HideKeys:        true,
		NoColors:        true,
		TimestampFormat: "15:04:05",
	},
	Level: logrus.DebugLevel,
}

// 初始化日志
func init() {
	// 尝试输出到文件
	file, err := os.OpenFile("acrnm.log", os.O_CREATE|os.O_WRONLY, 0666)
	if !checkErr(err) {
		return
	}
	writers := []io.Writer{
		file,
		os.Stdout,
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	log.SetOutput(fileAndStdoutWriter)
}

// 获取 log
func GetLog() *logrus.Logger {
	return log
}

// 错误处理
func checkErr(err error) bool {
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

// 爬虫
func Spider() (doc *goquery.Document) {
	resp, err := http.Get(config.Url)
	if !checkErr(err) {
		return nil
	}
	// 关闭链接
	defer resp.Body.Close()

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if !checkErr(err) {
		return nil
	}
	return
}
