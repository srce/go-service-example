FROM golang:1.15 as build

WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN go mod download
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /build/go-service-example ./cmd/go-service-example
