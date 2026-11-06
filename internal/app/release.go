package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

const ReleaseRepository = "marcioaltoe/roundfix"

var ErrNoReleases = errors.New("no releases published")

var runGitHubAPI = defaultRunGitHubAPI

type ReleaseAsset struct {
	Name        string
	APIURL      string
	DownloadURL string
	Size        int64
}

func LatestRelease(ctx context.Context) (string, []ReleaseAsset, error) {
	endpoint := "repos/" + ReleaseRepository + "/releases/latest"
	output, err := runGitHubAPI(ctx, endpoint)
	if err != nil {
		if isGitHubNotFound(err) {
			return "", nil, ErrNoReleases
		}
		return "", nil, fmt.Errorf("resolve latest release: %w", err)
	}

	var payload struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name        string `json:"name"`
			APIURL      string `json:"url"`
			DownloadURL string `json:"browser_download_url"`
			Size        int64  `json:"size"`
		} `json:"assets"`
	}
	if err := json.Unmarshal(output, &payload); err != nil {
		return "", nil, fmt.Errorf("parse latest release: %w", err)
	}
	if strings.TrimSpace(payload.TagName) == "" {
		return "", nil, fmt.Errorf("parse latest release: missing tag_name")
	}

	assets := make([]ReleaseAsset, 0, len(payload.Assets))
	for _, asset := range payload.Assets {
		assets = append(assets, ReleaseAsset{
			Name:        asset.Name,
			APIURL:      asset.APIURL,
			DownloadURL: asset.DownloadURL,
			Size:        asset.Size,
		})
	}
	return payload.TagName, assets, nil
}

// Release holds the published headings of one stable release.
type Release struct {
	Version string   `json:"version"`
	Titles  []string `json:"titles"`
	URL     string   `json:"url"`
}

const ReleaseListingPageSize = 30

// ListReleases reads one page of stable releases, newest first.
func ListReleases(ctx context.Context) ([]Release, error) {
	output, err := runGitHubAPI(ctx, fmt.Sprintf("repos/%s/releases?per_page=%d", ReleaseRepository, ReleaseListingPageSize))
	if err != nil {
		if isGitHubNotFound(err) {
			return nil, ErrNoReleases
		}
		return nil, fmt.Errorf("list releases: %w", err)
	}
	var payload []struct {
		TagName    string `json:"tag_name"`
		Body       string `json:"body"`
		Draft      bool   `json:"draft"`
		Prerelease bool   `json:"prerelease"`
		URL        string `json:"html_url"`
	}
	if err := json.Unmarshal(output, &payload); err != nil {
		return nil, fmt.Errorf("list releases: %w", err)
	}
	releases := make([]Release, 0, len(payload))
	for _, item := range payload {
		if item.Draft || item.Prerelease {
			continue
		}
		titles := []string{}
		for _, line := range strings.Split(item.Body, "\n") {
			if title, ok := strings.CutPrefix(line, "### "); ok {
				titles = append(titles, strings.TrimSpace(title))
			}
		}
		releases = append(releases, Release{Version: NormalizeVersion(item.TagName), Titles: titles, URL: item.URL})
	}
	if len(releases) == 0 {
		return nil, ErrNoReleases
	}
	sort.SliceStable(releases, func(i, j int) bool { return CompareVersions(releases[i].Version, releases[j].Version) > 0 })
	return releases, nil
}

func CompareVersions(left, right string) int {
	l := parseComparableVersion(left)
	r := parseComparableVersion(right)
	for i := 0; i < len(l.parts) || i < len(r.parts); i++ {
		var lp, rp int
		if i < len(l.parts) {
			lp = l.parts[i]
		}
		if i < len(r.parts) {
			rp = r.parts[i]
		}
		if lp > rp {
			return 1
		}
		if lp < rp {
			return -1
		}
	}
	if l.prerelease == r.prerelease {
		return 0
	}
	if l.prerelease {
		return -1
	}
	return 1
}

func NormalizeVersion(version string) string {
	trimmed := strings.TrimSpace(version)
	trimmed = strings.TrimPrefix(trimmed, "v")
	trimmed = strings.TrimPrefix(trimmed, "V")
	return trimmed
}

func IsDevelopmentVersion(version string) bool {
	normalized := strings.ToLower(NormalizeVersion(version))
	return normalized == "" || strings.Contains(normalized, "dev")
}

type comparableVersion struct {
	parts      []int
	prerelease bool
}

func parseComparableVersion(version string) comparableVersion {
	normalized := NormalizeVersion(version)
	if idx := strings.IndexByte(normalized, '+'); idx >= 0 {
		normalized = normalized[:idx]
	}
	main := normalized
	prerelease := false
	if idx := strings.IndexByte(normalized, '-'); idx >= 0 {
		prerelease = true
		main = normalized[:idx]
	}
	fields := strings.FieldsFunc(main, func(r rune) bool {
		return r == '.' || !unicode.IsDigit(r)
	})
	parts := make([]int, 0, len(fields))
	for _, field := range fields {
		if field == "" {
			continue
		}
		value, err := strconv.Atoi(field)
		if err != nil {
			continue
		}
		parts = append(parts, value)
	}
	return comparableVersion{parts: parts, prerelease: prerelease}
}

func defaultRunGitHubAPI(ctx context.Context, endpoint string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "gh", "api", endpoint)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("gh api %s failed: %w: %s", endpoint, err, boundedTail(output))
	}
	return bytes.TrimSpace(output), nil
}

func isGitHubNotFound(err error) bool {
	message := err.Error()
	return strings.Contains(message, "HTTP 404") || strings.Contains(message, "Not Found")
}

func boundedTail(output []byte) string {
	text := strings.TrimSpace(string(output))
	if len(text) <= 1200 {
		return text
	}
	return text[len(text)-1200:]
}
