# build stage
FROM golang AS build-env
WORKDIR /src
ADD . /src
RUN cd /src && go get && go mod tidy && go build -o ./authServer

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/authServer /app/
EXPOSE 5000
ENTRYPOINT ./authServer
