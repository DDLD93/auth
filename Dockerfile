# build stage
FROM ubuntu:jammy
WORKDIR /app
COPY auth .
EXPOSE 5000
CMD ["./auth"]
