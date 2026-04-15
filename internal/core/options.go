package core

type Options struct {
	Ingest   IngestOptions   `yaml:"ingest"`
	Encode   EncodeOptions   `yaml:"encode"`
	Dispatch DispatchOptions `yaml:"dispatch"`

	// ExposeMetrict bool `yaml:"expose_metrics"`
}

func NewOptions() *Options {
	return &Options{}
}

type IngestOptions struct {
	Cli  CliOptions  `yaml:"cli"`
	REST RestOptions `yaml:"rest"`
	// RabbitMQ bool
	// Kafka    bool
	// GRPC     bool
	// QSQ      bool
}

type CliOptions struct {
	Enabled bool `yaml:"enabled"`
}

type RestOptions struct {
	Enabled bool   `yaml:"enabled"`
	Port    string `yaml:"port" validate:"required_if=Enabled true"`
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
	Enabled bool   `yaml:"enabled"`
	StoreAt string `yaml:"store_at" validate:"required_if=Enabled true"`
}

type S3Options struct {
	Enabled         bool   `yaml:"enabled"`
	Endpoint        string `yaml:"endpoint"   validate:"required_if=Enabled true"`
	AccessKey       string `yaml:"access_key" validate:"required_if=Enabled true"`
	SecretAccessKey string `yaml:"secret_key" validate:"required_if=Enabled true"`
	Bucket          string `yaml:"bucket"     validate:"required_if=Enabled true"`
	UseSSL          bool   `yaml:"use_ssl"        validate:"required_if=Enabled true"`
}

type AzureBlobStorageOptions struct {
}

type GoogleCloudStorageOptions struct {
}
