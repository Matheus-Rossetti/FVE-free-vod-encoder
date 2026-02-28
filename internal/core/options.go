package core

type Options struct {
	Codec           string
	SegmentType     string
	SegmentDuration int

	UseTerminal bool
	UseREST     bool
	UseRabbitMQ bool

	StoreLocal        bool
	StoreS3           bool
	S3Endpoint        string
	S3AccessKey       string
	S3SecretAccessKey string
	S3UseSSL          bool

	Upload bool

	ExposeMetrics        bool
	ConcurrentEncodings  int
	Threads              int
	OutputFFmpegCommand  bool
	OutputEncodedVideoTo string
	MaxStoredVideos      int
}
