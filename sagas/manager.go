package sagas

import (
	"context"
	"github.com/novabankapp/saga.common/domain/base"
)

type Manager interface {
	Begin(ctx context.Context, sagaType string, sagaSteps []string, payload base.JSONB)
	Find(ctx context.Context, id string)
}
