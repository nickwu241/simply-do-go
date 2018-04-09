FROM alpine:latest

WORKDIR /app
RUN apk update && apk add ca-certificates

COPY public public
COPY simply-do .
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
ENTRYPOINT ["entrypoint.sh"]
CMD ["/app/simply-do"]

EXPOSE 8080
