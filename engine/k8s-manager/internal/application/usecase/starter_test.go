//go:build unit

package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service/mocks"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/usecase"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStartVersion(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	containerSvc := mocks.NewContainerService(t)

	version := testhelpers.NewVersionBuilder().Build()

	configName := "test-config-name"

	containerSvc.EXPECT().
		CreateVersionConfiguration(mock.Anything, version).
		Return(configName, nil).
		Once()

	mockCreateProcess(t, containerSvc, configName, version, *version.Workflows[0].Processes[0])

	starter := usecase.NewVersionStarter(logger, containerSvc)

	ctx := context.Background()
	err := starter.StartVersion(ctx, version)
	assert.NoError(t, err)
}

func TestStartVersion_WithMultipleProcesses(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	containerSvc := mocks.NewContainerService(t)

	processes := []*domain.Process{
		testhelpers.NewProcessBuilder().WithID("test-process-1").Build(),
		testhelpers.NewProcessBuilder().WithID("test-process-2").Build(),
		testhelpers.NewProcessBuilder().WithID("test-process-3").Build(),
	}

	workflows := []*domain.Workflow{
		testhelpers.NewWorkflowBuilder().WithProcesses(processes).Build(),
	}

	version := testhelpers.NewVersionBuilder().WithWorkflows(workflows).Build()

	configName := "test-config-name"

	containerSvc.EXPECT().
		CreateVersionConfiguration(mock.Anything, version).
		Return(configName, nil).
		Once()

	for _, w := range version.Workflows {
		for _, p := range w.Processes {
			mockCreateProcess(t, containerSvc, configName, version, *p)
		}
	}

	starter := usecase.NewVersionStarter(logger, containerSvc)

	ctx := context.Background()
	err := starter.StartVersion(ctx, version)
	assert.NoError(t, err)
}

func TestStartVersion_WithNetworking(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	containerSvc := mocks.NewContainerService(t)

	processes := []*domain.Process{
		testhelpers.NewProcessBuilder().
			WithType(domain.TriggerProcessType).
			WithNetworking(domain.Networking{
				SourcePort: 80,
				TargetPort: 80,
				Protocol:   "TCP",
			}).
			Build(),
	}
	workflows := []*domain.Workflow{
		testhelpers.NewWorkflowBuilder().WithProcesses(processes).Build(),
	}

	version := testhelpers.NewVersionBuilder().WithWorkflows(workflows).Build()

	configName := "test-config-name"

	containerSvc.EXPECT().
		CreateVersionConfiguration(mock.Anything, version).
		Return(configName, nil).
		Once()

	createProcessParams := service.CreateProcessParams{
		ConfigName: configName,
		Product:    version.Product,
		Version:    version.ID,
		Workflow:   version.Workflows[0].ID,
		Process:    version.Workflows[0].Processes[0],
	}

	containerSvc.EXPECT().
		CreateProcess(mock.Anything, createProcessParams).
		Return(nil).
		Once()

	createNetworkParams := service.CreateNetworkParams{
		Product:  version.Product,
		Version:  version.ID,
		Workflow: workflows[0].ID,
		Process:  processes[0],
	}

	containerSvc.EXPECT().
		CreateNetwork(mock.Anything, createNetworkParams).
		Return(nil).
		Once()

	starter := usecase.NewVersionStarter(logger, containerSvc)

	ctx := context.Background()
	err := starter.StartVersion(ctx, version)
	assert.NoError(t, err)
}

func TestStartVersion_ErrorCreatingConfig(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	containerSvc := mocks.NewContainerService(t)

	version := testhelpers.NewVersionBuilder().Build()

	expectedErr := errors.New("error creating configuration")

	configName := "test-config-name"

	containerSvc.EXPECT().
		CreateVersionConfiguration(mock.Anything, version).
		Return(configName, expectedErr).
		Once()

	starter := usecase.NewVersionStarter(logger, containerSvc)

	ctx := context.Background()
	err := starter.StartVersion(ctx, version)
	assert.ErrorIs(t, err, expectedErr)
}

func TestStartVersion_ErrorCreatingProcess(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: 0})
	containerSvc := mocks.NewContainerService(t)

	version := testhelpers.NewVersionBuilder().Build()

	expectedErr := errors.New("error creating process")

	configName := "test-config-name"

	containerSvc.EXPECT().
		CreateVersionConfiguration(mock.Anything, version).
		Return(configName, nil).
		Once()

	createProcessParams := service.CreateProcessParams{
		ConfigName: configName,
		Product:    version.Product,
		Version:    version.ID,
		Workflow:   version.Workflows[0].ID,
		Process:    version.Workflows[0].Processes[0],
	}

	containerSvc.EXPECT().
		CreateProcess(mock.Anything, createProcessParams).
		Return(expectedErr).
		Once()

	starter := usecase.NewVersionStarter(logger, containerSvc)

	ctx := context.Background()
	err := starter.StartVersion(ctx, version)
	assert.ErrorIs(t, err, expectedErr)
}

func TestStartVersion_ErrorCreatingNetwork(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	containerSvc := mocks.NewContainerService(t)

	expectedErr := errors.New("error creating network")

	version := getVersionWithNetrking(t)

	configName := "test-config-name"

	containerSvc.EXPECT().
		CreateVersionConfiguration(mock.Anything, version).
		Return(configName, nil).
		Once()

	createProcessParams := service.CreateProcessParams{
		ConfigName: configName,
		Product:    version.Product,
		Version:    version.ID,
		Workflow:   version.Workflows[0].ID,
		Process:    version.Workflows[0].Processes[0],
	}

	containerSvc.EXPECT().
		CreateProcess(mock.Anything, createProcessParams).
		Return(nil).
		Once()

	createNetworkParams := service.CreateNetworkParams{
		Product:  version.Product,
		Version:  version.ID,
		Workflow: version.Workflows[0].ID,
		Process:  version.Workflows[0].Processes[0],
	}

	containerSvc.EXPECT().
		CreateNetwork(mock.Anything, createNetworkParams).
		Return(expectedErr).
		Once()

	starter := usecase.NewVersionStarter(logger, containerSvc)

	ctx := context.Background()
	err := starter.StartVersion(ctx, version)
	assert.ErrorIs(t, err, expectedErr)
}

func getVersionWithNetrking(t *testing.T) domain.Version {
	t.Helper()
	processes := []*domain.Process{
		testhelpers.NewProcessBuilder().
			WithType(domain.TriggerProcessType).
			WithNetworking(domain.Networking{
				SourcePort: 80,
				TargetPort: 80,
				Protocol:   "TCP",
			}).
			Build(),
	}
	workflows := []*domain.Workflow{
		testhelpers.NewWorkflowBuilder().WithProcesses(processes).Build(),
	}

	return testhelpers.NewVersionBuilder().WithWorkflows(workflows).Build()
}

func mockCreateProcess(
	t *testing.T,
	containerSvc *mocks.ContainerService,
	configName string,
	version domain.Version,
	process domain.Process,
) {
	t.Helper()
	createProcessParams := service.CreateProcessParams{
		ConfigName: configName,
		Product:    version.Product,
		Version:    version.ID,
		Workflow:   version.Workflows[0].ID,
		Process:    &process,
	}

	containerSvc.EXPECT().
		CreateProcess(mock.Anything, createProcessParams).
		Return(nil).
		Once()
}
