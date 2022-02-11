# build stage
FROM ubuntu 
WORKDIR /app
COPY authServer .
EXPOSE 3000
CMD ["./authServer"]
