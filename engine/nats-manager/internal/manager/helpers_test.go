//go:build unit

package manager_test

import (
	"fmt"
	"reflect"

	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
)

type streamConfigMatcher struct {
	expectedStreamConfig *entity.StreamConfig
}

func newStreamConfigMatcher(expectedStreamConfig *entity.StreamConfig) *streamConfigMatcher {
	return &streamConfigMatcher{
		expectedStreamConfig: expectedStreamConfig,
	}
}

func (m streamConfigMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expectedStreamConfig)
}

func (m streamConfigMatcher) Matches(actual interface{}) bool {
	actualCfg, ok := actual.(*entity.StreamConfig)
	if !ok {
		return false
	}

	return reflect.DeepEqual(actualCfg, m.expectedStreamConfig)
}
