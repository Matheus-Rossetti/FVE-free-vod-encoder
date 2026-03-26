package core

type IngestQueue chan *Job
type EncodeQueue chan *Job
type DispatchQueue chan *Job

func (f *frevod) StartQueues() (IngestQueue, EncodeQueue, DispatchQueue) {

	ingestQueue := make(chan *Job, 999)
	encodeQueue := make(chan *Job, f.Options.Encode.ConcurrentEncodings)
	dispatchQueue := make(chan *Job, f.Options.Encode.ConcurrentEncodings)

	return ingestQueue, encodeQueue, dispatchQueue
}

func (i IngestQueue) PushJob(uriType URIType, uri, keyStarter string, source string) {

	job := NewJob()
	job.DownloadJob.Source = source
	job.DownloadJob.UriType = uriType
	job.DownloadJob.VideoUri = uri

	job.UploadJob.KeyStarter = keyStarter

	i <- job
}
