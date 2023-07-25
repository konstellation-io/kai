package usecase_test

import (
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/stretchr/testify/assert"
)

func TestRegisterProcess(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	processService := usecase.NewProcessService(logger)

	err := processService.RegisterProcess()
	assert.NoError(t, err)
}
