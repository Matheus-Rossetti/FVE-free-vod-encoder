package core

import (
	"fmt"
	"os/exec"
	"runtime"
	"slices"
	"strings"
	"syscall"
)

// This is quite a complex file, it operates multiple string concatenations
// Each part of the command is built in a different function
// We iterate through the same array in each function, that's on purpose
// Sacrifice a litte bit of performance for maintainability
// I tried my best not to make this ugly, okay? I'm Sorry

// THERE'S AN EXAMPLE OF THE OUTPUT OF THIS FUNCTION AT THE END OF THIS FILE
func BuildFFmpegCommand(video *Video, options *Options) *exec.Cmd {

	// --------- EACH BUILD FUNC RETURN A SLICE ---------
	input := []string{"-i", video.Source}
	filterComplex := buildFilterComplex(video)
	videoMaps := buildVideoMaps(video) // h.264
	audioMaps := buildAudioMaps(video) // aac | maybe switch to opus, heard sounds better at the same bitrate and handle 5.1 sound
	keyFramesAndQuality := getKeyFramesAndQuality()
	hlsOptions := buildHlsOptions(video)

	ffmpegArgs := slices.Concat(
		input,
		filterComplex,
		videoMaps,
		audioMaps, // this can be nil, but slices.Concat handles it
		keyFramesAndQuality,
		hlsOptions,
	)

	if options.OutputFFmpegCommand {
		// TODO some values need to be inclosed in double quotes "example"
		fmt.Println("\n", strings.Join(ffmpegArgs, " "))
	}

	// Returns a command that starts FFmpeg as a low-priority process, allowing it to use 100% of > SPARE < compute
	cmd := BuildForLowPrioExecution(ffmpegArgs)

	return cmd
}

func buildFilterComplex(video *Video) []string {

	flag := "-filter_complex"
	splitAmount := len(video.RenditionsToMake)

	// Example output:
	// [0:v]split=3[v1][v2]
	var splitSection strings.Builder
	fmt.Fprintf(&splitSection, "[0:v]split=%v", splitAmount)
	for index := range video.RenditionsToMake {
		fmt.Fprintf(&splitSection, "[v%v]", (index + 1))
	}

	// Example output:
	// ;[v1]scale=-2:720[v1out];[v2]scale=-2:480[v2out]
	for index, res := range video.RenditionsToMake {
		scale := fmt.Sprintf("-2:%v", res.name)
		if video.Orientation == "vertical" {
			scale = fmt.Sprintf("%v:-2", res.name)
		}

		fmt.Fprintf(&splitSection, ";[v%v]scale=%v[v%vout]",
			(index + 1),
			scale,
			(index + 1))
	}

	// Final string should be something like:
	// [0:v]split=3[v1][v2];[v1]scale=-2:720[v1out];[v2]scale=-2:480[v2out]

	filterComplex := []string{flag, splitSection.String()}
	return filterComplex
}

func buildVideoMaps(video *Video) []string {

	// TODO adapt to .env or config.yml in case
	// user wants custom bitrates

	// --- MAP OF BITRATE VALUES ---
	rates := map[int]struct {
		avg, max, buf string
	}{ // This is for 30 fps, if the video has 60fps we gotta double this values
		2160: {"25000k", "35000k", "50000k"},
		1440: {"12000k", "16000k", "24000k"},
		1080: {"6000k", "8000k", "12000k"},
		720:  {"3500k", "5000k", "7000k"},
		480:  {"1500k", "2000k", "3000k"},
	}

	mapFlag := "-map"
	var maps []string
	for index := range video.RenditionsToMake {

		// --- INDEX THE FLAGS AND RETREIVE VALUES FOR REFERENCE RESOLUTION ---
		version := fmt.Sprintf("[v%vout]", (index + 1))
		codecFlag := fmt.Sprintf("-c:v:%v", index)
		bitrateFlag := fmt.Sprintf("-b:v:%v", index)
		maxrateFlag := fmt.Sprintf("-maxrate:v:%v", index)
		bufsizeFlag := fmt.Sprintf("-bufsize:v:%v", index)
		bitrate := rates[video.RenditionsToMake[index].value]

		// --- TURN IT ALL INTO A SLICE ---
		maps = append(maps,
			mapFlag, version,
			codecFlag, "libx264",
			bitrateFlag, bitrate.avg,
			maxrateFlag, bitrate.max,
			bufsizeFlag, bitrate.buf,
		)
	}

	return maps
}

