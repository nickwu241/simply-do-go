FROM alpine:latest

WORKDIR /app
RUN apk update && apk add ca-certificates

COPY public public
COPY server entrypoint.sh ./

ENTRYPOINT ["/app/entrypoint.sh", "/app/server"]
CMD [""]

EXPOSE 8080
