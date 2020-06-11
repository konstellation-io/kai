// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"gitlab.com/konstellation/kre/admin-api/domain/entity"
)

type Alert struct {
	ID      string          `json:"id"`
	Type    AlertLevel      `json:"type"`
	Message string          `json:"message"`
	Runtime *entity.Runtime `json:"runtime"`
}

type ConfigurationVariable struct {
	Key   string                    `json:"key"`
	Value string                    `json:"value"`
	Type  ConfigurationVariableType `json:"type"`
}

type ConfigurationVariablesInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateRuntimeInput struct {
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
	AuthAllowedEmails     []string `json:"authAllowedEmails"`
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
}

type UpdateConfigurationInput struct {
	VersionID              string                         `json:"versionId"`
	ConfigurationVariables []*ConfigurationVariablesInput `json:"configurationVariables"`
}

type UsersInput struct {
	UserIds []string `json:"userIds"`
}

type AlertLevel string

const (
	AlertLevelError   AlertLevel = "ERROR"
	AlertLevelWarning AlertLevel = "WARNING"
)

var AllAlertLevel = []AlertLevel{
	AlertLevelError,
	AlertLevelWarning,
}

func (e AlertLevel) IsValid() bool {
	switch e {
	case AlertLevelError, AlertLevelWarning:
		return true
	}
	return false
}

func (e AlertLevel) String() string {
	return string(e)
}

func (e *AlertLevel) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AlertLevel(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AlertLevel", str)
	}
	return nil
}

func (e AlertLevel) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ConfigurationVariableType string

const (
	ConfigurationVariableTypeVariable ConfigurationVariableType = "VARIABLE"
	ConfigurationVariableTypeFile     ConfigurationVariableType = "FILE"
)

var AllConfigurationVariableType = []ConfigurationVariableType{
	ConfigurationVariableTypeVariable,
	ConfigurationVariableTypeFile,
}

func (e ConfigurationVariableType) IsValid() bool {
	switch e {
	case ConfigurationVariableTypeVariable, ConfigurationVariableTypeFile:
		return true
	}
	return false
}

func (e ConfigurationVariableType) String() string {
	return string(e)
}

func (e *ConfigurationVariableType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ConfigurationVariableType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ConfigurationVariableType", str)
	}
	return nil
}

func (e ConfigurationVariableType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type NodeStatus string

const (
	NodeStatusStarted NodeStatus = "STARTED"
	NodeStatusStopped NodeStatus = "STOPPED"
	NodeStatusError   NodeStatus = "ERROR"
)

var AllNodeStatus = []NodeStatus{
	NodeStatusStarted,
	NodeStatusStopped,
	NodeStatusError,
}

func (e NodeStatus) IsValid() bool {
	switch e {
	case NodeStatusStarted, NodeStatusStopped, NodeStatusError:
		return true
	}
	return false
}

func (e NodeStatus) String() string {
	return string(e)
}

func (e *NodeStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = NodeStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid NodeStatus", str)
	}
	return nil
}

func (e NodeStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type RuntimeStatus string

const (
	RuntimeStatusCreating RuntimeStatus = "CREATING"
	RuntimeStatusStarted  RuntimeStatus = "STARTED"
	RuntimeStatusError    RuntimeStatus = "ERROR"
)

var AllRuntimeStatus = []RuntimeStatus{
	RuntimeStatusCreating,
	RuntimeStatusStarted,
	RuntimeStatusError,
}

func (e RuntimeStatus) IsValid() bool {
	switch e {
	case RuntimeStatusCreating, RuntimeStatusStarted, RuntimeStatusError:
		return true
	}
	return false
}

func (e RuntimeStatus) String() string {
	return string(e)
}

func (e *RuntimeStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RuntimeStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RuntimeStatus", str)
	}
	return nil
}

func (e RuntimeStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type UserActivityType string

const (
	UserActivityTypeLogin                      UserActivityType = "LOGIN"
	UserActivityTypeLogout                     UserActivityType = "LOGOUT"
	UserActivityTypeCreateRuntime              UserActivityType = "CREATE_RUNTIME"
	UserActivityTypeCreateVersion              UserActivityType = "CREATE_VERSION"
	UserActivityTypePublishVersion             UserActivityType = "PUBLISH_VERSION"
	UserActivityTypeUnpublishVersion           UserActivityType = "UNPUBLISH_VERSION"
	UserActivityTypeStartVersion               UserActivityType = "START_VERSION"
	UserActivityTypeStopVersion                UserActivityType = "STOP_VERSION"
	UserActivityTypeUpdateSetting              UserActivityType = "UPDATE_SETTING"
	UserActivityTypeUpdateVersionConfiguration UserActivityType = "UPDATE_VERSION_CONFIGURATION"
	UserActivityTypeCreateUser                 UserActivityType = "CREATE_USER"
	UserActivityTypeRemoveUsers                UserActivityType = "REMOVE_USERS"
	UserActivityTypeUpdateAccessLevels         UserActivityType = "UPDATE_ACCESS_LEVELS"
	UserActivityTypeRevokeSessions             UserActivityType = "REVOKE_SESSIONS"
)

var AllUserActivityType = []UserActivityType{
	UserActivityTypeLogin,
	UserActivityTypeLogout,
	UserActivityTypeCreateRuntime,
	UserActivityTypeCreateVersion,
	UserActivityTypePublishVersion,
	UserActivityTypeUnpublishVersion,
	UserActivityTypeStartVersion,
	UserActivityTypeStopVersion,
	UserActivityTypeUpdateSetting,
	UserActivityTypeUpdateVersionConfiguration,
	UserActivityTypeCreateUser,
	UserActivityTypeRemoveUsers,
	UserActivityTypeUpdateAccessLevels,
	UserActivityTypeRevokeSessions,
}

func (e UserActivityType) IsValid() bool {
	switch e {
	case UserActivityTypeLogin, UserActivityTypeLogout, UserActivityTypeCreateRuntime, UserActivityTypeCreateVersion, UserActivityTypePublishVersion, UserActivityTypeUnpublishVersion, UserActivityTypeStartVersion, UserActivityTypeStopVersion, UserActivityTypeUpdateSetting, UserActivityTypeUpdateVersionConfiguration, UserActivityTypeCreateUser, UserActivityTypeRemoveUsers, UserActivityTypeUpdateAccessLevels, UserActivityTypeRevokeSessions:
		return true
	}
	return false
}

func (e UserActivityType) String() string {
	return string(e)
}

func (e *UserActivityType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UserActivityType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UserActivityType", str)
	}
	return nil
}

func (e UserActivityType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type VersionStatus string

const (
	VersionStatusStarting  VersionStatus = "STARTING"
	VersionStatusStarted   VersionStatus = "STARTED"
	VersionStatusPublished VersionStatus = "PUBLISHED"
	VersionStatusStopped   VersionStatus = "STOPPED"
)

var AllVersionStatus = []VersionStatus{
	VersionStatusStarting,
	VersionStatusStarted,
	VersionStatusPublished,
	VersionStatusStopped,
}

func (e VersionStatus) IsValid() bool {
	switch e {
	case VersionStatusStarting, VersionStatusStarted, VersionStatusPublished, VersionStatusStopped:
		return true
	}
	return false
}

func (e VersionStatus) String() string {
	return string(e)
}

func (e *VersionStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = VersionStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid VersionStatus", str)
	}
	return nil
}

func (e VersionStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
