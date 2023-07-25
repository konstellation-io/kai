package usecase

import (
	"context"
	"fmt"
	"io"

	"github.com/go-logr/logr"
)

type ProcessService struct {
	logger            logr.Logger
	processRegistry   ProcessRegistry
	processRepository ProcessRepository
}

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/${GOFILE} -package=mocks

type ProcessRegistry interface {
	RegisterProcess(ctx context.Context, product, version, process string, src io.Reader) (string, error)
}

type ProcessRepository interface {
	CreateProcess() error
}

type ProcessMetadata struct {
	Dockerfile string
}

func NewProcessService(
	logger logr.Logger,
	processRegistry ProcessRegistry,
	// processRepository ProcessRepository,
) *ProcessService {
	return &ProcessService{
		logger:          logger,
		processRegistry: processRegistry,
		//processRepository: processRepository,
	}
}

func (ps *ProcessService) RegisterProcess(
	ctx context.Context,
	product, version, process string,
	sources io.Reader,
) (string, error) {
	ps.logger.Info("Registering new process")

	processRef, err := ps.processRegistry.RegisterProcess(ctx, product, version, process, sources)
	if err != nil {
		return "", fmt.Errorf("registering process: %w", err)
	}

	ps.logger.Info("Registered process", "processRef", processRef)

	return processRef, nil
}

//
//func (ps *ProcessService) decompressSources(compressedSourcesFile *os.File, dstDir string) error {
//	compressedSources, err := os.Open(compressedSourcesFile.Name())
//	if err != nil {
//		return err
//	}
//
//	defer func() {
//		err := compressedSources.Close()
//		if err != nil {
//			ps.logger.Info("Error closing file %s: %s", compressedSources.Name(), err)
//		}
//	}()
//
//	sources, err := gzip.NewReader(compressedSources)
//	if err != nil {
//		return err
//	}
//
//	tarReader := tar.NewReader(sources)
//
//	for {
//		tarFile, err := tarReader.Next()
//		if err != io.EOF {
//			break
//		}
//
//		if err != nil {
//			return err
//		}
//
//		filePath := filepath.Join(dstDir, tarFile.Name)
//
//		if err := ps.processFile(tarReader, filePath, tarFile.Typeflag); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func (ps *ProcessService) processFile(tarReader *tar.Reader, filePath string, fileType byte) error {
//	switch fileType {
//	case tar.TypeDir:
//		if err := os.Mkdir(filePath, 0755); err != nil {
//			return fmt.Errorf("error creating krt dir %s: %w", filePath, err)
//		}
//
//	case tar.TypeReg:
//		outFile, err := os.Create(filePath)
//
//		if err != nil {
//			return fmt.Errorf("error creating krt file %s: %w", filePath, err)
//		}
//
//		if _, err := io.Copy(outFile, tarReader); err != nil {
//			return fmt.Errorf("error copying krt file %s: %w", filePath, err)
//		}
//
//		err = outFile.Close()
//		if err != nil {
//			return fmt.Errorf("error closing krt file %s: %w", filePath, err)
//		}
//
//	default:
//		return fmt.Errorf("error extracting krt files: uknown type [%v] in [%s]", fileType, filePath)
//	}
//
//	return nil
//}
//
//func copyStreamToTmpFile(sources io.Reader) (*os.File, error) {
//	tmpFile, err := os.CreateTemp("", "process-*")
//	if err != nil {
//		return nil, fmt.Errorf("creating temp file for process: %w", err)
//	}
//
//	_, err = io.Copy(tmpFile, sources)
//	if err != nil {
//		return nil, fmt.Errorf("copying sources to temp file: %w", err)
//	}
//
//	return tmpFile, nil
//}