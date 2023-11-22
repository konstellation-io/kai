package lokiclient

import (
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

func setFilterDefaults(lf *entity.LogFilters) {
	if lf.Limit == 0 {
		lf.Limit = 100
	}
}

func getQuery(lf entity.LogFilters) string {
	setFilterDefaults(&lf)

	query := fmt.Sprintf("{product_id=\"%s\"}", lf.ProductID)

	if lf.VersionID != "" {
		query += fmt.Sprintf(",version_id=\"%s\"", lf.VersionID)
	}

	return query
}

//query := `{product_id="rest-poc-2",version_id="v1.0.1"}`
//query := `{version="v1.0.0",name="mysql-backup"}`
