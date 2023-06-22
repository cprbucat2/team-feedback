# Build stage
FROM golang:alpine AS builder

RUN apk add make git

WORKDIR /app

COPY app/go.mod app/go.sum ./
RUN go mod download && go mod verify

COPY app .
COPY LICENSE LICENSE
RUN make build

# Run stage
FROM scratch as server

# Get latest CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/LICENSE /LICENSE
COPY --from=builder /app/tf-server /tf-server
COPY --from=builder /app/www /www

CMD [ "/tf-server" ]

LABEL Name=team-feedback Version=0.0.0
EXPOSE 8080
