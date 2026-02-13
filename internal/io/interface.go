package io

type Input interface {
	IngestVideo()
}

type Output interface {
	StoreOutput()
}
