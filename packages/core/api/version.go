package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dcdavidev/bastion/packages/core/version"
)

type VersionCheckResponse struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	NeedsUpdate    bool   `json:"needs_update"`
}

var (
	versionCache      *VersionCheckResponse
	versionCacheLock  sync.Mutex
	lastVersionCheck  time.Time
	versionCacheTTL   = 6 * time.Hour
)

func (h *Handler) VersionCheckHandler(w http.ResponseWriter, r *http.Request) {
	versionCacheLock.Lock()
	defer versionCacheLock.Unlock()

	if versionCache != nil && time.Since(lastVersionCheck) < versionCacheTTL {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(versionCache)
		return
	}

	// Fetch from GitHub
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/dcdavidev/bastion/releases/latest")
	if err != nil {
		// Fallback to current version if GitHub is unreachable
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(VersionCheckResponse{
			CurrentVersion: version.Version,
			LatestVersion:  version.Version,
			NeedsUpdate:    false,
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(VersionCheckResponse{
			CurrentVersion: version.Version,
			LatestVersion:  version.Version,
			NeedsUpdate:    false,
		})
		return
	}

	var githubRelease struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&githubRelease); err != nil {
		http.Error(w, "Failed to parse GitHub response", http.StatusInternalServerError)
		return
	}

	latest := strings.TrimPrefix(githubRelease.TagName, "v")
	current := strings.TrimPrefix(version.Version, "v")

	needsUpdate := isNewerVersion(latest, current)

	versionCache = &VersionCheckResponse{
		CurrentVersion: version.Version,
		LatestVersion:  latest,
		NeedsUpdate:    needsUpdate,
	}
	lastVersionCheck = time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versionCache)
}

func isNewerVersion(latest, current string) bool {
	if latest == current {
		return false
	}

	lParts := strings.Split(latest, ".")
	cParts := strings.Split(current, ".")

	for i := 0; i < len(lParts) && i < len(cParts); i++ {
		if lParts[i] > cParts[i] {
			return true
		}
		if lParts[i] < cParts[i] {
			return false
		}
	}

	return len(lParts) > len(cParts)
}
