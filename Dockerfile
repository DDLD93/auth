# build stage
FROM alpine 
WORKDIR /app
COPY authServer .
EXPOSE 3000
CMD ["./authServer"]
