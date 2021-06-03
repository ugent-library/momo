package routes

import (
	"fmt"
	"log"
	"net/http"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	gqlplayground "github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
	gqlgraph "github.com/ugent-library/momo/graph"
	gqlgenerated "github.com/ugent-library/momo/graph/generated"
	"github.com/ugent-library/momo/internal/auth/oidc"
	"github.com/ugent-library/momo/internal/controller"
	"github.com/ugent-library/momo/internal/ctx"
	"github.com/ugent-library/momo/internal/engine"
	mw "github.com/ugent-library/momo/internal/middleware"
	"github.com/ugent-library/momo/internal/render"
)

func Register(r chi.Router, e engine.Engine) {
	redirectURL := fmt.Sprintf("https://%s/auth/callback", viper.GetString("base-url"))

	auth, err := oidc.NewProvider(redirectURL)
	if err != nil {
		panic(err)
	}

	recs := controller.NewRecController(e)
	users := controller.NewUserController()
	requireUser := mw.RequireUser(auth.AuthCodeURL())
	setUser := mw.SetUser()
	setLocale := mw.SetLocale(e)

	// general middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// static file server
	r.Mount("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("static"))))

	// robots.txt
	r.Get("/robots.txt", controller.Robots(e))

	// graphql endpoint
	gqlserver := gqlhandler.NewDefaultServer(gqlgenerated.NewExecutableSchema(gqlgenerated.Config{Resolvers: &gqlgraph.Resolver{}}))
	r.Mount("/graphql", gqlserver)
	r.Mount("/graphql/playground", gqlplayground.Handler("GraphQL playground", "/graphql"))

	// OAI-PMH provider
	r.Mount("/oai", controller.OAI(e))

	// oEmbed endpoint
	r.Mount("/oembed", controller.OEmbed())

	// logout
	r.Get("/logout", users.Logout)
	// auth
	r.Get("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		var profile struct {
			Email                    string `json:"email"`
			EmailVerified            bool   `json:"email_verified"`
			Name                     string `json:"name"`
			FamilyName               string `json:"family_name"`
			GivenName                string `json:"given_name"`
			PreferredUsername        string `json:"preferred_username"`
			IdentityProvider         string `json:"identity_provider"`
			IdentityProviderIdentity string `json:"identity_provider_identity"`
		}
		err := auth.Exchange(r.URL.Query().Get("code"), &profile)
		if err != nil {
			log.Panic(err)
		}

		// TODO only store a remember token
		http.SetCookie(w, &http.Cookie{
			Name:     "user",
			Value:    profile.PreferredUsername,
			HttpOnly: true,
			Path:     "/",
		})

		http.Redirect(w, r, "/admin", http.StatusFound)
	})

	uiRoutes := func(r chi.Router) {
		for _, collection := range e.Collections() {
			r.Route("/collection/"+collection.Name, func(r chi.Router) {
				r.Use(setLocale)
				r.Use(setUser)
				r.Use(chimw.WithValue(ctx.CollectionKey, collection.Name))
				r.Use(chimw.WithValue(ctx.ThemeKey, collection.Theme))
				r.Get("/", recs.List)
				r.Get("/search", recs.Search)
				r.Get("/{id}", recs.Show)
				r.Get("/{id}/viewer", recs.Viewer)
			})
		}

		r.Route("/admin", func(r chi.Router) {
			r.Use(setLocale)
			r.Use(setUser)
			r.Use(requireUser)
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				render.Text(w, r, "Admin only")
			})
		})
	}

	for _, loc := range e.Locales() {
		r.Route("/"+loc.Language().String(), func(r chi.Router) {
			r.Use(chimw.WithValue(ctx.LocaleKey, loc))
			uiRoutes(r)
		})
	}

	uiRoutes(r)
}
