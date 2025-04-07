FROM docker.io/library/golang:1.24.2-bookworm AS golang-builder

WORKDIR /app
COPY . /app

RUN --mount=type=cache,target=$HOME/go \
    --mount=type=cache,target=$HOME/.cache/go-build \
    CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o ./dist/server ./cmd/server/...

FROM gcr.io/distroless/static AS final

LABEL maintainer="jochumdev"
USER nonroot:nonroot

COPY --from=golang-builder --chown=nonroot:nonroot /app/dist/server /helloworld

ENTRYPOINT [ "/helloworld" ]


