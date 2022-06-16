package base

import (
	"context"
	"errors"
	"fmt"
	"github.com/novabankapp/saga.common/domain"
	"github.com/novabankapp/saga.common/events"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type SagaSteps interface {
	GetCompensatingStepMessage(id string) domain.SagaStepMessage
	GetStepMessage(id string) domain.SagaStepMessage
}
type SagaProcessed interface {
	Processed(id string) bool
	AlreadyProcessed(id string) bool
}
type SagaProgress interface {
	GoBack(getCompensatingStepMessage func(id string) domain.SagaStepMessage,
		saveEvent func(event events.Event) bool,
		saveSagaState func(state domain.SagaState) bool) error
	Advance(getStepMessage func(id string) domain.SagaStepMessage,
		saveEvent func(event events.Event) bool,
		saveSagaState func(state domain.SagaState) bool) error
	OnStepEvent(eventType string, status domain.SagaStepStatus)
}
type Saga interface {
	SagaSteps
	SagaProcessed
	SagaProgress
}

type SagaAbstract struct {
	Type      string
	Steps     []string
	SagaState domain.SagaState
	Conn      *gorm.DB
}

func (s *SagaAbstract) SaveEvent(event events.Event) bool {
	var rEvent = event.(events.OutboxEvent)
	result := s.Conn.Create(&rEvent).WithContext(context.Background())
	if result.Error != nil && result.RowsAffected != 1 {
		return false
	}
	return true
}
func (s *SagaAbstract) SaveSagaState(state domain.SagaState) bool {
	result := s.Conn.Save(state).WithContext(context.Background())
	if result.Error != nil && result.RowsAffected != 1 {
		return false
	}
	return true
}
func (s *SagaAbstract) Processed(id string) bool {
	result := s.Conn.Create(domain.ConsumedMessage{
		Id:        id,
		CreatedAt: time.Now(),
	}).WithContext(context.Background())
	if result.Error != nil && result.RowsAffected != 1 {
		return false
	}
	return true
}

func (s *SagaAbstract) AlreadyProcessed(id string) bool {
	var consumed domain.ConsumedMessage
	result := s.Conn.First(&consumed, "id = ?", id).WithContext(context.Background())
	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false

		} else {
			return false

		}
		return false
	}
	return true
}

func (s *SagaAbstract) OnStepEvent(eventType string, status domain.SagaStepStatus) {
	stepStatus := s.SagaState.StepStatus
	stepStatus[eventType] = status

	if status == domain.SUCCEDED {
		//s.Advance()
	} else if status == domain.FAILED {
		//s.GoBack()
	}
	allStatus := make(map[string]bool)
	fields := reflect.ValueOf(stepStatus).MapKeys()
	for i := range fields {
		allStatus[fmt.Sprint(stepStatus[fields[i].String()])] = true
	}
	s.SagaState.SagaStatus = string(getSagaStatus(allStatus))
}
func (s *SagaAbstract) Advance(getStepMessage func(id string) domain.SagaStepMessage,
	saveEvent func(event events.Event) bool, saveSagaState func(state domain.SagaState) bool) error {
	nextStep := s.getNextStep()
	if nextStep == nil {
		s.SagaState.CurrentStep = ""
		return nil
	}
	stepEvent := getStepMessage(*nextStep)

	//event that saves in DB
	event := events.NewOutboxEvent(s.SagaState.Id, stepEvent.Type, stepEvent.EventType, stepEvent.Payload)
	saveEvent(event)

	s.SagaState.StepStatus[*nextStep] = domain.STARTED
	s.SagaState.CurrentStep = *nextStep
	//update sagaState
	saveSagaState(s.SagaState)

	return nil
}

func (s *SagaAbstract) GoBack(getCompensatingStepMessage func(id string) domain.SagaStepMessage,
	saveEvent func(event events.Event) bool, saveSagaState func(state domain.SagaState) bool) error {
	prevStep := s.getPreviousStep()
	if prevStep == nil {
		s.SagaState.CurrentStep = ""
		return nil
	}
	stepEvent := getCompensatingStepMessage(*prevStep)
	//event that saves in DB
	event := events.NewOutboxEvent(s.SagaState.Id, stepEvent.Type, stepEvent.EventType, stepEvent.Payload)
	saveEvent(event)

	s.SagaState.StepStatus[*prevStep] = domain.COMPENSATING
	s.SagaState.CurrentStep = *prevStep
	//update sagaState
	saveSagaState(s.SagaState)

	return nil
}
func (s *SagaAbstract) getPreviousStep() *string {
	indexFunc := func(arr []string, candidate string) int {
		for index, c := range arr {
			if c == candidate {
				return index
			}
		}
		return -1
	}
	index := indexFunc(s.Steps, s.SagaState.CurrentStep)

	if index == 0 {
		return nil
	}

	return &s.Steps[index-1]
}
func (s *SagaAbstract) getNextStep() *string {
	if s.SagaState.CurrentStep == "" {
		return &s.Steps[0]
	}
	indexFunc := func(arr []string, candidate string) int {
		for index, c := range arr {
			if c == candidate {
				return index
			}
		}
		return -1
	}
	index := indexFunc(s.Steps, s.SagaState.CurrentStep)
	if index == len(s.Steps)-1 || index == -1 {
		return nil
	}
	return &s.Steps[index+1]
}
func getSagaStatus(stepStates map[string]bool) domain.SagaStepStatus {
	if containsOnly(stepStates, domain.SUCCEDED) {
		return domain.COMPLETED
	} else if containsOnly2(stepStates, domain.STARTED, domain.SUCCEDED) {
		return domain.STARTED
	} else if containsOnly2(stepStates, domain.FAILED, domain.COMPENSATED) {
		return domain.ABORTED
	} else {
		return domain.ABORTING
	}
}
func containsOnly(stepStates map[string]bool, status domain.SagaStepStatus) bool {
	for k, _ := range stepStates {
		if k != string(status) {
			return false
		}
	}
	return true
}

func containsOnly2(stepStates map[string]bool, status1, status2 domain.SagaStepStatus) bool {
	for k, _ := range stepStates {
		if k != string(status1) && k != string(status2) {
			return false
		}
	}
	return true
}

/*func NewSagaAbstract(steps []string, sagaType string, sagaState domain.SagaState) Saga {
	return &SagaAbstract{
		steps,
		sagaType,
		sagaState,
	}
}*/

/*type ProcessPaymentSaga struct {
	Type string
	SagaAbstract
}

func Te() {

	p := ProcessPaymentSaga{
		Type: "Baba",
		SagaAbstract: SagaAbstract{
			Steps:[]string{"","",""},
			SagaState : domain.SagaState{

			},
		},
	}
	p.getNextStep()
}*/
