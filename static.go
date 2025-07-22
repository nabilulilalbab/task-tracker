package welcomesite

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:static
var embeddedFiles embed.FS

// StaticFS berisi file system untuk aset statis.
// Kita membuat sub-filesystem yang dimulai dari direktori 'static'.
var StaticFS http.FileSystem

func init() {
	fsys, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		panic(err)
	}
	StaticFS = http.FS(fsys)
}
