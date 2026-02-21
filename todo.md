
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

    - Make changes to it on the fly, like open or close an input or output method or spawn or despawn a worker.


    Actual todo:
    - There's a bug in the command builder that sets the same bitrate to all renditions, 480p rendition has bitrate for 1080p
    - Fixed bitrate for all videos of the same resolution isn't a good idea, check out VMAF from netflix
    - Or even 2-pass encoding
    - Also, set the -g to an actual real value, meaning we gotta get the framerate from ffprobe and do some simple math
    - after that, you can go back to:
    - adding retry logic to downloads
    - limiting the amount of videos inside downloaded-videos/ dir 
    - while not limiting new inputs, maybe a queue to download videos just like we have for the workers?
    - Dunno, good luck, future me!