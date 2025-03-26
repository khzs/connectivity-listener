FROM golang:latest AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o listener .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/listener .
EXPOSE 8080 50051
CMD ["./listener"]
