package saga_manager

import (
	"context"
	"github.com/google/uuid"
	"github.com/novabankapp/saga.common/domain"
	"github.com/novabankapp/saga.common/domain/base"
	"gorm.io/gorm"
)

type SagaManager[E base.SagaBase] struct {
	conn *gorm.DB
}

func NewSagaManager[E base.SagaBase](conn *gorm.DB) *SagaManager[E] {
	return &SagaManager[E]{
		conn,
	}
}
func (m *SagaManager[E]) Begin(ctx context.Context, sagaType string, payload base.JSONB) (*E, error) {
	id := uuid.New().String()
	state := domain.SagaState{
		Payload:     payload,
		Id:          id,
		Type:        sagaType,
		Version:     0,
		SagaStatus:  string(domain.STARTED),
		CurrentStep: "",
		StepStatus:  make(map[string]interface{}),
	}

}
