package skills

import (
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

func roundfixWordingBundle() fstest.MapFS {
	return fstest.MapFS{
		"roundfix/SKILL.md":           {Data: []byte(roundfixWordingFixture)},
		"roundfix/agents/openai.yaml": {Data: []byte(roundfixManifestFixture)},
	}
}

func TestRoundfixWordingIsSatisfiedByAReferenceFile(t *testing.T) {
	files := roundfixWordingBundle()
	const phrase = "Prefer `roundfix` commands over manual GitHub scraping."
	files["roundfix/SKILL.md"].Data = []byte(strings.ReplaceAll(roundfixWordingFixture, phrase, ""))
	files["roundfix/references/setup.md"] = &fstest.MapFile{Data: []byte(phrase)}
	if got := checkRoundfixWording(files); len(got) != 0 {
		t.Fatalf("wording diagnostics = %v, want none", got)
	}
}

func TestRoundfixWordingReportsAPhraseMissingFromEveryFile(t *testing.T) {
	files := roundfixWordingBundle()
	const phrase = "Prefer `roundfix` commands over manual GitHub scraping."
	files["roundfix/SKILL.md"].Data = []byte(strings.ReplaceAll(roundfixWordingFixture, phrase, ""))
	files["roundfix/references/setup.md"] = &fstest.MapFile{Data: []byte("other guidance")}
	want := []Diagnostic{{Path: "roundfix/SKILL.md", Message: "missing required wording \"" + phrase + "\""}}
	if got := checkRoundfixWording(files); !reflect.DeepEqual(got, want) {
		t.Fatalf("wording diagnostics = %v, want %v", got, want)
	}
}

func TestRoundfixWordingReportsBannedBrandingInAReference(t *testing.T) {
	for _, phrase := range []string{"reference project", "Reference Project"} {
		t.Run(phrase, func(t *testing.T) {
			files := roundfixWordingBundle()
			files["roundfix/references/setup.md"] = &fstest.MapFile{Data: []byte(phrase)}
			want := []Diagnostic{{Path: "roundfix/SKILL.md", Message: "contains banned reference branding \"" + phrase + "\""}}
			if got := checkRoundfixWording(files); !reflect.DeepEqual(got, want) {
				t.Fatalf("wording diagnostics = %v, want %v", got, want)
			}
		})
	}
}

const roundfixWordingFixture = "Roundfix\n" +
	"Prefer `roundfix` commands over manual GitHub scraping.\n" +
	"Report the Run ID\n" +
	"state whenever you summarize progress.\n" +
	"Review Runs (`fetch`, `resolve`, and `watch`) execute in the user's checkout\n" +
	"Branch Integrity Preflight runs before any fetch, Agent Session\n" +
	"`--skip-branch-integrity` is the only bypass\n" +
	"watch ends CleanUnverified, exits `3`\n" +
	"Roundfix publishes Outcome Comments\n" +
	"adapter: ok\n" +
	"profiles: ok\n" +
	"@agentclientprotocol/codex-acp\n" +
	"`--agent`, `--model`, and `--reasoning-effort` are all-or-none\n" +
	"no liveness signal\n" +
	"commit <path>\n" +
	"Settle surface: <path>\n" +
	"not advertised by runtime\n" +
	"two spaces followed by `reason: <one line>`\n" +
	"owner PID is provably dead\n" +
	"`done` becomes `completed`\n" +
	"This Verification Feedback retry never consumes a Round\n" +
	"Do not manually resolve CodeRabbit threads\n" +
	"Read every assigned Review Issue file completely\n" +
	"Update only assigned Review Issue statuses\n" +
	"Do not create commits inside an assigned Batch run.\n" +
	"Do not push inside an assigned Batch run.\n" +
	"Do not call GitHub, CodeRabbit, or other Review Source mutation APIs\n" +
	"Do not edit unassigned Review Issue files.\n" +
	"Do not mark any issue as `duplicated`\n" +
	"rtk bun run --cwd <package-dir> <script> [args...]\n"

const roundfixManifestFixture = `name: roundfix
entrypoint: SKILL.md
runtime_hints:
  command: roundfix watch --source coderabbit --pr <number> --until-clean
  complete_override_command: roundfix implement --spec <slug> --agent <agent> --model <model> --reasoning-effort <effort>
  profile_readiness_command: roundfix profiles validate --json
  agent_selection_contract: all-or-none
  review_run_contract: user's checkout
  watch_outcome_contract: CleanUnverified with exit code
  review_source_contract: Outcome Comments
  owns: assigned Review Issue files during Batch runs
  reports: Run state
`
