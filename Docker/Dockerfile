FROM golang:1.18-bullseye as builder
LABEL maintainer="Crypta Electrica <crypta@crypta.tech>"
RUN mkdir /build 
RUN git clone https://github.com/Crypta-Eve/truth /build
WORKDIR /build 
ENV GO111MODULE on
ENV GOPROXY https://proxy.golang.org
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main cmd/cli/main.go


FROM alpine
RUN adduser -S -D -H -h /app appuser && apk --no-cache add ca-certificates
USER appuser
COPY --from=builder /build/main /app/
WORKDIR /app
ENTRYPOINT ["./main"]
