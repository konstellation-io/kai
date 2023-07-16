package entity

type ServerInfo struct {
	Components []ComponentInfo
}

type ComponentInfo struct {
	Name    string
	Version string
}
