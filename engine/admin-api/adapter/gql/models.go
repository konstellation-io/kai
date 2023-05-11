// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
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

type CreateVersionInput struct {
	File      graphql.Upload `json:"file"`
	RuntimeID string         `json:"runtimeId"`
}

type LogPage struct {
	Cursor *string           `json:"cursor,omitempty"`
	Items  []*entity.NodeLog `json:"items"`
}

type PublishVersionInput struct {
	VersionName string `json:"versionName"`
	Comment     string `json:"comment"`
	RuntimeID   string `json:"runtimeId"`
}

type Settings struct {
	AuthAllowedDomains    []string `json:"authAllowedDomains"`
	SessionLifetimeInDays int      `json:"sessionLifetimeInDays"`
}

type SettingsInput struct {
	AuthAllowedDomains    []string `json:"authAllowedDomains,omitempty"`
	SessionLifetimeInDays *int     `json:"sessionLifetimeInDays,omitempty"`
}

type StartVersionInput struct {
	VersionName string `json:"versionName"`
	Comment     string `json:"comment"`
	RuntimeID   string `json:"runtimeId"`
}

type StopVersionInput struct {
	VersionName string `json:"versionName"`
	Comment     string `json:"comment"`
	RuntimeID   string `json:"runtimeId"`
}

type UnpublishVersionInput struct {
	VersionName string `json:"versionName"`
	Comment     string `json:"comment"`
	RuntimeID   string `json:"runtimeId"`
}

type UpdateAccessLevelInput struct {
	UserIds     []string    `json:"userIds"`
	AccessLevel AccessLevel `json:"accessLevel"`
	Comment     string      `json:"comment"`
}

type UpdateConfigurationInput struct {
	VersionName            string                         `json:"versionName"`
	RuntimeID              string                         `json:"runtimeId"`
	ConfigurationVariables []*ConfigurationVariablesInput `json:"configurationVariables"`
}

type UsersInput struct {
	UserIds []string `json:"userIds"`
	Comment string   `json:"comment"`
}

type AccessLevel string

const (
	AccessLevelViewer  AccessLevel = "VIEWER"
	AccessLevelManager AccessLevel = "MANAGER"
	AccessLevelAdmin   AccessLevel = "ADMIN"
)

var AllAccessLevel = []AccessLevel{
	AccessLevelViewer,
	AccessLevelManager,
	AccessLevelAdmin,
}

func (e AccessLevel) IsValid() bool {
	switch e {
	case AccessLevelViewer, AccessLevelManager, AccessLevelAdmin:
		return true
	}
	return false
}

func (e AccessLevel) String() string {
	return string(e)
}

func (e *AccessLevel) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AccessLevel(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AccessLevel", str)
	}
	return nil
}

func (e AccessLevel) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
