package service

import (
	"context"
	"encoding/json"

	"github.com/Nerzal/gocloak/v13"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/internal/errors"
)

type GocloakManagerClient struct {
	client *gocloak.GoCloak
	token  *gocloak.JWT
	ctx    context.Context
	cfg    *config.Config
}

func NewGocloakManager(cfg *config.Config) (*GocloakManagerClient, error) {
	wrapErr := errors.Wrapper("new gocloak manager: %w")

	client := gocloak.NewClient(cfg.Keycloak.URL)
	ctx := context.Background()
	token, err := client.LoginAdmin(
		ctx, cfg.Keycloak.AdminUsername,
		cfg.Keycloak.AdminPassword,
		cfg.Keycloak.MasterRealm,
	)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &GocloakManagerClient{
		client: client,
		token:  token,
		ctx:    ctx,
		cfg:    cfg,
	}, nil
}

func (gm *GocloakManagerClient) CreateUser(userData entity.UserGocloakData) error {
	wrapErr := errors.Wrapper("gocloak create user: %w")

	user := gocloak.User{
		FirstName: gocloak.StringP(userData.FirstName),
		LastName:  gocloak.StringP(userData.LastName),
		Email:     gocloak.StringP(userData.Email),
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP(userData.Username),
	}

	_, err := gm.client.CreateUser(gm.ctx, gm.token.AccessToken, gm.cfg.Keycloak.Realm, user)
	if err != nil {
		return wrapErr(err)
	}

	return nil
}

func (gm *GocloakManagerClient) GetUserByID(userID string) (entity.UserGocloakData, error) {
	wrapErr := errors.Wrapper("gocloak get user by id: %w")

	user, err := gm.client.GetUserByID(gm.ctx, gm.token.AccessToken, gm.cfg.Keycloak.Realm, userID)
	if err != nil {
		return entity.UserGocloakData{}, wrapErr(err)
	}

	return gocloakUserToUserData(user), nil
}

func (gm *GocloakManagerClient) UpdateUserProductPermissions(userID string, product string, permissions []string) error {
	wrapErr := errors.Wrapper("gocloak update user roles: %w")

	user, err := gm.client.GetUserByID(gm.ctx, gm.token.AccessToken, gm.cfg.Keycloak.Realm, userID)
	if err != nil {
		wrapErr(err)
	}

	rolesAttribute, ok := (*user.Attributes)["product_roles"]
	if !ok {
		rolesAttribute = make([]string, 0, 1)
	}

	userProductRoles := rolesAttribute[0]

	userRoles := make(map[string]interface{})
	json.Unmarshal([]byte(userProductRoles), &userRoles)

	if len(permissions) == 0 {
		delete(userRoles, product)
	} else {
		userRoles[product] = permissions
	}

	marshalledRoles, err := json.Marshal(userRoles)
	if err != nil {
		return wrapErr(err)
	}

	rolesAttribute[0] = string(marshalledRoles)

	(*user.Attributes)["product_roles"] = rolesAttribute

	err = gm.client.UpdateUser(gm.ctx, gm.token.AccessToken, gm.cfg.Keycloak.Realm, *user)
	if err != nil {
		return wrapErr(err)
	}

	return nil
}

func gocloakUserToUserData(user *gocloak.User) entity.UserGocloakData {
	return entity.UserGocloakData{
		ID:        *user.ID,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		Email:     *user.Email,
		Username:  *user.Username,
		Enabled:   *user.Enabled,
	}
}

func userDataToGocloak(userData entity.UserGocloakData) gocloak.User {
	return gocloak.User{
		ID:        gocloak.StringP(userData.ID),
		Username:  gocloak.StringP(userData.Username),
		Email:     gocloak.StringP(userData.Email),
		FirstName: gocloak.StringP(userData.FirstName),
		LastName:  gocloak.StringP(userData.LastName),
		Enabled:   gocloak.BoolP(userData.Enabled),
	}
}
