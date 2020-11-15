FROM golang:1.14 as build

WORKDIR /go/src/github.com/webdevops/azure-scheduledevents-exporter

# Get deps (cached)
COPY ./go.mod /go/src/github.com/webdevops/azure-scheduledevents-exporter
COPY ./go.sum /go/src/github.com/webdevops/azure-scheduledevents-exporter
RUN go mod download

# Compile
COPY ./ /go/src/github.com/webdevops/azure-scheduledevents-exporter
RUN make test
RUN make lint
RUN make build
RUN ./azure-scheduledevents-exporter --help

#############################################
# FINAL IMAGE
#############################################
FROM gcr.io/distroless/static
COPY --from=build /go/src/github.com/webdevops/azure-scheduledevents-exporter/azure-scheduledevents-exporter /
USER 1000
ENTRYPOINT ["/azure-scheduledevents-exporter"]
