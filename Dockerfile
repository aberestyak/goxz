ARG PROXY_REGISTRY=

FROM ${PROXY_REGISTRY:+$PROXY_REGISTRY/}golang:1.16-alpine3.13 as builder

ARG VERSION

COPY . /go/src/goxz
WORKDIR /go/src/goxz

ENV CGO_ENABLED=0

RUN go install \
    -installsuffix "static" \
    -ldflags "                                          \
      -X main.Version=${VERSION}                        \
      -X main.GoVersion=$(go version | cut -d " " -f 3) \
      -X main.Compiler=$(go env CC)                     \
      -X main.Platform=$(go env GOOS)/$(go env GOARCH) \
    " \
    ./...

FROM ${PROXY_REGISTRY:+$PROXY_REGISTRY/}alpine:3.13 as runtime

COPY --from=builder /go/bin/goxz /goxz

ENTRYPOINT ["/goxz"]
