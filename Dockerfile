# build stage
FROM ubuntu 
WORKDIR /app
COPY authServer .
EXPOSE 5000
CMD ["./authServer"]
