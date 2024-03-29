package registry

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
	"github.com/spf13/viper"
)

var _ service.ProcessRegistry = (*ProcessRegistry)(nil)

var (
	ErrFailedGetManifest = fmt.Errorf("failed to get image manifest")
	ErrFailedDeleteImage = fmt.Errorf("failed to delete image")
)

type ProcessRegistry struct {
}

func NewProcessRegistry() *ProcessRegistry {
	return &ProcessRegistry{}
}

func (c ProcessRegistry) DeleteProcess(ctx context.Context, imageName, version string) error {
	registryHost := viper.GetString(config.RegistryHostKey)
	authSecret := viper.GetString(config.RegistryAuthSecretKey)
	basicAuthHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(authSecret))

	getURL := formURL(registryHost, imageName, version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, http.NoBody)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", basicAuthHeader)
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrFailedGetManifest
	}

	digest := resp.Header.Get("Docker-Content-Digest")

	deleteURL := formURL(registryHost, imageName, digest)

	req, err = http.NewRequestWithContext(ctx, http.MethodDelete, deleteURL, http.NoBody)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", basicAuthHeader)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return ErrFailedDeleteImage
	}

	return nil
}

func formURL(registryHost, imageName, manifest string) string {
	if !strings.HasPrefix(registryHost, "http") {
		registryHost = "http://" + registryHost
	}

	return registryHost + "/v2/" + imageName + "/manifests/" + manifest
}
