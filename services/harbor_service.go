// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package services

import (
	"context"

	"github.com/tiagomelo/harbor-service/domain/harbor"
)

// HarborService is a service that provides operations on Harbors.
type HarborService struct {
	repository harbor.HarborRepository
}

// NewHarborService initializes a new HarborService.
func NewHarborService(repo harbor.HarborRepository) *HarborService {
	return &HarborService{repository: repo}
}

// UpsertHarbor upserts a Harbor.
func (s *HarborService) UpsertHarbor(ctx context.Context, h *harbor.Harbor) error {
	if err := s.repository.UpsertHarbor(ctx, h); err != nil {
		return err
	}
	return nil
}
