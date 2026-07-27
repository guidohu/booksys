FROM golang:1.25-alpine AS builder

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

EXPOSE 9090

CMD ["/bin/booksys", "--http_port", "9090"]