package manifest

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Load(path string) (*Manifest, error) {
	data, err := os.ReadFile(filepath.Join(path, "dep.json"))
	if err != nil {
		return nil, err
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func Save(path string, m *Manifest) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(path, "dep.json"), data, 0644)
}
