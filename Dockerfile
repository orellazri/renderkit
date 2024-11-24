FROM golang:1.23 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/renderkit .

FROM debian:bookworm
COPY --from=builder /bin/renderkit /bin/renderkit
ENTRYPOINT [ "/bin/renderkit" ]
