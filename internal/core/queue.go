package core

type IngestQueue chan *Job
type EncodeQueue chan *Job
type DispatchQueue chan *Job

var (
	ingestQueue   IngestQueue
	encodeQueue   EncodeQueue
	dispatchQueue DispatchQueue
)

func StartQueues(options *Options) (IngestQueue, EncodeQueue, DispatchQueue) {

	ingestQueue = make(chan *Job, 999)
	encodeQueue = make(chan *Job, options.Encode.ConcurrentEncodings)
	dispatchQueue = make(chan *Job, options.Encode.ConcurrentEncodings)

	return ingestQueue, encodeQueue, dispatchQueue
}
