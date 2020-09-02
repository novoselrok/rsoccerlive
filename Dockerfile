FROM golang:1.15-buster as build
WORKDIR /app
COPY . .
RUN go mod download

FROM build as serve
RUN go build -o main ./src
CMD ["/app/main"]
