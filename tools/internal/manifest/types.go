package manifest

type FileInfo struct {
	SHA256 string `json:"sha256"`
	Size   int64  `json:"size"`
}

type Security struct {
	NetworkAccess bool     `json:"networkAccess,omitempty"`
	FileAccess    []string `json:"fileAccess,omitempty"`
	UsesFFI       bool     `json:"usesFFI,omitempty"`
}

type Metadata struct {
	SourceURL  string   `json:"sourceUrl,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	Deprecated bool     `json:"deprecated,omitempty"`
}

type Manifest struct {
	ManifestVersion string              `json:"manifestVersion"`
	ID              string              `json:"id"`
	Name            string              `json:"name,omitempty"`
	Version         string              `json:"version"`
	Provides        []string            `json:"provides,omitempty"`
	Files           map[string]FileInfo `json:"files"`
	Dependencies    map[string]string   `json:"dependencies,omitempty"`
	Security        Security            `json:"security,omitempty"`
	Metadata        Metadata            `json:"metadata,omitempty"`
}
