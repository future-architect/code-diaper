# Build Stage
FROM golang:1.12.5 AS builder
ENV REPOSITORY github.com/future-architect/code-diaper
ADD . $GOPATH/src/$REPOSITORY
WORKDIR $GOPATH/src/$REPOSITORY
RUN GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-s -w' -a -installsuffix cgo -o /codediaper cmd/codediaper/codediaper.go

# Runtime Stage
FROM alpine:3.10.1
RUN apk add --no-cache ca-certificates
COPY --from=builder /codediaper .
CMD ["./codediaper"]
