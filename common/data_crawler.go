package common

import (
	"fmt"
	"strings"
)

//go:generate mockgen -source=$GOFILE -package=mocks -destination=mocks/$GOFILE

type DataCrawlerOption func(settings *DataCrawlerSettings)

type DataCrawlerSettings struct {
	key               string
	value             interface{}
	isJoinKeys        bool
	isSkipCollections bool
}

func WithJoinKeys() DataCrawlerOption {
	return func(settings *DataCrawlerSettings) {
		settings.isJoinKeys = true
	}
}

func WithSkipCollections() DataCrawlerOption {
	return func(settings *DataCrawlerSettings) {
		settings.isSkipCollections = true
	}
}

func WithPrefix(prefix string) DataCrawlerOption {
	return func(settings *DataCrawlerSettings) {
		isArray := strings.HasPrefix(settings.key, "[") && strings.HasSuffix(settings.key, "]")
		if len(prefix) > 0 && !strings.HasPrefix(settings.key, prefix) {
			if isArray {
				settings.key = fmt.Sprintf("%s%s", prefix, settings.key)
			} else {
				settings.key = fmt.Sprintf("%s.%s", prefix, settings.key)
			}
		}
	}
}

type DataCrawlerHandler func(k string, v interface{})

type DataCrawler interface {
	Walk(data interface{}, handler DataCrawlerHandler, opts ...DataCrawlerOption)
}

func NewFxDataCrawler() DataCrawler {
	return &dataCrawler{}
}

type dataCrawler struct {
}

func (m dataCrawler) Walk(data interface{}, handler DataCrawlerHandler, opts ...DataCrawlerOption) {
	settings := &DataCrawlerSettings{
		value: data,
	}

	for _, opt := range opts {
		opt(settings)
	}

	if !settings.isSkipCollections {
		handler(settings.key, settings.value)
	}

	switch body := data.(type) {
	case map[string]interface{}:
		for key, value := range body {
			m.walkItem(key, value, handler, opts...)
		}
	case map[interface{}]interface{}:
		for key, value := range body {
			m.walkItem(fmt.Sprint(key), value, handler, opts...)
		}
	case []interface{}:
		for i, item := range body {
			m.walkItem(fmt.Sprintf("[%d]", i), item, handler, opts...)
		}
	}
}

func (m dataCrawler) walkItem(key string, value interface{}, handler DataCrawlerHandler, opts ...DataCrawlerOption) {
	settings := &DataCrawlerSettings{
		key:   key,
		value: value,
	}

	for _, opt := range opts {
		opt(settings)
	}

	switch v := value.(type) {
	case map[string]interface{}, map[interface{}]interface{}, []interface{}:
		o := make([]DataCrawlerOption, 0)
		if settings.isJoinKeys {
			o = append(o, WithJoinKeys())
			o = append(o, WithPrefix(settings.key))
		}

		if settings.isSkipCollections {
			o = append(o, WithSkipCollections())
		}

		m.Walk(v, handler, o...)
	default:
		handler(settings.key, v)
	}
}
