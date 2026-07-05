package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
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
	main := normalized
	prerelease := false
	if idx := strings.IndexAny(normalized, "-+"); idx >= 0 {
		prerelease = strings.Contains(normalized[idx:], "-")
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
