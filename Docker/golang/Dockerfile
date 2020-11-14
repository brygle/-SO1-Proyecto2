FROM golang:1.14

WORKDIR /go/src/app
COPY go.mod .
RUN go mod download
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
EXPORT 5000
CMD ["app"]
