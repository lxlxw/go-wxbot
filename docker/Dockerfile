FROM alpine:latest
WORKDIR /app/

RUN apk add -U tzdata musl \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/shanghai" >> /etc/timezone

COPY ./go-wxbot /app/go-wxbot

ENTRYPOINT ["./go-wxbot"]
