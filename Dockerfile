FROM golang:1.24.2-bookworm AS golang-orbproxy-builder

WORKDIR /lib

COPY . /lib

RUN --mount=type=cache,target=$HOME/go \
    --mount=type=cache,target=$HOME/.cache/go-build \
    cd orbproxy; go build -o ../dist/orbproxy.so -buildmode=c-shared .

FROM envoyproxy/envoy:contrib-v1.33.2

COPY --from=golang-orbproxy-builder /lib/dist/orbproxy.so /lib/orbproxy.so