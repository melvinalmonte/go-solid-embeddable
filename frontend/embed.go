package frontend

import (
	"embed"
	"io/fs"
)

//go:generate npm install
//go:generate npm run build

//go:embed dist/*
var content embed.FS

func Contents() (fs.FS, error) {
	return fs.Sub(content, "dist")
}
