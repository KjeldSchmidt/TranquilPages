###############
# Build Stage #
###############
FROM golang:1.24-alpine AS build
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o tranquil-pages


#################
# Runtime Stage #
#################
FROM alpine:latest

RUN apk add --no-cache sqlite-libs

WORKDIR /app

COPY --from=build /app/tranquil-pages .

EXPOSE 8080

CMD ["./tranquil-pages"]
