package repository

type ObjectStorage interface {
	CreateBucket(name string) error
}
