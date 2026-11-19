package cli

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/preflight"
	"roundfix/internal/releaseplan"
)

func TestTheReleasePlanReadsTheManagedRecordsAtBaseAndTarget(t *testing.T) {
	t.Parallel()
	const adopterLock = `{"skills":{"roundfix":{},"write-tasks":{}}}`
	const manifest = `{"managedArtifacts":[{"path":"AGENTS.md"},{"path":"docs/agents/backend.md"}]}`
	for _, name := range []string{
		"the_fiscus_shape", "roundfix_itself", "a_skill_removed_in_the_range",
		"a_malformed_lock_makes_no_adopter", "project_config_is_not_read",
	} {
		t.Run(name, func(t *testing.T) {
			repo := newEmptyReleasePlanGitRepo(t)
			lock := adopterLock
			if name == "roundfix_itself" {
				lock = `{"skills":{"write-tasks":{}}}`
			}
			if name == "a_skill_removed_in_the_range" {
				lock = `{"skills":{"roundfix":{},"old-skill":{}}}`
				writeReleasePlanFile(t, repo, ".agents/skills/old-skill/SKILL.md", "old skill\n")
			}
			writeReleasePlanFile(t, repo, "skills-lock.json", lock)
			writeReleasePlanFile(t, repo, "docs/agents/setup-context.json", manifest)
			gitReleasePlan(t, repo, "add", "-A")
			commitReleasePlan(t, repo, "chore: seed managed records")
			gitReleasePlan(t, repo, "tag", "v0.1.0")

			var blocking []string
			wantState, wantExit := releaseplan.StateNoRelease, exitOK
			switch name {
			case "the_fiscus_shape", "roundfix_itself":
				commits := []struct {
					subject string
					files   map[string]string
				}{
					{"chore: update Roundfix", map[string]string{
						".agents/skills/write-tasks/SKILL.md": "updated skill\n",
						"AGENTS.md":                           "updated instructions\n",
						"docs/agents/setup-context.json":      manifest + "\n",
						"skills-lock.json":                    lock + "\n",
					}},
					{"chore: adopt Roundfix recommendations", map[string]string{
						".claude/agents/roundfix-tasks-author.md": "generated author\n",
						".roundfixrc.yml":                         "profiles: {}\n",
						"docs/agents/specific-repository.md":      "repository rules\n",
					}},
					{"docs: plan widget", map[string]string{
						"docs/specs/0097-widget/_prd.md":  "widget plan\n",
						"packages/backend/src/CONTEXT.md": "backend vocabulary\n",
					}},
					{"docs: plan move", map[string]string{
						"docs/specs/0099-move/_prd.md":  "move plan\n",
						"docs/user-guide/move.md":       "move guide\n",
						"packages/backend/.env.example": "EXAMPLE=value\n",
					}},
				}
				for index, commit := range commits {
					for path, content := range commit.files {
						writeReleasePlanFile(t, repo, path, content)
					}
					gitReleasePlan(t, repo, "add", "-A")
					sha := commitReleasePlan(t, repo, commit.subject)
					if name == "roundfix_itself" || index == len(commits)-1 {
						blocking = append(blocking, sha)
					}
				}
				wantState, wantExit = releaseplan.StateManualClassificationRequired, exitUnverified
			case "a_skill_removed_in_the_range":
				if err := os.Remove(filepath.Join(repo, ".agents/skills/old-skill/SKILL.md")); err != nil {
					t.Fatal(err)
				}
				writeReleasePlanFile(t, repo, "skills-lock.json", `{"skills":{"roundfix":{}}}`)
				gitReleasePlan(t, repo, "add", "-A")
				commitReleasePlan(t, repo, "chore: retire old skill")
			case "a_malformed_lock_makes_no_adopter":
				writeReleasePlanFile(t, repo, "skills-lock.json", "not JSON\n")
				writeReleasePlanFile(t, repo, "CLAUDE.md", "instructions\n")
				gitReleasePlan(t, repo, "add", "-A")
				blocking = append(blocking, commitReleasePlan(t, repo, "chore: change instructions"))
				wantState, wantExit = releaseplan.StateManualClassificationRequired, exitUnverified
			case "project_config_is_not_read":
				writeReleasePlanFile(t, repo, ".roundfixrc.yml", "invalid: [yaml\n")
				writeReleasePlanFile(t, repo, "packages/backend/CONTEXT.md", "backend vocabulary\n")
				gitReleasePlan(t, repo, "add", "-A")
				commitReleasePlan(t, repo, "docs: update configuration and vocabulary")
			}

			code, stdout, stderr := runReleasePlanCommandInRepo(t, repo, "--from", "v0.1.0", "--format", "json")
			if code != wantExit || stderr != "" {
				t.Fatalf("exit = %d, want %d; stdout=%s stderr=%s", code, wantExit, stdout, stderr)
			}
			plan := decodeReleasePlanJSON(t, stdout)
			if plan.State != wantState || !reflect.DeepEqual(plan.Classification.BlockingCommits, blocking) {
				t.Fatalf("state/blocking = %s/%v, want %s/%v", plan.State, plan.Classification.BlockingCommits, wantState, blocking)
			}
			for _, change := range plan.Changes {
				wantCrossing := false
				for _, sha := range blocking {
					wantCrossing = wantCrossing || change.CommitSHA == sha
				}
				if change.CrossesMaintenanceOnlyBoundary != wantCrossing {
					t.Errorf("%s crossing = %v, want %v", change.Subject, change.CrossesMaintenanceOnlyBoundary, wantCrossing)
				}
			}
		})
	}
}

