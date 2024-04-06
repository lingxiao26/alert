FROM golang:1.22-alpine as builder
WORKDIR /app
ARG VERSION
ENV GOPROXY=https://goproxy.cn
COPY ./go.mod ./
COPY ./go.sum ./
#RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o alert .

FROM alpine as certs
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk update && apk add ca-certificates tzdata

FROM busybox as runner
COPY --from=builder /app/alert /app
COPY --from=certs /etc/ssl/certs /etc/ssl/certs
COPY --from=certs /usr/share/zoneinfo /usr/share/zoneinfo
ENTRYPOINT ["/app"]