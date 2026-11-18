package config

import (
	"bytes"
	"path/filepath"
	"reflect"
	"testing"

	"roundfix/internal/gittest"
)

const recommendedAuthoringYAML = `skills:
  authoring:
    claude:
      idea:
        model: fable
        reasoning_effort: xhigh
      prd:
        model: fable
        reasoning_effort: high
      techspec:
        model: fable
        reasoning_effort: high
      tasks:
        model: opus
        reasoning_effort: high
      research:
        model: sonnet
        reasoning_effort: medium
`

func TestRecommendedAuthoringDeclarationRendersValidProjectConfig(t *testing.T) {
	t.Parallel()
	t.Run("exact_bytes", func(t *testing.T) {
		if got := RecommendedAuthoringYAML(); got != recommendedAuthoringYAML {
			t.Fatalf("rendering=%q want=%q", got, recommendedAuthoringYAML)
		}
	})
	t.Run("loads_as_project_config", func(t *testing.T) {
		home, root := t.TempDir(), t.TempDir()
		gittest.MarkRepository(t, root)
		path := filepath.Join(root, ".roundfixrc.yml")
		mustWrite(t, path, RecommendedAuthoringYAML())
		var warnings bytes.Buffer
		loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: root, Stderr: &warnings})
		if err != nil || warnings.Len() != 0 {
			t.Fatalf("load=%v warnings=%s", err, &warnings)
		}
		want := RecommendedAuthoringDeclaration()
		got := loaded.Config.Skills.Authoring.Claude
		if loaded.ProjectConfigPath != path || !reflect.DeepEqual(got, &want) || len(got.Research.Tools) != 0 {
			t.Fatalf("loaded=%+v declaration=%+v want=%+v", loaded, got, want)
		}
	})
	t.Run("returns_a_copy", func(t *testing.T) {
		want := RecommendedAuthoringDeclaration()
		changed := RecommendedAuthoringDeclaration()
		for _, entry := range []*AuthoringEntry{changed.Idea, changed.PRD, changed.TechSpec, changed.Tasks, &changed.Research.AuthoringEntry} {
			entry.Model = "changed"
			entry.ReasoningEffort = "low"
		}
		changed.Research.Tools = []string{"WebSearch"}
		if got := RecommendedAuthoringDeclaration(); !reflect.DeepEqual(got, want) {
			t.Fatalf("mutated recommendation: got=%+v want=%+v", got, want)
		}
		if got := RecommendedAuthoringYAML(); got != recommendedAuthoringYAML {
			t.Fatalf("mutated rendering=%q", got)
		}
	})
}
