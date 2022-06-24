package events

import (
	"github.com/google/uuid"
	"github.com/novabankapp/saga.common/domain/base"
	"time"
)

type Event interface {
}

type OutboxEvent struct {
	ID            string     `gorm:"primaryKey, column:id" json:"id"`
	AggregateId   string     `gorm:"column:aggregateid"`
	Payload       base.JSONB `gorm:"type:jsonb" sql:"type:jsonb" json:"payload"`
	Type          string     `gorm:"column:type"`
	AggregateType string     `gorm:"column:aggregatetype"`
	Timestamp     time.Time  `gorm:"timestamp"`
}

func (s *OutboxEvent) IsNoSQLEntity() bool {
	return true
}

func NewOutboxEvent(aggregateId, oType, aggregateType string, payload base.JSONB) Event {
	event := OutboxEvent{
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		Type:          oType,
		Payload:       payload,
	}
	event.FillDefaults()
	return &event
}
func (s *OutboxEvent) FillDefaults() {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	s.Timestamp = time.Now()
}
