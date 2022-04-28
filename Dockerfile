FROM golang:1.17 as builder

ARG VERSION=unknown

WORKDIR /app
COPY . .
RUN GO111MODULE=on GOOS=linux CGO_ENABLED=0 \
    go build -ldflags "-s -w -X main.version=${VERSION}" \
    -o /app/build/cmd/gophermart cmd/gophermart/*.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/build/cmd/gophermart /bin/cmd/gophermart

ENTRYPOINT ["/bin/cmd/gophermart"]