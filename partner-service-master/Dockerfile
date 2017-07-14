FROM golang:1.8.1-alpine

EXPOSE 80

COPY partner_service /go/bin/partner_service

CMD partner_service -grpcAddr 80