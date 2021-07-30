FROM golang:1.12-alpine3.10 AS build
COPY . /usr/local/src
WORKDIR /usr/local/src
RUN apk --no-cache add build-base && make

FROM alpine:3.10
WORKDIR /syncd
COPY --from=build /usr/local/src/output /syncd
EXPOSE 8878
CMD [ "/syncd/bin/syncd" ]
