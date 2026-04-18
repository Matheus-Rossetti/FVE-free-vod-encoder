# Frevod's Technical Breakdown.

## Objective

Frevod's objective is to facilitate turning a video into HLS segments for video on demand streaming of which, before Frevod, your options were:
- Pay cloud services.
- Learn FFmpeg and its highly complex flags.

Frevod aims to accomplish this while fitting in every encoding pipeline out there, so it needs to be agnostic, meaning: receive input from anywhere, anyhow, offer multiple encoding options and upload to anywhere, anyhow.
As of yet, it doesn't fully manage all this, but it will get there.

### Currently Supports:
**Ingest:** CLI and REST.
**Encode:** Concurrent encodings, fixed bitrate, ABR (From 4K to 480p) and perfectly handles different DARs.
**Dispatch:** Local and storages compatible with the S3 API.
> With more planned.

## Internal Architecture
Frevod uses a pipeline pattern for its internal architecture, it is composed of 5 different parts:
- Ingest Methods -> Ingestor -> Encoder -> Dispatcher -> Dispatch Methods.
### Visual Representation: 
<img width="1616" height="366" alt="Untitled-2026-03-30-1908" src="https://github.com/user-attachments/assets/d6463181-afc2-4b38-bf4b-948fea9ccfca" />
The Ingestor, Encoder and Dispatcher each listen to its own queue, those being `ingestQueue`, `encodeQueue` and `dispatchQueue`.

## Code
### Ingest and Dispatch Methods
The Go code is built in a way that's easy to add new *ingest* and *dispatch* methods.

**Ingest methods** call `core.PushJob(uriType, uri, keyStarter, methodName)` after processing the input. The decision of making *each ingest* method inject the uriType into the job was made with ease of feedback in mind
eg: REST can return a 422 (Unprocessable entity) if the uri is unsupported. By letting each *ingest method* call `core.CategorizeUri()` instead of calling it just once in the Ingestor, we add repeated code through all
*ingest methods* in exchange of **not** needing an error channel from the *Ingestor* to the respective methods. I much prefer repeating a function call than adding complexity, besides, all *ingest methods* need to validate their input anyway.

**Dispatch methods** work a little different since instead of pushing a job, they receive a job. Each *dispatch method* has to implement an interface with a `Dispatch(context, job)` function.
The thing with *dispatch methods* is that each one works very differently from one another, uploading output to S3 is very different than storing it locally. All the Dispatcher actually do is call the `Dispatch` function from
each *dispatch method*, log and clean local storage if the option to store locally is disabled.
The *dispatch methods* are initialized when the program starts and stored in an array, of which the Dispatcher iterates through.

### Encoder
The encoder gets metadata from a video, builds an ffmpeg command for it using this metadata and some options the user can choose with a config.yaml file, then run ffmpeg with that specific command.
*Note:* config.yaml is optional, Frevod does have default configs.
> It's worth to mention that the FFmpeg bin runs as low-priority on both Windows and Unix, so it will use 100% os spare compute but won't throttle the machine.
The specific implementation can be quite overwhelming at first, but upon close inspection, you'll notice it's quite simple.
An argument can be made about running the FFmpeg bin for each job is non-optimal, but in the grand scheme of things, a couple milliseconds added to a 10min encoding process is not significant enough as to make an improved (and more complex) implementation necessary.

## That's the gist of it!
If you have any questions, critiques, advice or just want to talk about it, please, join the [Discord Server](discord.gg/uVCpJTK6Dj)!

