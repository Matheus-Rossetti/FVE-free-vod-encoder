package core

type Options struct {
	Ingest   IngestOptions   `yaml:"input"`
	Encode   EncodeOptions   `yaml:"encode"`
	Dispatch DispatchOptions `yaml:"upload"`

	// ExposeMetrict bool `yaml:"expose_metrics"`
}

func NewOptions() *Options {
	return &Options{}
}

type IngestOptions struct {
	Cli  bool `yaml:"cli"`
	REST bool `yaml:"rest"`
	// RabbitMQ bool
	// Kafka    bool
	// GRPC     bool
	// QSQ      bool
}

type EncodeOptions struct {
	// Codec               string `yaml:"codec" validate:"oneof=h.264 h.265"`
	// SegmentType         string `yaml:"segment_type" validate:"oneof=fmp4 ts"`
	// SegmentDuration     int    `yaml:"segment_duration" validate:"min=0,max=10"`
	ConcurrentEncodings int  `yaml:"concurrent_encodings" validate:"min=1,max=100,required"`
	OutputFFmpegCommand bool `yaml:"output_ffmpeg_command"`
}

type DispatchOptions struct {
	S3 S3Options `yaml:"s3"`
	// UseAzureBlob    AzureBlobStorage
	// UseCloudStorage GoogleCloudStorage
	Local LocalOption `yaml:"local"`
}

type LocalOption struct {
	Enabled bool   `yaml:"use"`
	StoreAt string `yaml:"store_at" validate:"required_if=Use true"`
}

type S3Options struct {
	Enabled         bool   `yaml:"use"`
	Endpoint        string `yaml:"endpoint" validate:"required_if=Use true"`
	AccessKey       string `yaml:"access_key" validate:"required_if=Use true"`
	SecretAccessKey string `yaml:"secret_key" validate:"required_if=Use true"`
	Bucket          string `yaml:"bucket" validate:"required_if=Use true"`
	UseSSL          bool   `yaml:"ssl" validate:"required_if=Use true"`
}

type AzureBlobStorageOptions struct {
}

type GoogleCloudStorageOptions struct {
}
