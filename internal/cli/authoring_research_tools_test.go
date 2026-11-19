package cli

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestTheChecksNameTheRecommendedResearchToolsWhenResearchDeclaresNone(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, research string
	}{
		{"no_tools_key", "{model: sonnet}"},
		{"empty_tools_list", "{model: sonnet, tools: []}"},
		{"web_tools", "{model: sonnet, tools: [WebSearch, WebFetch]}"},
		{"mcp_tools_only", "{model: sonnet, tools: [mcp__kb__read]}"},
		{"a_problem_keeps_its_report", "{model: sonnet}"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home, root := withCLIWorkspace(t)
			writeAuthoringCommandFile(t, root, ".roundfixrc.yml", "skills:\n  authoring:\n    claude:\n      prd: {model: opus, reasoning_effort: xhigh}\n      research: "+tc.research+"\n")
			installAuthoringReportFixture(t)
			linkAuthoringReportFixture(t, root)
			wantDoctor := "authoring-models: ok (claude: prd, research applied)\n"
			wantJSON := `{"state":"applied","harness":"claude","declared":["prd","research"],"problems":[]}`
			wantText := ""
			switch tc.name {
			case "no_tools_key":
				wantDoctor = "authoring-models: warn (claude: prd, research applied; research declares no tools, so the Research Subagent cannot read the web; next: add tools: [WebSearch, WebFetch] under skills.authoring.claude.research in .roundfixrc.yml, then run roundfix skills install --target project)\n"
				wantJSON = `{"state":"applied","harness":"claude","declared":["prd","research"],"problems":[],"recommendedResearchTools":["WebSearch","WebFetch"]}`
				wantText = "Authoring models: research declares no tools, so the Research Subagent cannot read the web; add tools: [WebSearch, WebFetch] under skills.authoring.claude.research in .roundfixrc.yml, then apply with roundfix skills install --target project.\n"
			case "a_problem_keeps_its_report":
				writeAuthoringCommandFile(t, root, ".claude/agents/roundfix-tasks-author.md", "---\nname: roundfix-tasks-author\n---\n\n<!-- roundfix:generated skills.authoring.claude.tasks -->\n")
				wantDoctor = "authoring-models: warn (the installed files do not match the Authoring Model Declaration: .claude/agents/roundfix-tasks-author.md has no declared entry; next: roundfix skills install --target project)\n"
				wantJSON = `{"state":"unapplied","harness":"claude","declared":["prd","research"],"problems":[".claude/agents/roundfix-tasks-author.md has no declared entry"],"next":"roundfix skills install --target project"}`
				wantText = "Authoring models: the installed files do not match the declaration: .claude/agents/roundfix-tasks-author.md has no declared entry; apply with roundfix skills install --target project.\n"
			}
			before, beforeHome := snapshotDoctorPath(t, root), snapshotDoctorPath(t, home)
			configureAuthoringDoctorFixture(t)
			var stdout, stderr bytes.Buffer
			if code := runCLI(t, []string{"doctor"}, &stdout, &stderr); code != 0 || stderr.Len() != 0 {
				t.Fatalf("doctor exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
			}
			assertAuthoringRecommendationPosition(t, stdout.String(), "skills:", "authoring-models:", wantDoctor, "baseline-rules:")
			configureAuthoringUpdateFixture(t)
			raw := runUpdateCheckTest(t, "update", "check", "--format", "json")
			var fields map[string]json.RawMessage
			if err := json.Unmarshal([]byte(raw), &fields); err != nil {
				t.Fatal(err)
			}
			if got := string(fields["authoringModels"]); got != wantJSON {
				t.Fatalf("authoringModels=%s want=%s", got, wantJSON)
			}
			if string(fields["next"]) != `"none"` {
				t.Fatalf("changed top-level next: %s", raw)
			}
			output := runUpdateCheckTest(t, "update", "check")
			assertAuthoringRecommendationPosition(t, output, "Baseline:", "Authoring models:", wantText, "")
			if strings.Contains(output, "Next:") {
				t.Fatalf("changed text next: %s", output)
			}
			if !reflect.DeepEqual(before, snapshotDoctorPath(t, root)) || !reflect.DeepEqual(beforeHome, snapshotDoctorPath(t, home)) {
				t.Fatal("checks changed the repository or fixture home")
			}
		})
	}
}
