package core

type Options struct {
	ExposeMetrics bool

	// input
	Input InputOptions `yaml:"input"`

	// encode
	Encode EncodeOptions `yaml:"encode"`

	// upload
	Upload UploadOptions `yaml:"upload"`
}

type InputOptions struct {
	UseTerminal bool `yaml:"use_terminal"`
	UseREST     bool `yaml:"use_rest"`
	// UseRabbitMQ bool
	// UseKafka    bool
	// UseGRPC     bool
}

type EncodeOptions struct {
	Codec               string `yaml:"codec"`
	SegmentType         string `yaml:"segment_type"`
	SegmentDuration     int    `yaml:"segment_duration"`
	ConcurrentEncodings int    `yaml:"concurrent_encodings"`
	OutputFFmpegCommand bool   `yaml:"output_ffmpeg_command"`
}

type UploadOptions struct {
	S3 AmazonS3Options `yaml:"s3"`
	// UseAzureBlob    AzureBlobStorage
	// UseCloudStorage GoogleCloudStorage
	StoreLocal string `yaml:"store_local"`

	ConcurrentUploads int `yaml:"concurrent_uploads"`
}

type AmazonS3Options struct {
	Use               bool   `yaml:"use"`
	S3Endpoint        string `yaml:"s3_endpoint"`
	S3AccessKey       string `yaml:"s3_access_key"`
	S3SecretAccessKey string `yaml:"s3_secret_key"`
	S3UseSSL          bool   `yaml:"s3_use_ssl"`
}

type AzureBlobStorageOptions struct {
}

type GoogleCloudStorageOptions struct {
}
