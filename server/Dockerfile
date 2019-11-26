FROM golang:latest

WORKDIR .

COPY . .

RUN go build -o server .

EXPOSE 9999

CMD ["./server"]