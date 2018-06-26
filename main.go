package main

import (
	"net/http"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

var (
	listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9444").String()
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose Prometheus metrics.").Default("/metrics").String()
	emqURL        = kingpin.Flag("emq.uri", "HTTP API address of the EMQ node.").Default("http://127.0.0.1:8080").URL()
	emqUsername   = kingpin.Flag("emq.username", "EMQ username.").Default("admin").String()
	emqPassword   = kingpin.Flag("emq.password", "EMQ password.").Default("public").String()
	emqNodeName   = kingpin.Flag("emq.node", "Node name of the emq node to scrape.").Default("emq@127.0.0.1").String()
)

func init() {
	prometheus.MustRegister(version.NewCollector("emq_exporter"))
}

func main() {
	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("emq_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Infoln("Starting emq_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	httpClient := &http.Client{}
	nodeName := *emqNodeName
	username := *emqUsername
	password := *emqPassword
	prometheus.MustRegister(NewEMQCollector(httpClient, emqURL, nodeName, username, password))

	http.Handle(*metricsPath, promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
    <head><title>EMQ Exporter</title></head>
    <body>
    <h1>EMQ Exporter</h1>
    <p><a href="` + *metricsPath + `">Metrics</a></p>
    </body>
    </html>`))
	})

	log.Infoln("Listening on", *listenAddress)
	http.ListenAndServe(*listenAddress, nil)
}
