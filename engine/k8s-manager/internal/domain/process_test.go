//go:build unit

package domain_test

import (
	"testing"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestProcessTypeToString(t *testing.T) {
	testCases := []struct {
		name        string
		processType domain.ProcessType
		expected    string
	}{
		{"task process type", domain.TaskProcessType, "task"},
		{"trigger process type", domain.TriggerProcessType, "trigger"},
		{"unknown process type", domain.ExitProcessType, "exit"},
		{"unknown process type", domain.UnknownProcessType, "unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.processType.ToString())
		})
	}
}

func TestProcess_IsTrigger(t *testing.T) {
	testCases := []struct {
		name      string
		process   *domain.Process
		isTrigger bool
	}{
		{
			"task process",
			testhelpers.NewProcessBuilder().WithType(domain.TaskProcessType).Build(),
			false,
		},
		{
			"trigger process",
			testhelpers.NewProcessBuilder().WithType(domain.TriggerProcessType).Build(),
			true,
		},
		{
			"exit process",
			testhelpers.NewProcessBuilder().WithType(domain.ExitProcessType).Build(),
			false,
		},
		{
			"unknown process",
			testhelpers.NewProcessBuilder().WithType(domain.UnknownProcessType).Build(),
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isTrigger, tc.process.IsTrigger())
		})
	}
}
