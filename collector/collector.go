package collector

import (
	"log/slog"
	"strings"
	"time"

	"github.com/attilagyorffy/prometheus-exporter-omada-controller/omada"
	"github.com/prometheus/client_golang/prometheus"
)

// Instrumentation for the collector itself.
var (
	omadaCollectionsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "omada_collector_requests_total",
			Help: "number of requests to collect wireless data",
		})
	omadaErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "omada_collector_errors_total",
			Help: "number of errors while collecting wireless data",
		})
	omadaCollectorDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "omada_collector_duration_seconds",
			Help: "collection from the controller took this many seconds",
		})
)

func init() {
	prometheus.MustRegister(omadaCollectionsTotal)
	prometheus.MustRegister(omadaErrorsTotal)
	prometheus.MustRegister(omadaCollectorDuration)
}

// Collector implements the prometheus.Collector interface.
type Collector struct {
	controller *omada.Client
}

// New returns a prometheus.Collector that exports wifi station data.
func New(controller *omada.Client) *Collector {
	return &Collector{controller: controller}
}

// fixMAC adjusts address format from 00-00-00-FF-FF-FF to 00:00:00:ff:ff:ff.
func fixMAC(mac string) string {
	return strings.ToLower(strings.ReplaceAll(mac, "-", ":"))
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range stationMetrics {
		ch <- m.Desc
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	omadaCollectionsTotal.Inc()
	start := time.Now()
	defer func() {
		ms := time.Since(start).Milliseconds()
		omadaCollectorDuration.Set(float64(ms) / 1000.0)
	}()

	sites, err := c.controller.Sites()
	if err != nil {
		omadaErrorsTotal.Inc()
		slog.Error("failed to get sites", "error", err)
		return
	}

	for _, site := range sites {
		stations, err := c.controller.ConnectedClients(site.SiteID())
		if err != nil {
			omadaErrorsTotal.Inc()
			slog.Error("failed to get clients", "error", err)
			continue // try the next site
		}

		for _, sta := range stations {
			staMAC := fixMAC(sta.MAC)
			apMAC := fixMAC(sta.ApMAC)
			for _, metric := range stationMetrics {
				ch <- prometheus.MustNewConstMetric(
					metric.Desc,
					metric.ValueType,
					metric.Value(sta),
					staMAC,
					sta.Name,
					apMAC,
					sta.ApName,
					sta.SSID,
					site.Name)
			}
		}
	}
}
