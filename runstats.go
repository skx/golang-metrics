package runstats

import (
	"fmt"
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
		fmt.Printf("HOST: %s\n", host)
		var err error
		g, err = graphite.NewGraphite(host, 2003)
		fmt.Printf("%v\n", err)
		if err == nil {
			go runCollector()
		} else {
			fmt.Printf("Error setting up metrics: %s\n", err.Error())
		}
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
