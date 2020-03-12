FROM golang:1.13.5

WORKDIR /go/src/github.com/4726/game/services/matchmaking/queue
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o queue .

CMD ["./queue"]

EXPOSE 14000