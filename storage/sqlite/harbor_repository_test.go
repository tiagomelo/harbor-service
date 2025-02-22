// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/harbor-service/domain/harbor"
)

func TestUpsertHarbor(t *testing.T) {
	testCases := []struct {
		name          string
		input         *harbor.Harbor
		mockClosure   func() *sql.DB
		expectedError error
	}{
		{
			name: "happy path",
			input: &harbor.Harbor{
				UNLoc:       "unloc",
				Name:        "name",
				City:        "city",
				Country:     "country",
				Alias:       []string{"alias1", "alias2"},
				Regions:     []string{"region1", "region2"},
				Coordinates: []float64{1.0, 2.0},
				Province:    "province",
				Timezone:    "timezone",
				UNLocs:      []string{"unloc1", "unloc2"},
				Code:        "code",
			},
			mockClosure: func() *sql.DB {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mock.ExpectExec(regexp.QuoteMeta(upsertQuery)).WithArgs(
					"unloc",
					"name",
					"city",
					"country",
					"alias1,alias2",
					"region1,region2",
					"[1,2]",
					"province",
					"timezone",
					"unloc1,unloc2",
					"code",
				).WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
		},
		{
			name: "error",
			input: &harbor.Harbor{
				UNLoc:       "unloc",
				Name:        "name",
				City:        "city",
				Country:     "country",
				Alias:       []string{"alias1", "alias2"},
				Regions:     []string{"region1", "region2"},
				Coordinates: []float64{1.0, 2.0},
				Province:    "province",
				Timezone:    "timezone",
				UNLocs:      []string{"unloc1", "unloc2"},
				Code:        "code",
			},
			mockClosure: func() *sql.DB {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mock.ExpectExec(regexp.QuoteMeta(upsertQuery)).WithArgs(
					"unloc",
					"name",
					"city",
					"country",
					"alias1,alias2",
					"region1,region2",
					"[1,2]",
					"province",
					"timezone",
					"unloc1,unloc2",
					"code",
				).WillReturnError(errors.New("some error"))
				return db
			},
			expectedError: errors.New("upserting harbor: some error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := tc.mockClosure()
			repo := NewHarborRepository(db)
			err := repo.UpsertHarbor(context.TODO(), tc.input)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.expectedError)
				}
			}
		})
	}
}
