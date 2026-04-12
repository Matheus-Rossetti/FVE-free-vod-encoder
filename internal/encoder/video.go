package encoder

type Resolution struct {
	name  string
	value int
}

type Video struct {
	Height      int
	Width       int
	AspectRatio string
	Duration    string // Seconds

	HasAudio            bool // Useful when building the ffmpeg command, we dont include audio tags if hasAudio == false
	Orientation         string
	ReferenceResolution int          // in a vertical video, this is width value | for horizontal, height value
	Renditions          []Resolution // if a video is 1440p/QHD, this would be -> 1440p, 1080p, 720p, 480p
}

func NewVideo(metadata *Metadata) *Video {

	hasAudio := false
	videoHeight := 0
	mainStream := metadata.Streams[0]
	for index, stream := range metadata.Streams {

		// mainStream should be the one with the highest res
		if stream.Height > videoHeight {
			mainStream = metadata.Streams[index]
		}

		if stream.CodecType == "audio" {
			hasAudio = true
		}
	}

	// --------- CHECK VIDEO ORIENTATION ---------
	var orientation string
	var referenceResolution int
	if mainStream.Width > mainStream.Height {
		orientation = "horizontal"
		referenceResolution = mainStream.Height
	} else {
		orientation = "vertical" // vertical also covers square videos.
		referenceResolution = mainStream.Width
	}

	// --------- DECIDE WHICH RENDITIONS TO MAKE ---------
	resolutions := []Resolution{
		{"2160", 2160},
		{"1440", 1440},
		{"1080", 1080},
		{"720", 720},
		{"480", 480},
	}

	// only adds renditions of native res and lower, never upscale.
	var renditions []Resolution
	for _, res := range resolutions {
		if referenceResolution >= res.value {
			renditions = append(renditions, res)
		}
	}

	// in case a video res is lower than the lowest supported res
	if len(renditions) == 0 {
		renditions = append(renditions, resolutions[len(resolutions)-1]) // add the lowest possible res to the list
	}

	return &Video{
		Height:      mainStream.Height,
		Width:       mainStream.Width,
		AspectRatio: mainStream.AspectRatio,
		Duration:    mainStream.Duration,

		HasAudio:            hasAudio,
		Orientation:         orientation,
		ReferenceResolution: referenceResolution,
		Renditions:          renditions,
	}
}
