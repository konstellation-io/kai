package entity

type WorkflowKeyValueStores struct {
	WorkflowStore string
	Processes     map[string]string
}

type VersionKeyValueStores struct {
	ProjectStore    string
	WorkflowsStores map[string]*WorkflowKeyValueStores
}
