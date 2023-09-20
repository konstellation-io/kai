//go:build unit

package versionrepository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

var (
	userID          = "admin"
	creationDate    = time.Now().Add(-time.Hour)
	publicationDate = time.Now().Add(-time.Minute)
)

var domainVersion = &entity.Version{
	Tag:         "1.0.0",
	Description: "test description",
	Config: []entity.ConfigurationVariable{
		{
			Key:   "key1",
			Value: "value1",
		},
		{
			Key:   "key2",
			Value: "value2",
		},
	},

	CreationDate:      creationDate,
	CreationAuthor:    userID,
	PublicationDate:   &publicationDate,
	PublicationAuthor: &userID,
	Status:            entity.VersionStatusPublished,

	Workflows: []entity.Workflow{
		{
			Name: "workflow1",
			Type: entity.WorkflowTypeTraining,
			Config: []entity.ConfigurationVariable{
				{
					Key:   "key1",
					Value: "value1",
				},
			},
			Processes: []entity.Process{
				{
					Name:     "process1",
					Type:     entity.ProcessTypeTrigger,
					Image:    "image1",
					Replicas: 1,
					Config: []entity.ConfigurationVariable{
						{
							Key:   "key1",
							Value: "value1",
						},
					},
					ObjectStore: &entity.ProcessObjectStore{
						Name:  "objectStore1",
						Scope: entity.ObjectStoreScopeProduct,
					},
					Secrets: []entity.ConfigurationVariable{
						{
							Key:   "key1",
							Value: "value1",
						},
					},
					Subscriptions: []string{"subscription1", "subscription2"},
					Networking: &entity.ProcessNetworking{
						TargetPort:      8080,
						DestinationPort: 8080,
						Protocol:        "TCP",
					},
					ResourceLimits: &entity.ProcessResourceLimits{
						CPU: &entity.ResourceLimit{
							Request: "0.5",
							Limit:   "1.5",
						},
						Memory: &entity.ResourceLimit{
							Request: "400M",
							Limit:   "600M",
						},
					},
				},
				{
					Name:          "process2",
					Type:          entity.ProcessTypeTask,
					Image:         "image2",
					Replicas:      2,
					GPU:           true,
					Subscriptions: []string{"subscription3", "subscription4"},
				},
				{
					Name:          "process3",
					Type:          entity.ProcessTypeExit,
					Image:         "image3",
					Subscriptions: []string{"subscription5", "subscription6"},
				},
			},
		},
	},
}

var DTOVersion = &versionDTO{
	Tag:         "1.0.0",
	Description: "test description",
	Config: []configurationVariableDTO{
		{
			Key:   "key1",
			Value: "value1",
		},
		{
			Key:   "key2",
			Value: "value2",
		},
	},

	CreationDate:      creationDate,
	CreationAuthor:    userID,
	PublicationDate:   &publicationDate,
	PublicationAuthor: &userID,
	Status:            entity.VersionStatusPublished.String(),

	Workflows: []workflowDTO{
		{
			Name: "workflow1",
			Type: entity.WorkflowTypeTraining.String(),
			Config: []configurationVariableDTO{
				{
					Key:   "key1",
					Value: "value1",
				},
			},
			Processes: []processDTO{
				{
					Name:     "process1",
					Type:     entity.ProcessTypeTrigger.String(),
					Image:    "image1",
					Replicas: 1,
					Config: []configurationVariableDTO{
						{
							Key:   "key1",
							Value: "value1",
						},
					},
					ObjectStore: &processObjectStoreDTO{
						Name:  "objectStore1",
						Scope: entity.ObjectStoreScopeProduct.String(),
					},
					Secrets: []configurationVariableDTO{
						{
							Key:   "key1",
							Value: "value1",
						},
					},
					Subscriptions: []string{"subscription1", "subscription2"},
					Networking: &processNetworkingDTO{
						TargetPort:      8080,
						DestinationPort: 8080,
						Protocol:        "TCP",
					},
					ResourceLimits: &processResourceLimitsDTO{
						CPU: &resourceLimitDTO{
							Request: "0.5",
							Limit:   "1.5",
						},
						Memory: &resourceLimitDTO{
							Request: "400M",
							Limit:   "600M",
						},
					},
				},
				{
					Name:          "process2",
					Type:          entity.ProcessTypeTask.String(),
					Image:         "image2",
					Replicas:      2,
					GPU:           true,
					Subscriptions: []string{"subscription3", "subscription4"},
				},
				{
					Name:          "process3",
					Type:          entity.ProcessTypeExit.String(),
					Image:         "image3",
					Subscriptions: []string{"subscription5", "subscription6"},
				},
			},
		},
	},
}

func TestMapDTOToEntity(t *testing.T) {
	obtainedDomainVersion := mapDTOToEntity(DTOVersion)
	assert.Equal(t, domainVersion, obtainedDomainVersion)
}

func TestMapEntityToDTO(t *testing.T) {
	obtainedDTOVersion := mapEntityToDTO(domainVersion)
	assert.Equal(t, DTOVersion, obtainedDTOVersion)
}
