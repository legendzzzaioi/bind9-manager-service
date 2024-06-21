FROM golang:1.22.3-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 1
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata gcc musl-dev

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY ./etc /app/etc
RUN go build -ldflags="-s -w" -o /app/bind9-manager-service .

FROM internetsystemsconsortium/bind9:9.16
ENV TZ Asia/Shanghai
WORKDIR /app

RUN apt-get update && apt-get install -y musl-dev && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
COPY --from=builder /app/bind9-manager-service /app/bind9-manager-service
COPY --from=builder /app/etc /app/etc

VOLUME ["/etc/bind"]

EXPOSE 53/udp 53/tcp 953/tcp 8000/tcp

CMD ["./bind9-manager-service", "-f", "etc/bind9-api.yaml"]