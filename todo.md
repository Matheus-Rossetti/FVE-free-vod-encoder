
- Add the option for the user to use a custom ffmpeg path
- Get info from the video using ffprobe
- Build the ffmpeg command
- Accept flags for options like the destination of the encoded video
- Make FVE read a json, xml or watever file to get configs
- Add unit testing for the core package 
- Add metrics with prometheus
- Add optional metrics to expose how many videos were processed and things alike
- Make it be able to listen to a rabbitmq or kafka queue
- Add a REST API
- Add a gRCP API
- Make it able to download videos from somewhere
- Make it able to uplaod videos to S3 like services
- Containerize (bundle ffmpeg with it)
- Make a CLI tool
- Make a UI that consumes the REST API, native and web, with flutter for non power users.


Acceptable Beta version:
- Encode videos to HLS with fixed configs (segment_time, bitrate, etc...)
- Accepts terminal input with fairly nice TUI
- Accepts REST input with basic rate limiting and no auth
- Limit FFmpeg CPU usage
- Limit amount of concurrent encodings
- Basic tests 
-  


Acceptable 1.0 version:
    - Download videos
    - Limit ffmpeg CPU usage
    - Graceful shutdown with context package
    - REST API
    - Terminal input
    - RabbitMQ downloads from links on queue
    - Upload
    - Adaptable command builder based on video DAR
    - metrics for prometheus
    - docker
    - testing