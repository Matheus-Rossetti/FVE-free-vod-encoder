package core

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

// THERE'S AN EXAMPLE OF THE GENERATED COMMAND AT THE END OF THIS FILE

func BuildFFmpegCommand(video *Video, options *Options) *exec.Cmd {
	// I tried my best not to make this ugly, okay? I'm Sorry

	// segmentPattern := fmt.Sprintf("stream_%%v/segment_%%03d.m4s")
	// playlistPath := fmt.Sprintf("playlist_%%v.m3u8")
	// fmp4InitFilename := fmt.Sprintf("stream_%%v/init.mp4")

	// --------- EACH BUILD FUNC RETURN A SLICE ---------
	input := []string{"-i", video.Source}
	filterComplex := buildFilterComplex(video)
	// maps := buildMaps(video)

	args := slices.Concat(
		input,
		filterComplex,
	)

	// TODO make it an option to output the ffmpeg command
	fmt.Println("\n", strings.Join(args, " "))

	return exec.Command("ffmpeg", args...)
}

func buildFilterComplex(video *Video) []string {

	flag := "-filter_complex"
	splitAmount := len(video.RenditionsToMake)

	// Example output:
	// [0:v]split=3[v1][v2]
	var splitSection strings.Builder
	splitSection.WriteString(fmt.Sprintf("[0:v]split=%v", splitAmount))
	for index, _ := range video.RenditionsToMake {
		splitSection.WriteString(fmt.Sprintf("[v%v]", (index + 1)))
	}

	// Example output:
	// ;[v1]scale=-2:720[v1out];[v2]scale=-2:480[v2out]
	for index, res := range video.RenditionsToMake {
		scale := fmt.Sprintf("-2:%v", res.name)
		if video.Orientation == "vertical" {
			scale = fmt.Sprintf("%v:-2", res.name)
		}

		splitSection.WriteString(fmt.Sprintf(
			";[v%v]scale=%v[v%vout]",
			(index + 1),
			scale,
			(index + 1),
		))
	}

	// Final string should be something like:
	// [0:v]split=3[v1][v2];[v1]scale=-2:720[v1out];[v2]scale=-2:480[v2out]

	filterComplex := []string{flag, splitSection.String()}
	return filterComplex
}

/* args := []string{
	"-i", video.Source, // input

	"-filter_complex",
	"[0:v]split=3[v1][v2][v3];[v1]scale=-2:1080[v1out];[v2]scale=-2:720[v2out];[v3]scale=-2:480[v3out]",

	"-map", "[v1out]",
	"-c:v:0", "libx264", // video codec
	"-b:v:0", "5000k", // avarage video bitrate
	"-maxrate:v:0", "5350k", // max bitrate the video can get to
	"-bufsize:v:0", "7500k", // tells the player how much is safe to buffer

	"-map", "[v2out]",
	"-c:v:1", "libx264", // video codec
	"-b:v:1", "2800k", // avarage video bitrate
	"-maxrate:v:1", "2996k", // max bitrate the video can get to
	"-bufsize:v:1", "4200k", // tells the player how much is safe to buffer

	"-map", "[v3out]",
	"-c:v:2", "libx264", // video codec
	"-b:v:2", "1400k", // avarage video bitrate
	"-maxrate:v:2", "1498k", // max bitrate the video can get to
	"-bufsize:v:2", "2100k", // tells the player how much is safe to buffer

	// "-map", "a:0", "-c:a:0", "AAC", "-b:a:0", "192k",
	// "-map", "a:0", "-c:a:1", "AAC", "-b:a:1", "129k",
	// "-map", "a:0", "-c:a:2", "AAC", "-b:a:2", "96k",

	"-preset", "slow", // encodes slower but with better quality and compression
	"-pix_fmt", "yuv420p", // garantees yuv420p, which is widely adopted
	"-force_key_frames", "expr:gte(t, n_forced*1)", // force a keyframe every seconds
	"-sc_threshold", "0", // stops H.264 adding iframes at scene changes, we're forcing an iframe each 2 seconds
	"-g", "9999", // max frames between iframes, std value is 250, this option is just for safety

	// "-c:a", "aac", // audio codec
	// "-b:a", "128k", // audio bitrate

	"-f", "hls", // video format, in our case, either HLS or DASH
	"-hls_time", "2", // duration of each .ts segment
	"-hls_flags", "independent_segments", // throws and error if a segment doesn't start with an iframe
	"-hls_segment_type", "fmp4", // .mp4 segments instead of .ts
	"-hls_fmp4_init_filename", fmp4InitFilename,
	"-hls_playlist_type", "vod", // self explanatory
	"-hls_list_size", "0", // std value is 5, used for livestreams, we want all segments in the list so we input 0
	"-master_pl_name", "master.m3u8",
	// "-var_stream_map", "v:0,a:0 v:1,a:1 v:2,a:2",
	"-var_stream_map", "v:0 v:1 v:2",
	"-strftime_mkdir", "1",
	"-hls_segment_filename", segmentPattern, // name and dir for hls segments
	playlistPath, // name and dir for hls playlist
} */
