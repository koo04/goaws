FROM golang:1.15-alpine as builder
WORKDIR /goaws
COPY . .
RUN go build -o goaws ./app/cmd/

FROM alpine
EXPOSE 4100
COPY --from=builder /goaws/goaws /
COPY --from=builder /goaws/app/conf/goaws.yaml /conf/
ENTRYPOINT ["/goaws"]
