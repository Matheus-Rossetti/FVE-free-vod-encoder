package core

type Options struct {
	Codec           string
	SegmentType     string
	SegmentDuration int

	UseTerminal bool
	UseREST     bool
	UseRabbitMQ bool

	StoreLocal bool
	StoreS3    bool
	Upload     bool

	ExposeMetrics       bool
	ConcurrentEncodings int
	Threads             int
	OutputFFmpegCommand bool
}

// TODO set standard values and substitute them for the ones in .env
func ParseOptions() *Options {

	// threadAmount := (runtime.NumCPU() - 4)

	return &Options{
		Codec:               "h.264",
		SegmentType:         "fmp4",
		SegmentDuration:     2,
		UseTerminal:         true,
		UseREST:             false,
		UseRabbitMQ:         false,
		StoreLocal:          true,
		StoreS3:             false,
		Upload:              false,
		ExposeMetrics:       false,
		ConcurrentEncodings: 2,
		Threads:             1, // This is threads per rendition
		OutputFFmpegCommand: false,
	}
}
