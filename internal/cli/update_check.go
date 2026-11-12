package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"roundfix/internal/app"
	roundskills "roundfix/skills"
)

const (
	updateCheckSchema      = "roundfix/update-check/v1"
	updateCheckCacheSchema = "roundfix/update-check-cache/v1"
	updateCheckCacheFile   = "update-check.json"
	updateCheckInterval    = 24 * time.Hour
	updateCheckTimeout     = 5 * time.Second
	// A lookup lock older than this was left by a session that died mid-lookup.
	updateCheckLockStale = 2 * updateCheckTimeout
	updateCheckLockPoll  = 50 * time.Millisecond
	updateCheckBusyError = "another session's release lookup did not finish"
)

var updateCheckDeps = defaultUpdateCheckDependencies()

type updateCheckDependencies struct {
	now            func() time.Time
	currentVersion func() string
	listReleases   func(context.Context) ([]app.Release, error)
}

func defaultUpdateCheckDependencies() updateCheckDependencies {
	return updateCheckDependencies{now: time.Now, currentVersion: func() string { return app.Version }, listReleases: app.ListReleases}
}

type updateCheckCache struct {
	SchemaVersion string        `json:"schemaVersion"`
	CheckedAt     time.Time     `json:"checkedAt"`
	Releases      []app.Release `json:"releases"`
	Error         string        `json:"error"`
}

type updateCheckLookup struct {
	Source    string `json:"source"`
	CheckedAt string `json:"checkedAt,omitempty"`
	Error     string `json:"error"`
}

type updateCheckBaseline struct {
	State        string `json:"state"`
	FileChanges  int    `json:"fileChanges"`
	HistoryUnits int    `json:"historyUnits"`
	SkillChanges int    `json:"skillChanges"`
	Message      string `json:"message"`
}

type updateCheckReport struct {
	SchemaVersion     string                           `json:"schemaVersion"`
	Mode              string                           `json:"mode"`
	Installed         string                           `json:"installed"`
	Latest            string                           `json:"latest"`
	State             string                           `json:"state"`
	Lookup            updateCheckLookup                `json:"lookup"`
	Releases          []app.Release                    `json:"releases"`
	ReleasesTruncated bool                             `json:"releasesTruncated"`
	Baseline          updateCheckBaseline              `json:"baseline"`
	UserSkills        []roundskills.StaleUserSkillCopy `json:"userSkills,omitempty"`
	Next              string                           `json:"next"`
}

func runUpdateCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("update"))
		return exitOK
	}
	if len(args) == 0 || args[0] != "check" {
		printPreflightFailure("update", validationError{message: `expected subcommand "check"`}, stderr)
		return exitPreflight
	}
	flags := flag.NewFlagSet("update check", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	refresh := flags.Bool("refresh", false, "Read the release listing now instead of the cached one")
	format := flags.String("format", "text", "Output format: text (default) or json")
	err := flags.Parse(args[1:])
	if err == nil && len(flags.Args()) > 0 {
		err = validationError{message: fmt.Sprintf("unexpected argument %q", flags.Args()[0])}
	}
	if err == nil && *format != "text" && *format != "json" {
		err = validationError{message: fmt.Sprintf("format must be text or json, got %q", *format)}
	}
	if err != nil {
		printPreflightFailure("update check", err, stderr)
		return exitPreflight
	}
	loaded, err := loadRecoveryCommandConfig(environment, io.Discard)
	if err != nil {
		printPreflightFailure("update check", err, stderr)
		return exitPreflight
	}
	if err := warnRepositoryVersion(loaded, stderr); err != nil {
		printPreflightFailure("update check", err, stderr)
		return exitPreflight
	}
	report := buildUpdateCheckReport(ctx, loaded.HomeDir, loaded.Config.Update.Mode, *refresh)
	report.Baseline = previewUpdateCheckBaseline(ctx, loaded.GitRoot, environment)
	report.UserSkills = staleUserSkillCopies(ctx, loaded.HomeDir)
	if report.State == "update_available" || report.Baseline.State == "plan_ready" {
		report.Next = "apply"
	}
	if *format == "json" {
		// Reports have the same exit contract even when their output cannot be written.
		_ = json.NewEncoder(stdout).Encode(report)
	} else {
		printUpdateCheckReport(stdout, report)
	}
	return exitOK
}

