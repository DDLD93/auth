# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
ADD . /src
RUN go get
RUN cd /src && go build -o authServer

# final stage
FROM linux:alpine
WORKDIR /app
COPY --from=build-env /src/authServer /app/
EXPOSE 5000
ENTRYPOINT ./authServer
