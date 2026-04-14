/*
Copyright The Kubernetes Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/component-base/metrics/legacyregistry"
)

var (
	registerOnce                  sync.Once
	requestDurationSecondsBuckets = prometheus.ExponentialBuckets(0.05, 2, 9)
	draRequestsTotal              = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "nvidia_dra",
			Name:      "requests_total",
			Help:      "Total number of DRA prepare and unprepare requests.",
		},
		[]string{"driver", "operation"},
	)

	draRequestDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "nvidia_dra",
			Name:      "request_duration_seconds",
			Help:      "Duration of DRA prepare and unprepare requests.",
			Buckets:   requestDurationSecondsBuckets,
		},
		[]string{"driver", "operation"},
	)

	draRequestsInFlight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "nvidia_dra",
			Name:      "requests_inflight",
			Help:      "Number of in-flight DRA prepare and unprepare requests.",
		},
		[]string{"driver", "operation"},
	)

	draPreparedDevices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "nvidia_dra",
			Name:      "prepared_devices",
			Help:      "Current number of prepared devices by device type.",
		},
		[]string{"node", "driver", "device_type"},
	)

	// node_prepare_errors_total is scoped to kubelet node prepare failures.
	nodePrepareErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "nvidia_dra",
			Name:      "node_prepare_errors_total",
			Help:      "Total number of failures during DRA node prepare for kubelet plugins.",
		},
		[]string{"driver", "error_type"},
	)

	nodeUnprepareErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "nvidia_dra",
			Name:      "node_unprepare_errors_total",
			Help:      "Total number of failures during DRA node unprepare for kubelet plugins.",
		},
		[]string{"driver", "error_type"},
	)
)

func Register() {
	registerOnce.Do(func() {
		legacyregistry.RawMustRegister(
			draRequestsTotal,
			draRequestDurationSeconds,
			draRequestsInFlight,
			draPreparedDevices,
			nodePrepareErrorsTotal,
			nodeUnprepareErrorsTotal,
		)
	})
}

// InitializeDRARequestMetrics pre-creates zero-valued request series for a
// driver so they are visible on the first Prometheus scrape after process
// startup.
func InitializeDRARequestMetrics(driver string) {
	Register()
	initializeDRARequestMetrics(driver)
}

func initializeDRARequestMetrics(driver string) {
	if driver == "" {
		return
	}

	for _, operation := range []string{"prepare", "unprepare"} {
		draRequestsTotal.WithLabelValues(driver, operation)
		draRequestDurationSeconds.WithLabelValues(driver, operation)
		draRequestsInFlight.WithLabelValues(driver, operation)
	}
}

func TrackInFlight(driver, operation string) func() {
	Register()
	draRequestsInFlight.WithLabelValues(driver, operation).Inc()
	return func() {
		draRequestsInFlight.WithLabelValues(driver, operation).Dec()
	}
}

func ObserveRequest(driver, operation string, d time.Duration) {
	Register()
	draRequestsTotal.WithLabelValues(driver, operation).Inc()
	draRequestDurationSeconds.WithLabelValues(driver, operation).Observe(d.Seconds())
}

// IncNodePrepareError increments node_prepare_errors_total for a kubelet plugin.
func IncNodePrepareError(driver, errorType string) {
	Register()
	nodePrepareErrorsTotal.WithLabelValues(driver, errorType).Inc()
}

// IncNodeUnprepareError increments node_unprepare_errors_total for a kubelet plugin.
func IncNodeUnprepareError(driver, errorType string) {
	Register()
	nodeUnprepareErrorsTotal.WithLabelValues(driver, errorType).Inc()
}

// SetPreparedDevicesCounts sets prepared_devices gauges from a full snapshot (e.g. checkpoint-derived).
// Types that were previously non-zero but are absent or zero in counts are set back to 0.
func SetPreparedDevicesCounts(nodeName, driver, deviceType string, count int) {
	Register()
	draPreparedDevices.WithLabelValues(nodeName, driver, deviceType).Set(float64(count))
}
