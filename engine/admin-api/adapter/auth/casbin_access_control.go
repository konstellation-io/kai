package auth

import (
	"fmt"

	"github.com/casbin/casbin/v2"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
)

type Config struct {
	AdminRole string
}

type CasbinAccessControl struct {
	cfg      Config
	logger   logging.Logger
	enforcer *casbin.Enforcer
}

func NewCasbinAccessControl(cfg Config, logger logging.Logger, modelPath, policyPath string) (*CasbinAccessControl, error) {
	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, err
	}

	accessController := &CasbinAccessControl{
		cfg,
		logger,
		enforcer,
	}

	accessController.addCustomFunctions()

	return accessController, nil
}

func (a *CasbinAccessControl) addCustomFunctions() {
	a.enforcer.AddFunction("isAdmin", a.isAdminFunc)
	a.enforcer.AddFunction("hasGrantForResource", a.hasGrantsForResourceFunc)
}

func (a *CasbinAccessControl) CheckPermission(
	user *entity.User,
	product string,
	action auth.AccessControlAction,
) error {
	if !action.IsValid() {
		return InvalidAccessControlActionError
	}

	for _, realmRole := range user.Roles {
		allowed, err := a.enforcer.Enforce(realmRole, user.ProductGrants, product, action.String())
		if err != nil {
			a.logger.Errorf("error checking permission: %s", err)
			return err
		}

		a.logger.Infof(
			"Checking permission userID[%s] realmRole[%s] action[%s] product[%s] allowed[%t]",
			user.ID, realmRole, action, product, allowed,
		)

		if allowed {
			return nil
		}
	}

	//nolint:goerr113 // errors need to be wrapped
	return fmt.Errorf("you are not allowed to %s %s", action, product)
}

func (a *CasbinAccessControl) hasGrantsForResource(
	grants entity.ProductGrants,
	product,
	act string,
) bool {
	resGrants, ok := grants[product]
	if !ok {
		return false
	}

	for _, grant := range resGrants {
		if grant == act {
			return true
		}
	}

	return false
}

func (a *CasbinAccessControl) hasGrantsForResourceFunc(args ...interface{}) (interface{}, error) {
	grants := args[0].(entity.ProductGrants)
	resource := args[1].(string)
	act := args[2].(string)

	return a.hasGrantsForResource(grants, resource, act), nil
}

func (a *CasbinAccessControl) isAdminFunc(args ...interface{}) (interface{}, error) {
	role := args[0].(string)

	return role == a.cfg.AdminRole, nil
}
