package registry

import "time"

type Index struct {
	Version      string                  `json:"version"`
	LastUpdated  time.Time               `json:"lastUpdated"`
	Dependencies map[string]*Package     `json:"dependencies"`
	Scripts      map[string]*Package     `json:"scripts"`
}

type Package struct {
	Latest   string              `json:"latest"`
	Versions map[string]*Version `json:"versions"`
}

type Version struct {
	URL      string    `json:"url"`
	SHA256   string    `json:"sha256"`
	Size     int64     `json:"size"`
	Manifest *Manifest `json:"manifest"`
}

type Manifest struct {
	ManifestVersion string            `json:"manifestVersion"`
	ID              string            `json:"id"`
	Version         string            `json:"version"`
	Files           map[string]File   `json:"files"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	Security        Security          `json:"security,omitempty"`
	Metadata        Metadata          `json:"metadata,omitempty"`
}

type File struct {
	SHA256 string `json:"sha256"`
	Size   int64  `json:"size"`
}

type Security struct {
	NetworkAccess bool     `json:"networkAccess,omitempty"`
	FileAccess    []string `json:"fileAccess,omitempty"`
	UsesFFI       bool     `json:"usesFFI,omitempty"`
}

type Metadata struct {
	Tags      []string `json:"tags,omitempty"`
	SourceURL string   `json:"sourceUrl,omitempty"`
}

type DuplicateInfo struct {
	Exists          bool
	ExactMatch      bool
	ExistingVersion string
	AllVersions     []string
	PackageURL      string
}

