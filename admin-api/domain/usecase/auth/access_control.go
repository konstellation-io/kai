package auth

type AccessControlResource string

const ResMetrics AccessControlResource = "metrics"
const ResResourceMetrics AccessControlResource = "resource-metrics"
const ResRuntime AccessControlResource = "runtimes"
const ResVersion AccessControlResource = "versions"
const ResSettings AccessControlResource = "settings"
const ResUsers AccessControlResource = "users"
const ResAudit AccessControlResource = "audits"
const ResLogs AccessControlResource = "logs"

func (e AccessControlResource) IsValid() bool {
	switch e {
	case ResMetrics, ResResourceMetrics, ResRuntime, ResVersion, ResSettings, ResUsers, ResAudit, ResLogs:
		return true
	}
	return false
}

func (e AccessControlResource) String() string {
	return string(e)
}

type AccessControlAction string

const ActView AccessControlAction = "view"
const ActEdit AccessControlAction = "edit"

func (e AccessControlAction) IsValid() bool {
	switch e {
	case ActView, ActEdit:
		return true
	}
	return false
}

func (e AccessControlAction) String() string {
	return string(e)
}

type AccessControl interface {
	CheckPermission(userID string, resource AccessControlResource, action AccessControlAction) error
	ReloadUserRoles() error
}