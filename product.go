package acrnm

import (
	"bytes"
	"fmt"
	"sync"
)

// 商品
type Product struct {
	Name     string `cmps:"1"`
	Price    string `cmps:"2"`
	Variants sync.Map
}

func (p *Product) Variant() (r []string) {
	p.Variants.Range(func(key, value any) bool {
		r = append(r, fmt.Sprintf("%v %v", key, value))
		return true
	})
	return
}

func (p *Product) String() string {
	buf := bytes.NewBufferString("Product(")
	buf.WriteString(p.Name)
	buf.WriteString(", ")
	buf.WriteString(p.Price)
	p.Variants.Range(func(key, value any) bool {
		buf.WriteString(", ")
		buf.WriteString(key.(string))
		buf.WriteString(": ")
		buf.WriteString(value.(string))
		return true
	})
	buf.WriteString(")")
	return buf.String()
}
