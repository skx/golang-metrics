package runstats

import (
	"fmt"
	"net"
	"os"
	"path"
	"strconv"
	"time"

	graphite "github.com/marpaia/graphite-golang"
	"github.com/skx/golang-metrics/collector"
)

// Handle to the graphite object
var g *graphite.Graphite

// Should we be verbose?
var verbose bool

// The prefix of all our metrics
var prefix string

func init() {

	// Should we be verbose?
	v := os.Getenv("METRICS_VERBOSE")
	if v != "" {
		verbose = true
	}

	// Do we have a custom prefix?
	prefix = os.Getenv("METRICS_PREFIX")
	if prefix == "" {
		// If not use the default.
		prefix = "go." + path.Base(os.Args[0]) + "."
	}

	// Get the destination to send the metrics to.
	host := os.Getenv("METRICS")
	if host != "" {
		var err error

		// Split the into Host + Port
		h, p, err := net.SplitHostPort(host)
		if err != nil {
			// If that failed we assume the port was missing
			h = host
			p = "2003"
		}

		// port should be an integer
		port, err := strconv.Atoi(p)
		if err != nil {
			fmt.Printf("Error converting port %s to integer: %s\n", p, err.Error())
			return
		}

		//
		// We default to UDP, but if `METRICS_PROTOCOL` is set
		// to `tcp` we'll use that instead.
		//
		protocol := os.Getenv("METRICS_PROTOCOL")
		if protocol == "" {
			protocol = "udp"
		}

		//
		// Create the graphite-object.
		//
		g, err = graphite.GraphiteFactory(protocol, h, port, "")
		if err != nil {
			fmt.Printf("Error setting up metrics - skipping - %s\n", err.Error())
			return
		}

		go runCollector()
	}
}

func runCollector() {

	gaugeFunc := func(key string, val uint64) {
		if g != nil {
			if verbose {
				fmt.Printf("%s %d\n", prefix+key, val)
			}
			g.SimpleSend(prefix+key, strconv.FormatUint(val, 10))
		}
	}
	c := collector.New(gaugeFunc)

	//
	// A duration might be specified to alter our default
	// of ten seconds.
	//
	if os.Getenv("METRICS_DELAY") != "" {
		secs, err := strconv.Atoi(os.Getenv("METRICS_DELAY"))
		if err == nil {
			c.PauseDur = time.Second * time.Duration(secs)
		} else {
			fmt.Printf("Error processing $METRICS_DELAY: %s\n", err.Error())
		}
	}
	c.Run()
}
