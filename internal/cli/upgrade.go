package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
)

const (
	versionFreshnessCacheFile = "version-check.json"
	versionFreshnessInterval  = 24 * time.Hour
	versionFreshnessTimeout   = 2 * time.Second
)

var upgradeDeps = defaultUpgradeDependencies()
var versionFreshnessDeps = defaultVersionFreshnessDependencies()

type upgradeDependencies struct {
	latestRelease  func(context.Context) (string, []app.ReleaseAsset, error)
	downloadAsset  func(context.Context, app.ReleaseAsset) ([]byte, error)
	executablePath func() (string, error)
	currentVersion func() string
	goos           string
	goarch         string
}

type versionFreshnessDependencies struct {
	now            func() time.Time
	currentVersion func() string
	latestRelease  func(context.Context) (string, []app.ReleaseAsset, error)
}

type upgradeRequest struct {
	check bool
}

type versionFreshnessCache struct {
	CheckedAt     time.Time `json:"checked_at"`
	LatestVersion string    `json:"latest_version,omitempty"`
}

func defaultUpgradeDependencies() upgradeDependencies {
	return upgradeDependencies{
		latestRelease:  app.LatestRelease,
		downloadAsset:  defaultDownloadReleaseAsset,
		executablePath: os.Executable,
		currentVersion: func() string { return app.Version },
		goos:           runtime.GOOS,
		goarch:         runtime.GOARCH,
	}
}

func defaultVersionFreshnessDependencies() versionFreshnessDependencies {
	return versionFreshnessDependencies{
		now:            time.Now,
		currentVersion: func() string { return app.Version },
		latestRelease:  app.LatestRelease,
	}
}

func runUpgradeCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("upgrade"))
		return exitOK
	}
	req, err := parseUpgradeCommand(args)
	if err != nil {
		printUpgradeFailure(err, stderr)
		return exitPreflight
	}
	outcome, err := performUpgrade(ctx, req, commandDependenciesForContext(ctx).upgrade)
	if err != nil {
		fmt.Fprintf(stderr, "%s: upgrade failed: %v\n", app.Name, err)
		return exitRunFailed
	}
	fmt.Fprintln(stdout, outcome)
	return exitOK
}

