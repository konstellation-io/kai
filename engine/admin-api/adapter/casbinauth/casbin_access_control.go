package casbinauth

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
)

const _defaultResource = ""

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
	logger   logr.Logger
	enforcer *casbin.Enforcer
}

func NewCasbinAccessControl(
	logger logr.Logger,
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

func (a *CasbinAccessControl) CheckRoleGrants(user *entity.User, action auth.Action) error {
	err := a.checkGrants(user, _defaultResource, action)
	if err != nil {
		return err
	}

	return nil
}

func (a *CasbinAccessControl) CheckProductGrants(
	user *entity.User,
	product string,
	action auth.Action,
) error {
	err := a.checkGrants(user, product, action)
	if err != nil {
		return err
	}

	return nil
}

func (a *CasbinAccessControl) IsAdmin(user *entity.User) bool {
	for _, role := range user.Roles {
		if role == a.cfg.adminRole {
			return true
		}
	}

	return false
}

func (a *CasbinAccessControl) GetUserProductsWithViewAccess(user *entity.User) []string {
	if a.IsAdmin(user) {
		return nil
	}

	visibleProducts := make([]string, 0, len(user.ProductGrants))

	for prod := range user.ProductGrants {
		if err := a.CheckProductGrants(user, prod, auth.ActViewProduct); err == nil {
			visibleProducts = append(visibleProducts, prod)
		}
	}

	return visibleProducts
}

func (a *CasbinAccessControl) addCustomFunctions() {
	a.enforcer.AddFunction("isAdmin", a.isAdminFunc)
	a.enforcer.AddFunction("hasGrantsForResource", a.hasGrantsForResourceFunc)
	a.enforcer.AddFunction("isDefaultResource", a.isDefaultResourceFunc)
}

func (a *CasbinAccessControl) checkGrants(
	user *entity.User,
	product string,
	action auth.Action,
) error {
	if !action.IsValid() {
		return ErrInvalidAccessControlAction
	}

	for _, realmRole := range user.Roles {
		allowed, err := a.enforcer.Enforce(realmRole, user.ProductGrants, product, action.String())
		if err != nil {
			return err
		}

		a.logger.V(2).Info(
			"Checking grants", "user", user.ID, "role", realmRole, "action", action, "product", product, "allowed", allowed,
		)

		if allowed {
			return nil
		}
	}

	return auth.UnauthorizedError{
		Product: product,
		Action:  action,
	}
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

func (a *CasbinAccessControl) isDefaultResourceFunc(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidNumberOfArguments
	}

	res := args[0].(string)

	return res == _defaultResource, nil
}