func buildUpdateCheckReport(ctx context.Context, homeDir, mode string, refresh bool) updateCheckReport {
	deps := commandDependenciesForContext(ctx).updateCheck
	report := updateCheckReport{SchemaVersion: updateCheckSchema, Mode: mode, Installed: app.NormalizeVersion(deps.currentVersion()), State: "development", Lookup: updateCheckLookup{Source: "skipped"}, Releases: []app.Release{}, Next: "none"}
	if app.IsDevelopmentVersion(report.Installed) {
		return report
	}
	now := deps.now().UTC()
	path := filepath.Join(homeDir, ".roundfix", updateCheckCacheFile)
	cache, ok := readUpdateCheckCache(path, now)
	report.Lookup.Source = "cache"
	busy := false
	if refresh || !ok {
		// Sessions that start together take turns, so only the first reads the
		// network and the others report the cache it wrote.
		var release func()
		release, busy = acquireUpdateCheckLock(ctx, path+".lock")
		if release != nil {
			defer release()
		}
		var fresh bool
		cache, fresh = readUpdateCheckCache(path, now)
		if !refresh {
			ok = fresh
		}
		if busy {
			// Another session still holds the lookup: report the cache it last
			// wrote, even an expired one, instead of reading the network again.
			ok = cache.SchemaVersion == updateCheckCacheSchema
			if !ok {
				report.Lookup.Source = "skipped"
				report.Lookup.Error = updateCheckBusyError
				report.State = "unknown"
				return report
			}
		}
	}
	if (refresh && !busy) || !ok {
		report.Lookup.Source = "network"
		cache = updateCheckCache{SchemaVersion: updateCheckCacheSchema, CheckedAt: now, Releases: []app.Release{}}
		lookupCtx, cancel := context.WithTimeout(ctx, updateCheckTimeout)
		releases, err := deps.listReleases(lookupCtx)
		cancel()
		if err == nil && len(releases) == 0 {
			err = app.ErrNoReleases
		}
		if err != nil {
			if errors.Is(err, app.ErrNoReleases) {
				cache.Error = "no releases published"
			} else {
				cache.Error = err.Error()
			}
		} else {
			cache.Releases = releases
		}
		writeUpdateCheckCache(path, cache)
	}
	report.Lookup.CheckedAt = cache.CheckedAt.Format(time.RFC3339)
	report.Lookup.Error = cache.Error
	report.State = "unknown"
	if report.Lookup.Error != "" {
		return report
	}
	if len(cache.Releases) == 0 {
		report.Lookup.Error = "no releases published"
		return report
	}
	report.Latest = cache.Releases[0].Version
	report.State = "current"
	if app.CompareVersions(report.Latest, report.Installed) <= 0 {
		return report
	}
	report.State = "update_available"
	for _, release := range cache.Releases {
		if app.CompareVersions(release.Version, report.Installed) > 0 {
			if release.Titles == nil {
				release.Titles = []string{}
			}
			report.Releases = append(report.Releases, release)
		}
	}
	report.ReleasesTruncated = len(cache.Releases) == app.ReleaseListingPageSize && len(report.Releases) == len(cache.Releases)
	return report
}

func readUpdateCheckCache(path string, now time.Time) (updateCheckCache, bool) {
	content, err := os.ReadFile(path)
	if err != nil {
		return updateCheckCache{}, false
	}
	var cache updateCheckCache
	if json.Unmarshal(content, &cache) != nil || cache.SchemaVersion != updateCheckCacheSchema {
		return updateCheckCache{}, false
	}
	age := now.Sub(cache.CheckedAt)
	return cache, age >= 0 && age < updateCheckInterval
}

func writeUpdateCheckCache(path string, cache updateCheckCache) {
	// Cache failures never prevent a report or change its streams.
	content, err := json.Marshal(cache)
	if err != nil {
		return
	}
	if os.MkdirAll(filepath.Dir(path), 0o755) != nil {
		return
	}
	// A temporary file renamed into place keeps a concurrent reader from seeing
	// a partly written cache.
	file, err := os.CreateTemp(filepath.Dir(path), updateCheckCacheFile+".*.tmp")
	if err != nil {
		return
	}
	temporary := file.Name()
	_, writeErr := file.Write(content)
	closeErr := file.Close()
	if writeErr != nil || closeErr != nil || os.Chmod(temporary, 0o644) != nil || os.Rename(temporary, path) != nil {
		_ = os.Remove(temporary)
	}
}

