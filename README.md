<img width="1900" height="802" alt="Frevod logo" src="https://github.com/user-attachments/assets/47601dd5-6ce1-4723-ac6f-eca362d6a6da" />

---

## What is Frevod?
Frevod is a program that takes a video from _anywhere_, encodes it and stores HLS output _anywhere_ as well. It can be deployed on a container as a long-living worker acting as a microservice for your video pipeline, or be used as a CLI tool if you just need to encode a video here and there.

## Why just not use an FFmpeg script?
### Although Frevod uses ffmpeg under the hood, that's not all there is for this software:
* **It is agnostic:** Meaning it can download a video from a URL, encode it and push it to S3, or encode a local video, upload it to MinIO and also store the output locally.
* **Auto ABR:** Automatically makes multiple renditions (aka: multiple versions with different resolutions | 2160p, 1440p, 1080p...)
* **Handles DAR:** No matter if the original video is 16:9, 9:16 or 4:3, output will be the same.
* **Internal queue:** CPU won't explode when multiple videos are uploaded at the same time.
* **Very little RAM usage:** Downloads are stored directly into disc, meaning Frevod only uses a couple MBs of RAM.
* **Auto-cleanup:** If something goes wrong, it will automatically clean everything up, ensuring no ghost files are left anywhere, _even on remote storage!_.
* **Declarative config:** It accepts options through a simple `config.yaml` file, so you don't have to learn multiple complex FFmpeg flags, _which FFmpeg is famous for_.
* **Easy to deploy:** You can spin a container with a couple env vars, run it locally with a `config.yaml` file or rely on Frevod's fallback standard options.

> These are some of the things Frevod does right now, but it'll keep evolving, there's a "planned features" section later on this README.

## What options?
There's a couple options to choose from, divided into three categories.

### Ingest Options
> How will the videos get to Frevod? A REST API? Maybe an SQS queue.

**Input methods currently supported:**
* CLI
* REST
  
``` yaml
ingest:
  cli:
    enabled: true
  rest:
    enabled: false
    port: ":1137"
```

### Encode Options
> How will the videos be encoded? H.264 or H.265? How many videos do you want to encode at once? 
``` yaml
encode:
  concurrent_encodings: 2 
  output_ffmpeg_command: false
```

### Dispatch Options
> Where will Frevod store the output? A MinIO bucket? A local directory on your computer? Maybe both.

**Output methods currently supported:**
* Any storage compatible with the S3 API (Amazon S3, Cloudflare R2, MinIO...)
* Local storage
  
``` yaml
dispatch:
  s3:
    enabled: false
    bucket:
    endpoint: 
    access_key:
    secret_key:
    use_ssl: true
  local:
    enabled: true
    store_at: store/it/here/please
```

## How to use each input method?

---

### CLI:
Type the URI (a URL or local path) and the KEY (the name of the video after encoding) as arguments **after** running the executable, these are **not** flags.

Example: `http://coolvideos.com/funnycatvideo.mp4 funnycatvideo`

---

### REST:

Make a POST request to `/encode` with URI and KEY.

Example:
``` json
  {
    "uri": "/path/on/your/computer",
    "key": "making burgers with my brother",
  }
```

---

_Notice how you can get an input from the web or from local storage no matter what input method you're using_

## Nice, how can I use it?

### Most people will want to use this in Docker
### Here's a compose.yml example:
``` yaml
services:
  frevod:
    image: frevod:0.1.0 
    container_name: frevod-worker
    ports:
      - "1137:1137"
    environment:
      - BUCKET=    # Frevod automatically starts the REST API and the S3 connection when running in docker 
      - ENDPOINT=
      - ACCESS_KEY=
      - SECRET_ACCESS_KEY=
      - USE_SSL=true
      - CONCURRENT_ENCODINGS=2   # Encoding is heavy on the CPU, go easy on this number
    
    # if you want to encode videos from a local path
    # you can use a volume like in the example below
    # then input the docker dir as a uri. 
    volumes:
      - ./path/on/my/computer:/path/inside/docker
      # in this example, uri would be: "/path/inside/docker/{video-name.mp4}"
```

_You **can** use the CLI input method with Docker too, but... Are you sure you wanna do it this way?_ 🤨

## For local usage, check the [releases](https://github.com/Matheus-Rossetti/frevod/releases) page so you can download Frevod for your OS.
### Here's a config.yaml example:
``` yaml
ingest:
  cli:
    enabled: true
  rest:
    enabled: false
    port: ":1137"

encode:
  concurrent_encodings: 2 
  output_ffmpeg_command: false
  
dispatch:
  s3:
    enabled: false
    bucket:
    endpoint:
    access_key:
    secret_key:
    use_ssl: true
  local:
    enabled: true
    store_at: path/on/your/computer
```

_Docs end here, enjoy! <3_

# Planned features
This section will cover thing I want Frevod to support in the future, but also some things it does right now that I'm not fully satisfied with the result.

* **Save job state:** Right now, if something goes wrong during a job (like a network error, or maybe the user just closes the program mid job), Frevod just deletes everything about the specific job and jumps to the next. I made it this way because it was fast and I didn't want the initial version to leave trash around. The ideal would be to save the state of the job and retry from that point onward.
* **Built-in metrics:** Expose internal metrics like average job time, job completion progress, CPU usage, etc... To be consumed by Prometheus.
* **Multiple input / output methods:** Add support for RabbitMQ, Kafka, SQS, gRPC, Azure Blob, Google Cloud Storage, generic upload URLs and any other method required by the community.
* **Two pass encoding / per input bitrate** Bitrates are fixed for each resolution, it's enough for most videos but not ideal as the current bitrates may be too much for some videos and too little for others.
* **Multiple codecs:** Currently, Frevod supports H.264 and AAC, which is the standard for vod streaming, but additional support for H.265 and Opus would be nice, also, DASH support.
* **Notifier:** A notifier system to, amazingly, notify someone on Discord or Slack when a certain event occurs, like too many videos on the queue or a critical error.
* **Cool TUI:** Most users will run Frevod in a container, but CLI users need some love too. _Bubble Tea_

### I believe that, with the above features completed, a solid v1.0 can be released.
### _but I do have some plans for the far far future:_

* **Autoscalling:** Spawn or kill encoders based on queue size.
* **Zero disk mode:** From stream to ffmpeg, from ffmpeg to storage. Gotta make sure to add a RAM limit.
* **Custom FFmpeg:** Compile a fork of FFmpeg with only functions that Frevod uses and bundle it in the container.

## Some questions you might have

### Why does this exist?
At first, I was making a streaming service as a portfolio project but decided to just focus on the 'encode' part of it and make it as accessible as possible, you may have noticed that Frevod fits in every systems workflow, its built this way from the start. Currently, there's no solid option for a free, self hosted encoding service for vod streaming and I want to change that.

### Why 'Frevod'?
Frevod stands for "Free Encoder for Video on Demand", quite simplistic, I know.

### What's up with the logo?
I tried making something more common to this infra tool / cloud market, but man... let's just say it didn't come out right.. So I decided to make something different that breaks expectation. I've been playing Hades the last couple days so I based the logo on the game's art style. _Even considered calling 'features' -> 'boons' at some point, lol_.

Need help? Have any more questions? Join the [Discord server](discord.gg/uVCpJTK6Dj)!