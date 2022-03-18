# build stage
FROM ubuntu:jammy
WORKDIR /app
COPY auth .
EXPOSE 9000
CMD ["./auth"]
