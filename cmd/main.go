// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"github.com/tiagomelo/harbor-service/http/handlers"
	"github.com/tiagomelo/harbor-service/services"
	"github.com/tiagomelo/harbor-service/storage/sqlite"
)

// options represents the command line options.
type options struct {
	Port int `short:"p" long:"port" description:"server's port" required:"true"`
}

func run(port int, log *slog.Logger) error {
	ctx := context.Background()
	defer log.InfoContext(ctx, "completed")

	// setup database.
	const dbFilePath = "storage/sqlite/harbor.db"
	db, err := sqlite.Connect(dbFilePath)
	if err != nil {
		return errors.Wrapf(err, "opening database file %s", dbFilePath)
	}
	defer db.Close()

	// setup services and API.
	harborRepo := sqlite.NewHarborRepository(db)
	harborService := services.NewHarborService(harborRepo)
	apiMux := handlers.NewApiMux(&handlers.ApiMuxConfig{
		HarborService: harborService,
		Log:           log,
	})

	// setup server.
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: apiMux,
	}

	// manage server's lifecycle.
	return handleServerLifecycle(ctx, &srv, db, log)
}

// handleServerLifecycle manages the server's start, shutdown, and signal handling.
func handleServerLifecycle(ctx context.Context, srv *http.Server, db *sql.DB, log *slog.Logger) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	go func() {
		log.Info(fmt.Sprintf("API listening on %s", srv.Addr))
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")
	case sig := <-shutdown:
		log.InfoContext(ctx, fmt.Sprintf("Starting shutdown: %v", sig))
	}

	// graceful shutdown.
	const shutdownTimeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		srv.Close()
		return errors.Wrap(err, "could not stop server gracefully")
	}
	if err := db.Close(); err != nil {
		return errors.Wrap(err, "could not close database connection")
	}
	return nil
}

func main() {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(opts.Port, log); err != nil {
		log.Error("error", slog.Any("err", err))
		os.Exit(1)
	}
}
