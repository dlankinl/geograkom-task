FROM golang:1.22-bullseye

WORKDIR /app

ADD go.mod .

COPY . .

RUN go mod download

RUN go build -o /build cmd/main.go

CMD ["/build"]
