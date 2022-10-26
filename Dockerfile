## Build
FROM golang:1.18.1 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go .

RUN go build -o /go-paste

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /go-paste /go-paste

EXPOSE 8000

# USER nonroot:nonroot

CMD ["./go-paste"]