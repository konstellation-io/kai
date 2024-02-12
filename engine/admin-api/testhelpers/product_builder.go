package testhelpers

import (
	"github.com/bxcodec/faker/v3"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type ProductBuilder struct {
	product *entity.Product
}

func NewProductBuilder() *ProductBuilder {
	productID := faker.UUIDHyphenated()

	return &ProductBuilder{
		product: &entity.Product{
			ID:            productID,
			Name:          faker.Name(),
			Description:   "Test description",
			KeyValueStore: faker.UUIDHyphenated(),
			MinioConfiguration: entity.MinioConfiguration{
				Bucket: productID,
			},
			ServiceAccount: entity.ServiceAccount{
				Username: productID,
				Password: faker.Password(),
				Group:    productID,
			},
		},
	}
}

func (pb *ProductBuilder) WithID(id string) *ProductBuilder {
	pb.product.ID = id
	pb.product.MinioConfiguration.Bucket = id
	pb.product.ServiceAccount.Username = id
	pb.product.ServiceAccount.Group = id

	return pb
}

func (pb *ProductBuilder) WithPublishedVersion(publishedVersion *string) *ProductBuilder {
	pb.product.PublishedVersion = publishedVersion
	return pb
}

func (pb *ProductBuilder) WithServiceAccount(serviceAccount entity.ServiceAccount) *ProductBuilder {
	pb.product.ServiceAccount = serviceAccount
	return pb
}

func (pb *ProductBuilder) Build() *entity.Product {
	return pb.product
}
