package cmd

// DepManifest defines the structure of the dependency manifest file (v1 - legacy).
type DepManifest struct {
	ID               string            `json:"id"`
	Name             string            `json:"name,omitempty"`
	Version          string            `json:"version"`
	SourceURL        string            `json:"sourceUrl"`
	Files            []string          `json:"files"`
	Dependencies     map[string]string `json:"dependencies,omitempty"`
	Tags             []string          `json:"tags,omitempty"`
	CreatedFiles     []string          `json:"createdFiles,omitempty"`
	TouchedFiles     []string          `json:"touchedFiles,omitempty"`
	HasNetworkAccess bool              `json:"hasNetworkAccess,omitempty"`
}

// FileInfo contains metadata about a single file in the manifest.
type FileInfo struct {
	SHA256 string `json:"sha256"`
	Size   int64  `json:"size"`
}

// SecurityInfo contains security-related metadata.
type SecurityInfo struct {
	NetworkAccess bool     `json:"networkAccess,omitempty"`
	FileAccess    []string `json:"fileAccess,omitempty"`
	UsesFFI       bool     `json:"usesFFI,omitempty"`
}

// MetadataInfo contains additional metadata about the package.
type MetadataInfo struct {
	SourceURL  string   `json:"sourceUrl,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	Deprecated bool     `json:"deprecated,omitempty"`
}

// ManifestV2 defines the structure of the dependency manifest file (v2).
type ManifestV2 struct {
	ManifestVersion string              `json:"manifestVersion"`
	ID              string              `json:"id"`
	Name            string              `json:"name,omitempty"`
	Version         string              `json:"version"`
	Files           map[string]FileInfo `json:"files"`
	Dependencies    map[string]string   `json:"dependencies,omitempty"`
	Security        SecurityInfo        `json:"security,omitempty"`
	Metadata        MetadataInfo        `json:"metadata,omitempty"`
}