func buildAudioMaps(video *Video) []string {

	if video.HasAudio == false {
		return nil
	}

	// --- MAP OF BITRATE VALUES ---
	rates := map[int]string{
		2160: "320k",
		1440: "320k",
		1080: "192k",
		720:  "192k",
		480:  "128k",
	}

	mapFlag := "-map"
	var audioMaps []string
	for index := range video.RenditionsToMake {

		codecFlag := fmt.Sprintf("-c:a:%v", index)
		bitrateFlag := fmt.Sprintf("-b:a:%v", index)
		bitrate := rates[video.ReferenceResolution]

		audioMaps = append(audioMaps,
			mapFlag, "a:0",
			codecFlag, "aac",
			bitrateFlag, bitrate,
		)
	}
	// "-map", "a:0", "-c:a:0", "AAC", "-b:a:0", "192k",
	// "-map", "a:0", "-c:a:1", "AAC", "-b:a:1", "129k",
	// "-map", "a:0", "-c:a:2", "AAC", "-b:a:2", "96k",

	return audioMaps
}

func getKeyFramesAndQuality() []string {
	return []string{
		"-preset", "veryslow", // encodes slower but with better quality and compression
		"-pix_fmt", "yuv420p", // sets yuv420p, which is widely adopted
		"-force_key_frames", "expr:gte(t, n_forced*1)", // force a keyframe every seconds
		"-sc_threshold", "0", // stops H.264 from automatically adding iframes at scene changes
		"-g", "9999", // max frames between iframes, std value is 250, this option is just for safety
	}
}

func buildHlsOptions(video *Video) []string {

	segmentPattern := "stream_%v/segment_%03d.m4s"
	playlistPath := "playlist_%v.m3u8"
	fmp4InitFilename := "stream_%v/init.mp4"

	versionAmount := len(video.RenditionsToMake)

	var streamMap strings.Builder
	for index := range versionAmount {
		if video.HasAudio {
			fmt.Fprintf(&streamMap, "v:%v,a:%v ", index, index)
		} else {
			fmt.Fprintf(&streamMap, "v:%v ", index)
		}
	}

	return []string{
		"-f", "hls", // video format, in our case, either HLS or DASH
		"-hls_time", "2", // duration of each .ts segment
		"-hls_flags", "independent_segments", // throws and error if a segment doesn't start with an iframe
		"-hls_segment_type", "fmp4", // .mp4 segments instead of .ts
		"-hls_fmp4_init_filename", fmp4InitFilename,
		"-hls_playlist_type", "vod", // self explanatory
		"-hls_list_size", "0", // std value is 5, we want all segments in the list so we input 0 | weird, I know
		"-master_pl_name", "master.m3u8",
		"-var_stream_map", streamMap.String(),
		"-strftime_mkdir", "1",
		"-hls_segment_filename", segmentPattern, // name and dir for hls segments
		playlistPath, // name and dir for hls playlist
	}
}

func BuildForLowPrioExecution(ffmpegArgs []string) *exec.Cmd {
	// ------ SET LOW PRIO PROCESS FOR UNIX BASED ------
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {

		// Use nice to set low prio
		args := append([]string{"-n", "10", "ffmpeg"}, ffmpegArgs...)
		return exec.Command("nice", args...)
	}

	cmd := exec.Command("ffmpeg", ffmpegArgs...)

	// ------ SET LOW PRIO PROCESS FOR WINDOWS ------
	if runtime.GOOS == "windows" {
		const MICROSOFT_MAGIC_CONSTANT_THAT_STARTS_LOW_PRIORITY_PROCESSES = 0x00004000

		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: MICROSOFT_MAGIC_CONSTANT_THAT_STARTS_LOW_PRIORITY_PROCESSES,
		}
	}

	return cmd
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
