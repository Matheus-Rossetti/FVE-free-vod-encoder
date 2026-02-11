package core

import (
	"os/exec"
)

func BuildFFmpegCommand() *exec.Cmd {

	args := []string{
		"-i", "./video-test/video1.mov", // input
		"-c:v", "libx264", // video codec
		"-b:v", "4000k", // avarage video bitrate
		"-maxrate", "4800k", // max bitrate the video can get to
		"-bufsize", "9000k", // tells the player how much is safe to buffer
		"-preset", "slow", // encodes slower but with better quality and compression
		"-pix_fmt", "yuv420p", // garantees yuv420p, which is widely adopted
		"-force_key_frames", "expr:gte(t, n_forced*2)", // forcer a keyframe each 2 seconds
		"-sc_threshold", "0", // stops H.264 adding iframes at scene changes, we're forcing an iframe each 2 seconds
		"-g", "9999", // max frames between iframes, std value is 250, this option is just for safety
		"-c:a", "aac", // audio codec
		"-b:a", "128k", // audio bitrate
		"-f", "hls", // video format, in our case, either HLS or DASH
		"-hls_time", "2", // duration of each .ts segment
		"-hls_flags", "independent_segments", // throws and error if a segment doesn't start with an iframe
		"-hls_playlist_type", "vod", // self explanatory
		"-hls_list_size", "0", // std value is 5, used for livestreams, we want all segments in the list so we input 0
		"playlist.m3u8", // name of the hls playlist
	}

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
-hls_playlist_type vod
-hls_list_size 0

*/
