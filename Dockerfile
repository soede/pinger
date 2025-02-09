FROM golang:1.23

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

RUN go mod verify

ENV INFO=12
EXPOSE 8080

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/api

CMD ["app"]
