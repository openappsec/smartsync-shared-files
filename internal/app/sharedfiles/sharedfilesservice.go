// Copyright (C) 2022 Check Point Software Technologies Ltd. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sharedfiles

import (
	"context"

	"openappsec.io/smartsync-shared-files/internal/models"
)

// FileSystem exposes an interface for fs service operations
// To create mock for unittest please use this command
// mockgen -destination mocks/mock_FileSystem.go -package mocks openappsec.io/smartsync-shared-files/internal/app/sharedfiles FileSystem
type FileSystem interface {
	GetFilesList(ctx context.Context, pathPrefix string) ([]models.FileMetadata, error)
	GetFile(ctx context.Context, path string) ([]byte, error)
	PutFile(ctx context.Context, path string, content []byte, isTemp bool) error
}

// Service struct
type Service struct {
	fs FileSystem
}

// NewSharedFilesService returns a new instance of a demo service.
func NewSharedFilesService(fs FileSystem) (*Service, error) {
	return &Service{fs: fs}, nil
}
