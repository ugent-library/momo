package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
	"github.com/ugent-library/momo/internal/controller"
	"github.com/ugent-library/momo/internal/ctx"
	"github.com/ugent-library/momo/internal/engine"
	mw "github.com/ugent-library/momo/internal/middleware"
	"github.com/ugent-library/momo/internal/render"
	"golang.org/x/oauth2"
)

func Register(r chi.Router, e engine.Engine) {
	recs := controller.NewRecController(e)

	// general middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// static file server
	r.Mount("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("static"))))

	// OAI-PMH provider
	r.Mount("/oai", controller.OAI(e))

	// robots.txt
	r.Get("/robots.txt", controller.Robots(e))

	uiRoutes := func(r chi.Router) {
		for _, collection := range e.Collections() {
			r.Route("/collection/"+collection.Name, func(r chi.Router) {
				r.Use(mw.SetLocale(e))
				r.Use(chimw.WithValue(ctx.CollectionKey, collection.Name))
				r.Use(chimw.WithValue(ctx.ThemeKey, collection.Theme))
				r.Get("/", recs.List)
				r.Get("/search", recs.Search)
				r.Get("/{id}", recs.Show)
			})
		}
	}

	for _, loc := range e.Locales() {
		r.Route("/"+loc.Language().String(), func(r chi.Router) {
			r.Use(chimw.WithValue(ctx.LocaleKey, loc))
			uiRoutes(r)
		})
	}

	// test keycloak oidc -->

	provider, err := oidc.NewProvider(context.Background(), viper.GetString("oidc_url"))
	if err != nil {
		log.Panicf("oidc err: %v", err)
	}

	var redirectURL string
	if viper.GetBool("ssl") {
		redirectURL = fmt.Sprintf("https://%s/auth/callback", viper.GetString("host"))
	} else if viper.GetInt("port") != 80 {
		redirectURL = fmt.Sprintf("http://%s:%d/auth/callback", viper.GetString("host"), viper.GetInt("port"))
	} else {
		redirectURL = fmt.Sprintf("http://%s/auth/callback", viper.GetString("host"))
	}

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := oauth2.Config{
		ClientID:     viper.GetString("oidc_client_id"),
		ClientSecret: viper.GetString("oidc_client_secret"),
		RedirectURL:  redirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email", "phone", "address"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: viper.GetString("oidc_client_id")})

	authRequired := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, oauth2Config.AuthCodeURL(""), http.StatusFound)
	}

	r.Get("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		oauth2Token, err := oauth2Config.Exchange(context.Background(), r.URL.Query().Get("code"))
		if err != nil {
			log.Panic(err)
		}

		// Extract the ID Token from OAuth2 token.
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			// handle missing token
			log.Panic("token missing")
		}

		// Parse and verify ID Token payload.
		idToken, err := verifier.Verify(context.Background(), rawIDToken)
		if err != nil {
			// handle error
			log.Panic(err)
		}

		// Extract custom claims
		var claims struct {
			Email                    string `json:"email"`
			EmailVerified            bool   `json:"email_verified"`
			Name                     string `json:"name"`
			FamilyName               string `json:"family_name"`
			GivenName                string `json:"given_name"`
			PreferredUsername        string `json:"preferred_username"`
			IdentityProvider         string `json:"identity_provider"`
			IdentityProviderIdentity string `json:"identity_provider_identity"`
		}
		if err := idToken.Claims(&claims); err != nil {
			// handle error
			log.Panic(err)
		}

		render.JSON(w, r, &claims)
	})

	r.Route("/admin", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			authRequired(w, r)
		})
	})

	// <-- test keycloak oidc

	uiRoutes(r)
}
