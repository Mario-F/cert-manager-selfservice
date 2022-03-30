# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -o /cert-manager-selfservice

##
## Deploy
##
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /cert-manager-selfservice /cert-manager-selfservice

USER nonroot:nonroot

ENTRYPOINT ["/cert-manager-selfservice"]
