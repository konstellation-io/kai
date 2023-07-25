package registry

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/containers/buildah"
	"github.com/containers/buildah/define"
	"github.com/containers/buildah/imagebuildah"
	"github.com/containers/image/v5/transports"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"github.com/containers/storage"
	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/compression"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/spf13/viper"
)

const _defaultDockerfile = "Dockerfile"

type ProcessRegistry struct {
	logger logr.Logger
}

func NewProcessRegistry(logger logr.Logger) *ProcessRegistry {
	return &ProcessRegistry{
		logger: logger,
	}
}

func (pr *ProcessRegistry) RegisterProcess(ctx context.Context, product, version, process string, compressedSrc io.Reader) (string, error) {
	tmpDir, err := os.MkdirTemp("", "process-*")
	if err != nil {
		return "", fmt.Errorf("creating temporal directory: %w", err)
	}

	//defer func() {
	//	err := os.RemoveAll(tmpDir)
	//	if err != nil {
	//		pr.logger.Info("Error deleting temporal directory %q: %s", tmpDir, err)
	//	}
	//}()

	err = compression.UnpackFromReader(compressedSrc, tmpDir)
	if err != nil {
		return "", err
	}

	storeOptions, err := storage.DefaultStoreOptionsAutoDetectUID()
	if err != nil {
		return "", fmt.Errorf("error initializing store options: %s", err)
	}
	fmt.Printf("%+v\n", storeOptions)

	storeOptions.GraphDriverName = "vfs"
	//storeOptions.RunRoot = "/run/containers/storage"
	//storeOptions.GraphRoot = "/var/lib/containers/storage"
	fmt.Printf("%+v\n", storeOptions)

	store, err := storage.GetStore(storeOptions)
	if err != nil {
		return "", fmt.Errorf("initializing store: %w", err)
	}

	id, err := pr.buildImage(ctx, store, tmpDir)
	if err != nil {
		return "", fmt.Errorf("building image: %w", err)
	}

	imageName := fmt.Sprintf("%s-%s:%s", product, process, version)
	_, err = pr.pushImage(ctx, store, id, imageName)

	return imageName, nil
}

func (pr *ProcessRegistry) buildImage(ctx context.Context, buildStore storage.Store, srcPath string) (string, error) {
	dockerfilePath := path.Join(srcPath, _defaultDockerfile)

	buildOpts := define.BuildOptions{
		ContextDirectory: srcPath,
		Isolation:        buildah.IsolationOCIRootless,
		SystemContext: &types.SystemContext{
			SystemRegistriesConfPath: viper.GetString(config.RegistriesConfPathKey),
			SignaturePolicyPath:      viper.GetString(config.SignaturePolicyPathKey),
			//SystemRegistriesConfDirPath: ".",
		},
		//AddCapabilities: []string{"SYS_ADMIN"},
	}

	id, _, err := imagebuildah.BuildDockerfiles(ctx, buildStore, buildOpts, dockerfilePath)
	if err != nil {
		return "", fmt.Errorf("error building image: %s", err)
	}

	pr.logger.Info("Image successfully builded", "id", id)

	return id, nil

}

func (pr *ProcessRegistry) pushImage(
	ctx context.Context,
	store storage.Store,
	imageRef string,
	imageName string,
	// imageRef types.ImageReference,
) (string, error) {
	destSpec := fmt.Sprintf("docker://%s/%s", viper.GetString(config.RegistryURLKey), imageName)
	// destSpec := "docker://localhost:5000/kre-go:latest"
	dest, err := alltransports.ParseImageName(destSpec)
	// add the docker:// transport to see if they neglected it.
	if err != nil {
		destTransport := strings.Split(destSpec, ":")[0]
		if t := transports.Get(destTransport); t != nil {
			return "", err
		}

		if strings.Contains(destSpec, "://") {
			return "", err
		}

		destSpec = "docker://" + destSpec
		dest2, err2 := alltransports.ParseImageName(destSpec)
		if err2 != nil {
			return "", err2
		}
		dest = dest2
		pr.logger.Info("Assuming docker:// as the transport method for DESTINATION: %s", destSpec)
	}

	options := buildah.PushOptions{
		// SignaturePolicyPath: iopts.signaturePolicy,
		Store: store,
		SystemContext: &types.SystemContext{
			SystemRegistriesConfPath: viper.GetString(config.RegistriesConfPathKey),
			SignaturePolicyPath:      viper.GetString(config.SignaturePolicyPathKey),
			//SystemRegistriesConfDirPath: ".",
		},
	}

	pushedImageRef, digest, err := buildah.Push(ctx, imageRef, dest, options)
	if err != nil {
		return "", fmt.Errorf("pushing image: %s", err)
	}

	pr.logger.Info("Pushed Image", "ref", pushedImageRef, "digest", digest)

	return pushedImageRef.Name(), nil
}
