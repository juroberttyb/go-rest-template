# Modules caching
FROM golang:1.21-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Builder
FROM golang:1.21-alpine as builder
RUN apk --update add ca-certificates git make bash build-base
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN make app

# Runtime
FROM alpine:latest
RUN apk add --update --no-cache tzdata
COPY --from=builder /bin/app /app
COPY --from=builder /app/database/migrations /database/migrations
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENV TZ=Asia/Taipei
CMD ["/app"]
