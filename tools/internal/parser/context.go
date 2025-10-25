package parser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Deps-Tech/deps-registry/tools/internal/manifest"
)

type Context struct {
	PackageID       string
	PackagePath     string
	InternalModules map[string]bool
	Registry        *Registry
}

type Registry struct {
	packages map[string]*PackageInfo
	provides map[string]string
}

type PackageInfo struct {
	ID       string
	Version  string
	Provides []string
	Files    []string
}

func NewRegistry() *Registry {
	return &Registry{
		packages: make(map[string]*PackageInfo),
		provides: make(map[string]string),
	}
}

func (r *Registry) AddPackage(info *PackageInfo) {
	r.packages[info.ID] = info
	for _, alias := range info.Provides {
		r.provides[alias] = info.ID
	}
}

func (r *Registry) ResolveModule(modulePath string) string {
	parts := strings.Split(modulePath, ".")

	for i := len(parts); i > 0; i-- {
		testPath := strings.Join(parts[:i], ".")

		if _, ok := r.packages[testPath]; ok {
			return testPath
		}

		if pkgID, ok := r.provides[testPath]; ok {
			return pkgID
		}
	}

	return ""
}

func (r *Registry) GetPackage(id string) *PackageInfo {
	return r.packages[id]
}

func NewContext(packageID, packagePath string, registry *Registry) (*Context, error) {
	internalModules := make(map[string]bool)

	err := filepath.Walk(packagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".lua" {
			return nil
		}

		relPath, err := filepath.Rel(packagePath, path)
		if err != nil {
			return err
		}

		modulePath := strings.TrimSuffix(relPath, ".lua")
		modulePath = strings.ReplaceAll(modulePath, string(filepath.Separator), ".")

		internalModules[packageID+"."+modulePath] = true
		internalModules[modulePath] = true

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Context{
		PackageID:       packageID,
		PackagePath:     packagePath,
		InternalModules: internalModules,
		Registry:        registry,
	}, nil
}

func (c *Context) IsInternalModule(modulePath string) bool {
	if c.InternalModules[modulePath] {
		return true
	}

	if c.InternalModules[c.PackageID+"."+modulePath] {
		return true
	}

	parts := strings.Split(modulePath, ".")
	if len(parts) > 0 && parts[0] == c.PackageID {
		return true
	}

	return false
}

func LoadRegistryFromManifests(basePaths []string) (*Registry, error) {
	registry := NewRegistry()

	for _, basePath := range basePaths {
		items, err := os.ReadDir(basePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}

		for _, item := range items {
			if !item.IsDir() {
				continue
			}

			itemPath := filepath.Join(basePath, item.Name())
			versions, err := os.ReadDir(itemPath)
			if err != nil {
				continue
			}

			for _, version := range versions {
				if !version.IsDir() {
					continue
				}

				versionPath := filepath.Join(itemPath, version.Name())
				m, err := manifest.Load(versionPath)
				if err != nil {
					continue
				}

				files := make([]string, 0, len(m.Files))
				for file := range m.Files {
					files = append(files, file)
				}

				registry.AddPackage(&PackageInfo{
					ID:       m.ID,
					Version:  m.Version,
					Provides: m.Provides,
					Files:    files,
				})

				break
			}
		}
	}

	return registry, nil
}
