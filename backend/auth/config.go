package auth

import (
	"fmt"
	"os"
	"tranquil-pages/errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	OAuthConfig *oauth2.Config
)

func getRedirectURL() (string, error) {
	redicrect_route := "/auth/callback"

	if baseurl, ok := os.LookupEnv("BASE_URL"); ok {
		return baseurl + redicrect_route, nil
	}

	containerAppName, ok := os.LookupEnv("CONTAINER_APP_NAME")
	if !ok {
		return "", errors.ErrEnvNotSet("CONTAINER_APP_NAME")
	}

	containerAppEnvDnsSuffix, ok := os.LookupEnv("CONTAINER_APP_ENV_DNS_SUFFIX")
	if !ok {
		return "", errors.ErrEnvNotSet("CONTAINER_APP_ENV_DNS_SUFFIX")
	}

	return fmt.Sprintf("https://%s.%s/%s", containerAppName, containerAppEnvDnsSuffix, redicrect_route), nil
}

func InitOAuthConfig() error {
	clientId, ok := os.LookupEnv("OAUTH_CLIENT_ID")
	if !ok {
		return errors.ErrEnvNotSet("OAUTH_CLIENT_ID")
	}

	clientSecret, ok := os.LookupEnv("OAUTH_CLIENT_SECRET")
	if !ok {
		return errors.ErrEnvNotSet("OAUTH_CLIENT_SECRET")
	}

	redirectUrl, err := getRedirectURL()
	if err != nil {
		return err
	}

	OAuthConfig = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return nil
}
