package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	serverTimingHeader = "Server-Timing"
	appTimingMetric    = "app"
)

func setAppServerTiming(header http.Header, duration time.Duration) {
	metric := formatServerTiming(duration)
	existingValues := header.Values(serverTimingHeader)
	if len(existingValues) == 0 {
		header.Set(serverTimingHeader, metric)
		return
	}

	items := make([]string, 0, len(existingValues)+1)
	for _, value := range existingValues {
		for _, item := range strings.Split(value, ",") {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}
			if strings.EqualFold(serverTimingMetricName(item), appTimingMetric) {
				continue
			}
			items = append(items, item)
		}
	}
	items = append(items, metric)
	header.Set(serverTimingHeader, strings.Join(items, ", "))
}

func formatServerTiming(duration time.Duration) string {
	durationMS := float64(duration) / float64(time.Millisecond)
	return appTimingMetric + ";dur=" + strconv.FormatFloat(durationMS, 'f', 1, 64)
}

func serverTimingMetricName(item string) string {
	name := item
	if before, _, found := strings.Cut(item, ";"); found {
		name = before
	}
	return strings.TrimSpace(name)
}
