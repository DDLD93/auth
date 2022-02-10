# build stage
FROM golang 
WORKDIR /src
ADD . /src
RUN cd /src && go get && go mod tidy && go build -o ./auth
EXPOSE 5000
ENTRYPOINT ./auth
