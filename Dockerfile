FROM ubuntu:jammy 
WORKDIR /root/
COPY authServer ./
EXPOSE 5000
CMD ["./authServer"]
