// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

import (
	"github.com/99designs/gqlgen/graphql"
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

type PublishVersionInput struct {
	VersionTag string `json:"versionTag"`
	Comment    string `json:"comment"`
	ProductID  string `json:"productID"`
	Force      bool   `json:"force"`
}

type PublishedTrigger struct {
	Trigger string `json:"trigger"`
	URL     string `json:"url"`
}

type RegisterProcessInput struct {
	File        graphql.Upload `json:"file"`
	Version     string         `json:"version"`
	ProductID   string         `json:"productID"`
	ProcessID   string         `json:"processID"`
	ProcessType string         `json:"processType"`
}

type RegisterPublicProcessInput struct {
	File        graphql.Upload `json:"file"`
	Version     string         `json:"version"`
	ProcessID   string         `json:"processID"`
	ProcessType string         `json:"processType"`
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
