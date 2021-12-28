package common

import (
	"encoding/json"
	"encoding/xml"
	"github.com/pkg/errors"
)

var (
	ErrUnsupportedMediaType = errors.New("unsupported media type")
)

type MediaType string

const (
	MediaTypeJson MediaType = "application/json"
	MediaTypeXml  MediaType = "application/xml"
)

type MediaConverter interface {
	Marshal(mediaType MediaType, i interface{}) ([]byte, error)
	Unmarshal(mediaType MediaType, buf []byte) (interface{}, error)
}

func NewFxMediaConverter() MediaConverter {
	return &mediaConverter{
		marshalers: map[MediaType]func(obj interface{}) ([]byte, error){
			MediaTypeJson: func(obj interface{}) ([]byte, error) {
				return json.Marshal(obj)
			},
			MediaTypeXml: func(obj interface{}) ([]byte, error) {
				return xml.Marshal(obj)
			},
		},
		unmarshalers: map[MediaType]func(buf []byte) (interface{}, error){
			MediaTypeJson: func(buf []byte) (interface{}, error) {
				var obj interface{}
				err := json.Unmarshal(buf, &obj)
				if err != nil {
					return nil, err
				}

				return obj, nil
			},
			MediaTypeXml: func(buf []byte) (interface{}, error) {
				var obj interface{}
				err := xml.Unmarshal(buf, &obj)
				if err != nil {
					return nil, err
				}

				return obj, nil
			},
		},
	}
}

type mediaConverter struct {
	marshalers   map[MediaType]func(obj interface{}) ([]byte, error)
	unmarshalers map[MediaType]func(buf []byte) (interface{}, error)
}

func (c *mediaConverter) Marshal(mediaType MediaType, i interface{}) ([]byte, error) {
	marshaler, ok := c.marshalers[mediaType]
	if !ok {
		return nil, errors.Wrap(ErrUnsupportedMediaType, string(mediaType))
	}

	return marshaler(i)
}

func (c *mediaConverter) Unmarshal(mediaType MediaType, buf []byte) (interface{}, error) {
	unmarshaler, ok := c.unmarshalers[mediaType]
	if !ok {
		return nil, errors.Wrap(ErrUnsupportedMediaType, string(mediaType))
	}

	return unmarshaler(buf)
}