func parseUpgradeCommand(args []string) (upgradeRequest, error) {
	req := upgradeRequest{}
	fs := flag.NewFlagSet("upgrade", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.BoolVar(&req.check, "check", false, "Report the latest release without installing it")
	if err := fs.Parse(args); err != nil {
		return req, validationError{message: err.Error()}
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		return req, validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	return req, nil
}

func performUpgrade(ctx context.Context, req upgradeRequest, deps upgradeDependencies) (string, error) {
	tag, assets, err := deps.latestRelease(ctx)
	if err != nil {
		if errors.Is(err, app.ErrNoReleases) {
			return "no releases published", nil
		}
		return "", err
	}

	current := app.NormalizeVersion(deps.currentVersion())
	latest := app.NormalizeVersion(tag)
	if app.CompareVersions(latest, current) <= 0 {
		return fmt.Sprintf("already current %s", current), nil
	}
	if req.check {
		return fmt.Sprintf("upgrade available %s → %s", current, latest), nil
	}

	asset, ok := selectPlatformAsset(assets, deps.goos, deps.goarch)
	if !ok {
		return "", upgradeManualError{
			err: fmt.Errorf("release %s has no asset for %s/%s", latest, deps.goos, deps.goarch),
		}
	}
	executablePath, err := deps.executablePath()
	if err != nil {
		return "", upgradeManualError{err: fmt.Errorf("resolve current executable: %w", err), asset: asset}
	}
	if err := installReleaseAsset(ctx, deps, asset, assets, executablePath); err != nil {
		return "", upgradeManualError{err: err, asset: asset, executablePath: executablePath}
	}
	return fmt.Sprintf("upgraded %s → %s", current, latest), nil
}

type upgradeManualError struct {
	err            error
	asset          app.ReleaseAsset
	executablePath string
}

func (err upgradeManualError) Error() string {
	fallback := "manual fallback: open https://github.com/" + app.ReleaseRepository + "/releases"
	if err.asset.DownloadURL != "" {
		fallback = "manual fallback: download " + err.asset.DownloadURL
	}
	if err.executablePath != "" {
		fallback += " and replace " + err.executablePath
	}
	return fmt.Sprintf("%v; %s", err.err, fallback)
}

func (err upgradeManualError) Unwrap() error {
	return err.err
}

func installReleaseAsset(ctx context.Context, deps upgradeDependencies, asset app.ReleaseAsset, assets []app.ReleaseAsset, executablePath string) error {
	content, err := deps.downloadAsset(ctx, asset)
	if err != nil {
		return fmt.Errorf("download %s: %w", releaseAssetName(asset), err)
	}
	if asset.Size > 0 && int64(len(content)) != asset.Size {
		return fmt.Errorf("verify %s: size %d does not match expected %d", releaseAssetName(asset), len(content), asset.Size)
	}
	if checksumAsset, ok := selectChecksumAsset(assets); ok {
		checksumContent, err := deps.downloadAsset(ctx, checksumAsset)
		if err != nil {
			return fmt.Errorf("download checksum %s: %w", releaseAssetName(checksumAsset), err)
		}
		expected, err := checksumForAsset(checksumContent, asset.Name)
		if err != nil {
			return fmt.Errorf("verify %s: %w", releaseAssetName(asset), err)
		}
		actual := sha256.Sum256(content)
		if !bytes.Equal(actual[:], expected) {
			return fmt.Errorf("verify %s: checksum mismatch", releaseAssetName(asset))
		}
	}

	return replaceExecutableAtomically(executablePath, content)
}

func replaceExecutableAtomically(executablePath string, content []byte) error {
	info, err := os.Stat(executablePath)
	if err != nil {
		return fmt.Errorf("inspect %s: %w", executablePath, err)
	}
	dir := filepath.Dir(executablePath)
	base := filepath.Base(executablePath)
	temp, err := os.CreateTemp(dir, "."+base+".upgrade-*")
	if err != nil {
		return fmt.Errorf("create sibling temp file: %w", err)
	}
	tempPath := temp.Name()
	renamed := false
	defer func() {
		if !renamed {
			_ = os.Remove(tempPath)
		}
	}()
	if _, err := temp.Write(content); err != nil {
		_ = temp.Close()
		return fmt.Errorf("write temp binary: %w", err)
	}
	if err := temp.Chmod(info.Mode().Perm()); err != nil {
		_ = temp.Close()
		return fmt.Errorf("set temp binary permissions: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close temp binary: %w", err)
	}
	if err := os.Rename(tempPath, executablePath); err != nil {
		return fmt.Errorf("replace %s: %w", executablePath, err)
	}
	renamed = true
	return nil
}

func selectPlatformAsset(assets []app.ReleaseAsset, goos, goarch string) (app.ReleaseAsset, bool) {
	candidates := append([]app.ReleaseAsset(nil), assets...)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Name < candidates[j].Name
	})
	for _, asset := range candidates {
		name := strings.ToLower(asset.Name)
		if isChecksumAssetName(name) {
			continue
		}
		if strings.Contains(name, strings.ToLower(goos)) && strings.Contains(name, strings.ToLower(goarch)) {
			return asset, true
		}
	}
	return app.ReleaseAsset{}, false
}

func selectChecksumAsset(assets []app.ReleaseAsset) (app.ReleaseAsset, bool) {
	candidates := append([]app.ReleaseAsset(nil), assets...)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Name < candidates[j].Name
	})
	for _, asset := range candidates {
		if isChecksumAssetName(strings.ToLower(asset.Name)) {
			return asset, true
		}
	}
	return app.ReleaseAsset{}, false
}

func isChecksumAssetName(name string) bool {
	return strings.Contains(name, "checksum") || strings.Contains(name, "sha256")
}

func checksumForAsset(content []byte, assetName string) ([]byte, error) {
	if fields := strings.Fields(string(content)); len(fields) == 1 {
		return decodeSHA256(fields[0])
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		hash := fields[0]
		name := strings.TrimPrefix(fields[len(fields)-1], "*")
		if name == assetName || filepath.Base(name) == assetName {
			return decodeSHA256(hash)
		}
	}
	return nil, fmt.Errorf("checksum for %s not found", assetName)
}

