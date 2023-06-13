package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Nerzal/gocloak/v13"

	"github.com/konstellation-io/kai/engine/admin-api/internal/errors"
)

type KeycloakConfig struct {
	Realm         string
	MasterRealm   string
	AdminUsername string
	AdminPassword string
	AdminClientID string
}

type GocloakUserRegistry struct {
	client                *gocloak.GoCloak
	token                 *gocloak.JWT
	ctx                   context.Context
	cfg                   *KeycloakConfig
	tokenExpiresAt        time.Time
	refreshtokenExpiresAt time.Time
}

func WithClient(keycloakURL string) *gocloak.GoCloak {
	client := gocloak.NewClient(keycloakURL)
	return client
}

func NewGocloakUserRegistry(
	client *gocloak.GoCloak,
	cfg *KeycloakConfig,
) (*GocloakUserRegistry, error) {
	wrapErr := errors.Wrapper("new gocloak user registry: %w")

	ctx := context.Background()
	now := time.Now()

	token, err := client.LoginAdmin(
		ctx,
		cfg.AdminUsername,
		cfg.AdminPassword,
		cfg.MasterRealm,
	)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &GocloakUserRegistry{
		client:                client,
		token:                 token,
		ctx:                   ctx,
		cfg:                   cfg,
		tokenExpiresAt:        now.Add(time.Duration(token.ExpiresIn) * time.Second),
		refreshtokenExpiresAt: now.Add(time.Duration(token.RefreshExpiresIn) * time.Second),
	}, nil
}

func (gm *GocloakUserRegistry) refreshToken() error {
	wrapErr := errors.Wrapper("refresh token: %w")
	now := time.Now()

	if now.Before(gm.tokenExpiresAt) {
		return nil
	}

	var (
		token *gocloak.JWT
		err   error
	)

	if now.Before(gm.refreshtokenExpiresAt) {
		token, err = gm.client.RefreshToken(
			gm.ctx,
			gm.token.RefreshToken,
			gm.cfg.AdminClientID,
			"",
			gm.cfg.MasterRealm,
		)

		if err != nil {
			return wrapErr(err)
		}
	} else {
		token, err = gm.client.LoginAdmin(
			gm.ctx,
			gm.cfg.AdminUsername,
			gm.cfg.AdminPassword,
			gm.cfg.MasterRealm,
		)

		if err != nil {
			return wrapErr(err)
		}
	}

	gm.token = token
	gm.tokenExpiresAt = now.Add(time.Duration(token.ExpiresIn) * time.Second)
	gm.refreshtokenExpiresAt = now.Add(time.Duration(token.RefreshExpiresIn) * time.Second)

	return nil
}

func (gm *GocloakUserRegistry) UpdateUserProductGrants(userID, product string, grants []string) error {
	wrapErr := errors.Wrapper("gocloak update user roles: %w")

	err := gm.refreshToken()
	if err != nil {
		return wrapErr(err)
	}

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

	err = gm.refreshToken()
	if err != nil {
		return wrapErr(err)
	}

	err = gm.client.UpdateUser(gm.ctx, gm.token.AccessToken, gm.cfg.Realm, *user)
	if err != nil {
		return wrapErr(err)
	}

	return nil
}
