package storage

import (
	"context"
	"io"
)

// Storage defines the interface for file storage operations
type Storage interface {
	// Upload stores a file and returns the path
	Upload(ctx context.Context, path string, reader io.Reader) error

	// Download retrieves a file
	Download(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete removes a file
	Delete(ctx context.Context, path string) error

	// Exists checks if a file exists
	Exists(ctx context.Context, path string) (bool, error)

	// GetURL returns a URL to access the file
	GetURL(path string) string
}
