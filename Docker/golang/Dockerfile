FROM golang:1.12-alpine

RUN apk update && apk upgrade && \
	apk add --no-cache bash git openssh

WORKDIR /app

RUN go mod download
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/handlers
RUN go get github.com/streadway/amqp

COPY . .

EXPOSE 5000

CMD ["go", "run", "app.go"]
