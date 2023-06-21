package version

import (
	"time"

	"github.com/konstellation-io/krt/pkg/krt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type versionDTO struct {
	ID          string            `bson:"_id"`
	Name        string            `bson:"name"`
	Description string            `bson:"description"`
	Config      map[string]string `bson:"config,omitempty"`
	Workflows   []workflowDTO     `bson:"workflows"`

	CreationDate   time.Time `bson:"creationDate"`
	CreationAuthor string    `bson:"creationAuthor"`

	PublicationDate   *time.Time `bson:"publicationDate"`
	PublicationAuthor *string    `bson:"publicationUserId"`

	Status entity.VersionStatus `bson:"status"`

	Errors []string `bson:"errors"`
}

type workflowDTO struct {
	ID        string            `bson:"id"`
	Type      krt.WorkflowType  `bson:"type"`
	Config    map[string]string `bson:"config,omitempty"`
	Processes []processDTO      `bson:"processes"`
	Stream    string            `bson:"-"`
}

type processDTO struct {
	ID            string                 `bson:"id"`
	Name          string                 `bson:"name"`
	Type          krt.ProcessType        `bson:"type"`
	Image         string                 `bson:"image"`
	Replicas      int32                  `bson:"replicas"`
	GPU           bool                   `bson:"gpu"`
	Config        map[string]string      `bson:"config,omitempty"`
	ObjectStore   *processObjectStoreDTO `bson:"objectStore,omitempty"`
	Secrets       map[string]string      `bson:"secrets,omitempty"`
	Subscriptions []string               `bson:"subscriptions"`
	Networking    *processNetworkingDTO  `bson:"networking,omitempty"`
	Status        krt.ProcessStatus      `bson:"-"` // This field value is calculated in k8s
}

type processObjectStoreDTO struct {
	Name  string               `bson:"name"`
	Scope krt.ObjectStoreScope `bson:"scope"`
}

type processNetworkingDTO struct {
	TargetPort          int                    `bson:"targetPort"`
	TargetProtocol      krt.NetworkingProtocol `bson:"targetProtocol"`
	DestinationPort     int                    `bson:"destinationPort"`
	DestinationProtocol krt.NetworkingProtocol `bson:"destinationProtocol"`
}
