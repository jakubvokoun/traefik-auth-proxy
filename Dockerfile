FROM golang:1.22-alpine

RUN mkdir /app

COPY go.mod /app
COPY main.go /app

WORKDIR /app
RUN go build -o app

EXPOSE 8080

CMD [ "./app" ]
