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
}

type MinioConfiguration struct {
	User     string `bson:"user"`
	Group    string `bson:"group"`
	Password string `bson:"password"`
	Bucket   string `bson:"bucket"`
}

func (r *Product) Validate() error {
	return validate.Struct(r)
}
