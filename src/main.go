package main

import (
	"os"
	"fmt"
	"time"
	"strings"
	"net/url"
	"github.com/jessevdk/go-flags"
)

const (
	Author  = "webdevops.io"
	Version = "1.2.0"
)

var (
	argparser   *flags.Parser
	args        []string
	Logger      *DaemonLogger
	ErrorLogger *DaemonLogger
)

var opts struct {
	// general options
	ServerBind  string `       long:"bind"                env:"SERVER_BIND"   description:"Server address"                default:":8080"`
	ScrapeTime  time.Duration `long:"scrape-time"         env:"SCRAPE_TIME"   description:"Scrape time in seconds"        default:"1m"`
	Verbose []bool `           long:"verbose" short:"v"   env:"VERBOSE"       description:"Verbose mode"`

	// Api options
	ApiUrl      string `       long:"api-url"             env:"API_URL"       description:"Azure ScheduledEvents API URL" default:"http://169.254.169.254/metadata/scheduledevents?api-version=2017-08-01"`
	ApiTimeout  time.Duration `long:"api-timeout"         env:"API_TIMEOUT"   description:"Azure API timeout (seconds)"   default:"30s"`
	ApiErrorThreshold int `    long:"api-error-threshold" env:"API_ERROR_THRESHOLD"   description:"Azure API error threshold (after which app will panic)"   default:"5"`
}

func main() {
	initArgparser()

	// Init logger
	Logger = CreateDaemonLogger(0)
	ErrorLogger = CreateDaemonErrorLogger(0)

	// set verbosity
	Verbose = len(opts.Verbose) >= 1

	Logger.Messsage("Init Azure ScheduledEvents exporter v%s (written by %v)", Version, Author)

	Logger.Messsage("Starting metrics collection")
	Logger.Messsage("  API URL: %v", opts.ApiUrl)
	Logger.Messsage("  API timeout: %v", opts.ApiTimeout)
	Logger.Messsage("  scape time: %v", opts.ScrapeTime)
	Logger.Messsage("  error threshold: %v", opts.ApiErrorThreshold)
	setupMetricsCollection()
	startMetricsCollection()

	Logger.Messsage("Starting http server on %s", opts.ServerBind)
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

	// --api-url
	apiUrl, err := url.Parse(opts.ApiUrl)
	if err != nil {
		fmt.Println(err)
		fmt.Println()
		argparser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	// validate --api-url scheme
	switch (strings.ToLower(apiUrl.Scheme)) {
	case "http":
		break;
	case "https":
		break;
	default:
		fmt.Println("ApiURL scheme not allowed (must be http or https)")
		fmt.Println()
		argparser.WriteHelp(os.Stdout)
		os.Exit(1)
	}
}
