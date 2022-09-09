FROM golang:1.19.1-alpine3.16 AS builder

WORKDIR /build/
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" \
    rusprofile-wrapper/cmd/server

#-

FROM alpine:3.16 AS certificates

RUN apk --no-cache add ca-certificates

#-

FROM scratch

WORKDIR /api/

ENV PATH="/api/bin/:$PATH"
COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /build/internal/rpc_server/ /internal/rpc_server
COPY --from=builder /build/server /bin/gateway

EXPOSE 8000

CMD gateway