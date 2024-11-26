FROM golang:1.23.3
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o ./server ./cmd/server
CMD ["./server"]

