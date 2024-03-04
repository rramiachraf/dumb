FROM golang:1.19.4-alpine3.17

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

EXPOSE 5555/tcp

CMD ["/code/dumb"]
