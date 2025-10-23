package admin

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Odyssey-Classic/server/internal/data"
	"github.com/Odyssey-Classic/server/internal/services/admin/maps"
)

// API represents the main admin API structure
type API struct {
	router  chi.Router
	mapsAPI *maps.API
}

// New creates a new Admin API instance
func api(root data.Root) *API {
	api := &API{
		router:  chi.NewRouter(),
		mapsAPI: maps.NewFileBacked(root.MapsDir()),
	}

	api.setupMiddleware()
	api.setupRoutes()

	return api
}

// setupMiddleware configures common middleware for the admin API
func (a *API) setupMiddleware() {
	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.RealIP)
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)
	a.router.Use(middleware.SetHeader("Content-Type", "application/json"))
}

// setupRoutes configures all API routes
func (a *API) setupRoutes() {
	a.router.Route("/admin", func(r chi.Router) {
		// Mount maps API under /admin/maps
		r.Mount("/maps", a.mapsAPI.Routes())

		// Future admin endpoints can be added here
		// r.Mount("/users", a.usersAPI.Routes())
		// r.Mount("/settings", a.settingsAPI.Routes())
	})
}

// Routes returns the configured router
func (a *API) Routes() chi.Router {
	return a.router
}

// ServeHTTP implements http.Handler interface
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
