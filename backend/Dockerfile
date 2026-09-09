# Keep this in step with the go directive in api/v2/go.mod. The golang images
# set GOTOOLCHAIN=local, so a builder older than that directive fails at
# `go mod download` rather than fetching a newer toolchain.
FROM golang:1.26-alpine AS builder

ARG BOOKSYS_VERSION
ARG BOOKSYS_COMMIT
ARG BOOKSYS_DATE
ARG MODULE_PATH="server"

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY api/v2 api/v2
COPY api/v2/sample_config.yaml /sample_config.yaml
RUN echo ${BOOKSYS_COMMIT}
RUN cd api/v2; go mod download && go mod verify
RUN cd api/v2; go build -v -a -trimpath -o /build/app \
    -ldflags="-X '${MODULE_PATH}/version.Release=${BOOKSYS_VERSION}' -X '${MODULE_PATH}/version.Commit=${BOOKSYS_COMMIT}' -X '${MODULE_PATH}/version.BuildDate=${BOOKSYS_DATE}'" \
    ./main.go

FROM alpine:3.24 AS final

RUN addgroup -S booksys && adduser -S -G booksys booksys

COPY --from=builder /build/app /bin/booksys
COPY --from=builder /sample_config.yaml /etc/booksys/sample_config.yaml
RUN mkdir -p /var/lib/booksys/uploads && chown -R booksys:booksys /var/lib/booksys

USER booksys

# 9090 serves the API. 9091 is the debug listener (pprof, expvar) and is only
# opened when BOOKSYS_DEBUG_PORT is set, because a heap profile taken through
# it contains session secrets and the database password. Never publish it.
EXPOSE 9090 9091

CMD ["/bin/booksys", "--http_port", "9090"]