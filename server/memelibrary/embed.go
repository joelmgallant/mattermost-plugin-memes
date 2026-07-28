package memelibrary

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
)

// assetsFS carries the meme templates, fonts and per-template metadata into
// the binary. It replaces the go-bindata blob upstream generates, which cost
// 7.9 MB of generated Go to encode 2.1 MB of files.
//
//go:embed assets
var assetsFS embed.FS

// MustAsset returns the contents of the named asset, panicking when it is
// absent. Names are slash-separated and relative to the assets directory, for
// example "images/doge.jpg". The panic preserves go-bindata's contract, which
// memelibrary's package initialisation relies on.
func MustAsset(name string) []byte {
	data, err := assetsFS.ReadFile(path.Join("assets", name))
	if err != nil {
		panic(fmt.Sprintf("memelibrary: missing asset %q: %v", name, err))
	}
	return data
}

// AssetDir returns the base names of the entries in one asset directory, for
// example AssetDir("fonts").
func AssetDir(name string) ([]string, error) {
	entries, err := fs.ReadDir(assetsFS, path.Join("assets", name))
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names, nil
}
