package storage

import "context"

// Driver is the interface every storage backend must satisfy.
type Driver interface {
	// Upload stores data under the given folder and returns its public URL.
	Upload(ctx context.Context, folder string, data []byte) (publicURL string, err error)
}
