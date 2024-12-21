FROM docker.io/golang:1.22.10-alpine3.21 AS build

RUN apk add make git curl

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .git .
RUN make build

###############################################################

FROM docker.io/alpine:3.21

LABEL org.opencontainers.image.source="https://github.com/rramiachraf/dumb"
LABEL org.opencontainers.image.url="https://github.com/rramiachraf/dumb"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.description="Private alternative front-end for Genius."

COPY --from=build /code/dumb .

EXPOSE 5555/tcp

CMD ["./dumb"]

