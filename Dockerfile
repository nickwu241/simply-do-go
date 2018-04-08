FROM alpine:latest

WORKDIR /app
COPY public public
COPY server entrypoint.sh ./

ENTRYPOINT ["/app/entrypoint.sh"]
CMD [""]

EXPOSE 8080
