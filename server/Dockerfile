### builder stage ###

FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

### image stage ###

FROM alpine

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 9999

CMD ["./server"]