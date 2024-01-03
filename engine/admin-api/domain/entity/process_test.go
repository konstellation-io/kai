package entity_test

import (
	"testing"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestProcessType_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		processType entity.ProcessType
		expecterErr error
	}{
		{"valid trigger process type", entity.ProcessTypeTrigger, nil},
		{"valid task process type", entity.ProcessTypeTask, nil},
		{"valid exit process type", entity.ProcessTypeExit, nil},
		{"invalid process type", entity.ProcessType("invalid-process-type"), entity.ErrInvalidProcessType},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.ErrorIs(t, tc.processType.Validate(), tc.expecterErr)
		})
	}
}
