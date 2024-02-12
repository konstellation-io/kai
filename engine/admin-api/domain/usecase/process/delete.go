package process

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/spf13/viper"
)

// TODO
// Current config doesnt work, admin api is not able to get up
// Talk about usage of basic user pwd auth rather than token use
// Ensure config for both registry and admin-api are set up properly
// Add a test for this function
// Try it out in a local environment

func (ps *Service) DeleteProcess(
	ctx context.Context,
	user *entity.User,
	opts DeleteProcessOpts,
) (string, error) {
	ps.logger.Info("Deleting process", "Product", opts.Product, "Version", opts.Version, "Process", opts.Process, "IsPublic", opts.IsPublic)

	if err := opts.Validate(); err != nil {
		return "", err
	}

	if err := ps.checkDeleteGrants(user, opts.IsPublic, opts.Product); err != nil {
		return "", err
	}

	scope := ps.getProcessRegisterScope(opts.IsPublic, opts.Product)
	processID := ps.getProcessID(scope, opts.Process, opts.Version)

	_, err := ps.processRepository.GetByID(ctx, scope, processID)
	if err != nil {
		return "", err
	}

	if err := ps.deleteImageTag(opts, scope); err != nil {
		return "", err
	}

	if err := ps.processRepository.Delete(ctx, scope, processID); err != nil {
		return "", err
	}

	return processID, nil
}

func (ps *Service) deleteImageTag(opts DeleteProcessOpts, scope string) error {
	registryHost := viper.GetString(config.RegistryHostKey)
	authSecret := viper.GetString(config.RegistryAuthSecretKey)
	repositoryName := ps.getRepositoryName(scope, opts.Process)

	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://"+registryHost+"/v2/"+repositoryName+"/manifests/"+opts.Version, nil)
	if err != nil {
		return err
	}

	basicAuth := base64.StdEncoding.EncodeToString([]byte(authSecret))

	//username := "user"
	//password := "password"
	//basicAuth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	req.Header.Add("Authorization", "Basic "+authSecret)

	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	digest := resp.Header.Get("Docker-Content-Digest")

	req, err = http.NewRequest("DELETE", "http://"+registryHost+"/v2/"+repositoryName+"/manifests/"+digest, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Basic "+basicAuth)

	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to delete image: %s", resp.Status)
	}

	return nil
}
