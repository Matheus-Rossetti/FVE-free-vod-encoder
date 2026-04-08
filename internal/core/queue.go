package core

type IngestQueue chan *Job
type EncodeQueue chan *Job
type DispatchQueue chan *Job

var (
	ingestQueue   IngestQueue
	encodeQueue   EncodeQueue
	dispatchQueue DispatchQueue
)

func StartQueues(options *Options) (IngestQueue, EncodeQueue, DispatchQueue) {

	ingestQueue = make(chan *Job, 999)
	encodeQueue = make(chan *Job, options.Encode.ConcurrentEncodings)
	dispatchQueue = make(chan *Job, options.Encode.ConcurrentEncodings)

	return ingestQueue, encodeQueue, dispatchQueue
}

func PushJob(uriType URIType, uri, keyStarter string, source string) {

	job := NewJob()
	job.DownloadJob.Source = source
	job.DownloadJob.UriType = uriType
	job.DownloadJob.VideoUri = uri

	job.UploadJob.KeyStarter = keyStarter

	ingestQueue <- job
}
