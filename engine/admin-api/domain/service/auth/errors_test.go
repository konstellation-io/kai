//go:build unit

package auth_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/stretchr/testify/assert"
)

func TestUnauthorizedError_Error(t *testing.T) {
	err := auth.UnauthorizedError{
		Product: "product-01",
		Action:  auth.ActViewProduct,
	}

	stringFormat := "you don't have authorization to %s in product %s"
	expectedErrorMsg := fmt.Sprintf(stringFormat, err.Action.String(), err.Product)

	assert.Equal(t, expectedErrorMsg, err.Error())
}

func TestUnauthorizedError_Error_WithoutProduct(t *testing.T) {
	err := auth.UnauthorizedError{
		Action: auth.ActViewServerInfo,
	}

	stringFormat := "you don't have authorization to %s"
	expectedErrorMsg := fmt.Sprintf(stringFormat, err.Action.String())

	assert.Equal(t, expectedErrorMsg, err.Error())
}

func TestUnauthorizedError_Error_WithError(t *testing.T) {
	err := auth.UnauthorizedError{
		Action: auth.ActViewServerInfo,
		Err:    errors.New("wrapped error"),
	}

	stringFormat := "you don't have authorization to %s: %s"
	expectedErrorMsg := fmt.Sprintf(stringFormat, err.Action.String(), err.Err.Error())

	assert.Equal(t, expectedErrorMsg, err.Error())
}
