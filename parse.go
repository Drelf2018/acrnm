package acrnm

import (
	"sort"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// 单行商品
type Line struct {
	Name      string
	Price     string
	Selection *goquery.Selection
}

// 小于操作符
func (li Line) Less(lj Line) bool {
	if li.Name == lj.Name {
		return li.Price < lj.Price
	}
	return li.Name < lj.Name
}

// 多行商品
type Lines []Line

// Len 方法返回集合中的元素个数
func (lines Lines) Len() int {
	return len(lines)
}

// Less 方法报告索引 i 的元素是否比索引 j 的元素小
func (lines Lines) Less(i, j int) bool {
	return lines[i].Less(lines[j])
}

// Swap 方法交换索引 i 和 j 的两个元素
func (lines *Lines) Swap(i, j int) {
	(*lines)[i], (*lines)[j] = (*lines)[j], (*lines)[i]
}

// Sort 方法返回排序后梓神
func (lines *Lines) Sort() *Lines {
	sort.Sort(lines)
	return lines
}

// 获取文本
func Text(s *goquery.Selection, xpath string) string {
	return s.Find(xpath).Text()
}

// 执行函数
func Exec(s *goquery.Selection, xpath string, f func(int, *goquery.Selection)) {
	s.Find(xpath).Each(f)
}

// 获取商品细节变动
func (spider *Spider) GetDetail(product *Product, price string) func(int, *goquery.Selection) {
	return func(_ int, s *goquery.Selection) {
		variants := makeVariants(Text(s, config.XPath.Color), Text(s, config.XPath.Size), price)
		for _, v := range variants {
			variant, ok := product.Get(v)
			if !ok {
				product.Insert(v)
				log.Infof("上新 %v %v %v %v", product.Name, v.Color, v.Size, v.Price)
				go spider.UpdateHandler(product.Name, v)
			} else {
				if variant.Price != v.Price && variant.Soldout {
					log.Infof("价格变化 %v %v %v %v -> %v", product.Name, v.Color, v.Size, variant.Price, v.Price)
					go spider.ModifyHandler(product.Name, *variant, v)
					product.Modify(variant, v.Price)
				}
				variant.Soldout = false
			}
		}
	}
}

// 获取某一行商品数据
func (spider *Spider) GetRow(_ int, s *goquery.Selection) {
	price := Text(s, config.XPath.Price)
	if price == "" {
		return
	}
	name := Text(s, config.XPath.Name)
	spider.TotalLines = append(spider.TotalLines, Line{name, price, s})
}

// 获取全部商品数据
func (spider *Spider) GetData() {
	doc := Request()
	if doc == nil {
		return
	}

	// 先获取所有商品 再根据名称和价格排序
	spider.TotalLines = make(Lines, 0)
	Exec(doc.Selection, config.XPath.List, spider.GetRow)

	// 解析单行信息
	var product *Product
	for _, line := range *spider.TotalLines.Sort() {
		if product == nil || product.Name != line.Name {
			product = spider.Products.Get(line.Name)
		}
		Exec(line.Selection, config.XPath.Variants, spider.GetDetail(product, line.Price))
	}

	// 查找下架商品
	for name, product := range spider.Products {
		var removeList []int
		for i, variant := range product.Variants {
			if variant.Soldout {
				log.Infof("下架 %v %v %v %v", product.Name, variant.Color, variant.Size, variant.Price)
				removeList = append(removeList, i)
				go spider.SoldoutHandler(name, variant)
			} else {
				// 每次上新或者重新获取到的时候会把这个设置为 false 例如 Line40
				// 在获取完一次数据后就不会被上面 SoldoutHandler 检测到
				(*product).Variants[i].Soldout = true
			}
		}
		if len(removeList) != 0 {
			product.Remove(removeList)
		}
	}
	log.Infoln("更新完成")
}

type Handler func(string, ...Variant)

func None(string, ...Variant) {}

type Spider struct {
	UpdateHandler  Handler
	ModifyHandler  Handler
	SoldoutHandler Handler
	Products       ProductList
	TotalLines     Lines
}

// 设置处理器
func (spider *Spider) SetHandler(name string, handler Handler) {
	switch name {
	case "update":
		spider.UpdateHandler = handler
	case "modify":
		spider.ModifyHandler = handler
	case "soldout":
		spider.SoldoutHandler = handler
	}
}

// 轮询
func (spider *Spider) Interval(emit func()) {
	ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		spider.GetData()
		go emit()
	}
}

func NewSpider() *Spider {
	spider := Spider{None, None, None, make(ProductList), nil}
	spider.GetData()
	return &spider
}
