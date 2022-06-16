package domain

import (
	"github.com/novabankapp/saga.common/domain/base"
)

type SagaStatus string

const (
	SSTARTED  SagaStatus = "STARTED"
	ABORTING             = "ABORTING"
	ABORTED              = "ABORTED"
	COMPLETED            = "COMPLETED"
)

type SagaState struct {
	Payload     base.JSONB `gorm:"type:jsonb" sql:"type:jsonb" json:"payload"`
	Type        string     `gorm:"column:type" json:"type"`
	Id          string     `gorm:"primary_key" json:"id"`
	Version     int        ` json:"version"`
	CurrentStep string     `json:"current_step"`
	StepStatus  base.JSONB `gorm:"type:jsonb" sql:"type:jsonb" json:"step_status"`
	SagaStatus  string     `json:"saga_status"`
}
