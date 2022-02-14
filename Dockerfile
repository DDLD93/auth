# build stage
FROM ubuntu:jammy
WORKDIR /app
COPY authServer .
EXPOSE 5000
CMD ["./authServer"]
