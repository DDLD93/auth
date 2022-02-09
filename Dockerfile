# build stage
FROM golang AS build-env
ADD . /src
RUN go get
RUN cd /src && go build -o authServer

# final stage
FROM linux:alpine
WORKDIR /app
COPY --from=build-env /src/authServer /app/
EXPOSE 5000
ENTRYPOINT ./authServer
