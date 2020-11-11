# semtag

Tag your repository according to Semantic Versioning.

[pnikosis/semtag](https://github.com/pnikosis/semtag), but using Go.

## Installation

Using Homebrew

```bash
brew tap 1efty/tap/semtag
brew install 1efty/tap/semtag
```

## Usage

## Container

```bash
# pull the image from DockerHub
docker pull 1efty/semtag

# run on a git repository
docker run -it -v $PWD:/src -w /src --rm 1efty/semtag final
```
