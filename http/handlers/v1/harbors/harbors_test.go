// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package harbors

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/harbor-service/domain/harbor"
)

func TestHandleUpsert(t *testing.T) {
	testCases := []struct {
		name               string
		input              string
		mockClosure        func(rc *mockResponseController, hs *mockHarborService)
		expectedOutput     string
		expectedStatusCode int
	}{
		{
			name:               "valid input",
			input:              `{ "USLAX": { "name": "Los Angeles", "city": "Los Angeles", "country": "USA", "alias": [], "regions": [], "coordinates": [-118.2437, 34.0522], "province": "California", "timezone": "America/Los_Angeles", "unlocs": ["USLAX"], "code": "53001" } }`,
			mockClosure:        func(rc *mockResponseController, hs *mockHarborService) {},
			expectedOutput:     `{"message":"harbors upserted"}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "missing opening `{`",
			input:              `"USLAX": { "name": "Los Angeles", "city": "Los Angeles", "country": "USA" } }`,
			mockClosure:        func(rc *mockResponseController, hs *mockHarborService) {},
			expectedOutput:     `{"error":"invalid JSON: expected '{' at start"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "JSON key is not a string (forced non-string key)",
			input:              `{ true: { "name": "Los Angeles", "city": "Los Angeles", "country": "United States" }}`,
			mockClosure:        func(rc *mockResponseController, hs *mockHarborService) {},
			expectedOutput:     `{"error":"invalid JSON key"}`,
			expectedStatusCode: http.StatusBadRequest,
		},

		{
			name:               "JSON structure is broken (invalid harbor data format)",
			input:              `{ "USLAX": [ "Los Angeles", "USA" ] }`,
			mockClosure:        func(rc *mockResponseController, hs *mockHarborService) {},
			expectedOutput:     `{"error":"invalid JSON harbor structure"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "validation failure (missing required fields)",
			input:              `{ "USLAX": { "city": "Los Angeles" } }`,
			mockClosure:        func(rc *mockResponseController, hs *mockHarborService) {},
			expectedOutput:     `{"error":"[{\"field\":\"name\",\"error\":\"name is a required field\"},{\"field\":\"country\",\"error\":\"country is a required field\"}]"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "database failure on upsert",
			input: `{ "USLAX": { "name": "Los Angeles", "city": "Los Angeles", "country": "USA", "alias": [], "regions": [], "coordinates": [-118.2437, 34.0522], "province": "California", "timezone": "America/Los_Angeles", "unlocs": ["USLAX"], "code": "53001" } }`,
			mockClosure: func(rc *mockResponseController, hs *mockHarborService) {
				hs.UpsertHarborErr = errors.New("error upserting harbor")
			},
			expectedOutput:     `{"error":"error upserting harbor"}`,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "missing closing `}`",
			input:              `{ "USLAX": { "name": "Los Angeles", "city": "Los Angeles", "country": "USA" }`,
			mockClosure:        func(rc *mockResponseController, hs *mockHarborService) {},
			expectedOutput:     `{"error":"invalid JSON: expected '}' at end"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "error flushing response",
			input: `{ "UNLOC": {"name":"Los Angeles","city":"Los Angeles","country":"USA","alias":[],"regions":[],"coordinates":[-118.2437,34.0522],"province":"California","timezone":"America/Los_Angeles","unlocs":["USLAX"],"code":"53001"}}`,
			mockClosure: func(rc *mockResponseController, hs *mockHarborService) {
				rc.FlushErr = errors.New("error flushing response")
			},
			expectedOutput:     `{"error":"error flushing response"}`,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rc := new(mockResponseController)
			hs := new(mockHarborService)
			tc.mockClosure(rc, hs)
			newHttpResponseController = func(_ http.ResponseWriter) responseController {
				return rc
			}
			req, err := http.NewRequest(http.MethodPost, "/api/v1/harbors", bytes.NewBuffer([]byte(tc.input)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			h := New(hs)
			handler := http.HandlerFunc(h.HandleUpsert)
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedStatusCode, rr.Code)
			require.JSONEq(t, tc.expectedOutput, rr.Body.String())
		})
	}
}

type mockResponseController struct {
	FlushErr error
}

func (m *mockResponseController) Flush() error {
	return m.FlushErr
}

type mockHarborService struct {
	UpsertHarborErr error
}

func (m *mockHarborService) UpsertHarbor(ctx context.Context, harbor *harbor.Harbor) error {
	return m.UpsertHarborErr
}

func TestReadKey_NonStringToken(t *testing.T) {
	r := strings.NewReader("true")
	dec := json.NewDecoder(r)

	key, err := readKey(dec)
	require.Error(t, err)
	require.Equal(t, "", key)
	require.Equal(t, "invalid JSON key", err.Error())
}
