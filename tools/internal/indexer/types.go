package indexer

import "github.com/Deps-Tech/deps-registry/tools/internal/manifest"

type VersionInfo struct {
	URL      string            `json:"url"`
	SHA256   string            `json:"sha256"`
	Size     int64             `json:"size"`
	Manifest manifest.Manifest `json:"manifest"`
}

type PackageInfo struct {
	Latest   string                 `json:"latest"`
	Versions map[string]VersionInfo `json:"versions"`
}

type Index struct {
	Version      string                 `json:"version"`
	LastUpdated  string                 `json:"lastUpdated"`
	Dependencies map[string]PackageInfo `json:"dependencies"`
	Scripts      map[string]PackageInfo `json:"scripts"`
}
