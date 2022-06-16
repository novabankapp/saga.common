package implementations

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/novabankapp/saga.common/domain"
	"github.com/novabankapp/saga.common/domain/base"
	"gorm.io/gorm"
)

const (
	REQUEST = "REQUEST"
	CANCEL  = "CANCEL"
	//PAYMENT        = "PAYMENT"
	RESERVE_CREDIT = "RESERVE_CREDIT"
	BUY_TOKEN      = "BUY_TOKEN"
	TYPE           = "bill-payment"
)

type BillPaymentSaga struct {
	base.SagaAbstract
}

func NewSaga(sagaType string, state domain.SagaState, conn *gorm.DB) base.Saga {
	steps := []string{RESERVE_CREDIT, BUY_TOKEN}
	return &BillPaymentSaga{
		base.SagaAbstract{
			Type:      TYPE,
			Steps:     steps,
			SagaState: state,
			Conn:      conn,
		},
	}
}
func (s *BillPaymentSaga) GetCompensatingStepMessage(id string) domain.SagaStepMessage {
	payload := s.SagaState.Payload
	payload["type"] = CANCEL
	switch id {
	case RESERVE_CREDIT:
		return domain.SagaStepMessage{
			Type:      RESERVE_CREDIT,
			EventType: REQUEST,
		}
	case BUY_TOKEN:
		return domain.SagaStepMessage{
			Type:      BUY_TOKEN,
			EventType: REQUEST,
			Payload:   payload,
		}
	default:
		return domain.SagaStepMessage{
			Type:      RESERVE_CREDIT,
			EventType: REQUEST,
			Payload:   payload,
		}
	}

}
func (s *BillPaymentSaga) GetStepMessage(id string) domain.SagaStepMessage {
	switch id {
	case RESERVE_CREDIT:
		return domain.SagaStepMessage{
			Type:      RESERVE_CREDIT,
			EventType: REQUEST,
			Payload:   s.SagaState.Payload,
		}
	case BUY_TOKEN:
		return domain.SagaStepMessage{
			Type:      BUY_TOKEN,
			EventType: REQUEST,
			Payload:   s.SagaState.Payload,
		}
	default:
		return domain.SagaStepMessage{}
	}
}

//manager
func (s *BillPaymentSaga) Begin(ctx context.Context, payload base.JSONB) {
	id := uuid.New().String()

	state := domain.SagaState{
		Payload:     payload,
		Id:          id,
		Type:        TYPE,
		Version:     0,
		SagaStatus:  string(domain.STARTED),
		CurrentStep: "",
		StepStatus:  make(map[string]interface{}),
	}
	s.SagaState = state
	result := s.Conn.Create(&s.SagaState).WithContext(ctx)
	if result.Error != nil && result.RowsAffected != 1 {

	}
	s.Advance(s.GetStepMessage, s.SaveEvent, s.SaveSagaState)
}
func (s *BillPaymentSaga) Find(ctx context.Context, id string) (*domain.SagaState, error) {
	var state domain.SagaState
	result := s.Conn.First(&state, "id = ?", id).WithContext(context.Background())
	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("saga not found")

		}
		return nil, result.Error

	}
	return &state, nil
}
