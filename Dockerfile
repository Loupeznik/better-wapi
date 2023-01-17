FROM golang:1.19-alpine as builder

RUN apk update && apk add --no-cache alpine-sdk git tzdata

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@v1.8.0
RUN swag init

RUN go build -o ./app ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/app .
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENV GIN_MODE=release
ENV TZ=Europe/Prague
EXPOSE 8000

ENTRYPOINT ["./app"]
