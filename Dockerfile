FROM alpine:3.8

ADD main /

EXPOSE 8080

CMD ["./main"]