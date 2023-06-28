//go:build unit

package domain_test

import (
	"testing"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	testCases := []struct {
		name                      string
		version                   domain.Version
		expectedAmountOfProcesses int
	}{
		{
			name: "version with no processes",
			version: testhelpers.NewVersionBuilder().
				WithWorkflows([]*domain.Workflow{}).
				Build(),
			expectedAmountOfProcesses: 0,
		},
		{
			name: "version with one worklfow with 3 nodes",
			version: testhelpers.NewVersionBuilder().
				WithWorkflows([]*domain.Workflow{workflowWithThreeProcesses()}).
				Build(),
			expectedAmountOfProcesses: 3,
		},
		{
			name: "version with 3 worklfow with 3 nodes",
			version: testhelpers.NewVersionBuilder().
				WithWorkflows([]*domain.Workflow{
					workflowWithThreeProcesses(),
					workflowWithThreeProcesses(),
					workflowWithThreeProcesses(),
				}).
				Build(),
			expectedAmountOfProcesses: 9,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedAmountOfProcesses, tc.version.GetAmountOfProcesses())
		})
	}
}

func workflowWithThreeProcesses() *domain.Workflow {
	return testhelpers.NewWorkflowBuilder().
		WithProcesses([]*domain.Process{
			testhelpers.NewProcessBuilder().Build(),
			testhelpers.NewProcessBuilder().Build(),
			testhelpers.NewProcessBuilder().Build(),
		}).Build()
}
