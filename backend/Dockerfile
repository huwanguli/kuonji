# Stage 1: build
FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /kuonji .

# Stage 2: run
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

COPY --from=builder /kuonji /kuonji
COPY config/config.yaml /config/config.yaml

RUN mkdir -p /uploads /logs

EXPOSE 8080

CMD ["/kuonji"]
