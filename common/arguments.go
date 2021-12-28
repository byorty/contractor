package common

type Arguments struct {
	Mode         WorkerKind `short:"m" long:"mode" required:"true"`
	SpecFilename string     `short:"s" long:"spec" required:"true"`
	BaseUrl      string     `short:"b" long:"base-url" required:"true"`
}
