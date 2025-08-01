package cmd

// DepManifest defines the structure of the dependency manifest file.
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
