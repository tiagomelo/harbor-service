// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package harbor

import (
	"context"
)

// HarborRepository defines the interface that a harbor repository should implement.
type HarborRepository interface {
	// UpsertHarbor inserts or updates a harbor in the repository.
	UpsertHarbor(ctx context.Context, harbor *Harbor) error
}
