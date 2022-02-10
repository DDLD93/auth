# build stage
FROM ubuntu:18.04  
WORKDIR /app
COPY authServer ./
EXPOSE 3000
CMD ["./authServer"]
