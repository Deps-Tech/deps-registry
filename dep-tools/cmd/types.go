package cmd

// DepManifest defines the structure of the dependency manifest file.
type DepManifest struct {
	ID        string   `json:"id"`
	Version   string   `json:"version"`
	SourceURL string   `json:"sourceUrl"`
	Files     []string `json:"files"`
}
