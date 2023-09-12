FROM golang:alpine

COPY . /app
WORKDIR /app

ENV RUN_PORT=":8080"
ENV BASE_URL="127.0.0.1:8080"

RUN go build cmd/planning-poker/main.go

CMD ./main