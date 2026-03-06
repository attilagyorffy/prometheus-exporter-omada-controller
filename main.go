package main

import (
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/attilagyorffy/prometheus-exporter-omada-controller/collector"
	"github.com/attilagyorffy/prometheus-exporter-omada-controller/omada"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	port := envOrDefault("LISTEN_PORT", "6779")
	omadaURL := envOrDefault("OMADA_URL", "https://10.0.0.3:30077")
	omadaUser := os.Getenv("OMADA_USER")
	omadaPass := os.Getenv("OMADA_PASS")
	insecure := strings.ToLower(envOrDefault("OMADA_INSECURE", "true"))
	logLevel := envOrDefault("LOG_LEVEL", "info")

	setupLogging(logLevel)

	if omadaUser == "" || omadaPass == "" {
		slog.Error("OMADA_USER and OMADA_PASS environment variables are required")
		os.Exit(1)
	}

	skipVerify := insecure == "true" || insecure == "1"

	controller, err := omada.NewClient(omadaURL, omadaUser, omadaPass, skipVerify)
	if err != nil {
		slog.Error("failed to connect to Omada controller", "error", err)
		os.Exit(1)
	}

	prometheus.MustRegister(collector.New(controller))

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Omada Controller Exporter\n\nVisit /metrics for Prometheus metrics.\n")
	})

	addr := ":" + port
	slog.Info("starting exporter", "address", addr, "omada_url", omadaURL)
	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func setupLogging(level string) {
	var l slog.Level
	switch level {
	case "debug":
		l = slog.LevelDebug
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: l})))
}
