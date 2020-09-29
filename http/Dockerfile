FROM alpine:3.12

RUN apk add --no-cache python3=3.8.5-r0

ADD server_content /server_content

WORKDIR /server_content

EXPOSE 8080

ENTRYPOINT ["python3", "-m", "http.server"]