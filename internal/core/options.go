package core

type Options struct {
	Input  InputOptions  `yaml:"input"`
	Encode EncodeOptions `yaml:"encode"`
	Upload UploadOptions `yaml:"upload"`

	// ExposeMetrict bool `yaml:"expose_metrics"`
}

type InputOptions struct {
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
	ConcurrentEncodings int  `yaml:"concurrent_encodings" validate:"min=1,max=100"`
	OutputFFmpegCommand bool `yaml:"output_ffmpeg_command"`
}

type UploadOptions struct {
	S3 AmazonS3Options `yaml:"s3"`
	// UseAzureBlob    AzureBlobStorage
	// UseCloudStorage GoogleCloudStorage
	Local LocalOption `yaml:"local"`
}

type LocalOption struct {
	Use     bool   `yaml:"use"`
	StoreAt string `yaml:"store_at" validate:"required_if=Use true"`
}

type AmazonS3Options struct {
	Use             bool   `yaml:"use"`
	Endpoint        string `yaml:"endpoint" validate:"required_if=Use true"`
	AccessKey       string `yaml:"access_key" validate:"required_if=Use true"`
	SecretAccessKey string `yaml:"secret_key" validate:"required_if=Use true"`
	BucketName      string `yaml:"bucket_name" validate:"required_if=Use true"`
	UseSSL          bool   `yaml:"ssl" validate:"required_if=Use true"`
}

type AzureBlobStorageOptions struct {
}

type GoogleCloudStorageOptions struct {
}
