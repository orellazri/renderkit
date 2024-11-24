FROM golang:1.23 AS builder
ENV CGO_ENABLED=0
ENV GOOS=linux
ARG VERSION
ENV VERSION=$VERSION
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w -X main.version=$VERSION" -o /bin/renderkit .

FROM alpine:3
COPY --from=builder /bin/renderkit /bin/renderkit
ENTRYPOINT [ "/bin/renderkit" ]
