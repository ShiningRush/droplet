package middleware

import (
	"net/http"
	"testing"
	"time"
)

func TestSetAppServerTimingMergesExistingMetrics(t *testing.T) {
	header := http.Header{}
	header.Add(serverTimingHeader, "auth;dur=4.2")
	header.Add(serverTimingHeader, "db;dur=18.6")

	setAppServerTiming(header, 37*time.Millisecond+200*time.Microsecond)

	got := header.Get(serverTimingHeader)
	if got != "auth;dur=4.2, db;dur=18.6, app;dur=37.2" {
		t.Fatalf("Server-Timing header = %q", got)
	}
}

func TestSetAppServerTimingReplacesExistingAppMetric(t *testing.T) {
	header := http.Header{}
	header.Add(serverTimingHeader, "db;dur=18.6")
	header.Add(serverTimingHeader, "app;dur=40")

	setAppServerTiming(header, 10*time.Millisecond)

	got := header.Get(serverTimingHeader)
	if got != "db;dur=18.6, app;dur=10.0" {
		t.Fatalf("Server-Timing header = %q", got)
	}
}
