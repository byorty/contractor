package common

type SpecType string

const (
	SpecTypeOA2 SpecType = "oa2"
	SpecTypeOA3 SpecType = "oa3"
)

type Arguments struct {
	Mode           WorkerType        `short:"m" long:"mode" required:"true"`
	SpecLocation   string            `short:"s" long:"spec-location"`
	SpecType       SpecType          `short:"f" long:"spec-format" default:"oa2"`
	Url            string            `short:"u" long:"url"`
	Variables      map[string]string `short:"v" long:"variable"`
	Tags           []string          `short:"t" long:"tags"`
	CertFilename   string            `long:"cert-filename"`
	KeyFilename    string            `long:"key-filename"`
	ConfigFilename string            `short:"c" long:"config-filename"`
}
