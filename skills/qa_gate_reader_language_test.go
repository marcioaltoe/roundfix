package skills

import (
	"strings"
	"testing"
)

func TestTheQAGateReportsTextOutsideTheReaderLanguage(t *testing.T) {
	t.Parallel()
	data, err := embedded.ReadFile("qa-gate/SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	const heading = "### Reader Language"
	const paragraph = "When `docs/agents/agent-instructions.md` names a Reader Language, read the text a user reads on every surface the Spec declares: screen copy, command output, messages, e-mails, PDFs, user documentation, and the names, descriptions and parameters of MCP tools. Record each text written in another language as a finding with the surface, the text and the expected language, and classify it by user impact as any other finding. A business term the guide lets keep its original form is not a finding. When the guide names no Reader Language, report nothing for it."
	frontend := strings.Index(text, "### Frontend sweep\n")
	if frontend < 0 {
		t.Fatal("missing Frontend sweep")
	}
	next := text[frontend+len("### Frontend sweep\n"):]
	section := strings.Index(next, "\n### ")
	if section < 0 || !strings.HasPrefix(next[section:], "\n"+heading+"\n\n"+paragraph+"\n") || strings.Count(text, heading) != 1 {
		t.Fatal("Reader Language section must follow Frontend sweep with its exact paragraph")
	}
}
