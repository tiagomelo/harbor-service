// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package v1_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/harbor-service/domain/harbor"
	"github.com/tiagomelo/harbor-service/http/handlers"
	"github.com/tiagomelo/harbor-service/services"
	"github.com/tiagomelo/harbor-service/storage/sqlite"
)

var (
	testDb     *sql.DB
	testServer *httptest.Server
)

func TestMain(m *testing.M) {
	testDbFilePath := "../../../storage/sqlite/harbor_test.db"
	var err error
	testDb, err = sqlite.Connect(testDbFilePath)
	if err != nil {
		fmt.Println("error when connecting to the test database:", err)
		os.Exit(1)
	}
	harborRepo := sqlite.NewHarborRepository(testDb)
	harborService := services.NewHarborService(harborRepo)
	log := slog.New(slog.NewJSONHandler(io.Discard, nil))
	apiMux := handlers.NewApiMux(&handlers.ApiMuxConfig{
		HarborService: harborService,
		Log:           log,
	})
	testServer = httptest.NewServer(apiMux)
	defer testServer.Close()
	exitVal := m.Run()
	if err := testDb.Close(); err != nil {
		fmt.Println("error when closing test database:", err)
		os.Exit(1)
	}
	if err := os.Remove(testDbFilePath); err != nil {
		fmt.Println("error when deleting test database:", err)
		os.Exit(1)
	}
	os.Exit(exitVal)
}

func TestV1HandleUpsert(t *testing.T) {
	testCases := []struct {
		name            string
		inputFilePath   string
		outputFilePath  string
		expectedHarbors []harbor.Harbor
		expectedStatus  int
	}{
		{
			name:           "insert two new harbors",
			inputFilePath:  "../../../testdata/harbors/input/two_new_harbors.json",
			outputFilePath: "../../../testdata/harbors/output/harbors_upserted.json",
			expectedHarbors: []harbor.Harbor{
				{UNLoc: "AEAJM"},
				{UNLoc: "USLAX"},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Upsert: one new harbor, two existing harbors",
			inputFilePath:  "../../../testdata/harbors/input/one_new_two_existing_harbors.json",
			outputFilePath: "../../../testdata/harbors/output/harbors_upserted.json",
			expectedHarbors: []harbor.Harbor{
				{UNLoc: "AEAJM"},
				{UNLoc: "AEDXB"},
				{UNLoc: "USLAX"},
			},
			expectedStatus: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input, err := os.ReadFile(tc.inputFilePath)
			require.NoError(t, err)
			expectedOutput, err := os.ReadFile(tc.outputFilePath)
			require.NoError(t, err)

			resp, err := http.Post(testServer.URL+"/api/v1/harbors", "application/json", bytes.NewBuffer([]byte(input)))
			require.NoError(t, err)
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedStatus, resp.StatusCode)
			require.JSONEq(t, string(expectedOutput), string(body))

			rows, err := testDb.Query("SELECT unloc FROM harbors")
			require.NoError(t, err)
			defer rows.Close()

			var createdHarbors []harbor.Harbor
			for rows.Next() {
				var h harbor.Harbor
				require.NoError(t, rows.Scan(&h.UNLoc))
				createdHarbors = append(createdHarbors, h)
			}
			require.Len(t, createdHarbors, len(tc.expectedHarbors))
			require.ElementsMatch(t, tc.expectedHarbors, createdHarbors)
			require.NoError(t, rows.Err())
		})
	}
}

func TestSecurityMiddleware_RequestSmuggling(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/v1/harbors", bytes.NewBuffer([]byte("{}")))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Transfer-Encoding", "chunked")
	req.ContentLength = 1 // illegal when using chunked.

	rec := httptest.NewRecorder()
	testServer.Config.Handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "Invalid Transfer-Encoding header")
}
