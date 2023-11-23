package lokiclient

import (
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

var (
	productIDKey    = "product_id"
	versionTagKey   = "version_tag"
	workflowNameKey = "workflow_name"
	processNameKey  = "process_name"
	requestIDKey    = "request_id"
)

func setFilterDefaults(lf *entity.LogFilters) {
	if lf.Limit == 0 {
		lf.Limit = 100
	}
}

func getQuery(lf entity.LogFilters) string {
	setFilterDefaults(&lf)

	const madatoryQueryPart = "{%s=\"%s\", %s=\"%s\""
	const optionalQueryPart = ", %s=\"%s\""

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
