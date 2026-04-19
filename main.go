package main

import (
	"net/http"
	"os"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	lastTimestampMetric = "hyperbackup_last_backup_success_timestamp_seconds"
	lastVersionMetric   = "hyperbackup_last_backup_success_counter"
)

var (
	metrics = map[string]*prometheus.Desc{
		lastTimestampMetric: prometheus.NewDesc(
			lastTimestampMetric,
			"Unix Timestamp of the last time the task completed",
			[]string{"task"}, nil,
		),
		lastVersionMetric: prometheus.NewDesc(
			lastVersionMetric,
			"Last completed version",
			[]string{"task"}, nil,
		),
	}
)

type Exporter struct {
	backupLastPath, synobackupPath string
	mutex                          sync.RWMutex
	up                             prometheus.Gauge
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range metrics {
		ch <- m
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		metrics[lastTimestampMetric],
		prometheus.GaugeValue,
		0,
		"This is an example")

	ch <- prometheus.MustNewConstMetric(
		metrics[lastVersionMetric],
		prometheus.GaugeValue,
		0,
		"This is an example")
}

// https://stackoverflow.com/a/40326580
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
	var (
		backupLastPath = getEnv("BACKUP_LAST", "/var/packages/HyperBackup/var/last_result/backup.last")
		synobackupPath = getEnv("SYNOBACKUP", "/var/packages/HyperBackup/etc/synobackup.conf")
		port           = getEnv("PORT", "6533")
	)

	exporter := &Exporter{
		backupLastPath: backupLastPath,
		synobackupPath: synobackupPath,
	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(exporter)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(":"+port, nil)
}
