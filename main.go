package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/ini.v1"
)

const (
	lastTimestampMetric = "hyperbackup_last_backup_success_timestamp_seconds"
	lastVersionMetric   = "hyperbackup_last_backup_success_version"
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

type HyperBackupTask struct {
	Name          string
	LastUpdatedAt int64
	Version       int64
}

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
	e.mutex.Lock()
	defer e.mutex.Unlock()

	synoBackup, err := ini.Load(e.synobackupPath)
	if err != nil {
		// TODO: implement up metric
		log.Printf("Failed to load synobackup config: %v", err)
		return
	}

	backupLast, err := ini.Load(e.backupLastPath)
	if err != nil {
		log.Printf("Failed to read backup.last: %v", err)
		return
	}

	var tasks []HyperBackupTask

	// synobackup.conf lets us look up the numbers to names
	for _, section := range synoBackup.Sections() {
		if !strings.HasPrefix(section.Name(), "task_") {
			continue
		}

		lastBackupData, err := backupLast.GetSection(section.Name())
		if err != nil {
			tasks = append(tasks, HyperBackupTask{
				Name:          section.Key("name").String(),
				LastUpdatedAt: 0,
				Version:       0,
			})
			continue
		}

		tasks = append(tasks, HyperBackupTask{
			Name:          section.Key("name").String(),
			LastUpdatedAt: lastBackupData.Key("last_backup_success_time").MustInt64(),
			Version:       lastBackupData.Key("last_backup_success_version").MustInt64(),
		})
	}

	for _, task := range tasks {
		ch <- prometheus.MustNewConstMetric(
			metrics[lastTimestampMetric],
			prometheus.GaugeValue,
			float64(task.LastUpdatedAt),
			task.Name)

		ch <- prometheus.MustNewConstMetric(
			metrics[lastVersionMetric],
			prometheus.GaugeValue,
			float64(task.Version),
			task.Name)
	}
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
