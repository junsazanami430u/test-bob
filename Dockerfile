# Builder
FROM golang:1.18-alpine AS builder
ENV USER=appuser
ENV UID=10001
ENV GO111MODULE=on

RUN adduser --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" --no-create-home --uid "${UID}" "${USER}"
RUN apk update && apk add --no-cache ca-certificates tzdata && update-ca-certificates
WORKDIR /app

COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 go build -ldflags='-w -s -extldflags "-static"' -a -o /go/bin/gin-example
RUN chmod +x /go/bin/gin-example


# Minimalistic image
FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/bin/gin-example /usr/bin/gin-example

USER appuser:appuser

ENTRYPOINT [ "/usr/bin/gin-example" ]