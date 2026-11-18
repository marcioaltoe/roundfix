package skills

import (
	"io/fs"
	"strings"
	"testing"
)

func TestTheAuthoringSkillsPassTheirModelToThePhaseCommand(t *testing.T) {
	t.Parallel()
	for copy, source := range roundfixLayoutSources() {
		t.Run(copy, func(t *testing.T) {
			t.Parallel()
			var shared string
			for _, phase := range []string{"idea", "prd", "techspec", "tasks"} {
				data, err := fs.ReadFile(source, "write-"+phase+"/SKILL.md")
				if err != nil {
					t.Fatal(err)
				}
				text := string(data)
				original := "Once the Spec's slug is known, record the phase: run `roundfix phase start " + phase + " --spec <slug>` in the session or the Author Subagent that runs it. When the installed Roundfix does not know the command, say so in the report and continue."
				paragraph := "For `roundfix phase start " + phase + " --spec <slug>`, add `--model` with the exact model identifier the harness states for this session, never an alias and never a guess. Leave `--model` out when the session's instructions state none. Add `--author-subagent` when the phase's Author Subagent runs the phase. When the installed Roundfix refuses these flags, run the command without them and say so in the report."
				if strings.Count(text, original+"\n\n"+paragraph) != 1 {
					t.Errorf("write-%s: model instructions must follow the preserved phase-record paragraph exactly once", phase)
				}
				normalized := strings.ReplaceAll(paragraph, "phase start "+phase, "phase start <phase>")
				if shared != "" && normalized != shared {
					t.Errorf("write-%s: paragraph differs from the other authoring phases", phase)
				}
				shared = normalized
			}
		})
	}
}
