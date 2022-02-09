# build stage
FROM golang AS build-env
WORKDIR /src
ADD . /src
RUN cd /src && go get && go mod tidy && go build -o ./auth

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/auth /app/
ENTRYPOINT ./auth
