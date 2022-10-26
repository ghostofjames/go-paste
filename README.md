# go-paste

A tiny pastebin-like.

This was small project for learning Go.

## Installation

Download the latest [docker image](https://https://github.com/ghostofjames/go-paste/pkgs/container/go-paste) from the GitHub Container Registry.

```bash
docker pull ghcr.io/ghostofjames/go-paste:latest 
```

Run container, setting the `port` and `host`, as well as the `directory` to store pastes in.

```bash
docker run \
    -e HOST=<server-host>
    -p <port-on-host>:8000 \
    -v <path-on-host>:/files \
    ghcr.io/ghostofjames/go-paste
```

Alternatively, download the source, build and run using Go.

```bash
go mod download
go run .
```

### Configuration

Configuration can be done using enviromental variables when not using the docker image.

| Enviromental Variable | Description                      | Default     |
| ------------------    | ---------------------------------| ----------- |
| `HOST`                | Host for the server to listen on | `localhost` |
| `PORT`                | Port for the server to listen on | `8000`      |
| `FOLDER`              | Directory to store pastes        | `files`     |

## Usage

Upload a new file using a `curl` post request containing a `multipart/form-data` file.

```bash
cat file.txt | curl http://localhost:8000 -F file=@-
```

This will return a URL which can then be used to access the file.

```bash
http://localhost:8000/qs7kMtNVR
```

## Features

- [x] Write uploaded files to directory
- [x] Read files from directory and return as response
- [x] Configuration using enviromental variables
- [ ] Implement tests using [httptest](https://pkg.go.dev/net/http/httptest)
- [ ] Limit max file size
- [ ] Authentication?
