package usecase

import (
	"context"
	"fmt"
	"os"
	"path"
)

func (i *VersionInteractor) storeDashboards(ctx context.Context, dashboardsFolder string, runtimeID, version string) []error {
	i.logger.Infof("Storing dashboards for version %q in runtime %q", version, runtimeID)

	var errors []error = nil

	d, err := os.ReadDir(dashboardsFolder)
	if err != nil {
		errors = append([]error{fmt.Errorf("error listing dashboards files: %w", err)}, errors...)
	}

	for _, dashboard := range d {
		dashboardPath := path.Join(dashboardsFolder, dashboard.Name())

		err = i.dashboardService.Create(ctx, runtimeID, version, dashboardPath)
		if err != nil {
			errors = append([]error{fmt.Errorf("error creating dashboard: %w", err)}, errors...)
			continue
		}
	}

	return errors
}
