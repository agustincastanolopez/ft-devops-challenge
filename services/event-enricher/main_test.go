package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandleHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	handleHealthz(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var body map[string]string
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %s", body["status"])
	}
}

func TestHandleReadyz_NotReady(t *testing.T) {
	ready.Store(false)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()

	handleReadyz(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", res.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["status"] != "not_ready" {
		t.Errorf("expected status not_ready, got %s", body["status"])
	}
}

func TestHandleReadyz_Ready(t *testing.T) {
	ready.Store(true)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()

	handleReadyz(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["status"] != "ready" {
		t.Errorf("expected status ready, got %s", body["status"])
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	// Unset any env vars that might interfere.
	envVars := []string{
		"HTTP_PORT", "KAFKA_BROKERS", "KAFKA_CONSUMER_TOPIC",
		"KAFKA_PRODUCER_TOPIC", "KAFKA_CONSUMER_GROUP",
		"MYSQL_DSN", "REDIS_ADDR", "RABBITMQ_URL", "CLICKHOUSE_ADDR",
	}
	for _, k := range envVars {
		t.Setenv(k, "")
	}

	cfg := loadConfig()

	if cfg.HTTPPort != "8080" {
		t.Errorf("expected default HTTP_PORT 8080, got %s", cfg.HTTPPort)
	}
	if cfg.ConsumerTopic != "player.events.raw" {
		t.Errorf("expected default consumer topic player.events.raw, got %s", cfg.ConsumerTopic)
	}
	if cfg.ProducerTopic != "player.events.enriched" {
		t.Errorf("expected default producer topic player.events.enriched, got %s", cfg.ProducerTopic)
	}
	if cfg.ConsumerGroup != "event-enricher" {
		t.Errorf("expected default consumer group event-enricher, got %s", cfg.ConsumerGroup)
	}
}

func TestLoadConfig_FromEnv(t *testing.T) {
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("KAFKA_BROKERS", "broker1:9094,broker2:9094")

	cfg := loadConfig()

	if cfg.HTTPPort != "9090" {
		t.Errorf("expected HTTP_PORT 9090, got %s", cfg.HTTPPort)
	}
	if cfg.KafkaBrokers != "broker1:9094,broker2:9094" {
		t.Errorf("expected KAFKA_BROKERS broker1:9094,broker2:9094, got %s", cfg.KafkaBrokers)
	}
}

func TestEnvOrDefault(t *testing.T) {
	t.Run("returns fallback when env is empty", func(t *testing.T) {
		os.Unsetenv("TEST_ENV_OR_DEFAULT")
		if got := envOrDefault("TEST_ENV_OR_DEFAULT", "fallback"); got != "fallback" {
			t.Errorf("expected fallback, got %s", got)
		}
	})

	t.Run("returns env value when set", func(t *testing.T) {
		t.Setenv("TEST_ENV_OR_DEFAULT", "custom")
		if got := envOrDefault("TEST_ENV_OR_DEFAULT", "fallback"); got != "custom" {
			t.Errorf("expected custom, got %s", got)
		}
	})
}
