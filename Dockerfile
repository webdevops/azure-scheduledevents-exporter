FROM golang:1.10 as build

# golang deps
WORKDIR /tmp/app/
COPY ./src/glide.yaml /tmp/app/
COPY ./src/glide.lock /tmp/app/
RUN curl https://glide.sh/get | sh \
    && glide install

WORKDIR /go/src/azure-scheduledevents-exporter/src
COPY ./src /go/src/azure-scheduledevents-exporter/src
RUN mkdir /app/ \
    && cp -a /tmp/app/vendor ./vendor/ \
    && cp -a entrypoint.sh /app/ \
    && chmod 555 /app/entrypoint.sh \
    && go build -o /app/azure-scheduledevents-exporter

#############################################
# FINAL IMAGE
#############################################
FROM alpine
RUN apk add --no-cache \
        libc6-compat \
    	ca-certificates
COPY --from=build /app/ /app/
USER 1000
ENTRYPOINT ["/app/entrypoint.sh"]
