# Golang Prometheus

Monitoring a Golang project with Prometheus

## Init project

Init from the root directory project using the repository [https://github.com/cartapas/golang_prometheus](https://github.com/cartapas/golang_prometheus)

```shell
$ go mod init github.com/cartapas/golang_prometheus                                  1 ↵
    go: creating new go.mod: module github.com/cartapas/golang_prometheus
```

The `go.mod` file is created

```mod
module github.com/cartapas/golang_prometheus

go 1.24.4
```

## Installing dependencies

```shell
$ go get github.com/prometheus/client_golang/prometheus
    go get github.com/prometheus/client_golang/prometheus                                ✔
    go: downloading golang.org/x/sys v0.30.0
    go: added github.com/beorn7/perks v1.0.1
    go: added github.com/cespare/xxhash/v2 v2.3.0
    go: added github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822
    go: added github.com/prometheus/client_golang v1.22.0
    go: added github.com/prometheus/client_model v0.6.1
    go: added github.com/prometheus/common v0.62.0
    go: added github.com/prometheus/procfs v0.15.1
    go: added golang.org/x/sys v0.30.0
    go: added google.golang.org/protobuf v1.36.5

$ go get github.com/prometheus/client_golang/prometheus/promauto
$ go get github.com/prometheus/client_golang/prometheus/promhttp
```

The requirements are appended to the `go.mod` file and the `go.sum` is created with cheksum values.

```mod
require (
  github.com/beorn7/perks v1.0.1 // indirect
  github.com/cespare/xxhash/v2 v2.3.0 // indirect
  github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
  github.com/prometheus/client_golang v1.22.0 // indirect
  github.com/prometheus/client_model v0.6.1 // indirect
  github.com/prometheus/common v0.62.0 // indirect
  github.com/prometheus/procfs v0.15.1 // indirect
  golang.org/x/sys v0.30.0 // indirect
  google.golang.org/protobuf v1.36.5 // indirect
)
```
