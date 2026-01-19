package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LocalStorage implements Storage interface using local filesystem
type LocalStorage struct {
	basePath string
	baseURL  string
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(basePath, baseURL string) (*LocalStorage, error) {
	// Ensure base path exists
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve storage path: %w", err)
	}
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &LocalStorage{
		basePath: absPath,
		baseURL:  baseURL,
	}, nil
}

// safePath validates and returns a safe path within the base directory.
// Returns an error if the path attempts to escape the base directory.
func (s *LocalStorage) safePath(path string) (string, error) {
	// Clean the path to remove any .. or . components
	cleanPath := filepath.Clean(path)

	// Ensure the path doesn't start with / or contain ..
	if strings.HasPrefix(cleanPath, "/") || strings.HasPrefix(cleanPath, "..") || strings.Contains(cleanPath, "/../") {
		return "", fmt.Errorf("invalid path: path traversal not allowed")
	}

	// Join with base path and verify it's still within base
	fullPath := filepath.Join(s.basePath, cleanPath)

	// Verify the resolved path is within the base directory
	if !strings.HasPrefix(fullPath, s.basePath+string(os.PathSeparator)) && fullPath != s.basePath {
		return "", fmt.Errorf("invalid path: outside storage directory")
	}

	return fullPath, nil
}

// Upload stores a file to local filesystem
func (s *LocalStorage) Upload(ctx context.Context, path string, reader io.Reader) error {
	fullPath, err := s.safePath(path)
	if err != nil {
		return err
	}

	// Ensure parent directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy content
	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Download retrieves a file from local filesystem
func (s *LocalStorage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath, err := s.safePath(path)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return file, nil
}

// Delete removes a file from local filesystem
func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath, err := s.safePath(path)
	if err != nil {
		return err
	}
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// Exists checks if a file exists on local filesystem
func (s *LocalStorage) Exists(ctx context.Context, path string) (bool, error) {
	fullPath, err := s.safePath(path)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("failed to check file existence: %w", err)
}

// GetURL returns a URL to access the file
func (s *LocalStorage) GetURL(path string) string {
	return s.baseURL + "/" + path
}
