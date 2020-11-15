Azure ScheduledEvents Exporter
==============================

[![license](https://img.shields.io/github/license/webdevops/azure-scheduledevents-exporter.svg)](https://github.com/webdevops/azure-scheduledevents-exporter/blob/master/LICENSE)
[![DockerHub](https://img.shields.io/badge/DockerHub-webdevops%2Fazure--scheduledevents--exporter-blue)](https://hub.docker.com/r/webdevops/azure-scheduledevents-exporter/)
[![Quay.io](https://img.shields.io/badge/Quay.io-webdevops%2Fazure--scheduledevents--exporter-blue)](https://quay.io/repository/webdevops/azure-scheduledevents-exporter)

Prometheus exporter for [Azure ScheduledEvents](https://docs.microsoft.com/en-us/azure/virtual-machines/linux/scheduled-events) (planned VM maintenance) from the Azure API.

It fetches informations from `http://169.254.169.254/metadata/scheduledevents?api-version=2017-08-01`
and exports the parsed information as metric to Prometheus.


Hint: Interrested in automatic node draining for Kubernetes for Azure ScheduledEvents? see [azure-scheduledevents-manager](https://github.com/webdevops/azure-scheduledevents-manager)

Configuration
-------------

Normally no configuration is needed but can be customized using environment variables.

```
Usage:
  azure-scheduledevents-exporter [OPTIONS]

Application Options:
      --bind=                 Server address (default: :8080) [$SERVER_BIND]
      --scrape-time=          Scrape time in seconds (default: 1m)
                              [$SCRAPE_TIME]
  -v, --verbose               Verbose mode [$VERBOSE]
      --api-url=              Azure ScheduledEvents API URL (default:
                              http://169.254.169.254/metadata/scheduledevents?api-version=2017-11-01) [$API_URL]
      --api-timeout=          Azure API timeout (seconds) (default: 30s)
                              [$API_TIMEOUT]
      --api-error-threshold=  Azure API error threshold (after which app will
                              panic) (default: 0) [$API_ERROR_THRESHOLD]
      --metrics-requeststats  Enable request stats metrics
                              [$METRICS_REQUESTSTATS]

Help Options:
  -h, --help                  Show this help message
```

Metrics
-------

| Metric                                      | Description                                                                           |
|---------------------------------------------|---------------------------------------------------------------------------------------|
| `azure_scheduledevent_document_incarnation` | Document incarnation number (version)                                                 |
| `azure_scheduledevent_event`                | Fetched events from API                                                               |
| `azure_scheduledevent_request`              | Request histogram (count and request duration; disabled by default)                   |
| `azure_scheduledevent_request_error`        | Counter for failed requests                                                           |


Kubernetes Usage
----------------

```
---
apiVersion: app/v1
kind: DaemonSet
metadata:
  name: azure-scheduledevents
spec:
  updateStrategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 100%
  selector:
    matchLabels:
      app: azure-scheduledevents
  template:
    metadata:
      labels:
        app: azure-scheduledevents
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/path: /metrics
        prometheus.io/port: "8080"
    spec:
      terminationGracePeriodSeconds: 15
      nodeSelector:
        beta.kubernetes.io/os: linux
      tolerations:
      - effect: NoSchedule
        operator: Exists
      containers:
      - name: azure-scheduledevents
        image: webdevops/azure-scheduledevents-exporter
        securityContext:
          readOnlyRootFilesystem: true
          runAsNonRoot: true
          capabilities:
            drop: ['ALL']
        ports:
        - containerPort: 8080
          name: metrics
          protocol: TCP
        resources:
          limits:
            cpu: 100m
            memory: 50Mi
          requests:
            cpu: 1m
            memory: 50Mi
```
