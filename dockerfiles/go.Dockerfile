# first stage: build
FROM golang:1.18-alpine as builder

WORKDIR /usr/app

COPY ./go .

RUN go env -w GOPROXY="https://goproxy.io,direct"

RUN go env GOPROXY

RUN go mod tidy

RUN go build -o app ./cmd/api/main.go

# second stage: run app
FROM golang:1.18-alpine

WORKDIR /usr/app

COPY --from=builder /usr/app/app .

COPY --from=builder /usr/app/.env.production .

RUN apk --no-cache add curl

EXPOSE 3001

CMD [ "app" ]