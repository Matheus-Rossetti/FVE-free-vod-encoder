package rest

import "github.com/Matheus-Rossetti/frevod/internal/core"

func (rest *rest) pushJob(uriType core.URIType, uri, keyStarter string) {
	job := core.NewJob()

	job.DownloadJob.Source = "REST"
	job.DownloadJob.UriType = uriType
	job.DownloadJob.VideoUri = uri

	job.UploadJob.KeyStarter = keyStarter

	rest.downloadQueue <- job
}
