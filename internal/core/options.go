package core

type Options struct {
	Encoding    string
	SegmentType string

	UseTerminal bool
	UseREST     bool
	UseRabbitMQ bool

	StoreLocal bool
	StoreS3    bool
	Upload     bool

	ExposeMetrics       bool
	ConcurrentEncodings int
}

// TODO set standard values and substitute them for the ones in config.yml
func ParseOptions() *Options {
	return &Options{
		UseTerminal: true,
	}
}
