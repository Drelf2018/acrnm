package acrnm_test

import (
	"fmt"
	"testing"

	"github.com/Drelf2018/acrnm"
	"github.com/Drelf2018/request"
)

func TestSession(t *testing.T) {
	p := acrnm.GetProducts()[0]
	fmt.Printf("p: %#v\n", p)
	m := request.M{
		"title": p.Name + " " + p.Price,
		"desp":  p.Variant() + "\n\n" + p.MdImage(),
	}
	fmt.Printf("m: %v\n", m)
}
