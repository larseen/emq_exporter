package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	namespace     = "emq"
	defaultLabels = []string{"node", "otp_release", "version"}
	validID       = regexp.MustCompile(`\d{1,}[.]\d{1,}|\d{1,}`)
)

type metric struct {
	Type  prometheus.ValueType
	Desc  *prometheus.Desc
	Value func(values combinedResponse) float64
}

// Collector is the struct for the EMQ Collector
type Collector struct {
	client   *http.Client
	url      **url.URL
	node     string
	password string
	username string

	up                prometheus.Gauge
	totalScrapes      prometheus.Counter
	jsonParseFailures prometheus.Counter
	metrics           []*metric
}

//NewEMQCollector initializes every descriptor and returns a pointer to the collector
func NewEMQCollector(client *http.Client, url **url.URL, node string, username string, password string) *Collector {
	return &Collector{
		client:   client,
		url:      url,
		node:     node,
		username: username,
		password: password,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, "node", "up"),
			Help: "Was the last scrape of the EMQ node successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, "node", "total_scrapes"),
			Help: "Current total scrapes.",
		}),
		jsonParseFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, "node", "json_parse_failures"),
			Help: "Number of errors while parsing JSON.",
		}),
		metrics: []*metric{
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "cluster", "size"),
					"The total number of EMQ nodes in your cluster.",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.ClusterSize)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "node", "process_used"),
					"The amount of processes used by the EMQ node.",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.nodes.Result.ProcessesUsed)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "node", "process_available"),
					"The amount of processes available to the EMQ node.",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.nodes.Result.ProcessesAvailable)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "node", "max_fds"),
					"The amount of file descriptors available to the EMQ node.",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.nodes.Result.MaxFds)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "node", "memory_total"),
					"The max amount of memory used to the EMQ node.",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					str := validID.FindAllString(values.nodes.Result.MemoryTotal, -1)
					i, err := strconv.ParseFloat(str[0], 64)
					if err != nil {
						log.Error("error converting string into number")
					}
					return float64(i * 1000000)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "node", "memory_used"),
					"The amount of memory being used to the EMQ node.",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					str := validID.FindAllString(values.nodes.Result.MemoryUsed, -1)
					i, err := strconv.ParseFloat(str[0], 64)
					if err != nil {
						log.Error("error converting string into number")
					}
					return float64(i * 1000000)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_disconnected"),
					"The amount of packets disconnected",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsDisconnect)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "messages_qos2_received"),
					"The amount of packets QOS2 messages received",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.MessagesQos2Received)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_suback"),
					"The amount of packets suback",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsSuback)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_pubcomp_received"),
					"The amount of packets pubcomp received",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsPubcompReceived)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_unsuback"),
					"The amount of packets unsuback",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsUnsuback)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_pingresp"),
					"The amount of packets pingresp",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsPingresp)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_pingreq"),
					"The amount of packets pingreq",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsPingreq)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_pubrel_missed"),
					"The amount of packets pubrel missed",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsPubrelMissed)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_sent"),
					"The amount of packets sent",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsSent)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "messages_qos2_sent"),
					"The amount of QOS2 messages sent",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.MessagesQos2Sent)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_pubrec_missed"),
					"The amount of packets pubrec missed",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsPubrecMissed)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_unsubscribe"),
					"The amount of packets disconnected",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsUnsubscribe)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "bytes_received"),
					"The amount of bytes received",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.BytesReceived)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_connack"),
					"The amount of packets connack",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsConnack)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "messages_received"),
					"The amount of messages received",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.MessagesReceived)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "messages_dropped"),
					"The amount of messages dropped",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.MessagesDropped)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_pubrec_sent"),
					"The amount of packets pubrec sent",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsPubrecSent)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "messages_retained"),
					"The amount of messages retained",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.MessagesRetained)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_publish_received"),
					"The amount of packets publish received",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsPublishReceived)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_pubcomp_sent"),
					"The amount of packets pubcomp sent",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsPubcompSent)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_connect"),
					"The amount of packets connect",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsConnect)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_puback_received"),
					"The amount of packets puback received",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsPubackReceived)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "messages_sent"),
					"The amount of messages sent",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.MessagesSent)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_publish_sent"),
					"The amount of packets publish sent",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsPublishSent)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "bytes_sent"),
					"The amount of bytes sent",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.BytesSent)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_puback_sent"),
					"The amount of packets puback sent",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsPubackSent)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "messages_qos2_dropped"),
					"The amount of QOS2 messages dropped",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.MessagesQos2Dropped)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_pubrel_sent"),
					"The amount of packets pubrel sent",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsPubrelSent)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "messages_qos1_sent"),
					"The amount of QOS1 messages sent",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.MessagesQos1Sent)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_pubrel_received"),
					"The amount of packets pubrel received",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsPubrelReceived)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "messages_qos1_received"),
					"The amount of QOS1 messages received",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.MessagesQos1Received)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "messages_qos0_sent"),
					"The amount of QOS0 messages sent",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.MessagesQos0Sent)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_received"),
					"The amount of packets received",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsReceived)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_pubrec_received"),
					"The amount of packets pubrec received",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsPubrecReceived)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_pubcomp_missed"),
					"The amount of packets pubcomp missed",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsPubcompMissed)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "metric", "packets_puback_missed"),
					"The amount of packets puback missed",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.metrics.Result.PacketsPubackMissed)
				},
			},

			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "stats", "clients"),
					"The amount of clients using in the EMQ node.",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.stats.Result.ClientsCount)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "stats", "retained"),
					"The amount of retained messages in the EMQ node.",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.stats.Result.RetainedCount)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "stats", "routes"),
					"The amount of routes in use by the EMQ node.",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.stats.Result.RoutesCount)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "stats", "sessions"),
					"The amount of sessions in use by the EMQ node.",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.stats.Result.SessionsCount)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "stats", "subscribers"),
					"The amount of subscribers using the EMQ node.",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.stats.Result.SubscribersCount)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "stats", "subscriptions"),
					"The amount of subscriptions in use by the EMQ node.",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.stats.Result.SubscribersCount)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "stats", "topics"),
					"The amount of topics being used in the EMQ node.",
					defaultLabels, nil,
				),
				Value: func(values combinedResponse) float64 {
					return float64(values.stats.Result.TopicsCount)
				},
			},
		},
	}
}

