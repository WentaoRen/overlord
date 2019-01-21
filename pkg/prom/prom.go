package prom

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	statConns     = "overlord_proxy_conns"
	versionCounts = "overlord_version"
	statErr       = "overlord_proxy_err"

	statProxyTimer   = "overlord_proxy_timer"
	statHandlerTimer = "overlord_proxy_handler_timer"
)

var (
	conns        *prometheus.GaugeVec
	versions     *prometheus.GaugeVec
	gerr         *prometheus.GaugeVec
	proxyTimer   *prometheus.HistogramVec
	handlerTimer *prometheus.HistogramVec

	versionLabels        = []string{"appid", "version"}
	clusterLabels        = []string{"cluster"}
	clusterNodeErrLabels = []string{"cluster", "node", "cmd", "error"}
	clusterCmdLabels     = []string{"cluster", "cmd"}
	clusterNodeCmdLabels = []string{"cluster", "node", "cmd"}
	// On Prom switch
	On = true
)

// Init init prometheus.
func Init() {
	conns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: statConns,
			Help: statConns,
		}, clusterLabels)
	prometheus.MustRegister(conns)
	versions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: versionCounts,
			Help: versionCounts,
		}, versionLabels)
	prometheus.MustRegister(versions)
	gerr = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: statErr,
			Help: statErr,
		}, clusterNodeErrLabels)
	prometheus.MustRegister(gerr)
	proxyTimer = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    statProxyTimer,
			Help:    statProxyTimer,
			Buckets: prometheus.LinearBuckets(0, 10, 1),
		}, clusterCmdLabels)
	prometheus.MustRegister(proxyTimer)
	handlerTimer = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    statHandlerTimer,
			Help:    statHandlerTimer,
			Buckets: prometheus.LinearBuckets(0, 10, 1),
		}, clusterNodeCmdLabels)
	prometheus.MustRegister(handlerTimer)
	// metrics
	metrics()
}

func metrics() {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h := promhttp.Handler()
		h.ServeHTTP(w, r)
	})
}

// ProxyTime log timing information (in milliseconds).
func ProxyTime(cluster, node string, ts int64) {
	if proxyTimer == nil {
		return
	}
	proxyTimer.WithLabelValues(cluster, node).Observe(float64(ts))
}

// HandleTime log timing information (in milliseconds).
func HandleTime(cluster, node, cmd string, ts int64) {
	if handlerTimer == nil {
		return
	}
	handlerTimer.WithLabelValues(cluster, node, cmd).Observe(float64(ts))
}

// ErrIncr increments one stat error counter.
func ErrIncr(cluster, node, cmd, err string) {
	if gerr == nil {
		return
	}
	gerr.WithLabelValues(cluster, node, cmd, err).Inc()
}

// ConnIncr increments one stat error counter.
func ConnIncr(cluster string) {
	if conns == nil {
		return
	}
	conns.WithLabelValues(cluster).Inc()
}

// ConnDecr decrements one stat error counter.
func ConnDecr(cluster string) {
	if conns == nil {
		return
	}
	conns.WithLabelValues(cluster).Dec()
}

// VersionIncr incr version in use
func VersionIncr(appid, version string) {
	if versions == nil {
		return
	}
	versions.WithLabelValues(appid, version).Inc()
}

// VersionDecr decr version in use
func VersionDecr(appid, version string) {
	if versions == nil {
		return
	}
	versions.WithLabelValues(appid, version).Dec()
}
