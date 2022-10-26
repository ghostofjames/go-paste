# go-paste

A tiny pastebin-like.

This was small project for learning Go.

## Installation

Run using the docker [image](https://).

```bash
docker run -p <port-on-host>:8000 -v <path-on-host>:/files go-paste
```

## Usage

Upload a new file using a `curl` post request containing a `multipart/form-data` file.

```bash
cat file.txt | curl http://localhost:8000 -F file=@-
```

This will return a URL which can then be used to access the file.

```bash
localhost:8000/qs7kMtNVR
```

## Features

- [x] Write files
- [x] Read files
- [x] Configuration
