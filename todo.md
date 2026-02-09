- Install ffmpeg if it isn't found
- Add the option for the user to use a custom ffmpeg path
- Get info from the video using ffprobe
- Build the ffmpeg command
- Accept flags for options like the destination of the encoded video
- Make FVE read a json, xml or watever file to get configs
- Add unit testing for the core package 
- Add metrics with prometheus
- Make it be able to listen to a rabbitmq or kafka queue
- Add a REST API
- Add a GraphQL API
- Add a gRCP API
- Make it able to download videos from somewhere
- Make it able to uplaod videos to S3 like services
- Containerize it (bundle ffmpeg with it)
- Make a CLI tool to use it
- Make a UI that consumes the REST API, native and web, with flutter.


initially I wanted to call it FEVOD (Free Encoder for Video on Demand) but the 'FE'
sounded too much like 'Fee', and I didin't want it to be anywhere close to a paid product
so I changed it to FREVOD.




FFmpeg the way Redis is to caching or RabbitMQ is to messaging