# ACRNM 爬虫

使用了外部接口 [acrnm-api](https://github.com/Drelf2018/acrnm-api)

这是一个通过解析 [acrnm.com/index](https://acrnm.com/index?sort=default&filter=txt) 页面并返回数据的接口。

### 使用

```go
package main

import "github.com/Drelf2018/acrnm"

type Console struct{}

func (Console) Send(p *acrnm.Product) error {
	fmt.Println(p)
	return nil
}

func main() {
	acrnm.Run(Console{})
}
```