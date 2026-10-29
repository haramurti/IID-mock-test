// Package web embeds the built React app (frontend/dist) into the Go binary.
//
// `make build` (and the Dockerfile) copy frontend/dist into backend/web/dist
// before `go build`. Without that step only .gitkeep is embedded and the
// server answers "/" with a hint instead of the dashboard.
package web

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var dist embed.FS

// Dist returns the built frontend with the "dist/" prefix stripped.
func Dist() fs.FS {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err) // only fails if the embed directive is wrong
	}
	return sub
}

// Built reports whether a real frontend build was embedded.
func Built() bool {
	_, err := fs.Stat(Dist(), "index.html")
	return err == nil
}
