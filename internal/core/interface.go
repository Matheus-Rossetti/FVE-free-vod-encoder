package core

type Input interface {
	IngestVideo()
}

type Output interface {
	StoreOutput()
}
