# build stage
FROM unbuntu:jammy
WORKDIR /app
COPY auth .
EXPOSE 9000
CMD ["./auth"]
