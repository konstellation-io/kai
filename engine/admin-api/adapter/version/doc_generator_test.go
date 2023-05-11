package version_test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/version"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/stretchr/testify/require"
)

const versionName = "version1234"

func TestHTTPStaticDocGenerator_Generate(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	docFolder, err := os.MkdirTemp("", "test-version-doc")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(docFolder) // clean up

	storageFolder, err := os.MkdirTemp("", "test-api-storage")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(storageFolder) // clean up

	cfg := &config.Config{}
	cfg.Admin.BaseURL = "http://api.local"
	cfg.Admin.StoragePath = storageFolder

	readmeContent := []byte(`
# Example

## Image relative

This is an example:

![relative path image](./img/test.png)

## Image absolute

This is an example:

![absolute path image](https://absolute-url-example)

`)
	if err := os.WriteFile(filepath.Join(docFolder, "README.md"), readmeContent, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	expectedReadmeContent := []byte(`
# Example

## Image relative

This is an example:

![relative path image](http://api.local/static/version/version1234/docs/img/test.png)

## Image absolute

This is an example:

![absolute path image](https://absolute-url-example)

`)

	generator := version.NewHTTPStaticDocGenerator(cfg, logger)
	err = generator.Generate(versionName, docFolder)
	require.Nil(t, err)

	generatedReadme, err := os.ReadFile(path.Join(cfg.Admin.StoragePath, fmt.Sprintf("version/%s/docs/README.md", versionName)))
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, string(expectedReadmeContent), string(generatedReadme))
}

func TestHTTPStaticDocGenerator_GenerateWithNoContent(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mocks.NewMockLogger(ctrl)

	mocks.AddLoggerExpects(logger)

	cfg := &config.Config{}

	generator := version.NewHTTPStaticDocGenerator(cfg, logger)
	err := generator.Generate(versionName, "not-exists-folder")

	require.NotNil(t, err)
}
