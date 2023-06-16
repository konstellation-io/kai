package auth

import (
	"errors"

	"github.com/casbin/casbin/v2"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
)

type OptFunc func(*Opts)

type Opts struct {
	adminRole string
}

func defaultOpts() Opts {
	return Opts{
		adminRole: auth.DefaultAdminRole,
	}
}

func WithAdminRole(adminRole string) OptFunc {
	return func(opts *Opts) {
		opts.adminRole = adminRole
	}
}

type CasbinAccessControl struct {
	cfg      Opts
	logger   logging.Logger
	enforcer *casbin.Enforcer
}

func NewCasbinAccessControl(
	logger logging.Logger,
	modelPath,
	policyPath string,
	opts ...OptFunc,
) (*CasbinAccessControl, error) {
	o := defaultOpts()

	for _, fn := range opts {
		fn(&o)
	}

	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, err
	}

	accessController := &CasbinAccessControl{
		cfg:      o,
		logger:   logger,
		enforcer: enforcer,
	}

	accessController.addCustomFunctions()

	return accessController, nil
}

func (a *CasbinAccessControl) addCustomFunctions() {
	a.enforcer.AddFunction("isAdmin", a.isAdminFunc)
	a.enforcer.AddFunction("hasGrantsForResource", a.hasGrantsForResourceFunc)
}

func (a *CasbinAccessControl) checkGrants(
	user *entity.User,
	product string,
	action auth.AccessControlAction,
) error {
	if !action.IsValid() {
		return ErrInvalidAccessControlAction
	}

	for _, realmRole := range user.Roles {
		allowed, err := a.enforcer.Enforce(realmRole, user.ProductGrants, product, action.String())
		if err != nil {
			a.logger.Errorf("error checking grants: %s", err)
			return err
		}

		a.logger.Infof(
			"Checking grants userID[%s] role[%s] action[%s] product[%s] allowed[%t]",
			user.ID, realmRole, action, product, allowed,
		)

		if allowed {
			return nil
		}
	}

	return ErrNonAuthorized
}

func (a *CasbinAccessControl) CheckAdminGrants(
	user *entity.User,
	action auth.AccessControlAction,
) error {
	err := a.checkGrants(user, "", action)
	if errors.Is(err, ErrNonAuthorized) {
		return NonAdminAccess(action.String())
	} else if err != nil {
		return err
	}

	return nil
}

func (a *CasbinAccessControl) CheckProductGrants(
	user *entity.User,
	product string,
	action auth.AccessControlAction,
) error {
	err := a.checkGrants(user, product, action)
	if errors.Is(err, ErrNonAuthorized) {
		return NonAuthorizedForProductError(action.String(), product)
	} else if err != nil {
		return err
	}

	return nil
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

func (a *CasbinAccessControl) IsAdmin(user *entity.User) bool {
	for _, role := range user.Roles {
		if role == a.cfg.adminRole {
			return true
		}
	}

	return false
}

func (a *CasbinAccessControl) hasGrantsForResourceFunc(args ...interface{}) (interface{}, error) {
	if len(args) != 3 {
		return false, ErrInvalidNumberOfArguments
	}

	grants := args[0].(entity.ProductGrants)
	resource := args[1].(string)
	act := args[2].(string)

	return a.hasGrantsForResource(grants, resource, act), nil
}

func (a *CasbinAccessControl) isAdminFunc(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidNumberOfArguments
	}

	role := args[0].(string)

	return role == a.cfg.adminRole, nil
}
