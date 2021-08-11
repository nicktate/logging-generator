ARG DEV_IMAGE=ntate22/logging-generator:dev
FROM ${DEV_IMAGE} as builder

ENV GOPATH /go
WORKDIR /go/src/github.com/nicktate/logging-generator
COPY . .
RUN go build ./cmd/logging-generator

FROM debian:buster-slim
WORKDIR /app
COPY --from=builder /go/src/github.com/nicktate/logging-generator/logging-generator /app/logging-generator
CMD ["/app/logging-generator"]
