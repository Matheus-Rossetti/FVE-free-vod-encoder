package core

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
