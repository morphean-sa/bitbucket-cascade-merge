FROM golang:1.13-alpine AS builder
LABEL maintainer="Samuel Contesse <samuel.contesse@morphean.com>"
RUN apk add --update-cache alpine-sdk libgit2-dev && rm -rf /var/cache/apk/*
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine
RUN apk add --update-cache libgit2 && rm -rf /var/cache/apk/*
COPY --from=builder /build/main /
ENTRYPOINT ["/main"]
