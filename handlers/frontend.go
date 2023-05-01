package handlers

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
)

type customFs struct {
	fs       http.FileSystem
	fallback string
}

func (c *customFs) Open(name string) (http.File, error) {
	f, err := c.fs.Open(name)

	if errors.Is(err, os.ErrNotExist) {
		return c.fs.Open(c.fallback)
	}

	return f, err
}

func WithNotFoundFallbackFileSystem(root http.FileSystem, fallback string) http.FileSystem {
	return &customFs{
		fs:       root,
		fallback: fallback,
	}
}

func FileServer(contents fs.FS) http.Handler {
	return http.FileServer(WithNotFoundFallbackFileSystem(http.FS(contents), "index.html"))
}
