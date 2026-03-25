package cli

import "github.com/Matheus-Rossetti/frevod/internal/core"

func (c *cli) pushJob(uriType core.URIType, uri, keyStarter string) {

	job := core.NewJob()
	job.DownloadJob.Source = "cli"
	job.DownloadJob.UriType = uriType
	job.DownloadJob.VideoUri = uri

	job.UploadJob.KeyStarter = keyStarter

	c.downloadQueue <- job
}
