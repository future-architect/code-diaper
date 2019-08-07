# Build Stage
FROM golang:1.12.5 AS build
ENV REPOSITORY github.com/future-architect/code-diaper

# aiming layer cache(module download)
WORKDIR $GOPATH/src/$REPOSITORY
ENV GO111MODULE on
COPY go.mod go.sum ./
RUN go mod download
# If code was updated then execute below without using the cache
ADD . $GOPATH/src/$REPOSITORY
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-s -w' -a -installsuffix cgo -o /codediaper cmd/codediaper/codediaper.go

# Runtime Stage
FROM alpine:3.10.1
RUN apk add --no-cache ca-certificates
COPY --from=build /codediaper .
CMD ["./codediaper"]
