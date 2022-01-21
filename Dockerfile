FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY authServer ./
EXPOSE 5000
CMD ["./authServer"] 	//fmt.Println(token)
