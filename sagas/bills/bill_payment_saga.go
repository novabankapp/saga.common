package bills

import (
	"context"
	"github.com/google/uuid"
	baseRepository "github.com/novabankapp/common.data/repositories/base/cassandra"
	"github.com/novabankapp/saga.common/domain"
	"github.com/novabankapp/saga.common/domain/base"
	"github.com/novabankapp/saga.common/events"
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
	repo baseRepository.CassandraRepository[events.OutboxEvent]
	base.SagaAbstract
}

func NewBillPaymentSaga(
	state domain.SagaState,
	sagaRepo baseRepository.CassandraRepository[domain.SagaState],
	messageRepo baseRepository.CassandraRepository[domain.ConsumedMessage],
	repo baseRepository.CassandraRepository[events.OutboxEvent],
) base.Saga {
	steps := []string{RESERVE_CREDIT, BUY_TOKEN}
	return &BillPaymentSaga{
		repo,
		base.SagaAbstract{
			Type:        TYPE,
			Steps:       steps,
			SagaState:   state,
			MessageRepo: messageRepo,
			SagaRepo:    sagaRepo,
		},
	}
}
func (s *BillPaymentSaga) SaveEvent(event events.Event) bool {
	var rEvent = event.(events.OutboxEvent)
	result, err := s.repo.Create(context.Background(), rEvent)
	if err != nil {
		return false
	}
	return result
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
	_, err := s.SagaRepo.Create(context.Background(), s.SagaState)
	if err != nil {

	}
	s.Advance(s.GetStepMessage, s.SaveEvent, s.SaveSagaState)
}
func (s *BillPaymentSaga) Find(ctx context.Context, id string) (*domain.SagaState, error) {

	result, err := s.SagaRepo.GetById(context.Background(), id)
	if err != nil {

		return nil, err

	}
	return result, nil
}
