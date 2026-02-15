package core

import (
	"fmt"
	"os/exec"
)

// TODO using videoName for now, change to UUID later in case two video with the same name are encoded at once
func BuildFFmpegCommand(input, videoName string) *exec.Cmd {

	segmentPattern := fmt.Sprintf("stream_%%v/segment_%%03d.ts")
	playlistPath := fmt.Sprintf("playlist_%%v.m3u8")

	args := []string{
		"-i", input, // input

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
		// "-hls_segment_type", "fmp4", // .mp4 segments instead of .ts
		"-hls_playlist_type", "vod", // self explanatory
		"-hls_list_size", "0", // std value is 5, used for livestreams, we want all segments in the list so we input 0
		"-master_pl_name", "master.m3u8",
		// "-var_stream_map", "v:0,a:0 v:1,a:1 v:2,a:2",
		"-var_stream_map", "v:0 v:1 v:2",
		"-strftime_mkdir", "1",
		"-hls_segment_filename", segmentPattern, // name and dir for hls segments
		playlistPath, // name and dir for hls playlist
	}

	// TODO make it an option to output the ffmpeg command
	// fmt.Println("\n", strings.Join(args, " "))

	return exec.Command("ffmpeg", args...)
}

/*

ffmpeg
-i .\video-test\trailer.mp4

V ENCODING
-c:v libx264
-b:v 4000k // fullhd
-maxrate 4800k
-bufsize 9000k
-preset slow

COPATIBILITY
-pix_fmt yuv420p

KEY FRAMES FOR VOD
-force_key_frames "expr:gte(t, n_forced*2)"
-sc_threshold 0
-g 9999

A ENCONDING
-c:a aac
-b:a 128k

HLS
-f hls playlist.m3u8
-hls_time 2
-hls_flags independent_segments
-hls_segment_type fmp4
-hls_playlist_type vod
-hls_list_size 0
*/
