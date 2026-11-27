package cli

import (
	"strings"
	"testing"
)

func TestTheReviewPromptCarriesTheLanguageRuleInstruction(t *testing.T) {
	t.Parallel()
	prompt := buildReviewPrompt("base", "head", "candidate change")
	if strings.Count(prompt, reviewLanguageRuleInstruction) != 1 || !strings.Contains(prompt, absenceClaimInstruction+reviewLanguageRuleInstruction+"\n") {
		t.Fatal("language-rule instruction must occur once directly after the absence-claim instruction")
	}
	for _, line := range []string{
		"If there are findings, respond with Findings: on the first line, followed by each finding as a list item that starts with `- ` and names its file and line. If there are no findings, respond exactly: No findings. Include no other prose.\n\n",
		"Start each finding with its anchor, `path:line` or `path:start-end`, naming a line of the candidate diff, and state what breaks in a clause that starts with `Failure:`.\n",
	} {
		if !strings.Contains(prompt, line) {
			t.Errorf("missing finding grammar %q", line)
		}
	}
	block := "\nDelivery Conventions (" + deliveryConventionsVersion + "). A Roundfix delivery writes the following by design; do not report them as defects:\n"
	for _, convention := range deliveryConventions() {
		block += convention.ID + ". " + convention.Prompt + "\n"
	}
	if !strings.Contains(prompt, block) {
		t.Fatal("missing versioned Delivery Conventions block")
	}
}
