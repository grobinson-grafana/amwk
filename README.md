# amwk

`amwk` is a simple webhook server that receives alerts from Alertmanager. It
prints all alerts to stdout and tracks the number of times an alert has been
received.

## Installation

You can install `amwk` with `go install github.com/grobinson-grafana/amwk`,
or build it from source.

## Usage

```
Usage of amwk:
  -http-host string
    	The HTTP host (default "127.0.0.1")
  -http-port int
    	The HTTP port (default 8080)
```

### API

`amwk` has a simple HTTP API. You send it alerts via HTTP POST requests to `/`.
You can also get the number of times each alert has been received via an HTTP GET
request to `/fingerprints`.

```
{
  "839efbb46d0150df": {
    "2023-03-18T21:46:00Z": 1
  },
  "e2c03db4978aa696": {
    "2023-03-18T21:46:00Z": 1
  }
}
```