func (c *Collector) fetchAndDecodeNodes() (nodesResponse, error) {
	var chr nodesResponse

	u := *c.url
	u.Path = "/api/v2/monitoring/nodes/" + c.node
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return chr, fmt.Errorf("failed to get nodes response from %s://%s:%s%s: %s",
			u.Scheme, u.Hostname(), u.Port(), u.Path, err)
	}
	req.SetBasicAuth(c.username, c.password)
	res, err := c.client.Do(req)
	if err != nil {
		return chr, fmt.Errorf("failed to get nodes response from %s://%s:%s%s: %s",
			u.Scheme, u.Hostname(), u.Port(), u.Path, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return chr, fmt.Errorf("HTTP Request failed with code %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&chr); err != nil {
		c.jsonParseFailures.Inc()
		return chr, err
	}

	return chr, nil
}

func (c *Collector) fetchAndDecodeMetrics() (metricsResponse, error) {
	var chr metricsResponse

	u := *c.url
	u.Path = "/api/v2/monitoring/metrics/" + c.node
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return chr, fmt.Errorf("failed to get metrics from %s://%s:%s%s: %s",
			u.Scheme, u.Hostname(), u.Port(), u.Path, err)
	}
	req.SetBasicAuth(c.username, c.password)
	res, err := c.client.Do(req)
	if err != nil {
		return chr, fmt.Errorf("failed to get metrics from %s://%s:%s%s: %s",
			u.Scheme, u.Hostname(), u.Port(), u.Path, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return chr, fmt.Errorf("HTTP Request failed with code %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&chr); err != nil {
		c.jsonParseFailures.Inc()
		return chr, err
	}

	return chr, nil
}

func (c *Collector) fetchAndDecodeStats() (statsResponse, error) {
	var chr statsResponse

	u := *c.url
	u.Path = "/api/v2/monitoring/stats/" + c.node
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return chr, fmt.Errorf("failed to get stats from %s://%s:%s%s: %s",
			u.Scheme, u.Hostname(), u.Port(), u.Path, err)
	}
	req.SetBasicAuth(c.username, c.password)
	res, err := c.client.Do(req)
	if err != nil {
		return chr, fmt.Errorf("failed to get stats from %s://%s:%s%s: %s",
			u.Scheme, u.Hostname(), u.Port(), u.Path, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return chr, fmt.Errorf("HTTP Request failed with code %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&chr); err != nil {
		c.jsonParseFailures.Inc()
		return chr, err
	}

	return chr, nil
}

func (c *Collector) fetchAndDecodeManagment() (managementResponse, error) {
	var chr managementResponse

	u := *c.url
	u.Path = "/api/v2/management/nodes"
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return chr, fmt.Errorf("failed to get management info from %s://%s:%s%s: %s",
			u.Scheme, u.Hostname(), u.Port(), u.Path, err)
	}
	req.SetBasicAuth(c.username, c.password)
	res, err := c.client.Do(req)
	if err != nil {
		return chr, fmt.Errorf("failed to get management info from %s://%s:%s%s: %s",
			u.Scheme, u.Hostname(), u.Port(), u.Path, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return chr, fmt.Errorf("HTTP Request failed with code %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&chr); err != nil {
		c.jsonParseFailures.Inc()
		return chr, err
	}

	return chr, nil
}

// Describe is the describe fucntion function used by the prometheus package
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.Desc
	}

	ch <- c.up.Desc()
	ch <- c.totalScrapes.Desc()
	ch <- c.jsonParseFailures.Desc()
}

// Collect is the collect fucntion function used by the prometheus package
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.totalScrapes.Inc()
	defer func() {
		ch <- c.up
		ch <- c.totalScrapes
		ch <- c.jsonParseFailures
	}()

	nodes, err := c.fetchAndDecodeNodes()
	if err != nil {
		c.up.Set(0)
		log.Error(err)
		return
	}

	metrics, err := c.fetchAndDecodeMetrics()
	if err != nil {
		c.up.Set(0)
		log.Error(err)
		return
	}

	stats, err := c.fetchAndDecodeStats()
	if err != nil {
		c.up.Set(0)
		log.Error(err)
		return
	}

	management, err := c.fetchAndDecodeManagment()
	if err != nil {
		c.up.Set(0)
		log.Error(err)
		return
	}
	var ClusterSize = len(management.Result)
	var managementData ManagementResponseResult

	for _, v := range management.Result {
		if v.Name == c.node {
			managementData = v
		}
	}

	values := combinedResponse{
		nodes,
		metrics,
		stats,
		ClusterSize,
	}

	if values.nodes.Code == 0 {
		c.up.Set(1)
	} else {
		c.up.Set(0)
	}

	for _, metric := range c.metrics {
		ch <- prometheus.MustNewConstMetric(
			metric.Desc,
			metric.Type,
			metric.Value(values),
			values.nodes.Result.NodeName,
			values.nodes.Result.Release,
			managementData.Version,
		)
	}
}
