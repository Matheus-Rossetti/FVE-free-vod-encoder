package core

// TODO use fmt.Sprintf("") to concat strings

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
