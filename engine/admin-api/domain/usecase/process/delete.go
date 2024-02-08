package process

import (
	"context"
	"fmt"
	"net/http"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/spf13/viper"
)

// TODO
// Delete process registry from the database at the end of the function
// Define deleteProcessRegistry mutation then add this func to the resolver
// Ensure config for both registry and admin-api are set up properly
// Add a test for this function
// Try it out in a local environment

func (ps *Service) DeleteProcess(
	ctx context.Context,
	user *entity.User,
	opts DeleteProcessOpts,
) error {
	ps.logger.Info("Deleting process", "Product", opts.Product, "Version", opts.Version, "Process", opts.Process)

	if err := opts.Validate(); err != nil {
		return err
	}

	if err := ps.checkGrants(user, opts.IsPublic, opts.Product); err != nil {
		return err
	}

	if err := ps.checkIfProcessExists(opts); err != nil {
		return err
	}

	if err := ps.deleteImageTag(opts); err != nil {
		return err
	}

	return nil
}

func (ps *Service) checkIfProcessExists(opts DeleteProcessOpts) error {
	processID := ps.getProcessID(opts.Product, opts.Version, opts.Process)

	filter := repository.SearchFilter{
		ProcessID: processID,
	}

	if opts.IsPublic {
		processes, err := ps.processRepository.GlobalSearch(context.Background(), filter)
		if err != nil {
			return err
		}
		if len(processes) != 1 {
			return ErrRegisteredProcessNotFound
		}
	} else {
		processes, err := ps.processRepository.SearchByProduct(context.Background(), opts.Product, filter)
		if err != nil {
			return err
		}
		if len(processes) != 1 {
			return ErrRegisteredProcessNotFound
		}
	}

	return nil
}

func (ps *Service) deleteImageTag(opts DeleteProcessOpts) error {
	registryURL := viper.GetString(config.ImageRegistryURLKey)
	authSecret := viper.GetString(config.ImageRegistryAuthSecretKey)
	repositoryName := ps.getRepositoryName(opts.Product, opts.Process)

	client := &http.Client{}

	req, err := http.NewRequest("GET", registryURL+"/v2/"+repositoryName+"/manifests/"+opts.Version, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+authSecret)
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	digest := resp.Header.Get("Docker-Content-Digest")

	req, err = http.NewRequest("DELETE", registryURL+"/v2/"+repositoryName+"/manifests/"+digest, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+authSecret)

	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete image: %s", resp.Status)
	}

	return nil
}
