// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package v1

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tiagomelo/harbor-service/http/handlers/v1/harbors"
	"github.com/tiagomelo/harbor-service/http/middleware"
	"github.com/tiagomelo/harbor-service/services"
)

// Config struct holds the database connection and logger.
type Config struct {
	HarborService *services.HarborService
	Log           *slog.Logger
}

// Routes initializes and returns a new router with configured routes.
func Routes(c *Config) *mux.Router {
	router := mux.NewRouter()
	initializeRoutes(c.HarborService, router)
	router.Use(
		func(h http.Handler) http.Handler {
			return middleware.Logger(c.Log, h)
		},
		middleware.PreventRequestSmuggling,
	)
	return router
}

// initializeRoutes sets up the routes for harbor operations.
func initializeRoutes(harborService *services.HarborService, router *mux.Router) {
	harborHandlers := harbors.New(harborService)
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/harbors", harborHandlers.HandleUpsert).Methods(http.MethodPost)
}
