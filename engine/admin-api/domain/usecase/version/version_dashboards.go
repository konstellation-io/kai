package version

import (
	"context"
	"fmt"
	"os"
	"path"
)

func (h *Handler) storeDashboards(ctx context.Context, dashboardsFolder, product, version string) []error {
	h.logger.Info("Storing dashboards", "product", product, "version", version)

	var errors []error = nil

	d, err := os.ReadDir(dashboardsFolder)
	if err != nil {
		errors = append([]error{fmt.Errorf("error listing dashboards files: %w", err)}, errors...)
	}

	for _, dashboard := range d {
		dashboardPath := path.Join(dashboardsFolder, dashboard.Name())

		err = h.dashboardService.Create(ctx, product, version, dashboardPath)
		if err != nil {
			errors = append([]error{fmt.Errorf("error creating dashboard: %w", err)}, errors...)
			continue
		}
	}

	return errors
}
