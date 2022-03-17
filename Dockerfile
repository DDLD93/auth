# build stage
FROM alpine:3.15
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY auth .
EXPOSE 9000
CMD ["./auth"]
