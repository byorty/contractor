package common

import (
	"go.uber.org/fx"
)

var Constructors = fx.Provide(
	NewFxMediaConverter,
	NewFxWorkerContainer,
	NewFxDataCrawler,
	NewFxExpressionFactory,
	NewFxLoggerFactory,
	NewFxConfigProviderFactory,
)
