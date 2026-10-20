package spec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"roundfix/internal/preflight"
)

// GitLocator names a removed repository path at a reachable commit, or records
// the pull request from which a reader can recover it.
type GitLocator struct {
	Path        string
	Commit      string
	PullRequest string
}

// Text omits delimiters when the caller is already inside Markdown code.
func (l GitLocator) Text(inCode bool) string {
	path := l.Path
	if l.Commit != "" {
		path += "@" + l.Commit
	}
	if !inCode {
		path = "`" + path + "`"
	}
	if l.Commit != "" {
		return path
	}
	if l.PullRequest != "" {
		return path + " (pull request #" + l.PullRequest + ")"
	}
	return path + " (not reachable from any branch or tag)"
}

func gitLocatorAbsent(err error) bool {
	var exit *exec.ExitError
	return errors.As(err, &exit) && exit.ExitCode() == 1
}

func gitLocatorCommit(ctx context.Context, root, revision string) (string, error) {
	out, err := archiveGitRead(ctx, root, "rev-parse", "--verify", "--quiet", revision+"^{commit}")
	if gitLocatorAbsent(err) && ctx.Err() == nil {
		return "", nil
	}
	return strings.TrimSpace(string(out)), err
}

// GitLocatorHolds checks a file or directory without mistaking a missing object
// for a repository or process failure. Batch mode reports absent objects as data.
func GitLocatorHolds(ctx context.Context, repositoryRoot, path, commit string) (bool, error) {
	if commit == "" {
		return false, nil
	}
	command := exec.CommandContext(ctx, "git", "-C", repositoryRoot, "-c", "core.fsmonitor=false", "cat-file", "--batch-check")
	command.Env = authorizationGitEnvironment()
	command.Stdin = strings.NewReader(commit + ":" + path + "\n")
	var diagnostic bytes.Buffer
	command.Stderr = &diagnostic
	out, err := command.Output()
	if err != nil {
		return false, fmt.Errorf("inspect Git Locator %s@%s: %s: %w", path, commit, strings.TrimSpace(diagnostic.String()), err)
	}
	return !strings.HasSuffix(strings.TrimSpace(string(out)), " missing"), nil
}

// gitLocatorBranchRunner lets DetectDefaultBranch retain its fallback behavior
// for an absent symbolic ref while returning actual Git failures to the caller.
type gitLocatorBranchRunner struct{ err error }

func (r *gitLocatorBranchRunner) RunGit(ctx context.Context, root string, args ...string) (string, error) {
	if len(args) > 0 && args[0] == "symbolic-ref" {
		args = append([]string{"symbolic-ref", "--quiet"}, args[1:]...)
	}
	out, err := archiveGitRead(ctx, root, args...)
	if err != nil && (!gitLocatorAbsent(err) || ctx.Err() != nil) {
		r.err = err
	}
	return strings.TrimSpace(string(out)), err
}

// ResolveGitLocator uses Default Branch Head and tag reachability only. It
// never fetches a branch or relies on another work branch keeping a commit alive.
func ResolveGitLocator(ctx context.Context, repositoryRoot, path, preferred, pullRequest string) (GitLocator, error) {
	locator := GitLocator{Path: path, PullRequest: pullRequest}
	runner := &gitLocatorBranchRunner{}
	current, _ := runner.RunGit(ctx, repositoryRoot, "symbolic-ref", "--short", "HEAD")
	if runner.err != nil {
		return locator, fmt.Errorf("read Git Locator branch: %w", runner.err)
	}
	defaultHead, found, err := preflight.ResolveDefaultBranchHead(ctx, repositoryRoot, current, runner)
	if err != nil {
		return locator, fmt.Errorf("resolve Git Locator default branch: %w", err)
	}
	if runner.err != nil {
		return locator, fmt.Errorf("detect Git Locator default branch: %w", runner.err)
	}
	head := ""
	if found {
		head = defaultHead.Commit
	}
	if preferred != "" {
		commit, err := gitLocatorCommit(ctx, repositoryRoot, preferred)
		if err != nil {
			return locator, fmt.Errorf("resolve Git Locator provenance: %w", err)
		}
		holds, err := GitLocatorHolds(ctx, repositoryRoot, path, commit)
		if err != nil {
			return locator, err
		}
		if holds {
			reachable := false
			if head != "" {
				_, err := archiveGitRead(ctx, repositoryRoot, "merge-base", "--is-ancestor", commit, head)
				if err != nil && (!gitLocatorAbsent(err) || ctx.Err() != nil) {
					return locator, fmt.Errorf("check Git Locator ancestry: %w", err)
				}
				reachable = err == nil
			}
			if !reachable {
				tags, err := archiveGitRead(ctx, repositoryRoot, "tag", "--contains", commit)
				if err != nil {
					return locator, fmt.Errorf("check Git Locator tags: %w", err)
				}
				reachable = len(bytes.TrimSpace(tags)) > 0
			}
			if reachable {
				locator.Commit = commit
				return locator, nil
			}
		}
	}
	if head == "" {
		return locator, nil
	}
	out, err := archiveGitRead(ctx, repositoryRoot, "log", "-1", "--format=%H", head, "--", path)
	if err != nil {
		return locator, fmt.Errorf("find Git Locator path history: %w", err)
	}
	commit := strings.TrimSpace(string(out))
	if commit == "" {
		return locator, nil
	}
	holds, err := GitLocatorHolds(ctx, repositoryRoot, path, commit)
	if err != nil {
		return locator, err
	}
	if !holds {
		commit, err = gitLocatorCommit(ctx, repositoryRoot, commit+"^1")
		if err != nil {
			return locator, fmt.Errorf("resolve Git Locator deletion parent: %w", err)
		}
		holds, err = GitLocatorHolds(ctx, repositoryRoot, path, commit)
		if err != nil {
			return locator, err
		}
	}
	if holds {
		locator.Commit = commit
	}
	return locator, nil
}

