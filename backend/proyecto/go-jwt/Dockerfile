FROM golang:1.20

WORKDIR /backend/proyecto/go-jwt

COPY . .

RUN go mod tidy
RUN go build -o main .

EXPOSE 8090

ENTRYPOINT ["./main"]