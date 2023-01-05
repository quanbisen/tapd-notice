FROM golang:1.19.2-alpine3.15 as build
RUN echo -e http://mirrors.ustc.edu.cn/alpine/v3.15/main/ > /etc/apk/repositories
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
WORKDIR /opt/tapd-notice
RUN apk add --no-cache gcc musl-dev
COPY . /opt/tapd-notice
RUN go build -ldflags="-s -w" -o main main.go

FROM alpine:3.15.0
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' > /etc/timezone
WORKDIR /opt/tapd-notice
COPY --from=build /opt/tapd-notice/main /opt/tapd-notice/config/settings.yml ./
EXPOSE 8080
ENTRYPOINT ["./main"]
CMD ["-h", "0.0.0.0", "-p", "8080", "-c", "settings.yml"]