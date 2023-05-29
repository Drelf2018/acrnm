package acrnm

import (
	"sort"
	"strings"
)

// 种类
type Variant struct {
	Soldout bool
	Color   string
	Size    string
	Price   string
}

// 小于操作符
func (vi Variant) Less(vj Variant) bool {
	if vi.Color == vj.Color {
		return vi.Size < vj.Size
	}
	return vi.Color < vj.Color
}

// 中嘞生成函数
func makeVariants(color, size, price string) (variants []Variant) {
	for _, s := range strings.Split(size, "/") {
		for _, c := range strings.Split(color, "/") {
			variants = append(variants, Variant{false, c, s, price})
		}
	}
	return
}

// 商品
type Product struct {
	Name     string
	Variants []Variant
}

// Len 方法返回集合中的元素个数
func (product Product) Len() int {
	return len(product.Variants)
}

// Less 方法报告索引 i 的元素是否比索引 j 的元素小
func (product Product) Less(i, j int) bool {
	return product.Variants[i].Less(product.Variants[j])
}

// Swap 方法交换索引 i 和 j 的两个元素
func (product *Product) Swap(i, j int) {
	product.Variants[i], product.Variants[j] = product.Variants[j], product.Variants[i]
}

// Index 方法返回指定 variant 的序号 不存在返回 -1
func (product Product) Index(variant Variant) int {
	i := sort.Search(product.Len(), func(i int) bool {
		return !product.Variants[i].Less(variant)
	})
	if i < product.Len() {
		vi := product.Variants[i]
		if vi.Color == variant.Color && vi.Size == variant.Size {
			// x is present at data[i]
			return i
		}
	}
	// x is not present in data,
	// but i is the index where it would be inserted.
	return -1
}

// Get 方法返回指定 variant 相同数据的指针 不存在返回 nil
func (product *Product) Get(variant Variant) (*Variant, bool) {
	vid := product.Index(variant)
	if vid == -1 {
		return nil, false
	}
	return &(product.Variants[vid]), true
}

// Insert 方法用于顺序插入 variant
func (product *Product) Insert(variant Variant) {
	product.Variants = append(product.Variants, variant)
	sort.Sort(product)
}

// Modify 方法用于修改 variant 价格
func (product *Product) Modify(variant *Variant, price string) {
	variant.Price = price
	sort.Sort(product)
}

// Remove 方法用于删除一系列 variant
func (product *Product) Remove(removeList []int) {
	// 吗的 golang 排个倒序这么麻烦?
	sort.Sort(sort.Reverse(sort.IntSlice(removeList)))
	for _, i := range removeList {
		product.Variants = append(product.Variants[:i], product.Variants[i+1:]...)
	}
}

// 商品生成函数
func makeProduct(name string) *Product {
	return &Product{
		name,
		[]Variant{},
	}
}

// 总商品列表
type ProductList map[string]*Product

// 获取/新建商品
func (products *ProductList) Get(name string) (product *Product) {
	product, ok := (*products)[name]
	if ok {
		return
	}
	product = makeProduct(name)
	(*products)[name] = product
	return
}

// type CombinedProduct struct {
// 	Name  string
// 	Price string
// 	Color []string
// 	Size  []string
// }

// // 生成拼接字符串
// func (products *ProductList) Format() (cpl []CombinedProduct) {
// 	var cp CombinedProduct
// 	for name, product := range *products {
// 		cp = CombinedProduct{name, "", []string{}, []string{}}
// 		for _, variant := range product.Variants {
// 			if cp.Price == "" {
// 				cp.Price = variant.Price
// 			} else if cp.Price != variant.Price {
// 				cpl = append(cpl, cp)
// 				cp = CombinedProduct{name, variant.Price, []string{}, []string{}}
// 			}
// 		}
// 	}
// }
