package domain

type SagaStepStatus string

const (
	STARTED      SagaStepStatus = "STARTED"
	FAILED                      = "FAILED"
	SUCCEDED                    = "SUCCEEDED"
	COMPENSATING                = "COMPENSATING"
	COMPENSATED                 = "COMPENSATED"
)

type SagaStepState struct {
	Type string `json:"type"`
}
