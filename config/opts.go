package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"
)

type (
	Opts struct {
		// logger
		Logger struct {
			Debug   bool `           long:"debug"        env:"DEBUG"    description:"debug mode"`
			Verbose bool `short:"v"  long:"verbose"      env:"VERBOSE"  description:"verbose mode"`
			LogJson bool `           long:"log.json"     env:"LOG_JSON" description:"Switch log output to json format"`
		}

		// general options
		ServerBind string        `long:"bind"                env:"SERVER_BIND"   description:"Server address"                default:":8080"`
		ScrapeTime time.Duration `long:"scrape-time"         env:"SCRAPE_TIME"   description:"Scrape time in seconds"        default:"1m"`

		// Api options
		ApiUrl            string        `long:"api-url"             env:"API_URL"       description:"Azure ScheduledEvents API URL" default:"http://169.254.169.254/metadata/scheduledevents?api-version=2017-11-01"`
		ApiTimeout        time.Duration `long:"api-timeout"         env:"API_TIMEOUT"   description:"Azure API timeout (seconds)"   default:"30s"`
		ApiErrorThreshold int           `long:"api-error-threshold" env:"API_ERROR_THRESHOLD"   description:"Azure API error threshold (after which app will panic)"   default:"0"`

		Notification            []string `long:"notification"                 env:"NOTIFICATION"              description:"Shoutrrr url for notifications (https://containrrr.github.io/shoutrrr/)" env-delim:" "  json:"-"`
		NotificationMsgTemplate string   `long:"notification.messagetemplate" env:"NOTIFICATION_MESSAGE_TEMPLATE"  description:"Notification template" default:"%v"`

		// metrics
		MetricsRequestStats bool `long:"metrics-requeststats" env:"METRICS_REQUESTSTATS" description:"Enable request stats metrics"`
	}
)

func (o *Opts) GetJson() []byte {
	jsonBytes, err := json.Marshal(o)
	if err != nil {
		log.Panic(err)
	}
	return jsonBytes
}
