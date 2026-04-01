package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/fasttrack/event-enricher/consumer"
	"github.com/fasttrack/event-enricher/enricher"
)

// ready is flipped to true once the consumer loop has started successfully.
var ready atomic.Bool

func main() {
	// Structured JSON logging — no PII fields.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := loadConfig()

	brokers := strings.Split(cfg.KafkaBrokers, ",")

	slog.Info("starting event-enricher",
		"kafka_brokers", cfg.KafkaBrokers,
		"consumer_topic", cfg.ConsumerTopic,
		"producer_topic", cfg.ProducerTopic,
	)

	// Build the enricher pipeline (all stubs/no-ops for now).
	pipeline := enricher.NewPipeline(enricher.Config{
		MySQLDSN:       cfg.MySQLDSN,
		RedisAddr:      cfg.RedisAddr,
		RabbitMQURL:    cfg.RabbitMQURL,
		ClickHouseAddr: cfg.ClickHouseAddr,
		ProducerTopic:  cfg.ProducerTopic,
		KafkaBrokers:   brokers,
	})

	// Build the Kafka consumer.
	cons := consumer.New(consumer.Config{
		Brokers:       brokers,
		Topic:         cfg.ConsumerTopic,
		ConsumerGroup: cfg.ConsumerGroup,
		Pipeline:      pipeline,
	})

	// Context cancelled on SIGINT/SIGTERM for graceful shutdown.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Start health/readiness HTTP server.
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handleHealthz)
	mux.HandleFunc("/readyz", handleReadyz)

	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("http server listening", "port", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server error", "error", err)
			os.Exit(1)
		}
	}()

	// Start consuming (blocks until context is cancelled).
	go func() {
		ready.Store(true)
		slog.Info("consumer loop starting")
		if err := cons.Run(ctx); err != nil {
			slog.Error("consumer exited with error", "error", err)
		}
	}()

	// Block until shutdown signal.
	<-ctx.Done()
	slog.Info("shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
}

// config holds all environment-driven configuration.
type config struct {
	HTTPPort       string
	KafkaBrokers   string
	ConsumerTopic  string
	ProducerTopic  string
	ConsumerGroup  string
	MySQLDSN       string
	RedisAddr      string
	RabbitMQURL    string
	ClickHouseAddr string
}

func loadConfig() config {
	return config{
		HTTPPort:       envOrDefault("HTTP_PORT", "8080"),
		KafkaBrokers:   envOrDefault("KAFKA_BROKERS", "localhost:9092"),
		ConsumerTopic:  envOrDefault("KAFKA_CONSUMER_TOPIC", "player.events.raw"),
		ProducerTopic:  envOrDefault("KAFKA_PRODUCER_TOPIC", "player.events.enriched"),
		ConsumerGroup:  envOrDefault("KAFKA_CONSUMER_GROUP", "event-enricher"),
		MySQLDSN:       envOrDefault("MYSQL_DSN", "root:@tcp(localhost:3306)/players"),
		RedisAddr:      envOrDefault("REDIS_ADDR", "localhost:6379"),
		RabbitMQURL:    envOrDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		ClickHouseAddr: envOrDefault("CLICKHOUSE_ADDR", "http://localhost:8123"),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleReadyz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if ready.Load() {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
	json.NewEncoder(w).Encode(map[string]string{"status": "not_ready"})
}
