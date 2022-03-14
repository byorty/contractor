package mocker

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/byorty/contractor/common"
	"github.com/lucasjones/reggen"
	"math"
	"time"
)

type Generator interface {
	Generate() interface{}
}

func NewEqGenerator(value interface{}) Generator {
	return &eqGenerator{value: value}
}

type eqGenerator struct {
	value interface{}
}

func (g eqGenerator) Generate() interface{} {
	return g.value
}

func NewPositiveGenerator() Generator {
	return &minGenerator{min: 1}
}

func NewMinGenerator(min int) Generator {
	return &minGenerator{min: min}
}

type minGenerator struct {
	min int
}

func (g minGenerator) Generate() interface{} {
	return randomdata.Number(g.min, math.MaxInt64)
}

func NewMaxGenerator(max int) Generator {
	return &maxGenerator{max: max}
}

type maxGenerator struct {
	max int
}

func (g maxGenerator) Generate() interface{} {
	return randomdata.Number(math.MinInt32, g.max)
}

func NewRangeGenerator(min, max int) Generator {
	return &rangeGenerator{
		min: min,
		max: max,
	}
}

type rangeGenerator struct {
	min int
	max int
}

func (g rangeGenerator) Generate() interface{} {
	return randomdata.Number(g.min, g.max)
}

func NewRegexGenerator(expr string) Generator {
	return &regexGenerator{expr: expr}
}

type regexGenerator struct {
	expr string
}

func (g regexGenerator) Generate() interface{} {
	str, err := reggen.Generate(g.expr, randomdata.Number(10, 1000))
	if err != nil {
		return nil
	}

	return str
}

func NewEmptyGenerator() Generator {
	return &emptyGenerator{}
}

type emptyGenerator struct{}

func (g emptyGenerator) Generate() interface{} {
	return nil
}

func NewDateGenerator(layout string) Generator {
	switch layout {
	case common.ExpressionDateLayoutRFC3339:
		layout = time.RFC3339
	case common.ExpressionDateLayoutRFC3339Nano:
		layout = time.RFC3339Nano
	}

	return &dateGenerator{layout: layout}
}

type dateGenerator struct {
	layout string
}

func (g dateGenerator) Generate() interface{} {
	return time.Now().Format(g.layout)
}

func NewContainsGenerator(sub string) Generator {
	return &containsGenerator{
		sub: sub,
	}
}

type containsGenerator struct {
	sub string
}

func (g containsGenerator) Generate() interface{} {
	return fmt.Sprintf(randomdata.Paragraph(), g.sub)
}
