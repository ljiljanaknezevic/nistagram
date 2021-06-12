FROM golang:1.15-alpine3.13
WORKDIR app
RUN apk add --no-cache tzdata
ENV TZ=Europe/Belgrade
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apk update \
    && apk add socat \
    && apk add build-base
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o user-service
EXPOSE 8081
CMD ["./user-service"]