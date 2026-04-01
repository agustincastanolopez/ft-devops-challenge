package enricher

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	// Blank imports keep dependencies in go.mod so candidates can use them immediately.
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/rabbitmq/amqp091-go"
	_ "github.com/redis/go-redis/v9"
	_ "github.com/segmentio/kafka-go"
)

// Config holds connection strings for all downstream dependencies.
type Config struct {
	MySQLDSN       string
	RedisAddr      string
	RabbitMQURL    string
	ClickHouseAddr string
	ProducerTopic  string
	KafkaBrokers   []string
}

// RawEvent represents an incoming player event from Kafka.
type RawEvent struct {
	EventID    string          `json:"event_id"`
	PlayerID   string          `json:"player_id"`
	OperatorID string          `json:"operator_id"`
	EventType  string          `json:"event_type"`
	Timestamp  string          `json:"timestamp"`
	Payload    json.RawMessage `json:"payload"`
}

// EnrichedEvent is a RawEvent augmented with player profile metadata.
type EnrichedEvent struct {
	RawEvent
	PlayerTier  string `json:"player_tier"`
	PlayerSince string `json:"player_since"`
	Country     string `json:"country"`
	EnrichedAt  string `json:"enriched_at"`
}

// Pipeline orchestrates the enrichment stages. All methods are stubs (no-ops).
type Pipeline struct {
	cfg Config
}

// NewPipeline returns a Pipeline configured with the given settings.
// Connections are NOT established here — this is a stub.
func NewPipeline(cfg Config) *Pipeline {
	return &Pipeline{cfg: cfg}
}

// Process is the main entry point called by the consumer for each message.
// It runs the full enrichment pipeline. Currently all stages are no-ops.
func (p *Pipeline) Process(ctx context.Context, key, value []byte) error {
	start := time.Now()

	var raw RawEvent
	if err := json.Unmarshal(value, &raw); err != nil {
		slog.Error("failed to unmarshal raw event", "error", err)
		return err
	}

	// Security: log operator_id and event_type, but NOT player_id.
	slog.Info("processing event",
		"event_id", raw.EventID,
		"operator_id", raw.OperatorID,
		"event_type", raw.EventType,
	)

	// Stage 1: Deduplicate via Redis (60-second window).
	isDuplicate, err := p.deduplicate(ctx, raw)
	if err != nil {
		return err
	}
	if isDuplicate {
		slog.Info("duplicate event skipped", "event_id", raw.EventID)
		return nil
	}

	// Stage 2: Enrich from Aurora MySQL.
	enriched, err := p.enrich(ctx, raw)
	if err != nil {
		return err
	}

	// Stage 3: Publish enriched event to Kafka.
	if err := p.publishToKafka(ctx, enriched); err != nil {
		return err
	}

	// Stage 4: Enqueue high-priority events to RabbitMQ.
	if err := p.enqueueHighPriority(ctx, enriched); err != nil {
		return err
	}

	// Stage 5: Write to ClickHouse for analytics.
	if err := p.writeToClickHouse(ctx, enriched); err != nil {
		return err
	}

	slog.Info("event processed successfully",
		"event_id", raw.EventID,
		"operator_id", raw.OperatorID,
		"event_type", raw.EventType,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return nil
}

// deduplicate checks Redis for a recently-seen event_id.
// STUB: always returns false (not duplicate).
func (p *Pipeline) deduplicate(_ context.Context, event RawEvent) (bool, error) {
	slog.Debug("stub: redis dedup check", "event_id", event.EventID)
	// Real implementation: SET event_id NX EX 60
	return false, nil
}

// enrich reads player profile metadata from Aurora MySQL.
// STUB: returns the event with placeholder enrichment data.
func (p *Pipeline) enrich(_ context.Context, raw RawEvent) (EnrichedEvent, error) {
	slog.Debug("stub: mysql player lookup", "operator_id", raw.OperatorID)
	// Real implementation: SELECT tier, created_at, country FROM players WHERE id = ?
	return EnrichedEvent{
		RawEvent:    raw,
		PlayerTier:  "stub_gold",
		PlayerSince: "2024-01-01T00:00:00Z",
		Country:     "IE",
		EnrichedAt:  time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// publishToKafka writes the enriched event to player.events.enriched.
// STUB: logs and returns nil.
func (p *Pipeline) publishToKafka(_ context.Context, event EnrichedEvent) error {
	slog.Debug("stub: kafka produce",
		"topic", p.cfg.ProducerTopic,
		"event_id", event.EventID,
	)
	return nil
}

// enqueueHighPriority sends deposit and churn_risk events to RabbitMQ.
// STUB: logs and returns nil.
func (p *Pipeline) enqueueHighPriority(_ context.Context, event EnrichedEvent) error {
	switch event.EventType {
	case "deposit", "churn_risk":
		slog.Debug("stub: rabbitmq enqueue",
			"event_type", event.EventType,
			"event_id", event.EventID,
		)
	default:
		// Not a high-priority event; skip.
	}
	return nil
}

// writeToClickHouse inserts a denormalized event record for analytics.
// STUB: logs and returns nil.
func (p *Pipeline) writeToClickHouse(_ context.Context, event EnrichedEvent) error {
	slog.Debug("stub: clickhouse insert", "event_id", event.EventID)
	return nil
}
