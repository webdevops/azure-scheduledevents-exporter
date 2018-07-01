package main

import (
	"time"
	"log"
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type AzureScheduledEventResponse struct {
	DocumentIncarnation int `json:"DocumentIncarnation"`
	Events []AzureScheduledEvent `json:"events"`
}

type AzureScheduledEvent struct {
	EventId string `json:"EventId"`
	EventType string `json:"EventType"`
	ResourceType string `json:"ResourceType"`
	Resources []string `json:"Resources"`
	EventStatus string `json:"EventStatus"`
	NotBefore string `json:"NotBefore"`
}

var (
	scheduledEventCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_scheduled_event_count",
			Help: "Azure ScheduledEvent count",
		},
		[]string{},
	)

	scheduledDocumentIncarnation = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_scheduled_document_incarnation",
			Help: "Azure ScheduledEvent document incarnation",
		},
		[]string{},
	)


	scheduledEvent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_scheduled_event",
			Help: "Azure ScheduledEvent",
		},
		[]string{"EventID", "EventType", "ResourceType", "EventStatus", "NotBefore"},
	)

	scheduledEventCountdown = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_scheduled_event_countdown_seconds",
			Help: "Azure ScheduledEvent countdown",
		},
		[]string{"EventID"},
	)

	timeFormatList = []string{
		time.RFC3339,
		time.RFC1123,
		time.RFC822Z,
		time.RFC850,
	}

	httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
)


func initMetrics() {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(scheduledEvent)
	prometheus.MustRegister(scheduledDocumentIncarnation)
	prometheus.MustRegister(scheduledEventCount)
	prometheus.MustRegister(scheduledEventCountdown)

	go func() {
		for {
			probeCollect()
			time.Sleep(time.Duration(opts.ScrapeTime) * time.Second)
		}
	}()
}

func startHttpServer() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func probeCollect() {
	scheduledEvents, err := fetchApiUrl()
	if err != nil {
		panic(err.Error())
	}

	if len(scheduledEvents.Events) >= 1 {
		for _, event := range scheduledEvents.Events {
			scheduledEvent.With(prometheus.Labels{"EventID": event.EventId, "EventType": event.EventType, "ResourceType": event.ResourceType, "EventStatus": event.EventStatus, "NotBefore": event.NotBefore}).Set(1)

			notBefore, err := parseTime(event.NotBefore)
			if err == nil {
				scheduledEventCountdown.With(prometheus.Labels{"EventID": event.EventId}).Set(float64(time.Until(notBefore).Seconds()))
			} else {
				Logger.Error(fmt.Sprintf("Unable to parse time of eventid \"%v\"", event.EventId), err)
			}
		}
	} else {
		scheduledEvent.With(prometheus.Labels{}).Set(0)
		scheduledEventCountdown.With(prometheus.Labels{}).Set(0)
	}

	scheduledDocumentIncarnation.With(prometheus.Labels{}).Set(float64(scheduledEvents.DocumentIncarnation))
	scheduledEventCount.With(prometheus.Labels{}).Set(float64(len(scheduledEvents.Events)))
}

func fetchApiUrl() (*AzureScheduledEventResponse, error) {
	ret := &AzureScheduledEventResponse{}

	req, err := http.NewRequest("GET", opts.ApiUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func parseTime(value string) (parsedTime time.Time, err error) {
	for _, format := range timeFormatList {
		parsedTime, err = time.Parse(format, value)
		if err == nil {
			break
		}
	}

	return
}
