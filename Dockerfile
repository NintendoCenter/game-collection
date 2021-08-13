FROM golang:1.16 as builder

RUN go get github.com/go-delve/delve/cmd/dlv

ARG PROXY=
ARG SERVICE_PATH="."
WORKDIR /usr/src
COPY go.mod .
COPY go.sum .
RUN GOPROXY=${PROXY} go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -gcflags="all=-N -l" -installsuffix cgo -o service ./${SERVICE_PATH}

FROM alpine:latest
RUN apk add --no-cache ca-certificates libc6-compat tzdata

WORKDIR /usr/app
COPY --from=builder /go/bin/dlv .
COPY --from=builder /usr/src/service .
CMD ["./dlv", "--listen=:40000", "--headless", "--continue", "--api-version=2", "--accept-multiclient", "exec", "./service"]
