#### - DEV - ####
FROM golang:1.23.3 AS dev

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY src/go.mod go.mod
COPY src/go.sum go.sum
RUN go mod download

COPY src/ ./

RUN swag init
CMD ["go", "run", "."]

#### - TESTS - ####
FROM golang:1.23.3 AS tester

WORKDIR /app

COPY src/go.mod go.mod
COPY src/go.sum go.sum
RUN go mod download

COPY src/ ./

CMD ["go", "test", "-v", "./..."]

#### - BUILDER - ####
FROM golang:1.23.3 AS builder

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY src/go.mod go.mod
COPY src/go.sum go.sum
RUN go mod download

COPY src/ ./

RUN swag init && \
    go build -o /bin/main main.go


#### - SERVER - ####
FROM alpine:3.19.1 AS server

RUN apk add --no-cache gcompat=1.1.0-r4 libstdc++=13.2.1_git20231014-r0
# RUN apk add --no-cache gcompat libstdc++

WORKDIR /app

COPY --from=builder /bin/main ./main
COPY --from=builder /app/templates ./templates

RUN adduser --system --no-create-home nonroot
USER nonroot

ENV GIN_MODE=release

EXPOSE 8080

CMD ["./main"]
