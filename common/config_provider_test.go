package common_test

import (
	"github.com/byorty/contractor/common"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/config"
)

type Config struct {
	A         A
	Partition int
}

type A struct {
	B string
	C struct {
		D bool
		F int
	}
}

func TestConfigProviderSuite(t *testing.T) {
	suite.Run(t, new(ConfigProviderSuite))
}

type ConfigProviderSuite struct {
	suite.Suite
}

func (s *ConfigProviderSuite) TestPopulate() {
	reader := strings.NewReader("a: {b: bar, c: {d: true, f: 12}}")

	factory := common.NewFxConfigProviderFactory()
	provider, err := factory.CreateByOptions(config.Source(reader))
	s.Nil(err)

	var a A
	s.Nil(provider.PopulateByKey("a", &a))
	s.Equal("bar", a.B)
	s.Equal(true, a.C.D)
	s.Equal(12, a.C.F)

	var f int
	s.Nil(provider.PopulateByKey("a.c.f", &f))
	s.Equal(12, f)

	var cfg Config
	s.Nil(provider.Populate(&cfg))
	s.Equal(cfg.A, a)
	s.Equal(cfg.A.B, a.B)
}

func (s *ConfigProviderSuite) TestExpand() {
	var a int
	var b string
	varB := "hello world"
	err := os.Setenv("VAR_B", varB)
	if err != nil {
		s.Error(err)
	}

	reader := strings.NewReader(`
a: 1
b: "$VAR_B"
`)

	factory := common.NewFxConfigProviderFactory()
	provider, err := factory.CreateByOptions(config.Source(reader))
	s.Nil(err)

	s.Nil(provider.PopulateByKey("a", &a))
	s.Equal(1, a)

	s.Nil(provider.PopulateByKey("b", &b))
	s.Equal(varB, b)

	readerArray := strings.NewReader(`
rss:
  - name: forbes
    url: https://www.forbes.ru/newrss.xml
    num: 4
    m: 7
`)
	type Rss struct {
		Name string `yaml:"name"`
		Url  string `yaml:"url"`
		Num  int    `yaml:"num"`
		M    int    `yaml:"m"`
	}
	c := make([]Rss, 0)
	provider, err = factory.CreateByOptions(config.Source(readerArray))
	s.Nil(err)
	err = provider.PopulateByKey("rss", &c)
	s.Nil(err)
	s.Equal("forbes", c[0].Name)
	s.Equal("https://www.forbes.ru/newrss.xml", c[0].Url)
	s.Equal(4, c[0].Num)
	s.Equal(7, c[0].M)
}
