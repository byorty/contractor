package common

type SpecType string

const (
	SpecTypeOA2 SpecType = "oa2"
	SpecTypeOA3 SpecType = "oa3"
)

type EngineType string

const (
	EngineTypeOA  EngineType = "oa"
	EngineTypeE2E EngineType = "e2e"
)

type Arguments struct {
	Mode         WorkerType        `short:"m" long:"mode" required:"true"`
	SpecLocation string            `short:"s" long:"spec-location" required:"true"`
	SpecType     SpecType          `short:"f" long:"spec-format" required:"false" default:"oa2"`
	Url          string            `short:"u" long:"url" required:"true"`
	Variables    map[string]string `short:"v" long:"variable"`
	Tags         []string          `short:"t" long:"tags"`
	CertFilename string            `short:"c" long:"cert-filename"`
	KeyFilename  string            `short:"k" long:"key-filename"`
	Engine       EngineType        `short:"e" long:"engine"`
}
