# Base image
FROM alpine:latest

LABEL  maintainer="Jerome Doucet" version="0.1"

WORKDIR /app
COPY dahu-server /app/

EXPOSE 80
ENTRYPOINT /app/dahu-server