# golang-metrics

This is a simple library to automatically submit metrics to a remote
graphite/carbon server.

To use it you simply need to add the following import to your application:

    _ "github.com/skx/golang-metrics"

This works because the `init()` function will start a goroutine to collect
and submit metrics every 10 seconds.

The following example shows how little effort is required, the `import`
line literally does everything!

      package main

      import (
        "fmt"
        "time"

        _ "github.com/skx/golang-metrics"
      )

      func main() {
        for {
		   fmt.Printf("Alive\n")
		   time.Sleep(1 * time.Second)
        }
      }


Of course you need to tell the system _where_ to send your metrics.  So to
do that it reads the contents of the `METRICS` environmental variable.

Assuming you've saved the above example in `~/stats/main.go` you can start
it like so:

     $ METRICS=metrics.example.com go run ./main.go

The metrics submitted will be prefixed by the basename of your binary, and
 the prefix `go`.  i.e:

* `go.$(basename argv[0]).*`

As a concrete example an application binary called `overseer` would submit
metrics with names like these:

* `go.overseer.cpu.cgo_calls`
* `go.overseer.cpu.goroutines`
* `go.overseer.mem.alloc`
* `..`
