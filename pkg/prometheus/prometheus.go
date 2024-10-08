package prom

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type Config struct {
	Port int
	Path string
}

var Conf Config

func (C *Config) ExposeMetrics() {

	go func(port int, path string) {
		log.Printf("Prometheus metrics @ 0.0.0.0:%d%s", port, path)
		http.Handle(path, promhttp.Handler())

		listenAddress := fmt.Sprintf(":%d", port)
		http.ListenAndServe(listenAddress, nil)
	}(C.Port, C.Path)
}
