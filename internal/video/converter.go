package video

import (
	"fmt"
	"os/exec"
	"slices"

	"github.com/Matheus-Rossetti/video-converter-service/internal/app"
	"mvdan.cc/sh/v3/shell"
)

var resolutions = map[int]int{
	3840: 2160,
	2560: 1440,
	1920: 1080,
	1280: 720,
	854:  480,
}

var array = []int{3820, 2560, 1920, 1280, 854}

// TODO use fmt.Sprintf("") to concat strings
func Convert(video *Video, app *app.App) {
	command := "ffmpeg"
	keyFrameInterval := (video.Fps * app.Flags.SegmentDuration)

	index := slices.Index(array, video.Width) + 1
	versionsAmount := len(array) - index + 1

	flags := "-i \"" + video.Path + "\" "

	// SPLITS -------
	splits := "-filter_complex \"[0:v]split=" + fmt.Sprint(versionsAmount)
	for i := range versionsAmount {
		version := "[v" + fmt.Sprint(i+1) + "]"
		splits = splits + version
	}

	splits = splits + ";"
	flags = flags + splits
	// --------------

	// VERSIONS-----
	version := 1
	for _, value := range array {
		if video.Width >= value {
			res := fmt.Sprint(resolutions[value])
			flags = flags + "[v" + fmt.Sprint(version) + "]scale=w=-2:h=" + res + "[v" + res + "_scaled];[v" + res + "_scaled]pad=w=" + fmt.Sprint(value) + ":h=" + res + ":x=-1:y=-1:color=black[v" + res + "];"
			version++
		}
	}
	flags = flags + "\""
	//--------------

	// MAPS---------
	
	version = 0
	for _, value := range array {
		if video.Width >= value {
			flags = fmt.Sprintf("%v -map \"[v%v]\" -c:v:%v libx264 -preset medium -crf 23 -sc_threshold 0", flags, resolutions[value], version)
			version++
		}
	}
	flags = fmt.Sprintf("%v -map a:0 -c:a:0 aac -ac 2 -ar 48000", flags)
	// -------------

	// HLS ---------
	flags = fmt.Sprintf("%v -f hls -hls_time 3 -g %v -hls_playlist_type vod -hls_list_size 0 -hls_segment_filename \"%v/%%v/segment_%%03d.ts\" -hls_flags independent_segments -master_pl_name master.m3u8", flags, keyFrameInterval, video.Name)
	// -------------

	// STREAM MAP --
	flags = fmt.Sprintf("%v -var_stream_map \"", flags)
	version = 0 
	for _, value := range array {
		if video.Width >= value {
			flags = fmt.Sprintf("%vv:%v,name:%vp,agroup:main_audio ", flags, version, resolutions[value])
			version++
		}
	}

	flags = fmt.Sprintf("%va:0,name:audio,agroup:main_audio,default:yes\" \"%v/%%v/playlist.m3u8\"", flags, video.Name)
		
	// -------------
	fmt.Printf("%v %v\n", command, flags)

	args, err := shell.Fields(flags, nil)
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil{
		fmt.Printf("\n\nDEU ERRO: %v", err)
		panic("")
	}

	fmt.Print(string(output))
}

// FULL COMMAND
/* 
ffmpeg -i "video de teste.mp4" \
      -filter_complex "\
  [0:v]split=2[v1][v2]; \
  [v1]scale=w=-2:h=720[v720_scaled]; \
  [v720_scaled]pad=w=1280:h=720:x=-1:y=-1:color=black[v720]; \
  [v2]scale=w=-2:h=480[v480_scaled]; \
  [v480_scaled]pad=w=854:h=480:x=-1:y=-1:color=black[v480] \
  " \
      -map "[v720]" -c:v:0 libx264 -preset medium -crf 23 -sc_threshold 0 \

      -map "[v480]" -c:v:1 libx264 -preset medium -crf 23 -sc_threshold 0 \

      -map a:0 \
      -c:a:0 aac -ac 2 -ar 48000 \
      -f hls \
      -hls_time 3 \
      -g 72 \
      -hls_playlist_type vod \
      -hls_list_size 0 \
      -hls_segment_filename "processed-video/%v/segment_%03d.ts" \
      -hls_flags independent_segments \
      -master_pl_name master.m3u8 \
      -var_stream_map "v:0,name:720p,agroup:main_audio v:1,name:480p,agroup:main_audio a:0,name:audio,agroup:main_audio,default:yes" processed-video/%v/playlist.m3u8
*/
