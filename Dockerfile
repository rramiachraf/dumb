FROM golang:1.22.1-alpine3.19 as build

RUN apk add make git curl

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

###############################################################

FROM scratch

LABEL org.opencontainers.image.source="https://github.com/rramiachraf/dumb"
LABEL org.opencontainers.image.url="https://github.com/rramiachraf/dumb"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.description="Private alternative front-end for Genius."

COPY --from=build /code/dumb .
COPY --from=build /code/static static
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 5555/tcp

CMD ["./dumb"]

