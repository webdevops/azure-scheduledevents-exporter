FROM golang:1.12 as build

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
    && CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /app/azure-scheduledevents-exporter

#############################################
# FINAL IMAGE
#############################################
FROM scratch
COPY --from=build /app/azure-scheduledevents-exporter /
USER 1000
ENTRYPOINT ["/azure-scheduledevents-exporter"]
