package reader

import (
	"context"
	"github.com/byorty/contractor/tester"
	"github.com/getkin/kin-openapi/openapi3"
)

type Reader interface {
	Read(filename string) (tester.TestCase2List, error)
}

type Func func(ctx context.Context, filename string) (*openapi3.T, error)

type reader struct {
	ctx        context.Context
	readerFunc Func
}

func (r *reader) Read(filename string) (tester.TestCase2List, error) {
	spec, err := r.readerFunc(r.ctx, filename)
	if err != nil {
		return tester.TestCase2List{}, err
	}

}
