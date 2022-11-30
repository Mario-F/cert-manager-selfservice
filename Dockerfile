# syntax=docker/dockerfile:1

##
## Build Web
##
FROM node:18 AS build-web

WORKDIR /app

COPY ./web/package.json ./package.json
COPY ./web/yarn.lock ./yarn.lock
RUN yarn
COPY ./web ./
RUN yarn build

##
## Build Main
##
FROM golang:1.18-buster AS build-main

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
COPY --from=build-web /app/dist /app/web/dist

RUN go build -o /cert-manager-selfservice

##
## Deploy
##
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build-main /cert-manager-selfservice /cert-manager-selfservice

USER nonroot:nonroot

ENTRYPOINT ["/cert-manager-selfservice"]
