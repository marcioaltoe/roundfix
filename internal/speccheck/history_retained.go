package speccheck

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"roundfix/internal/spec"
)

// CodeHistoryRetained identifies retired files and citations that retain them.
const CodeHistoryRetained = "SC-HISTORY-RETAINED"

var retirementRegionOpening = regexp.MustCompile(`(?m)^<!-- setup-context-driven:begin id=guide\.docs-layout(?: [^\n]*)? -->\r?$`)

var historyLocatorPattern = regexp.MustCompile("`([^`\\s]+)@([0-9a-fA-F]{64}|[0-9a-fA-F]{40})`")

func historyGit(ctx context.Context, root string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", root, "-c", "core.fsmonitor=false"}, args...)...)
	cmd.Env = mechanicalGitEnvironment()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), strings.TrimSpace(string(out)), err)
	}
	return out, nil
}

func historyHasGit(ctx context.Context, root string) (bool, error) {
	_, err := historyGit(ctx, root, "rev-parse", "--is-inside-work-tree")
	if err != nil && strings.Contains(err.Error(), "not a git repository") {
		return false, nil
	}
	return err == nil, err
}

func retirementRulePresent(root string) (bool, error) {
	content, err := os.ReadFile(filepath.Join(root, "docs/agents/setup-context.json"))
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("read retirement manifest: %w", err)
	}
	var manifest struct {
		Artifacts []struct{ ID, Path string } `json:"managedArtifacts"`
	}
	if err := json.Unmarshal(content, &manifest); err != nil {
		return false, fmt.Errorf("parse retirement manifest: %w", err)
	}
	for _, artifact := range manifest.Artifacts {
		if artifact.ID != "guide.docs-layout" {
			continue
		}
		if _, ok, err := receiptPathInfo(root, artifact.Path); err != nil {
			return false, err
		} else if !ok {
			return false, nil
		}
		content, err = os.ReadFile(filepath.Join(root, artifact.Path))
		if err != nil {
			return false, fmt.Errorf("read retirement guide: %w", err)
		}
		text := string(content)
		opening := retirementRegionOpening.FindStringIndex(text)
		if opening == nil {
			return false, nil
		}
		text = text[opening[1]:]
		end := strings.Index(text, "<!-- setup-context-driven:end id=guide.docs-layout -->")
		if end < 0 {
			return false, nil
		}
		return strings.Contains(strings.Join(strings.Fields(text[:end]), " "), spec.RetirementRuleSentinel), nil
	}
	return false, nil
}

