FROM golang:1.17 AS builder

COPY . /src
WORKDIR /src

RUN mkdir bin && go build -o ./bin ./... && mv bin/$(ls bin) bin/server

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app

WORKDIR /app

EXPOSE 8000
VOLUME /etc/app-configs

CMD ["./server", "-conf", "/etc/app-configs/config.yaml"]
