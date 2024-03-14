FROM golang:1.22.1-alpine3.19 as build

RUN apk add make git curl

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM scratch

COPY --from=build /code/dumb .
COPY --from=build /code/static static
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 5555/tcp

CMD ["./dumb"]

