// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/tiagomelo/harbor-service/domain/harbor"
)

// upsertQuery is the SQL query to upsert a harbor into the database.
const upsertQuery = `
INSERT INTO harbors (unloc, name, city, country, alias, regions, coordinates, province, timezone, unlocs, code)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(unloc) DO UPDATE SET
	name = excluded.name,
	city = excluded.city,
	country = excluded.country,
	alias = excluded.alias,
	regions = excluded.regions,
	coordinates = excluded.coordinates,
	province = excluded.province,
	timezone = excluded.timezone,
	unlocs = excluded.unlocs,
	code = excluded.code;
`

// harborSQLiteRepository is a SQLite implementation of the harbor.Repository interface.
type harborSQLiteRepository struct {
	db *sql.DB
}

// NewHarborRepository creates a new harbor repository.
func NewHarborRepository(db *sql.DB) *harborSQLiteRepository {
	return &harborSQLiteRepository{db: db}
}

// UpsertHarbor upserts a harbor into the database.
func (r *harborSQLiteRepository) UpsertHarbor(ctx context.Context, harbor *harbor.Harbor) error {
	flattenAliases := strings.Join(harbor.Alias, ",")
	flattenRegions := strings.Join(harbor.Regions, ",")
	flattenCoordinates := strings.Join(strings.Fields(fmt.Sprint(harbor.Coordinates)), ",")
	flattenUNLocs := strings.Join(harbor.UNLocs, ",")
	if _, err := r.db.ExecContext(ctx, upsertQuery,
		harbor.UNLoc,
		harbor.Name,
		harbor.City,
		harbor.Country,
		flattenAliases,
		flattenRegions,
		flattenCoordinates,
		harbor.Province,
		harbor.Timezone,
		flattenUNLocs,
		harbor.Code,
	); err != nil {
		return errors.Wrap(err, "upserting harbor")
	}
	return nil
}
