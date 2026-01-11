## Base Image Go
FROM golang:1.25.5 AS base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

## 本番用のバイナリ
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o main .

## for debug
FROM golang:1.25.5 AS debug

WORKDIR /app

RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY . .

EXPOSE 1323
EXPOSE 2345

CMD ["dlv","debug","--headless","--listen=:2345","--api-version=2","--accept-multiclient","/app/main.go"]

FROM golang:1.25.5 AS final


COPY --from=base /app/main /app/main

#RUN chmod 777 /app/main


EXPOSE 1323
CMD ["/app/main"]