func decodeSHA256(value string) ([]byte, error) {
	if len(value) != sha256.Size*2 {
		return nil, fmt.Errorf("checksum has length %d, want %d", len(value), sha256.Size*2)
	}
	decoded, err := hex.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("decode checksum: %w", err)
	}
	return decoded, nil
}

func defaultDownloadReleaseAsset(ctx context.Context, asset app.ReleaseAsset) ([]byte, error) {
	if strings.TrimSpace(asset.APIURL) == "" {
		return nil, fmt.Errorf("asset %s has no GitHub API URL", releaseAssetName(asset))
	}
	endpoint := strings.TrimPrefix(asset.APIURL, "https://api.github.com/")
	cmd := exec.CommandContext(ctx, "gh", "api", "-H", "Accept: application/octet-stream", endpoint)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("gh api download failed: %w: %s", err, boundedCommandTail(output))
	}
	return output, nil
}

func releaseAssetName(asset app.ReleaseAsset) string {
	if asset.Name != "" {
		return asset.Name
	}
	if asset.DownloadURL != "" {
		return asset.DownloadURL
	}
	return asset.APIURL
}

func maybeReportVersionFreshness(ctx context.Context, loaded roundconfig.Loaded, stderr io.Writer) {
	deps := commandDependenciesForContext(ctx).versionFreshness
	if deps.now == nil {
		deps.now = time.Now
	}
	if deps.currentVersion == nil {
		deps.currentVersion = func() string { return app.Version }
	}
	if deps.latestRelease == nil {
		deps.latestRelease = app.LatestRelease
	}

	current := app.NormalizeVersion(deps.currentVersion())
	if app.IsDevelopmentVersion(current) || strings.TrimSpace(loaded.HomeDir) == "" {
		return
	}

	now := deps.now()
	cachePath := versionFreshnessCachePath(loaded.HomeDir)
	cache, fresh := readVersionFreshnessCache(cachePath, now)
	latest := cache.LatestVersion
	if !fresh {
		checkCtx, cancel := context.WithTimeout(ctx, versionFreshnessTimeout)
		tag, _, err := deps.latestRelease(checkCtx)
		cancel()
		cache = versionFreshnessCache{CheckedAt: now}
		if err == nil {
			cache.LatestVersion = app.NormalizeVersion(tag)
			latest = cache.LatestVersion
		} else {
			latest = ""
		}
		writeVersionFreshnessCache(cachePath, cache)
	}

	if latest != "" && app.CompareVersions(latest, current) > 0 {
		fmt.Fprintf(stderr, "%s %s is behind latest %s; run %s upgrade\n", app.Name, current, latest, app.Name)
	}
}

func versionFreshnessCachePath(homeDir string) string {
	return filepath.Join(homeDir, ".roundfix", versionFreshnessCacheFile)
}

func readVersionFreshnessCache(path string, now time.Time) (versionFreshnessCache, bool) {
	content, err := os.ReadFile(path)
	if err != nil {
		return versionFreshnessCache{}, false
	}
	var cache versionFreshnessCache
	if err := json.Unmarshal(content, &cache); err != nil {
		return versionFreshnessCache{}, false
	}
	if cache.CheckedAt.IsZero() {
		return versionFreshnessCache{}, false
	}
	age := now.Sub(cache.CheckedAt)
	if age < 0 || age >= versionFreshnessInterval {
		return cache, false
	}
	return cache, true
}

func writeVersionFreshnessCache(path string, cache versionFreshnessCache) {
	content, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return
	}
	content = append(content, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	_ = os.WriteFile(path, content, 0o644)
}

func printUpgradeFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: upgrade failed: %v\n", app.Name, err)
	fmt.Fprintf(stderr, "Run '%s upgrade --help' for usage.\n", app.Name)
}

func boundedCommandTail(output []byte) string {
	text := strings.TrimSpace(string(output))
	if len(text) <= 1200 {
		return text
	}
	return text[len(text)-1200:]
}
