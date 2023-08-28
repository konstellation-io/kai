// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type CreateProductInput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateVersionInput struct {
	File      graphql.Upload `json:"file"`
	ProductID string         `json:"productID"`
}

type LogPage struct {
	Cursor *string              `json:"cursor,omitempty"`
	Items  []*entity.ProcessLog `json:"items"`
}

type PublishVersionInput struct {
	VersionTag string `json:"versionTag"`
	Comment    string `json:"comment"`
	ProductID  string `json:"productID"`
}

type RegisterProcessInput struct {
	File        graphql.Upload `json:"file"`
	Version     string         `json:"version"`
	ProductID   string         `json:"productID"`
	ProcessID   string         `json:"processID"`
	ProcessType string         `json:"processType"`
}

type RegisteredImage struct {
	ProcessedImageID string `json:"processedImageID"`
}

type RevokeUserProductGrantsInput struct {
	TargetID string  `json:"targetID"`
	Product  string  `json:"product"`
	Comment  *string `json:"comment,omitempty"`
}

type StartVersionInput struct {
	VersionTag string `json:"versionTag"`
	Comment    string `json:"comment"`
	ProductID  string `json:"productID"`
}

type StopVersionInput struct {
	VersionTag string `json:"versionTag"`
	Comment    string `json:"comment"`
	ProductID  string `json:"productID"`
}

type UnpublishVersionInput struct {
	VersionTag string `json:"versionTag"`
	Comment    string `json:"comment"`
	ProductID  string `json:"productID"`
}

type UpdateUserProductGrantsInput struct {
	TargetID string   `json:"targetID"`
	Product  string   `json:"product"`
	Grants   []string `json:"grants"`
	Comment  *string  `json:"comment,omitempty"`
}
