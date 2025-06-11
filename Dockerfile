## Base Image Go
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

#RUN GOOS=linux GOARCH=amd64 go build -o server .
RUN go build -o main .

# 軽量な実行環境
FROM golang:1.23 AS final

#RUN apk --no-cache add ca-certificates

COPY --from=builder /app /app

RUN chmod 777 /app/main

EXPOSE 1323

CMD ["/app/main"]