ARG GOLANG_VERSION=1.22

FROM golang:$GOLANG_VERSION-alpine as go-base

## set the environment variables.
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn

WORKDIR ./web-digger

# Copy depencies and modules' list.
COPY go.mod go.sum ./

RUN go mod tidy

## Copy the source code.
COPY . .

## Build the source code
RUN CGO_ENABLED=0 GOOS=linux go build -o /web-digger

## Expose the port
EXPOSE 8000

