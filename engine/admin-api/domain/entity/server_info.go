package entity

type ComponentStatus string

const (
	ComponentStatusKO = "KO"
	ComponentStatusOK = "OK"
)

type ComponentInfo struct {
	Component string
	Version   string
	Status    ComponentStatus
}

type ServerInfo struct {
	Components []ComponentInfo
}
