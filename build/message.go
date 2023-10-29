package main

import (
	"github.com/Drelf2018/asyncio"
	"github.com/Drelf2018/request"
	"github.com/Drelf2020/utils"
)

type Sender interface {
	Send(string)
}

type Senders []Sender

func (s Senders) Send(msg string) {
	asyncio.ForEach(s, func(s Sender) {
		s.Send(msg)
	})
}

type FangTang struct {
	token string
}

func (f FangTang) Send(msg string) {
	resp := request.Post("https://sctapi.ftqq.com/"+f.token+".send", request.Data("title", msg))
	utils.LogErr(resp.Error())
}

type Guild struct {
	url string
}

func (g Guild) Send(msg string) {
	resp := request.Get(g.url, request.Data("msg", msg))
	utils.LogErr(resp.Error())
}

type QQ struct {
	url string
	uid string
}

func (q QQ) Send(msg string) {
	resp := request.Get(q.url, request.Data("user_id", q.uid, "message", msg))
	utils.LogErr(resp.Error())
}
