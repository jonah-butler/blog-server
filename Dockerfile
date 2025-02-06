FROM golang:alpine

WORKDIR /blog_api
COPY . .

RUN go build -o ./bin/api ./cmd/api

CMD ["/blog_api/bin/api"]
EXPOSE 8080