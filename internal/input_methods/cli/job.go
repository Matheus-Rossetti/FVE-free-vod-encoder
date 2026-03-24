package cli

import "github.com/Matheus-Rossetti/frevod/internal/core"

func (c *cli) pushJob(uri string, uriType core.URIType, keyStarter string) {
	if c.err != nil {
		return
	}

	job := core.NewJob()
	job.DownloadJob.Source = "cli"
	job.DownloadJob.UriType = uriType
	job.DownloadJob.VideoUri = uri

	job.UploadJob.KeyStarter = keyStarter

	c.downloadQueue <- job
}
