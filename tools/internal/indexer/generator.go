package indexer

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/Deps-Tech/deps-registry/tools/internal/filesystem"
	"github.com/Deps-Tech/deps-registry/tools/internal/manifest"
)

func Generate(distPath, cdnURL string) (*Index, error) {
	deps, err := generateForType("deps", distPath, cdnURL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate deps index: %w", err)
	}

	scripts, err := generateForType("scripts", distPath, cdnURL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate scripts index: %w", err)
	}

	return &Index{
		Version:      "1.0",
		LastUpdated:  time.Now().UTC().Format(time.RFC3339),
		Dependencies: deps,
		Scripts:      scripts,
	}, nil
}

func generateForType(itemType, distPath, cdnURL string) (map[string]PackageInfo, error) {
	result := make(map[string]PackageInfo)
	itemsPath := filepath.Join(distPath, itemType)

	if _, err := os.Stat(itemsPath); os.IsNotExist(err) {
		return result, nil
	}

	files, err := os.ReadDir(itemsPath)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`^(.+)-([\d\.]+)\.zip$`)
	packageVersions := make(map[string][]string)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		matches := re.FindStringSubmatch(file.Name())
		if len(matches) == 3 {
			pkgName := matches[1]
			version := matches[2]
			packageVersions[pkgName] = append(packageVersions[pkgName], version)
		}
	}

	for pkgName, versions := range packageVersions {
		sort.Strings(versions)
		latest := versions[len(versions)-1]

		pkgInfo := PackageInfo{
			Latest:   latest,
			Versions: make(map[string]VersionInfo),
		}

		for _, version := range versions {
			fileName := fmt.Sprintf("%s-%s.zip", pkgName, version)
			filePath := filepath.Join(itemsPath, fileName)

			hash, err := filesystem.SHA256File(filePath)
			if err != nil {
				continue
			}

			info, err := os.Stat(filePath)
			if err != nil {
				continue
			}

			m, err := readManifestFromZip(filePath)
			if err != nil {
				continue
			}

			url := fmt.Sprintf("%s/%s/%s", cdnURL, itemType, fileName)
			pkgInfo.Versions[version] = VersionInfo{
				URL:      url,
				SHA256:   hash,
				Size:     info.Size(),
				Manifest: *m,
			}
		}

		result[pkgName] = pkgInfo
	}

	return result, nil
}

func readManifestFromZip(zipPath string) (*manifest.Manifest, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "dep.json" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				return nil, err
			}

			var m manifest.Manifest
			if err := json.Unmarshal(content, &m); err != nil {
				return nil, err
			}
			return &m, nil
		}
	}

	return nil, fmt.Errorf("dep.json not found in %s", zipPath)
}
