package core

import "os"

type Job struct {
	IngestJob   IngestJob
	EncodeJob   EncodeJob
	DispatchJob DispatchJob
}

func NewJob() *Job {
	return &Job{
		IngestJob:   IngestJob{},
		EncodeJob:   EncodeJob{},
		DispatchJob: DispatchJob{},
	}
}

func PushJob(uriType URIType, uri, keyStarter string, source string) {

	job := NewJob()
	job.IngestJob.Source = source
	job.IngestJob.UriType = uriType
	job.IngestJob.VideoUri = uri

	job.DispatchJob.KeyStarter = keyStarter

	ingestQueue <- job
}

type IngestJob struct {
	VideoUri     string
	UriType      URIType
	Source       string
	UploadMethod string
}

type EncodeJob struct {
	AbsoluteVideoPath string
	File              *os.File
	DownloadedFile    bool
}

type DispatchJob struct {
	FromDir    string
	KeyStarter string
}
