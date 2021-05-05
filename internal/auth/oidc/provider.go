package oidc

import (
	"context"
	"errors"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type Provider struct {
	oauthClient   *oauth2.Config
	tokenVerifier *oidc.IDTokenVerifier
}

func NewProvider(callbackURL string) (*Provider, error) {
	oidcProvider, err := oidc.NewProvider(context.Background(), viper.GetString("oidc_url"))
	if err != nil {
		return nil, err
	}

	// Configure an OpenID Connect aware OAuth2 client.
	oauthClient := &oauth2.Config{
		ClientID:     viper.GetString("oidc_client_id"),
		ClientSecret: viper.GetString("oidc_client_secret"),
		RedirectURL:  callbackURL,
		// Discovery returns the OAuth2 endpoints.
		Endpoint: oidcProvider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email", "phone", "address"},
	}

	tokenVerifier := oidcProvider.Verifier(&oidc.Config{ClientID: viper.GetString("oidc_client_id")})

	p := Provider{
		oauthClient:   oauthClient,
		tokenVerifier: tokenVerifier,
	}

	return &p, nil
}

func (p *Provider) AuthCodeURL() string {
	return p.oauthClient.AuthCodeURL("")
}

func (p *Provider) Exchange(code string, profile interface{}) error {
	ctx := context.Background()
	oauthToken, err := p.oauthClient.Exchange(ctx, code)
	if err != nil {
		return err
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauthToken.Extra("id_token").(string)
	if !ok {
		// handle missing token
		return errors.New("token missing")
	}

	// Parse and verify ID Token payload.
	idToken, err := p.tokenVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return err
	}

	return idToken.Claims(&profile)
}
