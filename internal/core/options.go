package core

type Options struct {
	ExposeMetrics bool

	// input
	Input InputOptions `mapstructure:"input"`

	// encode
	Encode EncodeOptions `mapstructure:"encode"`

	// upload
	Upload UploadOptions `mapstructure:"upload"`
}

type InputOptions struct {
	UseTerminal bool `mapstructure:"use_terminal"`
	UseREST     bool `mapstructure:"use_rest"`
	// UseRabbitMQ bool
	// UseKafka    bool
	// UseGRPC     bool
}

type EncodeOptions struct {
	Codec               string `mapstructure:"codec"`
	SegmentType         string `mapstructure:"segment_type"`
	SegmentDuration     int    `mapstructure:"segment_duration"`
	ConcurrentEncodings int    `mapstructure:"concurrent_encodings"`
	OutputFFmpegCommand bool   `mapstructure:"output_ffmpeg_command"`
}

type UploadOptions struct {
	S3 AmazonS3Options `mapstructure:"s3"`
	// UseAzureBlob    AzureBlobStorage
	// UseCloudStorage GoogleCloudStorage
	Local LocalOption `mapstructure:"local"`

	ConcurrentUploads int `mapstructure:"concurrent_uploads"`
}

type LocalOption struct {
	Use     bool   `mapstructure:"use"`
	StoreAt string `mapstructure:"store_at"`
}

type AmazonS3Options struct {
	Use               bool   `mapstructure:"use"`
	S3Endpoint        string `mapstructure:"s3_endpoint"`
	S3AccessKey       string `mapstructure:"s3_access_key"`
	S3SecretAccessKey string `mapstructure:"s3_secret_key"`
	S3BucketName      string `mapstructure:"s3_bucket_name"`
	S3UseSSL          bool   `mapstructure:"s3_use_ssl"`
}

type AzureBlobStorageOptions struct {
}

type GoogleCloudStorageOptions struct {
}
