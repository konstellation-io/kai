package repository

import "context"

//go:generate mockery --name ObjectStorage --output ../../mocks --filename object_storage.go --structname MockObjectStorage

type ObjectStorage interface {
	CreateBucket(ctx context.Context, bucket string) error
	CreateBucketPolicy(ctx context.Context, bucket string) (string, error)
}
