package version

import (
	"time"
)

type ConfigurationVariable struct {
	Key   string
	Value string
}

type versionDTO struct {
	ID          string                  `bson:"_id"`
	Name        string                  `bson:"name"`
	Description string                  `bson:"description"`
	Config      []ConfigurationVariable `bson:"config,omitempty"`
	Workflows   []workflowDTO           `bson:"workflows"`

	CreationDate   time.Time `bson:"creationDate"`
	CreationAuthor string    `bson:"creationAuthor"`

	PublicationDate   *time.Time `bson:"publicationDate"`
	PublicationAuthor *string    `bson:"publicationAuthor"`

	Status string `bson:"status"`

	Errors []string `bson:"errors"`
}

type workflowDTO struct {
	ID        string                  `bson:"id"`
	Name      string                  `bson:"name"`
	Type      string                  `bson:"type"`
	Config    []ConfigurationVariable `bson:"config,omitempty"`
	Processes []processDTO            `bson:"processes"`
}

type processDTO struct {
	ID            string                  `bson:"id"`
	Name          string                  `bson:"name"`
	Type          string                  `bson:"type"`
	Image         string                  `bson:"image"`
	Replicas      int32                   `bson:"replicas"`
	GPU           bool                    `bson:"gpu"`
	Config        []ConfigurationVariable `bson:"config,omitempty"`
	ObjectStore   *processObjectStoreDTO  `bson:"objectStore,omitempty"`
	Secrets       []ConfigurationVariable `bson:"secrets,omitempty"`
	Subscriptions []string                `bson:"subscriptions"`
	Networking    *processNetworkingDTO   `bson:"networking,omitempty"`
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
