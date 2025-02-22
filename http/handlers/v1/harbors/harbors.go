// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package harbors

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/tiagomelo/harbor-service/domain/harbor"
	"github.com/tiagomelo/harbor-service/http/web"
	"github.com/tiagomelo/harbor-service/validate"
)

// UpsertHarborRequest represents the request to upsert a harbor.
type UpsertHarborRequest struct {
	UNLoc       string    `json:"unloc"`
	Name        string    `json:"name" validate:"required"`
	City        string    `json:"city" validate:"required"`
	Country     string    `json:"country" validate:"required"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	UNLocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
}

// ToHarbor converts the UpsertHarborRequest to a harbor.Harbor.
func (u *UpsertHarborRequest) ToHarbor() *harbor.Harbor {
	return &harbor.Harbor{
		UNLoc:       u.UNLoc,
		Name:        u.Name,
		City:        u.City,
		Country:     u.Country,
		Alias:       u.Alias,
		Regions:     u.Regions,
		Coordinates: u.Coordinates,
		Province:    u.Province,
		Timezone:    u.Timezone,
		UNLocs:      u.UNLocs,
		Code:        u.Code,
	}
}

// UpsertHarborResponse represents the response to upsert a harbor.
type UpsertHarborResponse struct {
	Message string `json:"message"`
}

// responseController is an interface that wraps the Flush method.
type responseController interface {
	Flush() error
}

// harborService defines the interface for the harbor service.
type harborService interface {
	UpsertHarbor(ctx context.Context, harbor *harbor.Harbor) error
}

// handlers defines the handlers for managing harbors.
type handlers struct {
	service harborService
}

// New creates the handlers for managing harbors.
func New(service harborService) *handlers {
	return &handlers{
		service: service,
	}
}

// maxBufferedReaderSize is the maximum size of the buffered reader.
const maxBufferedReaderSize = 32 * 1024

// For ease of unit testing.
var (
	// newHttpResponseController is a function that creates a new response controller.
	newHttpResponseController = func(rw http.ResponseWriter) responseController {
		return http.NewResponseController(rw)
	}
)

// handlerError is a custom error type that carries an HTTP status code.
type handlerError struct {
	code int
	msg  string
}

func (he handlerError) Error() string {
	return he.msg
}

// readKey reads the next token and returns it as a string.
func readKey(dec *json.Decoder) (string, error) {
	tok, err := dec.Token()
	if err != nil {
		return "", err
	}
	if key, ok := tok.(string); ok {
		return key, nil
	}
	return "", errors.New("invalid JSON key")
}

// HandleUpsert handles the upsert of harbors in a streaming fashion.
func (h *handlers) HandleUpsert(w http.ResponseWriter, r *http.Request) {
	ctr := newHttpResponseController(w)
	bufReader := bufio.NewReaderSize(r.Body, maxBufferedReaderSize)
	dec := json.NewDecoder(bufReader)
	// check for opening '{'.
	if err := h.readExpectedToken(dec, json.Delim('{')); err != nil {
		web.RespondWithError(w, http.StatusBadRequest, "invalid JSON: expected '{' at start")
		return
	}
	// process each harbor in the JSON object.
	if herr := h.processHarbors(r.Context(), dec); herr != nil {
		web.RespondWithError(w, herr.code, herr.Error())
		return
	}
	// check for closing '}'.
	if err := h.readExpectedToken(dec, json.Delim('}')); err != nil {
		web.RespondWithError(w, http.StatusBadRequest, "invalid JSON: expected '}' at end")
		return
	}
	// flush response and finalize.
	if err := ctr.Flush(); err != nil {
		web.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	web.RespondAfterFlush(w, UpsertHarborResponse{Message: "harbors upserted"})
}

// readExpectedToken reads the next token and verifies it matches the expected delimiter.
func (h *handlers) readExpectedToken(dec *json.Decoder, expected json.Delim) error {
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	if tok != expected {
		return fmt.Errorf("unexpected token: got %v, expected %v", tok, expected)
	}
	return nil
}

// processHarbors iterates through each harbor entry in the JSON and delegates processing.
func (h *handlers) processHarbors(ctx context.Context, dec *json.Decoder) *handlerError {
	for dec.More() {
		if herr := h.processSingleHarbor(ctx, dec); herr != nil {
			return herr
		}
	}
	return nil
}

// processSingleHarbor handles processing of a single harbor entry.
func (h *handlers) processSingleHarbor(ctx context.Context, dec *json.Decoder) *handlerError {
	key, err := readKey(dec)
	if err != nil {
		return &handlerError{http.StatusBadRequest, "invalid JSON key"}
	}
	var req UpsertHarborRequest
	if err = dec.Decode(&req); err != nil {
		return &handlerError{http.StatusBadRequest, "invalid JSON harbor structure"}
	}
	if err = validate.Check(req); err != nil {
		return &handlerError{http.StatusBadRequest, err.Error()}
	}
	req.UNLoc = key
	if err = h.service.UpsertHarbor(ctx, req.ToHarbor()); err != nil {
		return &handlerError{http.StatusInternalServerError, "error upserting harbor"}
	}
	return nil
}
