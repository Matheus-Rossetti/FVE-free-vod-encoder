package core

type VideoJob struct {
	// cli | rabbitmq | REST API | etc...
	Source            string
	AbsoluteVideoPath string
}
