package usecase

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/pkg/compensator"
	"golang.org/x/net/context"
)

type VersionStarter struct {
	logger           logr.Logger
	containerService service.ContainerService
}

func NewVersionStarter(logger logr.Logger, orchStarter service.ContainerService) VersionStarterService {
	return &VersionStarter{
		logger,
		orchStarter,
	}
}

func (s *VersionStarter) StartVersion(ctx context.Context, version domain.Version) error {
	s.logger.Info("Running version starter", "product", version.Product, "version", version.Tag)

	compensations := compensator.New()

	if err := s.createVersionResources(ctx, version, compensations); err != nil {
		if compensationsErrors := compensations.Execute(); compensationsErrors != nil {
			s.logger.Error(compensationsErrors, "Error(s) executing compensations")
		}

		return err
	}

	return nil
}

func (s *VersionStarter) createVersionResources(ctx context.Context, version domain.Version, compensations *compensator.Compensator) error {
	configName, err := s.containerService.CreateVersionConfiguration(ctx, version)
	if err != nil {
		return fmt.Errorf("create version configuration: %w", err)
	}

	compensations.AddCompensation(s.deleteConfigurationFunc(version))
	compensations.AddCompensation(s.deleteProcessesFunc(version))

	for _, workflow := range version.Workflows {
		for _, process := range workflow.Processes {
			err := s.containerService.CreateProcess(ctx, service.CreateProcessParams{
				ConfigName: configName,
				Product:    version.Product,
				Version:    version.Tag,
				Workflow:   workflow.Name,
				Process:    process,
			})
			if err != nil {
				return fmt.Errorf("create process %q: %w", process.Name, err)
			}

			if process.IsTrigger() && process.Networking != nil {
				err := s.containerService.CreateNetwork(ctx, service.CreateNetworkParams{
					Product:  version.Product,
					Version:  version.Tag,
					Workflow: workflow.Name,
					Process:  process,
				})
				if err != nil {
					return fmt.Errorf("create network: %w", err)
				}
			}

			compensations.AddCompensation(s.deleteNetworkFunc(version))
		}
	}

	err = s.containerService.WaitProcesses(ctx, version)
	if err != nil {
		return err
	}

	return nil
}

func (s *VersionStarter) deleteConfigurationFunc(version domain.Version) compensator.Compensation {
	return func() error {
		return s.containerService.DeleteConfiguration(context.Background(), version.Product, version.Tag)
	}
}

func (s *VersionStarter) deleteProcessesFunc(version domain.Version) compensator.Compensation {
	return func() error {
		return s.containerService.DeleteProcesses(context.Background(), version.Product, version.Tag)
	}
}

func (s *VersionStarter) deleteNetworkFunc(version domain.Version) compensator.Compensation {
	return func() error {
		return s.containerService.DeleteNetwork(context.Background(), version.Product, version.Tag)
	}
}
