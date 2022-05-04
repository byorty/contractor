package reader

import (
	"github.com/byorty/contractor/common"
	"go.uber.org/fx"
)

type Constructor func() Reader

type FactoryFxIn struct {
	fx.In
	Descriptors []common.Descriptor[Constructor] `group:"open_api_reader"`
}

func NewFxFactory(in FactoryFxIn) common.Factory[Constructor] {
	return common.NewFactory[Constructor](in.Descriptors)
}
