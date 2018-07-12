FROM golang:1.8-alpine AS compile
COPY hello-docker-web.go /go
RUN go build hello-docker-web.go

FROM alpine:latest
COPY --from=compile /go/hello-docker-web /
USER nobody:nobody
EXPOSE 8080
ENTRYPOINT ["/hello-docker-web"]
