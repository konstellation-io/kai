package auth

import (
	"fmt"

	"github.com/casbin/casbin/v2"

	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
)

type CasbinAccessControl struct {
	logger   logging.Logger
	enforcer *casbin.Enforcer
}

func NewCasbinAccessControl(logger logging.Logger, modelPath, policyPath string) (*CasbinAccessControl, error) {
	e, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, err
	}

	accessControl := &CasbinAccessControl{
		logger,
		e,
	}

	return accessControl, nil
}

// change input params to ones obtained from jwt token.
func (a *CasbinAccessControl) CheckGrant(userID string, resource auth.AccessControlResource, action auth.AccessControlAction) error {
	if !resource.IsValid() {
		return invalidAccessControlResourceError
	}

	if !action.IsValid() {
		return invalidAccessControlActionError
	}

	allowed, err := a.enforcer.Enforce(userID, resource.String(), action.String())
	if err != nil {
		a.logger.Errorf("error checking grant: %s", err)
		return err
	}

	a.logger.Infof("Checking grant userID[%s] resource[%s] action[%s] allowed[%t]", userID, resource, action, allowed)

	if !allowed {
		//nolint:goerr113 // errors need to be wrapped
		return fmt.Errorf("you are not allowed to %s %s", action, resource)
	}

	return nil
}
