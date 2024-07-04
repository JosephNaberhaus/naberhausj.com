# naberhausj.com
My personal webiste

## Development
I've [Dockerized](./Dockerfile) dependencies you need for serving the project. The only two things you need installed locally are [Docker](https://www.docker.com/) and [Make](https://www.gnu.org/software/make/). The Make target will automatically build a Docker image with the required dependencies.

Run `make help` to see the available commands.

## Scripts
Some scripts I've used when writing articles.

### Encode video to WebM
Encode a video to WebM. Make sure to adjust the dimensions to match the aspect ratio of the original video.

```shell
ffmpeg -i <in> -s 720x1280 -vcodec libvpx -acodec libvorbis output.webm
```

Use the `--b:v` flag to adjust the average bit rate (e.g. `--b:v 1000K`)
Use the `-crf` flag to specify a quality that the encode will attempt to maintain (e.g. `--crf 10`)

### Encode video to GIF

```shell
ffmpeg -i <in> -vf "fps=10,scale=<width>:-1:flags=lanczos,split[s0][s1];[s0]palettegen[p];[s1][p]paletteuse" -loop 0 output.gif
```