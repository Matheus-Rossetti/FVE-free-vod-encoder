package core

// func TestPushJob(t *testing.T) {
// 	testCases := []struct {
// 		name        string
// 		source      string
// 		uriType     URIType
// 		uri         string
// 		keyStarter  string
// 		expectedJob Job
// 	}{
// 		{
// 			name:       "correct job from cli",
// 			uriType:    Url,
// 			uri:        "https://mycoolvideo.com",
// 			keyStarter: "key/starter",
// 			expectedJob: Job{
// 				DownloadJob: DownloadJob{
// 					Source:   "cli",
// 					UriType:  Url,
// 					VideoUri: "https://mycoolvideo.com",
// 				},
// 				UploadJob: UploadJob{
// 					KeyStarter: "key/starter",
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range testCases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			queue := make(chan *Job, 2)
// 			cli := cli{
// 				ingestQueue: queue,
// 			}

// 			cli.ingestQueue.PushJob(tt.uriType, tt.uri, tt.keyStarter, "cli")
// 			job := <-queue

// 			if tt.expectedJob.DownloadJob.Source != job.DownloadJob.Source {
// 				t.Errorf("expected source to be %v, got %v",
// 					tt.expectedJob.DownloadJob.Source,
// 					job.DownloadJob.Source)
// 			}

// 			if tt.expectedJob.DownloadJob.UriType != job.DownloadJob.UriType {
// 				t.Errorf("expected uri type to be %v, got %v",
// 					tt.expectedJob.DownloadJob.UriType,
// 					job.DownloadJob.UriType)
// 			}

// 			if tt.expectedJob.DownloadJob.VideoUri != job.DownloadJob.VideoUri {
// 				t.Errorf("expected uri to be %v, got %v",
// 					tt.expectedJob.DownloadJob.VideoUri,
// 					job.DownloadJob.VideoUri)
// 			}

// 			if tt.expectedJob.UploadJob.KeyStarter != job.UploadJob.KeyStarter {
// 				t.Errorf("expected key starter to be %v, got %v",
// 					tt.expectedJob.UploadJob.KeyStarter,
// 					job.UploadJob.KeyStarter)
// 			}
// 		})
// 	}
// }
