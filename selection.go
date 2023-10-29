package acrnm

import (
	"io"

	"github.com/Drelf2018/asyncio"
	"github.com/PuerkitoBio/goquery"
)

type Selection struct {
	*goquery.Selection
}

// 获取文本
func (s *Selection) Text(xpath string) string {
	return s.Selection.Find(xpath).Text()
}

// 获取子元素
func (s *Selection) XPath(xpath string, f func(*Selection)) {
	each := make([]Selection, 0)
	s.Selection.Find(xpath).Each(func(i int, s *goquery.Selection) {
		each = append(each, Selection{s})
	})
	asyncio.ForEachPtr(each, f)
}

// 新建根节点
func FromReader(r io.Reader) *Selection {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil
	}
	return &Selection{doc.Selection}
}
