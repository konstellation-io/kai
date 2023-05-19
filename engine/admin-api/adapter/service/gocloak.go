package service

import (
	"context"
	"encoding/json"

	"github.com/Nerzal/gocloak/v13"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/internal/errors"
)

type KeycloakConfig struct {
	URL           string
	Realm         string
	MasterRealm   string
	AdminUsername string
	AdminPassword string
}

type GocloakService struct {
	client *gocloak.GoCloak
	token  *gocloak.JWT
	ctx    context.Context
	cfg    *KeycloakConfig
}

func WithClient(keycloakURL string) *gocloak.GoCloak {
	client := gocloak.NewClient(keycloakURL)
	return client
}

func NewGocloakService(
	client *gocloak.GoCloak,
	cfg *KeycloakConfig,
) (*GocloakService, error) {
	wrapErr := errors.Wrapper("new gocloak service: %w")

	ctx := context.Background()
	token, err := client.LoginAdmin(
		ctx,
		cfg.AdminUsername,
		cfg.AdminPassword,
		cfg.MasterRealm,
	)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &GocloakService{
		client: client,
		token:  token,
		ctx:    ctx,
		cfg:    cfg,
	}, nil
}

func (gm *GocloakService) GetUserByID(userID string) (*entity.User, error) {
	wrapErr := errors.Wrapper("gocloak get user by id: %w")

	user, err := gm.client.GetUserByID(gm.ctx, gm.token.AccessToken, gm.cfg.Realm, userID)
	if err != nil {
		return nil, wrapErr(err)
	}

	return gocloakUserToUser(user), nil
}

func (gm *GocloakService) UpdateUserProductGrants(userID, product string, grants []string) error {
	wrapErr := errors.Wrapper("gocloak update user roles: %w")

	user, err := gm.client.GetUserByID(gm.ctx, gm.token.AccessToken, gm.cfg.Realm, userID)
	if err != nil {
		return wrapErr(err)
	}

	if user.Attributes == nil {
		user.Attributes = &map[string][]string{}
	}
	rolesAttribute, ok := (*user.Attributes)["product_roles"]
	if !ok {
		rolesAttribute = make([]string, 1)
	}

	userProductGrants := rolesAttribute[0]
	userGrantsByProduct := make(map[string]interface{})

	if userProductGrants == "" {
		userProductGrants = "{}"
	}
	if err = json.Unmarshal([]byte(userProductGrants), &userGrantsByProduct); err != nil {
		return wrapErr(err)
	}

	if len(grants) == 0 {
		delete(userGrantsByProduct, product)
	} else {
		userGrantsByProduct[product] = grants
	}

	marshalledRoles, err := json.Marshal(userGrantsByProduct)
	if err != nil {
		return wrapErr(err)
	}

	rolesAttribute[0] = string(marshalledRoles)

	(*user.Attributes)["product_roles"] = rolesAttribute

	err = gm.client.UpdateUser(gm.ctx, gm.token.AccessToken, gm.cfg.Realm, *user)
	if err != nil {
		return wrapErr(err)
	}

	return nil
}

func gocloakUserToUser(user *gocloak.User) *entity.User {
	return &entity.User{
		ID:        *user.ID,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		Email:     *user.Email,
		Username:  *user.Username,
		Enabled:   *user.Enabled,
	}
}
