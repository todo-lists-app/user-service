package user

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/go-resty/resty/v2"
	ConfigBuilder "github.com/keloran/go-config"
	"strings"
	"time"
)

type GocloakInterface interface {
	LoginClient(ctx context.Context, clientID, clientSecret, realm string) (*gocloak.JWT, error)
	GetUserByID(ctx context.Context, token, realm, userID string) (*gocloak.User, error)
	RetrospectToken(ctx context.Context, token, clientID, clientSecret, realm string) (*gocloak.IntroSpectTokenResult, error)
}

func DeleteKeyCloakUser(ctx context.Context, cfg ConfigBuilder.Config, userId, accessToken string) error {
	client := gocloak.NewClient(cfg.Keycloak.Host)
	cond := func(resp *resty.Response, err error) bool {
		if resp != nil && resp.IsError() {
			if e, ok := resp.Error().(*gocloak.HTTPErrorResponse); ok {
				msg := e.String()
				return strings.Contains(msg, "Cached clientScope not found") || strings.Contains(msg, "unknown_error")
			}
		}
		return false
	}
	rest := client.RestyClient()
	rest.SetRetryCount(10).SetRetryWaitTime(2 * time.Second).AddRetryCondition(cond)
	token, err := client.LoginClient(ctx, cfg.Keycloak.Client, cfg.Keycloak.Secret, cfg.Keycloak.Realm)
	if err != nil {
		return logs.Errorf("error logging in: %v", err)
	}

	logs.Infof("user would have been deleted: %v", userId)

	if err := client.DeleteUser(ctx, token.AccessToken, cfg.Keycloak.Realm, userId); err != nil {
		return logs.Errorf("error deleting user: %v", err)
	}

	return nil
}
