package main

import (
	"fmt"
	"strings"

	"github.com/Drelf2018/acrnm"
	"github.com/Drelf2018/event"
)

func main() {
	conf := acrnm.ReadConfig()
	spider := acrnm.New(conf)
	sender := Senders{}
	spider.OnCommand("new", func(e *event.Event, v ...*acrnm.Product) {
		p := v[0]
		msg := fmt.Sprintf("%v %v\n%v", p.Name, p.Price, strings.Join(p.Variant(), "\n"))
		sender.Send(msg)
	})
	spider.Run(conf.Interval)
}
