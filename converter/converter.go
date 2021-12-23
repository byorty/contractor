package converter

import (
	"context"
	"github.com/byorty/contractor/common"
)

type Converter interface {
	Convert(ctx context.Context, arguments common.Arguments) (common.TemplateContainer, error)
}
