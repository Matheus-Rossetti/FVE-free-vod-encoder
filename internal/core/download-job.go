package core

type DownloadJob struct {
	VideoUri string
	Source   string
}

func NewDownloadJob(videoUri, source string) DownloadJob {
	return DownloadJob{
		VideoUri: videoUri,
		Source:   source,
	}
}
