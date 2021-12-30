# syntax=docker/dockerfile:1

##
## Build backend
##
FROM golang:1.16-alpine AS backend-build

WORKDIR /app

COPY ./backend/go.mod ./
COPY ./backend/go.sum ./
RUN go mod download

COPY ./backend ./

RUN ls

RUN go build -o /backend-server cmd/main.go

##
## Build UI
##
FROM node:14-alpine3.13 AS ui-build

WORKDIR /app

COPY ./ui/package.json ./
COPY ./ui/yarn.lock ./
COPY ./ui/public ./public
COPY ./ui/src ./src

RUN ls

RUN yarn
RUN yarn build

##
## Deploy
##
FROM alpine:3.15.0

WORKDIR /app

COPY --from=backend-build /backend-server /app/backend-server
COPY --from=ui-build /build-server /app/ui/build

EXPOSE 4000

ENTRYPOINT ["/app/backend-server"]
