package core

import (
	"fmt"
	"os/exec"
	"strings"
)

// TODO use fmt.Sprintf("") to concat strings

func ConvertVideo(inputPath string) {

	command := "ffmpeg"
	flags := fmt.Sprintf("ffmpeg -i %v -map 0:v:0 -c:v:0 libx264 -b:v:0 4500k -maxrate:v:0 4500k -bufsize:v:0 9000k -s:v:0 1920x1080 -g 120 -keyint_min 120 -sc_threshold 0 -map 0:v:0 -c:v:1 libx264 -b:v:1 2500k -maxrate:v:1 2500k -bufsize:v:1 5000k -s:v:1 1280x720 -g 120 -keyint_min 120 -sc_threshold 0 -map 0:a:0 -c:a aac -b:a 128k -ac 2 -f hls -var_stream_map \"v:0,a:0 v:1,a:0\" -master_pl_name master.m3u8 -hls_time 4 -hls_playlist_type vod -hls_segment_filename \"stream_%%v/data%%03d.ts\" \"stream_%%v.m3u8\"", inputPath)
	args := strings.Fields(flags)

	cmd := exec.Command(command, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}

	fmt.Printf("Command output: %s\n", string(output))
}
