#### - DEV - ####
FROM golang:1.24.2 AS dev

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY ./go.mod go.mod
COPY ./go.sum go.sum
RUN go mod download

COPY ./ ./

RUN swag init -d cmd
CMD ["go", "run", "cmd/main.go"]

#### - TESTS - ####
FROM golang:1.24.2 AS tester

WORKDIR /app

COPY ./go.mod go.mod
COPY ./go.sum go.sum
RUN go mod download

COPY ./ ./

CMD ["go", "test", "-v", "./..."]

#### - BUILDER - ####
FROM golang:1.24.2 AS builder

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY ./go.mod go.mod
COPY ./go.sum go.sum
RUN go mod download

COPY ./ ./

RUN swag init -d cmd && \
    go build -o /bin/main cmd/main.go


#### - SERVER - ####
FROM alpine:3.19.1 AS server

RUN apk add --no-cache gcompat=1.1.0-r4 libstdc++=13.2.1_git20231014-r0
# RUN apk add --no-cache gcompat libstdc++

WORKDIR /app

COPY --from=builder /bin/main ./main
COPY --from=builder /app/internal/templates ./internal/templates

RUN adduser --system --no-create-home nonroot
USER nonroot

ENV GIN_MODE=release

EXPOSE 8080

CMD ["./main"]
