package runstats

import (
	"fmt"
	"net"
	"os"
	"path"
	"strconv"

	graphite "github.com/marpaia/graphite-golang"
	"github.com/skx/golang-metrics/collector"
)

var g *graphite.Graphite

func init() {
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

		g, err = graphite.NewGraphite(h, port)
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
			g.SimpleSend(prefix+key, strconv.FormatUint(val, 10))
		}
	}
	c := collector.New(gaugeFunc)
	c.Run()
}
