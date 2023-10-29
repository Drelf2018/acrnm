package acrnm

import (
	"io"
	"sync"

	"github.com/Drelf2018/asyncio"
	"github.com/Drelf2018/cmps"
	"github.com/Drelf2018/event"
	"github.com/Drelf2020/utils"

	_ "unsafe"
)

// 初始化日志
var log = utils.SetOutputFile("acrnm.log")

type Parser struct {
	XPath
	event.AsyncEventS[*Product]
	Pre cmps.SafeSlice[*Product]
	New cmps.SafeSlice[*Product]
}

func (p *Parser) Get(name, price string) *Product {
	product := &Product{Name: name, Price: price, Variants: sync.Map{}}
	if o := p.New.Search(product); o != nil {
		return o
	}
	p.New.Insert(product)
	return product
}

// 获取全部商品数据
func (p *Parser) Parse(r io.Reader) {
	root := FromReader(r)
	if root == nil {
		return
	}
	// 先获取所有商品 再根据名称和价格排序
	root.XPath(p.List, func(list *Selection) {
		price := list.Text(p.Price)
		if price == "" {
			return
		}
		name := list.Text(p.Name)
		product := p.Get(name, price)
		list.XPath(p.Variants, func(s *Selection) {
			product.Variants.Store(s.Text(p.Color), s.Text(p.Size))
		})
	})

	asyncio.ForEach(p.New.I, func(product *Product) {
		o := p.Pre.Search(product)
		if o == nil {
			p.Dispatch("new", product)
			log.Infof("上新 %v", product)
		}
	})

	p.Pre.I, p.New.I = p.New.I, make([]*Product, 0)
	log.Infoln("更新完成")
}

func (p *Parser) Run(interval float64) {
	defer Close()
	event.Heartbeat(interval, interval, func(e *event.Event, count int) {
		p.Pre.Delete(&Product{Name: "SAC-J6010", Price: "1,883.00 EUR"})
		r, err := Refresh()
		if utils.LogErr(err) {
			return
		}
		p.Parse(r)
	})
}

func New(c *Config) *Parser {
	start(c.Chrome)
	spider := Parser{
		XPath:       c.XPath,
		AsyncEventS: event.Default[*Product](),
		Pre:         cmps.SafeSlice[*Product]{I: make([]*Product, 0)},
		New:         cmps.SafeSlice[*Product]{I: make([]*Product, 0)},
	}
	r, err := Get(c.Url)
	if utils.LogErr(err) {
		panic(err)
	}
	spider.Parse(r)
	return &spider
}
