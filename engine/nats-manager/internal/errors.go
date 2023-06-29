package internal

import "errors"

var ErrInvalidObjectStoreName = errors.New("invalid object store name")
var ErrInvalidObjectStoreScope = errors.New("invalid object store scope")
var ErrInvalidKeyValueStoreScope = errors.New("invalid key-value store scope")
var ErrEmptyWorkflow = errors.New("workflow name cannot be empty")
var ErrEmptyProcessID = errors.New("process id cannot be empty")
var ErrNoWorkflowsDefined = errors.New("no workflow defined")
var ErrNoOptFilter = errors.New("optFilter param accepts 0 or 1 value")
