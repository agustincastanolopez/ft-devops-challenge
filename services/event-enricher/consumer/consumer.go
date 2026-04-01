package consumer

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

// Pipeline is the interface the enricher must satisfy.
type Pipeline interface {
	Process(ctx context.Context, key, value []byte) error
}

// Config holds Kafka consumer configuration.
type Config struct {
	Brokers       []string
	Topic         string
	ConsumerGroup string
	Pipeline      Pipeline
}

// Consumer reads from a Kafka topic and passes messages to the enricher pipeline.
type Consumer struct {
	reader   *kafka.Reader
	pipeline Pipeline
}

// New creates a Consumer. It does NOT start consuming — call Run() for that.
func New(cfg Config) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topic,
		GroupID:  cfg.ConsumerGroup,
		MinBytes: 1e3,  // 1 KB
		MaxBytes: 10e6, // 10 MB
	})

	return &Consumer{
		reader:   r,
		pipeline: cfg.Pipeline,
	}
}

// Run starts the consume loop. It blocks until ctx is cancelled.
func (c *Consumer) Run(ctx context.Context) error {
	defer func() {
		if err := c.reader.Close(); err != nil {
			slog.Error("failed to close kafka reader", "error", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			slog.Info("consumer context cancelled, stopping")
			return ctx.Err()
		default:
		}

		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			slog.Error("failed to read message",
				"error", err,
				"topic", c.reader.Config().Topic,
			)
			continue
		}

		// Security: log event metadata but never player_id or PII.
		slog.Info("message received",
			"topic", msg.Topic,
			"partition", msg.Partition,
			"offset", msg.Offset,
		)

		if err := c.pipeline.Process(ctx, msg.Key, msg.Value); err != nil {
			slog.Error("pipeline processing failed",
				"error", err,
				"topic", msg.Topic,
				"partition", msg.Partition,
				"offset", msg.Offset,
			)
			// TODO: implement DLQ / retry strategy
			continue
		}
	}
}
