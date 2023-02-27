FROM alpine:edge AS builder
LABEL maintainer="Samuel Contesse <samuel.contesse@morphean.com>"
RUN apk add --update-cache alpine-sdk go libgit2-dev=1.5.1-r0 && rm -rf /var/cache/apk/*
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:edge
RUN apk add --update-cache libgit2=1.5.1-r0 && rm -rf /var/cache/apk/*
COPY --from=builder /build/main /
ENTRYPOINT ["/main"]
