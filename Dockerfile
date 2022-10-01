# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

WORKDIR /app

COPY . .

RUN go build -o /tz

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

ARG MONGODB_URI_ARG
ENV MONGODB_URI=$MONGODB_URI_ARG

COPY --from=build /tz /tz

EXPOSE 8080

ENTRYPOINT ["/tz"]
