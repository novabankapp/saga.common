package payloads

import (
	"github.com/novabankapp/saga.common/domain/base"
	"github.com/novabankapp/saga.common/sagas/bills"
)

func MakePayload() base.JSONB {
	payload := make(base.JSONB)
	payload["type"] = bills.REQUEST

	return payload
}
