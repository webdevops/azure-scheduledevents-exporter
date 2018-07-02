package main

import (
	"os"
	"fmt"
	"github.com/jessevdk/go-flags"
	"azure-scheduledevents-exporter/src/logger"
	"net/url"
)

const (
	Author  = "webdevops.io"
	Version = "0.1.0"
)

var (
	argparser   *flags.Parser
	args        []string
	Logger      *logger.DaemonLogger
	ErrorLogger *logger.DaemonLogger
)

var opts struct {
	ApiUrl      string `long:"api-url"       env:"APIURL"        description:"Azure ScheduledEvents API URL" default:"http://169.254.169.254/metadata/scheduledevents?api-version=2017-08-01"`
	apiUrl      *url.URL
	ScrapeTime  int `   long:"scrape-time"   env:"SCRAPE_TIME"   description:"Scrape time in seconds"        default:"360"`
}

func main() {
	initArgparser()

	// Init logger
	Logger = logger.CreateDaemonLogger(0)
	ErrorLogger = logger.CreateDaemonErrorLogger(0)

	Logger.Messsage("init")

	u, err := url.Parse(opts.ApiUrl)
	if err != nil {
		panic(err)
	}
	opts.apiUrl = u

	Logger.Messsage("starting metrics collection")
	initMetrics()

	Logger.Messsage("starting http server")
	startHttpServer()
}

func initArgparser() {
	argparser = flags.NewParser(&opts, flags.Default)
	_, err := argparser.Parse()

	// check if there is an parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println(err)
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}

}
