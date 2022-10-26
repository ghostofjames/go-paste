# go-paste

A tiny pastebin-like.

This was small project for learning Go.

## Installation

docker...

## Usage

Upload a new file using a `curl` post request containing a `multipart/form-data` file.

```
cat file.txt | curl http://localhost:8000 -F file=@-
```

This will return a URL which can then be used to access the file.

```
localhost:8000/qs7kMtNVR
```


## Features 
- [x] Write files
- [x] Read files
- [x] Configuration
