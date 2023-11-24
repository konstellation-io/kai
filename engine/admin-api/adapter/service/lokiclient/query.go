package lokiclient

import (
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

const (
	productIDKey    = "product_id"
	versionTagKey   = "version_tag"
	workflowNameKey = "workflow_name"
	processNameKey  = "process_name"
	requestIDKey    = "request_id"
)

func getQuery(lf entity.LogFilters) string {
	const (
		madatoryQueryPart = "{%s=\"%s\", %s=\"%s\""
		optionalQueryPart = ", %s=\"%s\""
	)

	// mandatory part of the query
	query := fmt.Sprintf(madatoryQueryPart, productIDKey, lf.ProductID, versionTagKey, lf.VersionID)

	// optional part of the query
	if lf.WorkflowName != "" {
		query += fmt.Sprintf(optionalQueryPart, workflowNameKey, lf.WorkflowName)
	}

	if lf.ProcessName != "" {
		query += fmt.Sprintf(optionalQueryPart, processNameKey, lf.ProcessName)
	}

	if lf.RequestID != "" {
		query += fmt.Sprintf(optionalQueryPart, requestIDKey, lf.RequestID)
	}

	query += "}"

	return query
}
