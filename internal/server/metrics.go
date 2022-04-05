package server

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func StartMetricsExporter(port int) {
	endpointServer := fmt.Sprintf(":%d", port)
	endpointHandler := "/metrics"

	log.Infof("Starting metrics exporter on %s%s", endpointServer, endpointHandler)

	http.Handle(endpointHandler, promhttp.Handler())
	err := http.ListenAndServe(endpointServer, nil)
	if err != nil {
		log.Error("Failed to start metrics exporter endpoint")
	}
}
