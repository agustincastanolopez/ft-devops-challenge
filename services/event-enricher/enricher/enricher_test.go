package enricher

import (
	"context"
	"encoding/json"
	"testing"
)

func TestProcess_ValidEvent(t *testing.T) {
	p := NewPipeline(Config{
		ProducerTopic: "player.events.enriched",
		KafkaBrokers:  []string{"localhost:9092"},
	})

	raw := RawEvent{
		EventID:    "evt-001",
		PlayerID:   "player-123",
		OperatorID: "operator-456",
		EventType:  "login",
		Timestamp:  "2024-01-01T00:00:00Z",
		Payload:    json.RawMessage(`{"session_id":"abc"}`),
	}

	value, err := json.Marshal(raw)
	if err != nil {
		t.Fatalf("failed to marshal test event: %v", err)
	}

	if err := p.Process(context.Background(), []byte("key"), value); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestProcess_InvalidJSON(t *testing.T) {
	p := NewPipeline(Config{})

	err := p.Process(context.Background(), []byte("key"), []byte("not-json"))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestProcess_HighPriorityEventTypes(t *testing.T) {
	p := NewPipeline(Config{
		ProducerTopic: "player.events.enriched",
		KafkaBrokers:  []string{"localhost:9092"},
	})

	for _, eventType := range []string{"deposit", "churn_risk", "login", "game_session"} {
		t.Run(eventType, func(t *testing.T) {
			raw := RawEvent{
				EventID:    "evt-" + eventType,
				PlayerID:   "player-123",
				OperatorID: "operator-456",
				EventType:  eventType,
				Timestamp:  "2024-01-01T00:00:00Z",
				Payload:    json.RawMessage(`{}`),
			}
			value, _ := json.Marshal(raw)

			if err := p.Process(context.Background(), nil, value); err != nil {
				t.Errorf("expected no error for event type %s, got %v", eventType, err)
			}
		})
	}
}

func TestEnrich_ReturnsEnrichedEvent(t *testing.T) {
	p := NewPipeline(Config{})

	raw := RawEvent{
		EventID:    "evt-001",
		PlayerID:   "player-123",
		OperatorID: "operator-456",
		EventType:  "login",
		Timestamp:  "2024-01-01T00:00:00Z",
	}

	enriched, err := p.enrich(context.Background(), raw)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if enriched.EventID != raw.EventID {
		t.Errorf("expected EventID %s, got %s", raw.EventID, enriched.EventID)
	}
	if enriched.OperatorID != raw.OperatorID {
		t.Errorf("expected OperatorID %s, got %s", raw.OperatorID, enriched.OperatorID)
	}
	if enriched.EnrichedAt == "" {
		t.Error("expected EnrichedAt to be set")
	}
	if enriched.PlayerTier == "" {
		t.Error("expected PlayerTier to be set")
	}
}

func TestDeduplicate_StubReturnsFalse(t *testing.T) {
	p := NewPipeline(Config{})

	raw := RawEvent{EventID: "evt-001"}
	isDup, err := p.deduplicate(context.Background(), raw)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if isDup {
		t.Error("stub deduplicate should always return false")
	}
}
