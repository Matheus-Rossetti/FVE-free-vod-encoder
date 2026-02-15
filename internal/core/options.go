package core

type Options struct {
}

/*
OPTIONS TO OFFER:
	Encoding and containerization:
	- H.264 or H.265
	- .ts or .m4s output

	Inputs:
	- Terminal with path or download link
	- Accept HTTP requests with download link
	- Listen to a RabbitMQ queue

	Outputs:
	- Local | specify path
	- An storage service compatible with S3 API
	- Maybe any storage, as long as the user informs the upload URL
*/