// CitationRewrite shares destination scanning and bounded path replacement
// between Archive Relink and text copied out of a removed folder.
type CitationRewrite struct {
	Moves    map[string]string
	Relative bool // Resolve ./ and ../ prose only for history conversion.
	Locate   func(path string) (GitLocator, bool, error)
}

func (rw CitationRewrite) locate(path string) (GitLocator, bool, error) {
	if rw.Locate == nil {
		return GitLocator{}, false, nil
	}
	return rw.Locate(path)
}

// Apply changes only destinations and bounded path text, retaining all other bytes.
func (rw CitationRewrite) Apply(content []byte, from, to string) ([]byte, int, error) {
	sources := make([]string, 0, len(rw.Moves))
	for source := range rw.Moves {
		if source != "" {
			sources = append(sources, source)
		}
	}
	sort.Slice(sources, func(i, j int) bool {
		if len(sources[i]) != len(sources[j]) {
			return len(sources[i]) > len(sources[j])
		}
		return sources[i] < sources[j]
	})
	// Whole-link replacements can enclose an image destination. Plan edits
	// before writing so nested destinations cannot move the output cursor back.
	type edit struct {
		start, end int
		text       string
		count      int
	}
	var edits []edit
	for _, span := range scanArchiveMarkdown(content) {
		destination := string(content[span.start:span.end])
		target, suffix, relative := archiveLinkTarget(destination, filepath.Dir(from))
		target = filepath.ToSlash(target)
		valid := relative && !filepath.IsAbs(target) && target != ".." && !strings.HasPrefix(target, "../")
		next, moved := rw.Moves[target]
		var locator GitLocator
		located := false
		if valid && !moved {
			var err error
			locator, located, err = rw.locate(target)
			if err != nil {
				return nil, 0, err
			}
		}
		if located {
			label := content[span.labelStart:span.labelEnd]
			nested := 0
			if !span.reference && len(scanArchiveMarkdown(label)) > 0 {
				var err error
				label, nested, err = rw.Apply(label, from, to)
				if err != nil {
					return nil, 0, err
				}
			}
			replacement := string(label) + " (" + locator.Text(false) + ")"
			if span.reference {
				replacement = `\[` + string(label) + "]: " + locator.Text(false)
			}
			edits = append(edits, edit{span.linkStart, span.linkEnd, replacement, 1 + nested})
			continue
		}
		replacement := destination
		if (relative && moved) || (valid && from != to) {
			if !moved {
				next = target
			}
			rel, err := filepath.Rel(filepath.Dir(to), filepath.FromSlash(next))
			if err != nil {
				// RelinkContent historically preserves an unrelativizable destination.
				if rw.Locate != nil || from != to {
					return nil, 0, fmt.Errorf("rebase citation %q: %w", destination, err)
				}
			} else {
				replacement = (&url.URL{Path: filepath.ToSlash(rel)}).EscapedPath() + suffix
			}
		}
		changed := 0
		if replacement != destination {
			changed = 1
		}
		edits = append(edits, edit{span.start, span.end, replacement, changed})
	}
	sort.SliceStable(edits, func(i, j int) bool {
		if edits[i].start != edits[j].start {
			return edits[i].start < edits[j].start
		}
		return edits[i].end > edits[j].end
	})
	code := archiveMarkdownCode(content)
	var out bytes.Buffer
	count, end := 0, 0
	for _, e := range edits {
		if e.start < end {
			continue
		} // an enclosing whole-link edit owns these bytes
		text, n, err := rw.paths(content[end:e.start], code[end:e.start], sources, from)
		if err != nil {
			return nil, 0, err
		}
		out.Write(text)
		out.WriteString(e.text)
		count += n + e.count
		end = e.end
	}
	text, n, err := rw.paths(content[end:], code[end:], sources, from)
	if err != nil {
		return nil, 0, err
	}
	out.Write(text)
	return out.Bytes(), count + n, nil
}

