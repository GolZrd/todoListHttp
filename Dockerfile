FROM golang:1.24

COPY . .
RUN go build -o main cmd/main.go

CMD ["./main"] 