# golang-metrics

This is a simple library to automatically submit metrics to a remote
graphite/carbon server.

Using this library is as simple as adding the following import to your
application:

    _ "github.com/skx/golang-metrics"

Importing the library is all that you need to do, because the `init()`
function will launch a goroutine to collect and submit metrics.


## Sample Application

The following example shows how little effort is required, as documented
above the `import` carries out the setup and launches the metric-collection
automatically.

      package main

      import (
        "fmt"
        "time"

        _ "github.com/skx/golang-metrics"
      )

      func main() {
        for {
		   fmt.Printf("This process is alive, and submitting metrics!\n")
		   time.Sleep(1 * time.Second)
        }
      }


## Configuration

There is only one thing you need to configure, which is the address of
the host to submit your metrics to.   The environmental variable
`METRICS` will be used for that purpose.

* If there is no `METRICS` variable defined then no metrics will be collected.
   * Because collecting things is pointless without somewhere to send them.

Assuming you've saved the example included earlier in `~/main.go` you can
start it like so:

     $ METRICS=metrics.example.com go run ./main.go

The metrics will be submitted to the host `metrics.example.com`, and each
metric will prefixed by the basename of your binary, and the prefix `go`.  i.e:

* `go.$(basename argv[0]).*`

As a concrete example an application binary called `overseer` would submit
metrics with names like these:

* `go.overseer.cpu.cgo_calls`
* `go.overseer.cpu.goroutines`
* `go.overseer.mem.alloc`
* `..`

If you need to submit to a non-standard port you can include that in your `$METRICS` setup:

     $ METRICS=metrics.example.com:2233 ./application


## Advanced Configuration

Although we've discussed the only mandatory setting, `METRICS`, there
are a couple more environmental variables you might wish to set.

| Setting            | Purpose                                                  |
| -------------------|----------------------------------------------------------|
| `METRICS_VERBOSE`  | If this is non-empty metrics will be echoed to STDOUT    |
| `METRICS_DELAY`    | If this is set to an integer then it will be used to     |
|                    | control how often metrics are sent.  The default is      |
|                    | `10` meaning metrics will be submitted every ten seconds |
| `METRICS_PROTOCOL` | If this is set to `tcp` then TCP will be used, instead  |
|                    | of the default of UDP updates                           |


## Systemd

If you're launching your application under the control of systemd you can
configure the destination in your `.service` file by adding an `Environment` setting.  For example the following service:

     [Unit]
     Description=My service..

     [Service]
     WorkingDirectory=/srv/blah
     User=blah
     Environment="METRICS=metrics.example.com:2003"
     ExecStart=/srv/blah/bin/application arugments
     Restart=always
     StartLimitInterval=2
     StartLimitBurst=20
     PrivateTmp=yes
     RestrictAddressFamilies=AF_INET AF_INET6 AF_UNIX


## Licensing / Credits

The metric-collection code was copied from the library written by Brian Hatfield, which is available here:

* https://github.com/bmhatfield/go-runtime-metrics

That code was licensed under the MIT license, which has been included in this
repository.
