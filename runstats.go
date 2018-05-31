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

func init() {

	// Should we be verbose?
	v := os.Getenv("METRICS_VERBOSE")
	if v != "" {
		verbose = true
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
			fmt.Printf("Error convering port %s to integer: %s\n", p, err.Error())
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
	prefix := "go." + path.Base(os.Args[0]) + "."

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
