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

RUN yarn
RUN yarn build

##
## Deploy
##
FROM alpine:3.15.0

WORKDIR /app

COPY --from=backend-build /backend-server /app/backend/backend-server
COPY --from=ui-build /app/build /app/ui/build

EXPOSE 4000

ENV VES_API_KEY ""
ENV MOT_HISTORY_API_KEY ""
ENV JWT_SIGNING_SECRET ""
ENV MONGO_CONNECTION_STRING ""

CMD /app/backend/backend-server -vesapi-key=${VES_API_KEY} -mothistoryapi-key=${MOT_HISTORY_API_KEY} -jwt-signing-secret=${JWT_SIGNING_SECRET} -mongo-connection-string=${MONGO_CONNECTION_STRING}
