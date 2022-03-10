package mocker

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/lucasjones/reggen"
	"math"
	"time"
)

type Generator interface {
	Generate() interface{}
}

type eqGenerator struct {
	value interface{}
}

func (g eqGenerator) Generate() interface{} {
	return g.value
}

type minGenerator struct {
	min int
}

func (g minGenerator) Generate() interface{} {
	return randomdata.Number(g.min)
}

type regexGenerator struct {
	expr string
}

func (g regexGenerator) Generate() interface{} {
	str, err := reggen.Generate(g.expr, math.MaxInt)
	if err != nil {
		return nil
	}

	return str
}

type emptyGenerator struct{}

func (g emptyGenerator) Generate() interface{} {
	return nil
}

type dateGenerator struct {
	layout string
}

func (g dateGenerator) Generate() interface{} {
	return time.Now().Format(g.layout)
}
