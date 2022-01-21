FROM ubuntu:18.04  
WORKDIR /root/
COPY authServer ./
EXPOSE 5000
CMD ["./authServer"]
