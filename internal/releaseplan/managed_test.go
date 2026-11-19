package releaseplan

import (
	"context"
	"errors"
	"testing"
)

func TestAnAdopterManagedPathIsMaintenanceOnly(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name      string
		paths     []string
		crosses   bool
		noAdopter bool
	}{
		{name: "the_setup_manifest", paths: []string{"docs/agents/setup-context.json"}},
		{name: "a_recorded_guide", paths: []string{"docs/agents/backend.md"}},
		{name: "skills_lock", paths: []string{"skills-lock.json"}},
		{name: "a_locked_skill_file", paths: []string{".agents/skills/write-tasks/references/qa.md"}},
		{name: "a_claude_skill_link", paths: []string{".claude/skills/roundfix"}},
		{name: "claude_md", paths: []string{"CLAUDE.md"}},
		{name: "specific_repository", paths: []string{"docs/agents/specific-repository.md"}},
		{name: "a_generated_author_subagent", paths: []string{".claude/agents/roundfix-tasks-author.md"}},
		{name: "project_config", paths: []string{".roundfixrc.yml"}},
		{name: "a_nested_context_md", paths: []string{"packages/backend/src/contexts/tax/CONTEXT.md"}},
		{name: "an_unlocked_local_skill_crosses", paths: []string{".agents/skills/revisao/scripts/run.ts"}, crosses: true},
		{name: "another_claude_agent_crosses", paths: []string{".claude/agents/reviewer.md"}, crosses: true},
		{name: "context_map_crosses", paths: []string{"CONTEXT-MAP.md"}, crosses: true},
		{name: "an_unrecorded_guide_crosses", paths: []string{"docs/agents/frontend.md"}, crosses: true},
		{name: "a_mixed_commit_crosses", paths: []string{"AGENTS.md", "packages/backend/.env.example"}, crosses: true},
		{name: "not_an_adopter_changes_nothing", paths: []string{"CLAUDE.md"}, crosses: true, noAdopter: true},
		{name: "an_empty_commit_still_crosses", crosses: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ClassifyChanges(ClassifyRequest{
				Commits: []Commit{{SHA: "change", Subject: "chore: maintain adopter", ChangedPaths: tt.paths}},
				Managed: ManagedPaths{Adopter: !tt.noAdopter,
					Files: []string{"AGENTS.md", "docs/agents/backend.md"}, Skills: []string{"roundfix", "write-tasks"}},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Changes) != 1 || result.Changes[0].CrossesMaintenanceOnlyBoundary != tt.crosses {
				t.Fatalf("changes = %+v, want crossing %v", result.Changes, tt.crosses)
			}
		})
	}
}

func TestManagedPathsRejectInvalidPathsAndSkillNames(t *testing.T) {
	t.Parallel()
	managed := ManagedPaths{Adopter: true, Files: []string{"../outside"}, Skills: []string{"", ".", "..", "nested/skill", `nested\skill`, "roundfix"}}
	for _, name := range []string{"", ".", "..", "/CONTEXT.md", "../CONTEXT.md", "../outside",
		".agents/skills/nested/skill/SKILL.md", ".agents/skills/unlocked/SKILL.md",
		".agents/skills/roundfix-extra/SKILL.md", ".claude/agents/nested/roundfix-author.md"} {
		t.Run(name, func(t *testing.T) {
			if managed.Contains(name) {
				t.Fatalf("invalid or unrecorded path %q is contained", name)
			}
		})
	}
}

func TestManagedMaintenanceIgnoresConventionalImpact(t *testing.T) {
	t.Parallel()
	for _, subject := range []string{"feat: update instructions", "fix: correct instructions", "feat!: replace instructions"} {
		t.Run(subject, func(t *testing.T) {
			commit := Commit{SHA: "instructions", Subject: subject, ChangedPaths: []string{"CLAUDE.md", "docs/specs/0001-x/_prd.md"}}
			result, err := ClassifyChanges(ClassifyRequest{Commits: []Commit{commit}, Managed: ManagedPaths{Adopter: true}})
			if err != nil || result.State != StateNoRelease || result.Classification.Impact != ImpactNone {
				t.Fatalf("classification = %+v, err = %v, want no_release", result, err)
			}
			if !ClassifyCommit(commit).CrossesMaintenanceOnlyBoundary {
				t.Fatal("ClassifyCommit changed its legacy boundary")
			}
		})
	}
}

type managedBuildSource struct {
	ceilingBuildSource
	managed    ManagedPaths
	managedErr error
}

func (source *managedBuildSource) ManagedPaths(context.Context, Range) (ManagedPaths, error) {
	return source.managed, source.managedErr
}

func TestBuildReadsManagedPathsBeforeCheckingReleasePreparation(t *testing.T) {
	t.Parallel()
	source := &managedBuildSource{
		ceilingBuildSource: ceilingBuildSource{commits: []Commit{{SHA: "maintenance", Subject: "chore: instructions", ChangedPaths: []string{"CLAUDE.md"}}}},
		managed:            ManagedPaths{Adopter: true},
	}
	plan, err := Build(t.Context(), Request{}, source)
	if err != nil || plan.State != StateNoRelease || source.calls != 0 {
		t.Fatalf("Build = %+v, err = %v, grant reads = %d; want no_release without reading grants", plan, err, source.calls)
	}
	errRead := errors.New("cannot read managed records")
	source.managedErr = errRead
	if _, err := Build(t.Context(), Request{}, source); !errors.Is(err, errRead) {
		t.Fatalf("Build error = %v, want managed source error", err)
	}
}
