FROM golang:1.19.4-alpine3.17 AS builder
RUN apk --no-cache add gcc musl-dev git curl bash
RUN go install github.com/revel/cmd/revel@latest
WORKDIR /work
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . ./
ENV APP_PACKAGE="github.com/pottava/gpu-node-manager" \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN revel build --application-path . --target-path /tmp/app -m prod \
    -X "${APP_PACKAGE}/src/app/util.Date=$( date +%Y-%m-%d --utc )"

FROM alpine:3.17
RUN apk --no-cache add bash ca-certificates tzdata
COPY --from=builder /tmp/app /app
ENTRYPOINT ["/app/run.sh"]
