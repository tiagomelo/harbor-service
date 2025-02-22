// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Logger is a middleware that logs the start and end of each HTTP request along with
// some additional information.
func Logger(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		log.Info("request started",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remoteaddr", r.RemoteAddr),
		)
		next.ServeHTTP(w, r)
		log.Info("request completed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remoteaddr", r.RemoteAddr),
			slog.Duration("since", time.Since(start)),
		)
	})
}

// PreventRequestSmuggling blocks smuggling attempts by validating Transfer-Encoding headers.
func PreventRequestSmuggling(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Transfer-Encoding") == "chunked" && r.ContentLength > 0 {
			http.Error(w, "Invalid Transfer-Encoding header", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
