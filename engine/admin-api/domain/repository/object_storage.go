package repository

//go:generate mockery --name ObjectStorage --output ../../mocks --filename object_storage.go --structname MockObjectStorage

type ObjectStorage interface {
	CreateFolder(name string) error
}