FROM alpine:3.4

RUN apk add --no-cache ca-certificates

ADD main main
RUN chmod +x main

CMD ["./main"]