// acquireUpdateCheckLock creates the lookup lock beside the cache and returns
// the function that removes it. It waits up to updateCheckLockStale for another
// session's lock, which outlasts that session's bounded lookup, and replaces a
// stale one. busy reports that another session still held the lock, so the
// caller must not read the network. A lock that cannot be created at all
// returns neither, because a cache beside it could not be written either.
func acquireUpdateCheckLock(ctx context.Context, path string) (release func(), busy bool) {
	if os.MkdirAll(filepath.Dir(path), 0o755) != nil {
		return nil, false
	}
	deadline := time.NewTimer(updateCheckLockStale)
	defer deadline.Stop()
	for {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			_ = file.Close()
			return func() { _ = os.Remove(path) }, false
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, false
		}
		if info, statErr := os.Stat(path); statErr == nil {
			if age := time.Since(info.ModTime()); age > updateCheckLockStale || age < -updateCheckLockStale {
				_ = os.Remove(path)
				continue
			}
		}
		select {
		case <-ctx.Done():
			return nil, true
		case <-deadline.C:
			return nil, true
		case <-time.After(updateCheckLockPoll):
		}
	}
}

func previewUpdateCheckBaseline(ctx context.Context, gitRoot string, environment commandEnvironment) updateCheckBaseline {
	if gitRoot == "" {
		return updateCheckBaseline{State: "not_checked"}
	}
	var stdout bytes.Buffer
	runBaselineUpdateCommand(ctx, []string{"--repo", gitRoot, "--format", "json"}, &stdout, io.Discard, environment)
	var result baselineUpdateResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return updateCheckBaseline{State: "unknown", Message: updateCheckLine(fmt.Sprintf("decode Baseline preview: %v", err), 0)}
	}
	return updateCheckBaseline{State: result.State, FileChanges: len(result.FileChanges), HistoryUnits: len(result.History.Units), SkillChanges: len(result.Skills.Outdated) + len(result.Skills.Missing) + len(result.Skills.Retired), Message: updateCheckLine(result.Message, 0)}
}

func updateCheckLine(value string, limit int) string {
	line := strings.Join(strings.Fields(value), " ")
	if limit > 0 {
		runes := []rune(line)
		if len(runes) > limit {
			line = string(runes[:limit])
		}
	}
	return line
}

func printUpdateCheckReport(stdout io.Writer, report updateCheckReport) {
	switch report.State {
	case "update_available":
		fmt.Fprintf(stdout, "Roundfix %s → %s is available (update.mode %s).\n", report.Installed, report.Latest, report.Mode)
		fmt.Fprintf(stdout, "Changes since %s:\n", report.Installed)
		for _, release := range report.Releases {
			titles := strings.Join(release.Titles, "; ")
			if len(release.Titles) == 0 {
				titles = release.URL
			}
			fmt.Fprintf(stdout, "- %s: %s\n", release.Version, titles)
		}
		if report.ReleasesTruncated {
			fmt.Fprintln(stdout, "- earlier releases: https://github.com/"+app.ReleaseRepository+"/releases")
		}
	case "current":
		fmt.Fprintf(stdout, "Roundfix %s is current (latest release %s; update.mode %s).\n", report.Installed, report.Latest, report.Mode)
	case "unknown":
		fmt.Fprintf(stdout, "Roundfix %s: latest release not reached (%s); the session continues.\n", report.Installed, updateCheckLine(report.Lookup.Error, 300))
	case "development":
		fmt.Fprintf(stdout, "Roundfix %s is a development build; no release comparison (update.mode %s).\n", report.Installed, report.Mode)
	}
	switch report.Baseline.State {
	case "current":
		fmt.Fprintln(stdout, "Baseline: current with the installed Roundfix.")
	case "plan_ready":
		fmt.Fprintf(stdout, "Baseline: plan ready with the installed Roundfix: %d file change(s), %d history unit(s), %d skill change(s).\n", report.Baseline.FileChanges, report.Baseline.HistoryUnits, report.Baseline.SkillChanges)
	case "not_checked":
		fmt.Fprintln(stdout, "Baseline: not checked (outside a Git repository).")
	default:
		if report.Baseline.Message != "" {
			fmt.Fprintf(stdout, "Baseline: %s: %s\n", report.Baseline.State, report.Baseline.Message)
		} else {
			fmt.Fprintf(stdout, "Baseline: %s.\n", report.Baseline.State)
		}
	}
	printUpdateCheckUserSkills(stdout, report.UserSkills)
	if report.Next == "apply" {
		if report.Mode == "auto" {
			fmt.Fprintln(stdout, "Next: update.mode auto: apply it now with the Roundfix skill's --update mode.")
		} else {
			fmt.Fprintln(stdout, "Next: apply it with the Roundfix skill's --update mode; update.mode notify applies nothing on its own.")
		}
	}
}
