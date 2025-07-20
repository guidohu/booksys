FROM golang:1.24

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY api/v2 api/v2
COPY api/v2/sample_config.yaml /config.yaml
RUN cd api/v2; go mod download && go mod verify

RUN ls ./api/v2
RUN cd api/v2; go build -v -o /usr/local/bin/app ./main.go
EXPOSE 9090

CMD ["app", "--http_port", "9090"]