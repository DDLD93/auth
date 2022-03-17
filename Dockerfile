# build stage
FROM alpine:3.15
RUN apk add --no-cache ca-certificates
RUN apt install golang -y
WORKDIR /app
COPY . . 
RUN go get
RUN go build -o auth .
EXPOSE 9000
CMD ["./auth"]
