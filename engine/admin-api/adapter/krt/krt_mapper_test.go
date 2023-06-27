//go:build unit

package krt_test

import (
	"testing"

	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/stretchr/testify/require"

	krtapp "github.com/konstellation-io/kai/engine/admin-api/adapter/krt"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

func getExpectedVersion() *entity.Version {
	commonObjectStore := &entity.ProcessObjectStore{
		Name:  "emails",
		Scope: "workflow",
	}

	return &entity.Version{
		ID:          "", // ID is not defined in KRT YAML
		Name:        "email-classificator",
		Description: "Email classificator for branching features.",
		Version:     "v1.0.0",
		Config: []entity.ConfigurationVariable{
			{
				Key:   "keyA",
				Value: "value1",
			},
		},
		Workflows: []entity.Workflow{
			{
				Name: "go-classificator",
				Type: "data",
				Config: []entity.ConfigurationVariable{
					{
						Key:   "keyA",
						Value: "value1",
					},
				},
				Processes: []entity.Process{
					{
						Name:          "entrypoint",
						Type:          "trigger",
						Image:         "konstellation/kai-grpc-trigger:latest",
						Replicas:      krt.DefaultNumberOfReplicas,
						GPU:           krt.DefaultGPUValue,
						Subscriptions: []string{"exitpoint"},
						Networking: &entity.ProcessNetworking{
							TargetPort:      9000,
							DestinationPort: 9000,
							Protocol:        "TCP",
						},
					},
					{
						Name:          "etl",
						Type:          "task",
						Image:         "konstellation/kai-etl-task:latest",
						Replicas:      krt.DefaultNumberOfReplicas,
						GPU:           krt.DefaultGPUValue,
						ObjectStore:   commonObjectStore,
						Subscriptions: []string{"entrypoint"},
					},
					{
						Name:          "email-classificator",
						Type:          "task",
						Image:         "konstellation/kai-ec-task:latest",
						Replicas:      krt.DefaultNumberOfReplicas,
						GPU:           krt.DefaultGPUValue,
						ObjectStore:   commonObjectStore,
						Subscriptions: []string{"etl"},
					},
					{
						Name:          "exitpoint",
						Type:          "exit",
						Image:         "konstellation/kai-exitpoint:latest",
						Replicas:      krt.DefaultNumberOfReplicas,
						GPU:           krt.DefaultGPUValue,
						ObjectStore:   commonObjectStore,
						Subscriptions: []string{"etl", "stats-storer"},
					},
				},
			},
		},
	}
}

func TestKrtYmlMapper(t *testing.T) {
	// GIVEN an expected version
	expectedVersion := getExpectedVersion()

	// GIVEN a KRT YAML file with a valid format
	krtYml, err := krtapp.ParseFile("../../test_assets/classificator_krt.yaml")
	require.NoError(t, err)

	err = krtYml.Validate()
	require.NoError(t, err)

	// WHEN the KRT YAML is mapped to a Version entity
	version := krtapp.MapKrtYamlToVersion(krtYml)

	// THEN the Version entity is the expected
	require.Equal(t, expectedVersion, version)
}
