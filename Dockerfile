FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/renderkit cmd/renderkit/main.go

ENTRYPOINT ["/app/renderkit"]
