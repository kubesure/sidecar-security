FROM golang:1.13.3-alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /go/src/sidecar-security
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -mod vendor

FROM scratch
WORKDIR /opt
COPY --from=builder /go/bin/sidecar-security .
EXPOSE 8000
ENTRYPOINT ["/opt/sidecar-security"]