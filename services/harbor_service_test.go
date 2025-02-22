// Copyright (c) 2025 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/harbor-service/domain/harbor"
)

func TestUpsertHarbor(t *testing.T) {
	testCases := []struct {
		name          string
		mockClosure   func(*mockHarborRepository)
		expectedError error
	}{
		{
			name:        "happy path",
			mockClosure: func(m *mockHarborRepository) {},
		},
		{
			name: "error",
			mockClosure: func(m *mockHarborRepository) {
				m.upsertHarborErr = errors.New("updser harbor error")
			},
			expectedError: errors.New("updser harbor error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mockHarborRepository)
			tc.mockClosure(mockRepo)
			s := NewHarborService(mockRepo)
			err := s.UpsertHarbor(context.Background(), &harbor.Harbor{})
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

type mockHarborRepository struct {
	upsertHarborErr error
}

func (m *mockHarborRepository) UpsertHarbor(ctx context.Context, harbor *harbor.Harbor) error {
	return m.upsertHarborErr
}
