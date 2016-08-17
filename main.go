package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"strings"
	"time"
)

const defaultLogLevel = log.InfoLevel

func initLogger() {
	log.SetLevel(getLogLevel())
}

func getLogLevel() log.Level {
	lvl := strings.ToLower(os.Getenv("LOG_LEVEL"))
	level, err := log.ParseLevel(lvl)
	if err != nil {
		level = defaultLogLevel
	}
	return level
}

func main() {
	initLogger()

	var (
		listenAddress = flag.String("web.listen-address", ":9108", "Address to listen on for web interface and telemetry.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		esURI         = flag.String("es.uri", "http://localhost:9200", "HTTP API address of an Elasticsearch node.")
		esTimeout     = flag.Duration("es.timeout", 5*time.Second, "Timeout for trying to get stats from Elasticsearch.")
		esAllNodes    = flag.Bool("es.all", false, "Export stats for all nodes in the cluster.")
	)
	flag.Parse()

	exporter := NewExporter(*esURI, *esTimeout, *esAllNodes)
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Elasticsearch Exporter</title></head>
             <body>
             <h1>Elasticsearch Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	log.WithFields(log.Fields{
		"Version":       Version,
		"Revision":      Revision,
		"Branch":        Branch,
		"BuildDate":     BuildDate,
		"ListenAddress": *listenAddress,
	}).Info("Starting Server:")

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
