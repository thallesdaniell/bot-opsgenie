FROM golang:1.20 as builder

WORKDIR /app
COPY . ./
RUN cd /app/ && CGO_ENABLED=0 GOOS=linux go build -o opslink .

FROM alpine:latest

RUN apk add tzdata
COPY --from=builder /app/opslink /app/opslink

ENTRYPOINT [ "/app/opslink" ]