func TestManagedRecordsUseCommittedEndpoints(t *testing.T) {
	t.Parallel()
	repo := newEmptyReleasePlanGitRepo(t)
	writeReleasePlanFile(t, repo, "skills-lock.json", `{"skills":{"roundfix":{},"old-skill":{}}}`)
	writeReleasePlanFile(t, repo, "docs/agents/setup-context.json", `{"managedArtifacts":[{"path":"docs/agents/old.md"}]}`)
	gitReleasePlan(t, repo, "add", "-A")
	commitReleasePlan(t, repo, "chore: seed records")
	gitReleasePlan(t, repo, "tag", "v0.1.0")
	writeReleasePlanFile(t, repo, "skills-lock.json", `{"skills":{"roundfix":{},"new-skill":{}}}`)
	writeReleasePlanFile(t, repo, "docs/agents/setup-context.json", `{"managedArtifacts":[{"path":"docs/agents/new.md"}]}`)
	gitReleasePlan(t, repo, "add", "-A")
	commitReleasePlan(t, repo, "chore: replace recorded paths")
	source := newReleasePlanGitSource(repo, nil)
	releaseRange, err := source.ResolveRange(t.Context(), "v0.1.0", "")
	if err != nil {
		t.Fatal(err)
	}
	// Neither the target records nor the base records come from working-tree bytes.
	writeReleasePlanFile(t, repo, "skills-lock.json", `{"skills":{"uncommitted":{}}}`)
	writeReleasePlanFile(t, repo, "docs/agents/setup-context.json", `{"managedArtifacts":[{"path":"uncommitted.md"}]}`)
	managed, err := source.ManagedPaths(t.Context(), releaseRange)
	want := releaseplan.ManagedPaths{Adopter: true,
		Files: []string{"docs/agents/new.md", "docs/agents/old.md"}, Skills: []string{"new-skill", "old-skill", "roundfix"}}
	if err != nil || !reflect.DeepEqual(managed, want) {
		t.Fatalf("ManagedPaths = %+v, err = %v, want %+v", managed, err, want)
	}
	for _, operation := range []string{"rev-parse", "ls-tree", "show"} {
		t.Run(operation+"_failure", func(t *testing.T) {
			errGit := errors.New("managed Git read failed")
			source := newReleasePlanGitSource(repo, managedReadFailingGit{operation: operation, err: errGit})
			_, err := source.ManagedPaths(t.Context(), releaseRange)
			var sourceError releaseplan.GitSourceError
			if !errors.As(err, &sourceError) || !errors.Is(err, errGit) {
				t.Fatalf("ManagedPaths error = %v, want GitSourceError wrapping read failure", err)
			}
		})
	}
}

type managedReadFailingGit struct {
	operation string
	err       error
}

func (runner managedReadFailingGit) RunGit(ctx context.Context, dir string, args ...string) (string, error) {
	if len(args) > 0 && args[0] == runner.operation {
		return "", runner.err
	}
	return (preflight.ExecGitRunner{}).RunGit(ctx, dir, args...)
}

func TestInvalidManagedRecordsContributeNothing(t *testing.T) {
	t.Parallel()
	for name, lock := range map[string]string{
		"missing": "", "malformed": "{", "non_object_skills": `{"skills":["roundfix"]}`, "null_skills": `{"skills":null}`,
	} {
		t.Run(name, func(t *testing.T) {
			repo := newEmptyReleasePlanGitRepo(t)
			writeReleasePlanFile(t, repo, "seed.txt", "seed\n")
			if lock != "" {
				writeReleasePlanFile(t, repo, "skills-lock.json", lock)
			}
			writeReleasePlanFile(t, repo, "docs/agents/setup-context.json", "not JSON\n")
			gitReleasePlan(t, repo, "add", "-A")
			sha := commitReleasePlan(t, repo, "chore: seed invalid records")
			source := newReleasePlanGitSource(repo, nil)
			managed, err := source.ManagedPaths(t.Context(), releaseplan.Range{
				Base: releaseplan.VersionRef{CommitSHA: sha}, Target: releaseplan.RevisionRef{CommitSHA: sha},
			})
			if err != nil || managed.Adopter || len(managed.Files) != 0 || len(managed.Skills) != 0 {
				t.Fatalf("ManagedPaths = %+v, err = %v, want empty non-adopter", managed, err)
			}
			if strings.TrimSpace(gitReleasePlanOutput(t, repo, "status", "--porcelain")) != "" {
				t.Fatal("managed record read changed the repository")
			}
		})
	}
}
