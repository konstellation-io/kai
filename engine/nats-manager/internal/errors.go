package internal

import "errors"

var ErrInvalidObjectStoreName = errors.New("invalid object store name")
var ErrInvalidObjectStoreScope = errors.New("invalid object store scope")
var ErrInvalidKeyValueStoreScope = errors.New("invalid key-value store scope")
var ErrEmptyWorkflowName = errors.New("workflow name cannot be empty")
var ErrEmptyProcessName = errors.New("process name cannot be empty")
var ErrNoWorkflowsDefined = errors.New("no workflows defined")
var ErrNoOptFilter = errors.New("optFilter param accepts 0 or 1 value")
