FROM golang:1.21-alpine as builder

WORKDIR /app

# COPY go.mod, go.sum and download the dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# COPY All things inside the project and build

COPY . ./

RUN go build -o main cmd/main.go

FROM alpine:3.19.0 as development

WORKDIR /app

COPY --from=builder /app/main main
COPY --from=builder /app/env.docker .env



EXPOSE 8080
CMD [ "./main" ]

FROM alpine:3.19.0 as production

WORKDIR /app

COPY --from=builder /app/main main

EXPOSE 8080
CMD [ "./main" ]