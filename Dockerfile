# build stage
FROM golang:1.13-alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
ADD . /src
RUN cd /src && go build -o goapp

# final stage
FROM alpine
WORKDIR /app
RUN apk --no-cache add mysql-client tzdata
COPY --from=build-env /src/goapp /app/
ENTRYPOINT ./goapp