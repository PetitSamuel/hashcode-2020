FROM golang:latest
ADD . /
COPY . /
WORKDIR /
CMD go run hash.go