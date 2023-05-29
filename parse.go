package acrnm

import (
	"sort"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var products = make(ProductList)

type Handler func(string, ...Variant)

var UpdateHandler, ModifyHandler, SoldoutHandler Handler

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

var TotalLines Lines

// 获取文本
func Text(s *goquery.Selection, xpath string) string {
	return s.Find(xpath).Text()
}

// 执行函数
func Exec(s *goquery.Selection, xpath string, f func(int, *goquery.Selection)) {
	s.Find(xpath).Each(f)
}

// 获取商品细节变动
func GetDetail(product *Product, price string) func(int, *goquery.Selection) {
	return func(_ int, s *goquery.Selection) {
		variants := makeVariants(Text(s, config.XPath.Color), Text(s, config.XPath.Size), price)
		for _, v := range variants {
			variant, ok := product.Get(v)
			if !ok {
				product.Insert(v)
				log.Infof("上新 %v %v %v %v", product.Name, v.Color, v.Size, v.Price)
				if UpdateHandler != nil {
					go UpdateHandler(product.Name, v)
				}
			} else {
				if variant.Price != v.Price && variant.Soldout {
					log.Infof("价格变化 %v %v %v %v -> %v", product.Name, v.Color, v.Size, variant.Price, v.Price)
					if ModifyHandler != nil {
						go ModifyHandler(product.Name, *variant, v)
					}
					product.Modify(variant, v.Price)
				}
				variant.Soldout = false
			}
		}
	}
}

// 获取某一行商品数据
func GetRow(_ int, s *goquery.Selection) {
	price := Text(s, config.XPath.Price)
	if price == "" {
		return
	}
	name := Text(s, config.XPath.Name)
	TotalLines = append(TotalLines, Line{name, price, s})
}

// 获取全部商品数据
func GetData() {
	doc := Spider()
	if doc == nil {
		return
	}

	// 先获取所有商品 再根据名称和价格排序
	TotalLines = make(Lines, 0)
	Exec(doc.Selection, config.XPath.List, GetRow)

	// 解析单行信息
	var product *Product
	for _, line := range *TotalLines.Sort() {
		if product == nil || product.Name != line.Name {
			product = products.Get(line.Name)
		}
		Exec(line.Selection, config.XPath.Variants, GetDetail(product, line.Price))
	}

	// 查找下架商品
	for name, product := range products {
		var removeList []int
		for i, variant := range product.Variants {
			if variant.Soldout {
				log.Infof("下架 %v %v %v %v", product.Name, variant.Color, variant.Size, variant.Price)
				removeList = append(removeList, i)
				if SoldoutHandler != nil {
					go SoldoutHandler(name, variant)
				}
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

// 设置处理器
func SetHandler(name string, handler Handler) {
	switch name {
	case "update":
		UpdateHandler = handler
	case "modify":
		ModifyHandler = handler
	case "soldout":
		SoldoutHandler = handler
	}
}

func None(*ProductList) {}

// 轮询
func Interval(initial func(*ProductList), finish func(*ProductList)) {
	ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
	defer ticker.Stop()
	GetData()
	initial(&products)
	for range ticker.C {
		GetData()
		finish(&products)
	}
}
