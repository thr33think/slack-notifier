# Intermediate Build Image 
FROM lushdigital/docker-golang-dep as builder
COPY ./ /go/src/slack-notifier/
WORKDIR /go/src/slack-notifier
RUN dep ensure -v && \
  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o slack-notifier ./*.go
RUN apk --update add ca-certificates

# Main Image
FROM scratch
COPY --from=builder /go/src/slack-notifier/slack-notifier /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 8080
ENTRYPOINT ["/slack-notifier"]