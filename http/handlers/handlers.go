// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package handlers

import (
	"log/slog"

	"github.com/gorilla/mux"
	v1 "github.com/tiagomelo/harbor-service/http/handlers/v1"
	"github.com/tiagomelo/harbor-service/services"
)

// ApiMuxConfig struct holds the configuration for the API.
type ApiMuxConfig struct {
	HarborService *services.HarborService
	Log           *slog.Logger
}

// NewApiMux creates and returns a new mux.Router configured with version 1 (v1) routes.
func NewApiMux(c *ApiMuxConfig) *mux.Router {
	return v1.Routes(&v1.Config{
		HarborService: c.HarborService,
		Log:           c.Log,
	})
}
