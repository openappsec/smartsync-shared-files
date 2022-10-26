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

package filesystem

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"openappsec.io/errors"
	"openappsec.io/log"
	"openappsec.io/smartsync-shared-files/internal/models"
)

const (
	fsBaseConfig = "filesystem_db"
	fsConfigRoot = fsBaseConfig + ".root"
	fsConfigTTL  = fsBaseConfig + ".ttl"
)

// Adapter for filesystem ops on local drive
type Adapter struct {
	root string
	ttl  time.Duration
}

// Configuration service interface for fetching config
type Configuration interface {
	GetDuration(key string) (time.Duration, error)
	GetString(key string) (string, error)
}

// NewAdapter creates new adapter
func NewAdapter(conf Configuration) (*Adapter, error) {
	root, err := conf.GetString(fsConfigRoot)
	if err != nil {
		return &Adapter{}, err
	}
	ttl, err := conf.GetDuration(fsConfigTTL)
	if err != nil {
		return &Adapter{}, err
	}
	err = os.MkdirAll(root, 0755)
	if err != nil {
		return &Adapter{}, err
	}
	cleanup(root, ttl)
	return &Adapter{root: root, ttl: ttl}, nil
}

func cleanup(root string, ttl time.Duration) {
	tempFiles := make([]string, 0)
	err := filepath.WalkDir(
		root,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil || d == nil {
				log.Errorf("fail to cleanup. dirEntry: %+v, err: %v", d, err)
				return err
			}
			if d.IsDir() {
				return nil
			}
			if models.IsTempFile(path) {
				tempFiles = append(tempFiles, path)
			}
			return nil
		},
	)
	if err != nil {
		log.Warnf("failed to cleanup old temp files")
	}
	go func(root string, ttl time.Duration) {
		time.Sleep(ttl)
		for _, filePath := range tempFiles {
			err := os.Remove(filePath)
			if err != nil && !os.IsNotExist(err) {
				log.Warnf("failed to cleanup file: %v", filePath)
			}
		}
	}(root, ttl)
}

// GetFilesList return a list of files with their last modified time
func (a *Adapter) GetFilesList(ctx context.Context, pathPrefix string) ([]models.FileMetadata, error) {
	log.WithContext(ctx).Infof("list files with prefix: %v", pathPrefix)
	matches, err := filepath.Glob(a.root + pathPrefix + "*")
	if err != nil {
		return []models.FileMetadata{}, err
	}
	log.WithContext(ctx).Infof("matched: %v", matches)
	var files []models.FileMetadata
	for _, match := range matches {
		log.WithContext(ctx).Debugf("walk dir: %v", match)
		err = filepath.WalkDir(
			match,
			func(path string, d fs.DirEntry, err error) error {
				if d.IsDir() {
					return nil
				}
				fileInfo, err := d.Info()
				if err != nil {
					return err
				}
				log.WithContext(ctx).Debugf("adding file: %v to response", path)
				path = path[len(a.root):]
				log.WithContext(ctx).Infof("%v", path)
				files = append(files, models.FileMetadata{Path: path, LastModified: fileInfo.ModTime()})
				return nil
			},
		)
	}
	return files, nil
}

// GetFile return file content
func (a *Adapter) GetFile(ctx context.Context, path string) ([]byte, error) {
	log.WithContext(ctx).Debugf("get file: %v", path)
	data, err := os.ReadFile(a.root + path)
	if err != nil {
		if os.IsNotExist(err) {
			log.WithContext(ctx).Warnf("file %v not found", path)
			err = errors.Wrap(err, "file not found").SetClass(errors.ClassNotFound)
		} else {
			log.WithContext(ctx).Errorf("failed to read file %v", path)
		}
		return []byte{}, err
	}
	log.WithContext(ctx).Debugf("got file, file length %v", len(data))
	return data, nil
}

// PutFile write a file, set ttl if isTemp is true
func (a *Adapter) PutFile(ctx context.Context, path string, content []byte, isTemp bool) error {
	log.WithContext(ctx).Debugf("put file: %v, length: %v", path, len(content))

	if err := os.MkdirAll(a.root+filepath.Dir(path), 0750); err != nil && !os.IsExist(err) {
		return err
	}
	if err := os.WriteFile(a.root+path, content, 0666); err != nil {
		log.WithContext(ctx).Errorf("failed to put file: %v", err)
		return err
	}
	if isTemp {
		go func() {
			time.Sleep(a.ttl)
			log.WithContext(ctx).Debugf("ttl expired for file: %v", path)
			err := os.Remove(a.root + path)
			if err != nil {
				log.WithContext(ctx).Warnf("failed to remove file %v. err: %v", path, err)
			}
		}()
	}
	return nil
}
