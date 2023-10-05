package entity

type KeyValueStoreScope string

const (
	KVScopeGlobal   KeyValueStoreScope = "global"
	KVScopeVersion  KeyValueStoreScope = "version"
	KVScopeWorkflow KeyValueStoreScope = "workflow"
	KVScopeProcess  KeyValueStoreScope = "process"
)

type WorkflowKeyValueStores struct {
	WorkflowStore string
	Processes     map[string]string
}

type VersionKeyValueStores struct {
	ProjectStore    string
	WorkflowsStores map[string]*WorkflowKeyValueStores
}
