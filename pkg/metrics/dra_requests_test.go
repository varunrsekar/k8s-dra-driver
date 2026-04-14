package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"k8s.io/component-base/metrics/legacyregistry"
)

func resetDRARequestMetricsForTest() {
	legacyregistry.Reset()
	registerOnce = sync.Once{}
	draRequestsTotal.Reset()
	draRequestDurationSeconds.Reset()
	draRequestsInFlight.Reset()
	draPreparedDevices.Reset()
	nodePrepareErrorsTotal.Reset()
	nodeUnprepareErrorsTotal.Reset()
}

func TestInitializeDRARequestMetricsExposesZeroValuedSeries(t *testing.T) {
	resetDRARequestMetricsForTest()

	InitializeDRARequestMetrics("gpu.nvidia.com")

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	NewLegacyPrometheusHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	for _, expected := range []string{
		`nvidia_dra_requests_total{driver="gpu.nvidia.com",operation="prepare"} 0`,
		`nvidia_dra_requests_total{driver="gpu.nvidia.com",operation="unprepare"} 0`,
		`nvidia_dra_request_duration_seconds_count{driver="gpu.nvidia.com",operation="prepare"} 0`,
		`nvidia_dra_requests_inflight{driver="gpu.nvidia.com",operation="prepare"} 0`,
	} {
		if !strings.Contains(body, expected) {
			t.Fatalf("metrics output missing %s\nfull output:\n%s", expected, body)
		}
	}
}
