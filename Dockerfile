# Build stage
FROM golang:alpine AS builder

RUN apk add make git

WORKDIR /app

COPY app/go.mod app/go.sum ./
RUN go mod download && go mod verify

COPY app .
COPY LICENSE LICENSE
RUN make -B build

# Run stage
FROM scratch as server

# Get latest CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/LICENSE /LICENSE
COPY --from=builder /app/tf-server /tf-server
COPY --from=builder /app/www /www

EXPOSE 8080
ENTRYPOINT [ "/tf-server" ]

LABEL Name=team-feedback Version=0.0.0

# Test stage
FROM golang:alpine as test
RUN apk add make git build-base

WORKDIR /go/test
COPY app/go.mod app/go.sum ./
RUN go mod download && go mod verify

COPY app/Makefile ./
RUN CGO_ENABLED=0 make install-golangci

COPY app .
RUN CC=gcc CGO_ENABLED=1 make test
