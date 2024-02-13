package registry

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
	"github.com/spf13/viper"
)

var _ service.ProcessRegistry = (*ProcessRegistry)(nil)

type ProcessRegistry struct {
}

func NewProcessRegistry() *ProcessRegistry {
	return &ProcessRegistry{}
}

func (c ProcessRegistry) DeleteProcess(imageName, version string) error {
	registryHost := viper.GetString(config.RegistryHostKey)
	authSecret := viper.GetString(config.RegistryAuthSecretKey)

	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://"+registryHost+"/v2/"+imageName+"/manifests/"+version, nil)
	if err != nil {
		return err
	}

	basicAuth := base64.StdEncoding.EncodeToString([]byte(authSecret))

	req.Header.Add("Authorization", "Basic "+basicAuth)
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	digest := resp.Header.Get("Docker-Content-Digest")

	req, err = http.NewRequest("DELETE", "http://"+registryHost+"/v2/"+imageName+"/manifests/"+digest, nil)
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
