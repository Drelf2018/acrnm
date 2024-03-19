package acrnm

import (
	"fmt"
	"strings"

	"github.com/Drelf2018/request"
)

type Variant struct {
	Color string
	Size  string
}

// 商品
type Product struct {
	Name     string
	Href     string
	Price    string
	Variants []Variant
}

func (p *Product) Variant() string {
	var r []string
	for _, variant := range p.Variants {
		r = append(r, "- **"+variant.Color+"**: "+variant.Size)
	}
	return strings.Join(r, "\n\n")
}

func (p *Product) MapKey() string {
	return p.Name
}

func (p *Product) Equal(n *Product) bool {
	if p.Name != n.Name {
		return false
	}
	if p.Price != n.Price {
		return false
	}
	if len(p.Variants) != len(n.Variants) {
		return false
	}
	return true
}

func (p *Product) Image() string {
	s := request.Get(fmt.Sprintf("%s/image%s", url, p.Href)).Text()
	return s[1 : len(s)-1]
}

func (p *Product) MdImage() string {
	return fmt.Sprintf("![%s](%s)", p.Name, p.Image())
}

func (p *Product) String() string {
	return fmt.Sprintf("Product(%s, %s)", p.Name, p.Price)
}
