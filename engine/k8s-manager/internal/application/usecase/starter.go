package usecase

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"golang.org/x/net/context"
)

type VersionStarter struct {
	logger           logr.Logger
	containerStarter service.ContainerStarter
}

func NewVersionStarter(logger logr.Logger, orchStarter service.ContainerStarter) *VersionStarter {
	return &VersionStarter{
		logger,
		orchStarter,
	}
}

func (s *VersionStarter) StartVersion(ctx context.Context, version domain.Version) error {
	s.logger.Info("Running version starter", "product", version.Product, "version", version.Name)

	configName, err := s.containerStarter.CreateVersionConfiguration(ctx, version)
	if err != nil {
		return fmt.Errorf("create version configuration: %w", err)
	}

	for _, workflow := range version.Workflows {
		for _, process := range workflow.Processes {
			err := s.containerStarter.CreateProcess(ctx, service.CreateProcessParams{
				ConfigName: configName,
				Product:    version.Product,
				Version:    version.Name,
				Workflow:   workflow.Name,
				Process:    process,
			})
			if err != nil {
				return fmt.Errorf("create process %q: %w", process.Name, err)
			}

			if process.IsTrigger() && process.Networking != nil {
				err := s.containerStarter.CreateNetwork(ctx, service.CreateNetworkParams{
					Product:  version.Product,
					Version:  version.Name,
					Workflow: workflow.Name,
					Process:  process,
				})
				if err != nil {
					return fmt.Errorf("create network: %w", err)
				}
			}
		}
	}

	return nil
}
