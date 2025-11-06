package flags

import "flag"

type Flags struct {
	InputPath string
	OutputName string
	SegmentDuration int
}

func New() *Flags{
	return &Flags{}
}

func (f *Flags) Bind() {
	flag.StringVar(&f.InputPath, "path", "", "Video Path")
	flag.StringVar(&f.OutputName, "name", "processed-video", "Video Name")
	flag.IntVar(&f.SegmentDuration, "segment_duration", 3, "Duration of each .ts segment")

	flag.Parse()
}
