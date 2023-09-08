package usecase

import (
	"context"
	"fmt"
	"os"
	"path"
)

func (i *VersionInteractor) storeDashboards(ctx context.Context, dashboardsFolder, product, version string) []error {
	i.logger.Info("Storing dashboards", "product", product, "version", version)

	var errors []error = nil

	d, err := os.ReadDir(dashboardsFolder)
	if err != nil {
		errors = append([]error{fmt.Errorf("error listing dashboards files: %w", err)}, errors...)
	}

	for _, dashboard := range d {
		dashboardPath := path.Join(dashboardsFolder, dashboard.Name())

		err = i.dashboardService.Create(ctx, product, version, dashboardPath)
		if err != nil {
			errors = append([]error{fmt.Errorf("error creating dashboard: %w", err)}, errors...)
			continue
		}
	}

	return errors
}
