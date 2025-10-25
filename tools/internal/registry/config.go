package registry

import (
	"os"
	"time"
)

const (
	DefaultCDNURL = "https://cdn.depscian.tech"
	CacheTTL      = 5 * time.Minute
	IndexPath     = "/index.json"
)

func GetCDNURL() string {
	if url := os.Getenv("DEPS_CDN_URL"); url != "" {
		return url
	}
	return DefaultCDNURL
}

