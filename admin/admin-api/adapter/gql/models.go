// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/konstellation-io/kre/admin/admin-api/domain/entity"
)

type ConfigurationVariablesInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateRuntimeInput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateUserInput struct {
	Email       string             `json:"email"`
	AccessLevel entity.AccessLevel `json:"accessLevel"`
}

type CreateVersionInput struct {
	File      graphql.Upload `json:"file"`
	RuntimeID string         `json:"runtimeId"`
}

type DeleteAPITokenInput struct {
	ID string `json:"id"`
}

type GenerateAPITokenInput struct {
	Name string `json:"name"`
}

type LogPage struct {
	Cursor *string           `json:"cursor"`
	Items  []*entity.NodeLog `json:"items"`
}

type PublishVersionInput struct {
	VersionID string `json:"versionId"`
	Comment   string `json:"comment"`
}

type SettingsInput struct {
	AuthAllowedDomains    []string `json:"authAllowedDomains"`
	SessionLifetimeInDays *int     `json:"sessionLifetimeInDays"`
}

type StartVersionInput struct {
	VersionID string `json:"versionId"`
	Comment   string `json:"comment"`
}

type StopVersionInput struct {
	VersionID string `json:"versionId"`
	Comment   string `json:"comment"`
}

type UnpublishVersionInput struct {
	VersionID string `json:"versionId"`
	Comment   string `json:"comment"`
}

type UpdateAccessLevelInput struct {
	UserIds     []string           `json:"userIds"`
	AccessLevel entity.AccessLevel `json:"accessLevel"`
	Comment     string             `json:"comment"`
}

type UpdateConfigurationInput struct {
	VersionID              string                         `json:"versionId"`
	ConfigurationVariables []*ConfigurationVariablesInput `json:"configurationVariables"`
}

type UsersInput struct {
	UserIds []string `json:"userIds"`
	Comment string   `json:"comment"`
}
