package acrnm

import (
	"fmt"
	"io"
	"strings"

	"github.com/Drelf2020/utils"
	"github.com/tebeka/selenium"
)

var (
	service *selenium.Service
	driver  selenium.WebDriver
)

type Chrome struct {
	Port int64  `default:"9000"`
	Path string `default:"./chromedriver.exe"`
}

func Close() {
	driver.Quit()
	service.Stop()
}

func start(c Chrome) {
	var err error
	// 启动 Chrome 浏览器
	service, err = selenium.NewChromeDriverService(c.Path, int(c.Port))
	if utils.LogErr(err) {
		panic(err)
	}

	// 添加无头参数
	caps := selenium.Capabilities{"browserName": "chrome"}
	// caps.AddChrome(chrome.Capabilities{Args: []string{"--headless"}})

	// 设置全局变量
	driver, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", c.Port))
	if utils.LogErr(err) {
		panic(err)
	}
}

func reader() (io.Reader, error) {
	html, err := driver.PageSource()
	if err != nil {
		return nil, err
	}
	return strings.NewReader(html), nil
}

func Get(url string) (io.Reader, error) {
	err := driver.Get(url) // "https://postman-echo.com/get"
	if err != nil {
		return nil, err
	}
	return reader()
}

func Refresh() (io.Reader, error) {
	err := driver.Refresh()
	if err != nil {
		return nil, err
	}
	return reader()
}
