FROM golang:latest

WORKDIR /app

COPY ./ /app

RUN go get -v ./...

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon \
                        -color=true \
                        -exclude-dir="./tools/dev/db_data" \
                        -build="go build -v -o ./tools/dev/bin/hydric-linux-x86_64 ./cmd/hydric" \
                        -command="./tools/dev/bin/hydric-linux-x86_64 -config tools/dev/config.json"