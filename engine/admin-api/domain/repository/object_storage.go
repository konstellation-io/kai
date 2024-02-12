package repository

import "context"

//go:generate mockery --name ObjectStorage --output ../../mocks --filename object_storage.go --structname MockObjectStorage

type ObjectStorage interface {
	CreateBucket(ctx context.Context, bucket string) error
	DeleteBucket(ctx context.Context, bucket string) error
	CreateBucketPolicy(ctx context.Context, bucket string) (string, error)
	DeleteBucketPolicy(ctx context.Context, policyName string) error
	UploadImageSources(ctx context.Context, product, image string, sources []byte) error
	DeleteImageSources(ctx context.Context, product, image string) error
}
