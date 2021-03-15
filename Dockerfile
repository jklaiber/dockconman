FROM golang:1.15-alpine AS build_base

RUN apk add --no-cache git

WORKDIR /tmp/dockconman

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/dockconman .

FROM alpine:3.9 
RUN apk add ca-certificates
RUN apk add --update docker openrc
RUN rc-update add docker boot

COPY --from=build_base /tmp/dockconman/out/dockconman /app/dockconman
RUN ln -s /app/dockconman dockconman

CMD ["./dockconman"]