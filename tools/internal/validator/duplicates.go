package validator

import (
	"crypto/sha256"
	"fmt"
	"sort"

	"github.com/Deps-Tech/deps-registry/tools/internal/manifest"
)

type DuplicateSet struct {
	Packages  []string
	Signature string
}

func (d *DuplicateSet) Error() string {
	return fmt.Sprintf("duplicate packages detected: %v (same files)", d.Packages)
}

func DetectDuplicates(manifests map[string]*manifest.Manifest) []*DuplicateSet {
	signatureMap := make(map[string][]string)

	for id, m := range manifests {
		sig := generateFileSignature(m)
		signatureMap[sig] = append(signatureMap[sig], id)
	}

	duplicates := []*DuplicateSet{}
	for sig, packages := range signatureMap {
		if len(packages) > 1 {
			sort.Strings(packages)
			duplicates = append(duplicates, &DuplicateSet{
				Packages:  packages,
				Signature: sig,
			})
		}
	}

	return duplicates
}

func generateFileSignature(m *manifest.Manifest) string {
	fileNames := []string{}
	for fileName := range m.Files {
		fileNames = append(fileNames, fileName)
	}
	sort.Strings(fileNames)

	h := sha256.New()
	for _, name := range fileNames {
		h.Write([]byte(name))
		h.Write([]byte{0})
		h.Write([]byte(m.Files[name].SHA256))
		h.Write([]byte{0})
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

