FROM golang:1.11 AS builder
ENV GO111MODULE=on

WORKDIR /opt

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -race -a -installsuffix nocgo -o /opt/test ./main.go &&\
    chmod +x /opt/test

# race detector requries glibc
FROM debian:stretch
COPY --from=builder /opt/test /opt/test

ENTRYPOINT ["/opt/test"]
