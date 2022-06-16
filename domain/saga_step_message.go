package domain

import "github.com/novabankapp/saga.common/domain/base"

type SagaStepMessage struct {
	Type      string
	EventType string
	Payload   base.JSONB
}

func NewSagaStepMessage(stepType, eventType string, payload base.JSONB) *SagaStepMessage {
	return &SagaStepMessage{
		stepType,
		eventType,
		payload,
	}
}
