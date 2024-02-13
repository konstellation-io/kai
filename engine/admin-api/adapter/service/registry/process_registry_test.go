//go:build integration

package registry_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/registry"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteProcess(t *testing.T) {
	const (
		imageName = "productID_processID"
		version   = "versionID"
		basicAuth = "user:password"
		digest    = "sha256:1234567890"
	)

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		callCount++

		switch callCount {
		case 1:
			expectedURL := `/v2/productID_processID/manifests/versionID`

			actualQuery, err := url.QueryUnescape(req.URL.String())
			require.NoError(t, err)

			assert.Equal(t, expectedURL, actualQuery)

			rw.Header().Set("Docker-Content-Digest", digest)
			rw.Write([]byte(`{}`))

		case 2:
			expectedURL := `/v2/productID_processID/manifests/sha256:1234567890`

			actualQuery, err := url.QueryUnescape(req.URL.String())
			require.NoError(t, err)

			assert.Equal(t, expectedURL, actualQuery)

			rw.WriteHeader(http.StatusAccepted)

		default:
			t.Error("Unexpected call")
		}
	}))
	defer server.Close()

	viper.Set(config.RegistryHostKey, server.URL)
	viper.Set(config.RegistryAuthSecretKey, basicAuth)

	processRegistry := registry.NewProcessRegistry()

	err := processRegistry.DeleteProcess(imageName, version)
	require.NoError(t, err)
}

func TestDeleteProcess_GetManifestError(t *testing.T) {
	const (
		imageName = "productID_processID"
		version   = "versionID"
		basicAuth = "user:password"
	)

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		callCount++

		switch callCount {
		case 1:
			expectedURL := `/v2/productID_processID/manifests/versionID`

			actualQuery, err := url.QueryUnescape(req.URL.String())
			require.NoError(t, err)

			assert.Equal(t, expectedURL, actualQuery)

			rw.WriteHeader(http.StatusInternalServerError)

		default:
			t.Error("Unexpected call")
		}
	}))
	defer server.Close()

	viper.Set(config.RegistryHostKey, server.URL)
	viper.Set(config.RegistryAuthSecretKey, basicAuth)

	processRegistry := registry.NewProcessRegistry()

	err := processRegistry.DeleteProcess(imageName, version)
	require.Error(t, err)
}

func TestDeleteProcess_DeleteManifestError(t *testing.T) {
	const (
		imageName = "productID_processID"
		version   = "version"
		basicAuth = "user:password"
		digest    = "sha256:1234567890"
	)

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		callCount++

		switch callCount {
		case 1:
			expectedURL := `/v2/productID_processID/manifests/version`

			actualQuery, err := url.QueryUnescape(req.URL.String())
			require.NoError(t, err)

			assert.Equal(t, expectedURL, actualQuery)

			rw.Header().Set("Docker-Content-Digest", digest)
			rw.Write([]byte(`{}`))

		case 2:
			expectedURL := `/v2/productID_processID/manifests/sha256:1234567890`

			actualQuery, err := url.QueryUnescape(req.URL.String())
			require.NoError(t, err)

			assert.Equal(t, expectedURL, actualQuery)

			rw.WriteHeader(http.StatusInternalServerError)

		default:
			t.Error("Unexpected call")
		}
	}))
	defer server.Close()

	viper.Set(config.RegistryHostKey, server.URL)
	viper.Set(config.RegistryAuthSecretKey, basicAuth)

	processRegistry := registry.NewProcessRegistry()

	err := processRegistry.DeleteProcess(imageName, version)
	require.Error(t, err)
}
