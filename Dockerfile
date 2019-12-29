FROM golang:1.13.4-alpine3.10 as builder
RUN apk add --no-cache --update bash git
WORKDIR /src
ADD ./go.mod ./go.sum ./
RUN go mod download
ADD ./ ./
RUN go build -o /dist/tbox_backend main.go

FROM alpine:latest
RUN apk add --update ca-certificates && \
    rm -rf /var/cache/apk/*

COPY --from=builder /dist/tbox_backend /app/bin/tbox_backend
COPY --from=builder /src/db/migrations /app/bin/db/migrations

WORKDIR /app/bin
CMD ["/app/bin/tbox_backend"]
