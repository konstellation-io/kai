package entity

type ComponentStatus string

const (
	ComponentStatusKO = "KO"
	ComponentStatusOK = "OK"
)

type ComponentInfo struct {
	Name    string
	Version string
}

type ServerInfo struct {
	Components []ComponentInfo
	Status     string
}
