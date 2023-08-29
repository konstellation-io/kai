package versionrepository

import (
	"time"
)

type configurationVariableDTO struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

type versionDTO struct {
	ID          string                     `bson:"_id"`
	Tag         string                     `bson:"tag"`
	Description string                     `bson:"description"`
	Config      []configurationVariableDTO `bson:"config,omitempty"`
	Workflows   []workflowDTO              `bson:"workflows"`

	CreationDate   time.Time `bson:"creationDate"`
	CreationAuthor string    `bson:"creationAuthor"`

	PublicationDate   *time.Time `bson:"publicationDate"`
	PublicationAuthor *string    `bson:"publicationAuthor"`

	Status string `bson:"status"`

	Errors []string `bson:"errors"`
}

type workflowDTO struct {
	ID        string                     `bson:"id"`
	Name      string                     `bson:"name"`
	Type      string                     `bson:"type"`
	Config    []configurationVariableDTO `bson:"config,omitempty"`
	Processes []processDTO               `bson:"processes"`
}

type processDTO struct {
	ID             string                     `bson:"id"`
	Name           string                     `bson:"name"`
	Type           string                     `bson:"type"`
	Image          string                     `bson:"image"`
	Replicas       int32                      `bson:"replicas"`
	GPU            bool                       `bson:"gpu"`
	Config         []configurationVariableDTO `bson:"config,omitempty"`
	ObjectStore    *processObjectStoreDTO     `bson:"objectStore,omitempty"`
	Secrets        []configurationVariableDTO `bson:"secrets,omitempty"`
	Subscriptions  []string                   `bson:"subscriptions"`
	Networking     *processNetworkingDTO      `bson:"networking,omitempty"`
	ResourceLimits *processResourceLimitsDTO  `bson:"resourceLimits,omitempty"`
}

type processObjectStoreDTO struct {
	Name  string `bson:"name"`
	Scope string `bson:"scope"`
}

type processNetworkingDTO struct {
	TargetPort      int    `bson:"targetPort"`
	DestinationPort int    `bson:"destinationPort"`
	Protocol        string `bson:"protocol"`
}

type resourceLimitDTO struct {
	Request string `bson:"request"`
	Limit   string `bson:"limit"`
}

type processResourceLimitsDTO struct {
	CPU    *resourceLimitDTO `bson:"cpu,omitempty"`
	Memory *resourceLimitDTO `bson:"memory,omitempty"`
}
