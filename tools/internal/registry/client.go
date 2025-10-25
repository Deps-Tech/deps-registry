package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Client struct {
	cdnURL    string
	index     *Index
	cacheTime time.Time
	cacheTTL  time.Duration
	mu        sync.RWMutex
	httpClient *http.Client
}

func NewClient(cdnURL string) *Client {
	if cdnURL == "" {
		cdnURL = GetCDNURL()
	}
	
	return &Client{
		cdnURL:   cdnURL,
		cacheTTL: CacheTTL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) fetchIndex() error {
	url := c.cdnURL + IndexPath
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch index: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	
	var index Index
	if err := json.Unmarshal(body, &index); err != nil {
		return fmt.Errorf("failed to parse index: %w", err)
	}
	
	c.mu.Lock()
	c.index = &index
	c.cacheTime = time.Now()
	c.mu.Unlock()
	
	return nil
}

func (c *Client) getIndex() (*Index, error) {
	c.mu.RLock()
	needsRefresh := c.index == nil || time.Since(c.cacheTime) > c.cacheTTL
	c.mu.RUnlock()
	
	if needsRefresh {
		if err := c.fetchIndex(); err != nil {
			c.mu.RLock()
			hasCache := c.index != nil
			c.mu.RUnlock()
			
			if hasCache {
				return c.index, nil
			}
			return nil, err
		}
	}
	
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.index, nil
}

func (c *Client) GetLatestVersion(itemType, id string) (string, error) {
	index, err := c.getIndex()
	if err != nil {
		return "", err
	}
	
	var pkg *Package
	switch itemType {
	case "deps", "dependencies":
		pkg = index.Dependencies[id]
	case "scripts":
		pkg = index.Scripts[id]
	default:
		return "", fmt.Errorf("unknown item type: %s", itemType)
	}
	
	if pkg == nil {
		return "", fmt.Errorf("package not found: %s", id)
	}
	
	return pkg.Latest, nil
}

func (c *Client) CheckDuplicate(itemType, id, version string) (*DuplicateInfo, error) {
	index, err := c.getIndex()
	if err != nil {
		return nil, err
	}
	
	var pkg *Package
	switch itemType {
	case "deps", "dependencies":
		pkg = index.Dependencies[id]
	case "scripts":
		pkg = index.Scripts[id]
	default:
		return nil, fmt.Errorf("unknown item type: %s", itemType)
	}
	
	info := &DuplicateInfo{
		Exists: pkg != nil,
	}
	
	if !info.Exists {
		return info, nil
	}
	
	versions := make([]string, 0, len(pkg.Versions))
	for v := range pkg.Versions {
		versions = append(versions, v)
	}
	sort.Strings(versions)
	
	info.AllVersions = versions
	info.ExistingVersion = pkg.Latest
	info.ExactMatch = pkg.Versions[version] != nil
	
	if info.ExactMatch && pkg.Versions[version] != nil {
		info.PackageURL = pkg.Versions[version].URL
	}
	
	return info, nil
}

func (c *Client) GetAllDependencies() ([]string, error) {
	index, err := c.getIndex()
	if err != nil {
		return nil, err
	}
	
	deps := make([]string, 0, len(index.Dependencies))
	for id := range index.Dependencies {
		deps = append(deps, id)
	}
	sort.Strings(deps)
	
	return deps, nil
}

func (c *Client) GetAllScripts() ([]string, error) {
	index, err := c.getIndex()
	if err != nil {
		return nil, err
	}
	
	scripts := make([]string, 0, len(index.Scripts))
	for id := range index.Scripts {
		scripts = append(scripts, id)
	}
	sort.Strings(scripts)
	
	return scripts, nil
}

func (c *Client) IsAvailable() bool {
	_, err := c.getIndex()
	return err == nil
}

