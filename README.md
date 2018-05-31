# golang-metrics

This is a simple library to automatically submit metrics to a remote
graphite/carbon server.

To use it add the following to your application:

      import ( "github.com/skx/golang-metrics")

For example the following complete program will submit metrics every
few seconds:

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


To specify your server simply export the environmental variable `METRICS`:

     $ METRICS=metrics.example.com go run ./main.go