func citationPathBoundary(c byte) bool {
	return strings.ContainsRune(" \t\r\n`\"'()[]<>#.,;:", rune(c))
}

func (rw CitationRewrite) paths(content []byte, code []bool, sources []string, from string) ([]byte, int, error) {
	var out bytes.Buffer
	count := 0
	for i := 0; i < len(content); {
		matched := false
		if i == 0 || citationPathBoundary(content[i-1]) {
			if rw.Locate != nil {
				end := i
				for end < len(content) && !strings.ContainsRune(" \t\r\n`\"'()[]<>#", rune(content[end])) {
					end++
				}
				token := string(content[i:end])
				suffix := strings.TrimLeft(string(content[end:]), "`")
				if token != "" && (strings.Contains(token, "@") || strings.HasPrefix(suffix, " (not reachable from any branch or tag)") || strings.HasPrefix(suffix, " (pull request #")) {
					out.Write(content[i:end])
					i = end
					continue
				}
			}
			for _, source := range sources {
				if !bytes.HasPrefix(content[i:], []byte(source)) {
					continue
				}
				end := i + len(source)
				if end < len(content) && content[end] == '/' {
					if end+1 < len(content) && !citationPathBoundary(content[end+1]) {
						continue
					}
					end++
				} else if end < len(content) && !citationPathBoundary(content[end]) {
					continue
				}
				out.WriteString(rw.Moves[source])
				count++
				i = end
				matched = true
				break
			}
			if !matched && (rw.Locate != nil || rw.Relative) {
				// Paths can contain punctuation (including dots); strip sentence punctuation
				// only at the end, and never reinterpret an existing path@commit locator.
				end := i
				for end < len(content) && !strings.ContainsRune(" \t\r\n`\"'()[]<>#", rune(content[end])) {
					end++
				}
				path := strings.TrimRight(string(content[i:end]), ".,;:")
				stop := i + len(path)
				if path != "" {
					target := strings.TrimSuffix(path, "/")
					if rw.Relative {
						target = citationProseTarget(target, from)
					}
					if rw.Relative {
						if destination, moved := rw.Moves[target]; moved {
							out.WriteString(destination)
							count++
							i = stop
							matched = true
						}
					}
					locator, located, err := rw.locate(target)
					if err != nil {
						return nil, 0, err
					}
					if located && !matched {
						out.WriteString(locator.Text(code[i]))
						count++
						i = stop
						matched = true
					}
				}
			}
		}
		if !matched {
			out.WriteByte(content[i])
			i++
		}
	}
	return out.Bytes(), count, nil
}

// ApplyText changes bounded paths only. It never parses links or adds Markdown
// delimiters, so comments and source files receive plain Git Locators.
func (rw CitationRewrite) ApplyText(content []byte, path string) ([]byte, int, error) {
	sources := make([]string, 0, len(rw.Moves))
	for source := range rw.Moves {
		if source != "" {
			sources = append(sources, source)
		}
	}
	sort.Slice(sources, func(i, j int) bool {
		if len(sources[i]) != len(sources[j]) {
			return len(sources[i]) > len(sources[j])
		}
		return sources[i] < sources[j]
	})
	code := make([]bool, len(content))
	for i := range code {
		code[i] = true
	}
	return rw.paths(content, code, sources, path)
}

// Only explicit relative prose paths resolve against the citing directory.
func citationProseTarget(path, source string) string {
	if strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") {
		return filepath.ToSlash(filepath.Clean(filepath.Join(filepath.Dir(source), filepath.FromSlash(path))))
	}
	return path
}
