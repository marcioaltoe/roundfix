// Suite: auditing-binary age evidence.
// Invariant: evidence is resolved only where it is meaningful, and its absence is never read as current.
// Boundary IN: the audited repository's Git objects and declared Roundfix version.
// Boundary OUT: the comparison itself, owned by internal/app.
package spec

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"roundfix/internal/app"
)

// AuditorEvidence is what one repository can say about the auditing binary's
// age. Both fields are optional: an absent signal means the question could not
// be answered here, never that the binary is current.
type AuditorEvidence struct {
	TreeVersion string
	Ancestry    app.AncestryResult
}

// roundfixVersionManifest is the file the Roundfix repository declares its own
// released version in. Its absence is the ordinary case: an adopting repository
// declares no Roundfix version, so the version signal cannot answer there.
var roundfixVersionManifest = filepath.Join("dist", "npm", "roundfix", "package.json")

// ResolveAuditorEvidence collects what the audited repository can say about the
// auditing binary's age.
//
// Both signals are only meaningful when the audited tree is the repository the
// binary was built from. Across repositories — the ordinary fleet case, where
// Roundfix audits an adopting repository — a Roundfix build commit is not an
// object in that tree and the tree declares no Roundfix version, so neither
// answers and the caller reports unknown. That is the honest result, not a
// degraded one: the alternative is comparing a commit against history it never
// belonged to.
func ResolveAuditorEvidence(ctx context.Context, repoRoot string, binary app.AuditingBinary) AuditorEvidence {
	evidence := AuditorEvidence{Ancestry: app.AncestryUnknown}
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot == "" {
		return evidence
	}

	commit := auditorBuildCommit(binary)
	if commit != "" && gitObjectExists(ctx, repoRoot, commit) {
		evidence.Ancestry = auditorAncestry(ctx, repoRoot, commit)
	}
	if evidence.Ancestry == app.AncestryUnknown {
		// The version signal is the fallback a released build needs: it leaves
		// its build commit empty by design, so ancestry can never answer for
		// the binaries most readers run.
		evidence.TreeVersion = declaredRoundfixVersion(repoRoot)
	}
	return evidence
}

// auditorBuildCommit strips the dirty marker the Makefile appends. The marker
// says the build carried uncommitted changes; the commit beside it is still the
// best available anchor, and `auditing_binary` shows the marker to the reader.
func auditorBuildCommit(binary app.AuditingBinary) string {
	return strings.TrimSuffix(strings.TrimSpace(binary.Commit), "-dirty")
}

func gitObjectExists(ctx context.Context, repoRoot string, commit string) bool {
	command := exec.CommandContext(ctx, "git", "-C", repoRoot, "cat-file", "-e", commit+"^{commit}")
	return command.Run() == nil
}

func auditorAncestry(ctx context.Context, repoRoot string, commit string) app.AncestryResult {
	head, err := gitOutput(ctx, repoRoot, "rev-parse", "HEAD")
	if err != nil {
		return app.AncestryUnknown
	}
	resolved, err := gitOutput(ctx, repoRoot, "rev-parse", commit+"^{commit}")
	if err != nil {
		return app.AncestryUnknown
	}
	if resolved == head {
		return app.AncestryNotOlder
	}
	command := exec.CommandContext(ctx, "git", "-C", repoRoot, "merge-base", "--is-ancestor", resolved, head)
	if err := command.Run(); err == nil {
		return app.AncestryOlder
	} else {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			// Not an ancestor and not equal: built from divergent history, so
			// it does not predate this tree even though it differs from it.
			return app.AncestryNotOlder
		}
	}
	return app.AncestryUnknown
}

func gitOutput(ctx context.Context, repoRoot string, args ...string) (string, error) {
	command := exec.CommandContext(ctx, "git", append([]string{"-C", repoRoot}, args...)...)
	output, err := command.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func declaredRoundfixVersion(repoRoot string) string {
	content, err := os.ReadFile(filepath.Join(repoRoot, roundfixVersionManifest))
	if err != nil {
		return ""
	}
	var manifest struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(content, &manifest); err != nil {
		return ""
	}
	return strings.TrimSpace(manifest.Version)
}
