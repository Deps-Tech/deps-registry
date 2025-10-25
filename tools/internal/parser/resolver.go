package parser

import (
	"strings"
)

type ResolvedDependency struct {
	PackageID    string
	OriginalPath string
	Resolved     bool
}

func ResolveDependencies(ctx *Context, rawModules []string) map[string]*ResolvedDependency {
	resolved := make(map[string]*ResolvedDependency)

	for _, module := range rawModules {
		if ctx.IsInternalModule(module) {
			continue
		}

		parts := strings.Split(module, ".")
		rootModule := parts[0]

		if strings.EqualFold(rootModule, ctx.PackageID) {
			continue
		}

		if isBuiltin(rootModule) {
			continue
		}

		pkgID := ctx.Registry.ResolveModule(module)
		if pkgID == "" {
			pkgID = rootModule
		}

		if pkgID != "" && !strings.EqualFold(pkgID, ctx.PackageID) {
			resolved[pkgID] = &ResolvedDependency{
				PackageID:    pkgID,
				OriginalPath: module,
				Resolved:     ctx.Registry.GetPackage(pkgID) != nil,
			}
		}
	}

	return resolved
}

var builtinModules = map[string]bool{
	"bit":       true,
	"bit32":     true,
	"math":      true,
	"string":    true,
	"table":     true,
	"os":        true,
	"io":        true,
	"debug":     true,
	"coroutine": true,
	"package":   true,
	"utf8":      true,
}

func isBuiltin(module string) bool {
	return builtinModules[strings.ToLower(module)]
}
