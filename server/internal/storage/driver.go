package storage

import "context"

// Driver is the interface every storage backend must satisfy.
type Driver interface {
	// Upload stores data and returns its public URL.
	// ext must include the leading dot, e.g. ".webp", ".mp4".
	// contentType is the MIME type, e.g. "image/webp", "video/mp4".
	Upload(ctx context.Context, folder, ext, contentType string, data []byte) (publicURL string, err error)
}
