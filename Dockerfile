FROM golang:1.20.4-alpine3.18 AS build-stage

WORKDIR /home/app

COPY ./ /home/app

RUN mkdir -p /home/build
RUN go mod download
RUN go build -v -o /home/build/api ./cmd/api


FROM alpine

WORKDIR /home/build

COPY --from=build-stage /home/build/api /home/build/api
COPY --from=build-stage /home/app/views /home/build/views
CMD ["/home/build/api"]