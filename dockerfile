FROM alpine
COPY ./example /app/example
WORKDIR /app
EXPOSE 8080

CMD ./example