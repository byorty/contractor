package common

type SpecType string

const (
	SpecTypeOA2 SpecType = "oa2"
	SpecTypeOA3 SpecType = "oa3"
)

type Arguments struct {
	Mode         WorkerType `short:"m" long:"mode" required:"true"`
	SpecFilename string     `short:"s" long:"spec-filename" required:"true"`
	SpecType     SpecType   `short:"t" long:"spec-type" required:"false" default:"oa2"`
	BaseUrl      string     `short:"b" long:"base-url" required:"true"`
}
