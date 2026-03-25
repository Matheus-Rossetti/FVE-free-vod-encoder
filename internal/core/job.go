package core

import "os"

type Job struct {
	DownloadJob DownloadJob
	EncodeJob   EncodeJob
	UploadJob   UploadJob
}

func NewJob() *Job {
	return &Job{
		DownloadJob: DownloadJob{},
		EncodeJob:   EncodeJob{},
		UploadJob:   UploadJob{},
	}
}

type DownloadJob struct {
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

type UploadJob struct {
	FromDir    string
	KeyStarter string
}