func detectHistoryRetained(result *Result, root, specsRoot string) error {
	enabled, err := retirementRulePresent(root)
	if err != nil {
		return err
	}
	if !enabled {
		addSkip(result, CodeHistoryRetained, "recorded guide.docs-layout region does not carry "+spec.RetirementRuleSentinel)
		return nil
	}
	ctx := context.Background()
	git, err := historyHasGit(ctx, root)
	if err != nil {
		return err
	}
	if !git {
		addSkip(result, CodeHistoryRetained, "Git Locator resolution and validation: repository is outside Git")
	}
	specRel, err := filepath.Rel(root, specsRoot)
	if err != nil {
		return err
	}
	specRel = filepath.ToSlash(specRel)
	archiveRoot := spec.ArchiveSpecRoot(specsRoot, specRel == "docs/specs")
	archiveRel, err := filepath.Rel(root, archiveRoot)
	if err != nil {
		return err
	}
	archiveRel = filepath.ToSlash(archiveRel)
	pending, err := spec.PendingTaskCreates(specsRoot)
	if err != nil {
		return err
	}
	candidates, err := spec.RemovedHistoryPaths(ctx, root)
	if err != nil {
		return err
	}
	historyRoot := "docs/history"
	report := func(path string, line int, summary, fix string) {
		result.Findings = append(result.Findings, Finding{
			Code:     CodeHistoryRetained,
			Severity: SeverityError,
			Summary:  path + " " + summary,
			Where:    []Location{{Path: path, Line: line}},
			Fix:      fix,
		})
	}
	// A single walk sees both tracked files and newly planted files. Symlinks are
	// not followed, so the check cannot read outside its repository through them.
	roots := []string{"docs", specRel}
	seen := map[string]bool{}
	locators := map[string]spec.GitLocator{}
	locate := func(path, preferred, pr string) (spec.GitLocator, error) {
		key := path + "\x00" + preferred + "\x00" + pr
		if l, ok := locators[key]; ok {
			return l, nil
		}
		l, err := spec.ResolveGitLocator(ctx, root, path, preferred, pr)
		if err == nil {
			locators[key] = l
		}
		return l, err
	}
	missingLinks := map[spec.HistoryCitation]bool{}
	var citations []historyLocatorCitation
	for _, scanRoot := range roots {
		err := filepath.WalkDir(filepath.Join(root, scanRoot), func(full string, entry fs.DirEntry, walkErr error) error {
			if errors.Is(walkErr, os.ErrNotExist) {
				return nil
			}
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() {
				return nil
			}
			path, err := filepath.Rel(root, full)
			if err != nil {
				return err
			}
			path = filepath.ToSlash(path)
			if seen[path] {
				return nil
			}
			seen[path] = true
			for _, family := range []string{"findings", "backlog", "adr", "reviews", "handoffs"} {
				if strings.HasPrefix(path, historyRoot+"/"+family+"/") {
					report(path, 1, "retains a retired history file", "Delete "+path+"; Git keeps its text.")
					return nil
				}
			}
			if !entry.Type().IsRegular() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			content, err := os.ReadFile(full)
			if err != nil {
				return fmt.Errorf("read retained history %s: %w", path, err)
			}
			if filepath.ToSlash(filepath.Dir(path)) == "docs/findings" {
				fm, err := parseFindingFrontmatter(content)
				if err != nil {
					return err
				}
				if terminalHistoryFinding(fm.status.value) {
					report(path, fm.status.line+1, "retains a terminal Finding", "Delete "+path+"; Git keeps its text.")
				}
			}
			record := spec.ArchiveRecord{}
			isRecord := filepath.ToSlash(filepath.Dir(path)) == archiveRel && spec.IsArchiveRecord(entry.Name(), content)
			if isRecord {
				record, err = spec.ParseArchiveRecord(content)
				if err != nil {
					return err
				}
			}
			if !isRecord && !spec.IsLiveDocument(path, specRel, archiveRel) && !strings.HasPrefix(path, "docs/user-guide/") {
				return nil
			}
			for _, link := range spec.MissingRelativeLinks(root, path, content) {
				if pending[link.Target] || (link.Self && !isRecord) {
					continue
				}
				fix := "Restore " + link.Target + " or replace its link; Git Locator resolution is skipped outside Git."
				copyPath, promoted := spec.RecordLinkPromotion(root, link, record)
				if isRecord && promoted {
					relative, err := filepath.Rel(filepath.Dir(path), filepath.FromSlash(copyPath))
					if err != nil {
						return err
					}
					fix = "Replace the link with " + filepath.ToSlash(relative) + ", the copy the archive promoted to " + copyPath + "."
				} else if git {
					var l spec.GitLocator
					if isRecord {
						l, err = spec.ResolveRecordLink(ctx, root, link, record)
					} else {
						l, err = locate(link.Target, "", "")
					}
					if err != nil {
						return err
					}
					if !isRecord && l.Commit == "" {
						if candidate, ok := spec.ResolveMissingLinkByName(link.Target, candidates); ok {
							l, err = locate(candidate, "history-full", "")
							if err != nil {
								return err
							}
						}
					}
					fix = "Replace the link with " + l.Text(false) + "."
				}
				report(path, link.Line, "links to a missing file or itself: "+link.Target, fix)
				missingLinks[spec.HistoryCitation{Path: path, Line: link.Line, Target: link.Target}] = true
			}
			if !git {
				return nil
			}

			for _, match := range historyLocatorPattern.FindAllSubmatchIndex(content, -1) {
				citations = append(citations, historyLocatorCitation{
					path: path, line: bytes.Count(content[:match[0]], []byte("\n")) + 1,
					target: string(content[match[2]:match[3]]), commit: strings.ToLower(string(content[match[4]:match[5]])), pullRequest: record.PullRequest,
				})
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("check retained history: %w", err)
		}
	}
	// The new pass reads only tracked and untracked, unignored Citation Sources
	// and local Git history. Nothing is sent anywhere. The legacy walk above
	// retains its findings and wording; a missing-link finding owns its line
	// and target, so a removed citation never reports that pair twice.
	if git && len(candidates) > 0 {
		paths, err := spec.CitationSources(ctx, root, true)
		if err != nil {
			return err
		}
		for _, path := range paths {
			// History files, including Archive Records, keep their own checks.
			if strings.HasPrefix(path, historyRoot+"/") {
				continue
			}
			content, err := os.ReadFile(filepath.Join(root, path))
			if err != nil {
				return fmt.Errorf("read history citation source %s: %w", path, err)
			}
			found, err := spec.HistoryCitationsInSource(content, path, candidates)
			if err != nil {
				return err
			}
			for _, citation := range found {
				if pending[citation.Target] || missingLinks[citation] {
					continue
				}
				locator, err := locate(citation.Target, "", "")
				if err != nil {
					return err
				}
				plain := true
				switch strings.ToLower(filepath.Ext(path)) {
				case ".md", ".markdown", ".mdx":
					plain = false
				}
				report(path, citation.Line, "cites a removed history file without a Git Locator: "+citation.Target,
					"Replace the citation with "+locator.Text(plain)+".")
			}
		}
	}
	invalid, err := invalidHistoryLocators(ctx, root, citations)
	if err != nil {
		return err
	}
	for _, c := range invalid {
		replacement, err := locate(c.target, c.commit, c.pullRequest)
		if err != nil {
			return err
		}
		report(c.path, c.line, "names an invalid Git Locator "+c.target+"@"+c.commit, "Replace the Git Locator with "+replacement.Text(false)+".")
	}
	return nil
}

// Validate all objects in one Git process and prove reachability once per
// distinct commit, even when many files cite different paths at that commit.
type historyLocatorCitation struct {
	path, target, commit, pullRequest string
	line                              int
}

func invalidHistoryLocators(ctx context.Context, root string, citations []historyLocatorCitation) ([]historyLocatorCitation, error) {
	if len(citations) == 0 {
		return nil, nil
	}
	var input strings.Builder
	keys := []string{}
	holds := map[string]bool{}
	for _, c := range citations {
		key := c.commit + ":" + c.target
		if _, seen := holds[key]; seen {
			continue
		}
		holds[key] = false
		keys = append(keys, key)
		input.WriteString(key + "\n")
	}
	command := exec.CommandContext(ctx, "git", "-C", root, "-c", "core.fsmonitor=false", "cat-file", "--batch-check")
	command.Env = mechanicalGitEnvironment()
	command.Stdin = strings.NewReader(input.String())
	var diagnostic bytes.Buffer
	command.Stderr = &diagnostic
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("validate Git Locator objects: %s: %w", strings.TrimSpace(diagnostic.String()), err)
	}
	lines := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	if len(lines) != len(keys) {
		return nil, fmt.Errorf("validate Git Locator objects: got %d responses for %d paths", len(lines), len(keys))
	}
	for i, key := range keys {
		holds[key] = !strings.HasSuffix(lines[i], " missing")
	}
	reachable := map[string]bool{}
	var invalid []historyLocatorCitation
	for _, c := range citations {
		held := holds[c.commit+":"+c.target]
		if held {
			if _, known := reachable[c.commit]; !known {
				locator, err := spec.ResolveGitLocator(ctx, root, c.target, c.commit, "")
				if err != nil {
					return nil, err
				}
				reachable[c.commit] = locator.Commit == c.commit
			}
		}
		if !held || !reachable[c.commit] {
			invalid = append(invalid, c)
		}
	}
	return invalid, nil
}

func terminalHistoryFinding(status string) bool {
	switch status {
	case "done", "deferred", "deprecated", "superseded", "closed", "cancelled":
		return true
	}
	return false
}
