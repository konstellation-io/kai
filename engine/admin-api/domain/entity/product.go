package entity

import (
	"time"

	"github.com/go-playground/validator/v10"
)

//nolint:gochecknoglobals // validate have to be a global variable.
var (
	validate = validator.New()
)

type Product struct {
	ID                 string             `bson:"_id" validate:"required"`
	Name               string             `bson:"name" validate:"required,lte=40"`
	Description        string             `bson:"description" validate:"required,lte=500"`
	CreationDate       time.Time          `bson:"creationDate"`
	Owner              string             `bson:"owner"`
	MinioConfiguration MinioConfiguration `bson:"minioConfiguration"`
	KeyValueStore      string             `bson:"keyValueStore"`
	PublishedVersion   *string            `bson:"publishedVersion"`
	ServiceAccount     ServiceAccount     `bson:"serviceAccount"`
}

type MinioConfiguration struct {
	Bucket string `bson:"bucket"`
}

type ServiceAccount struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
	Group    string `bson:"group"`
}

func (p *Product) Validate() error {
	return validate.Struct(p)
}

func (p *Product) HasVersionPublished() bool {
	return p.PublishedVersion != nil
}

func (p *Product) UpdatePublishedVersion(publishedVersion string) {
	p.PublishedVersion = &publishedVersion
}

func (p *Product) RemovePublishedVersion() {
	p.PublishedVersion = nil
}
