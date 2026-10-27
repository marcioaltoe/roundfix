package skills

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func roundfixSection(text, heading string) string {
	start := strings.Index(text, "## "+heading)
	if start < 0 {
		return ""
	}
	body := text[start+len("## "+heading):]
	if next := strings.Index(body, "\n## "); next >= 0 {
		body = body[:next]
	}
	return strings.TrimSpace(body)
}

func roundfixNumberedItem(section string, number int) string {
	prefix := fmt.Sprintf("%d. ", number)
	item := ""
	for _, line := range strings.Split(section, "\n") {
		if strings.HasPrefix(line, prefix) {
			item = line
		} else if item != "" {
			if regexp.MustCompile(`^[1-8]\. `).MatchString(line) {
				break
			}
			item += "\n" + line
		}
	}
	return strings.TrimSpace(item)
}

func roundfixWhitespace(text string) string {
	return strings.Join(strings.Fields(text), " ")
}

func TestTheFollowUpModeRunsItsStepsInOrder(t *testing.T) {
	t.Parallel()
	for name, source := range roundfixLayoutSources() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			data, err := fs.ReadFile(source, "roundfix/references/modes.md")
			if err != nil {
				t.Fatal(err)
			}
			text := string(data)
			section := roundfixSection(text, "--follow-up")
			for i := 1; i <= 8; i++ {
				if !regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(fmt.Sprintf("%d. ", i))).MatchString(section) {
					t.Errorf("follow-up section omits item %d", i)
				}
			}
			items := make([]string, 0, 8)
			current := ""
			for _, line := range strings.Split(section, "\n") {
				if regexp.MustCompile(`^[1-8]\. `).MatchString(line) {
					if current != "" {
						items = append(items, current)
					}
					current = line
				} else if current != "" {
					current += "\n" + line
				}
			}
			if current != "" {
				items = append(items, current)
			}
			if len(items) != 8 {
				t.Fatalf("follow-up items = %d, want 8", len(items))
			}
			for i, item := range items {
				if !strings.HasPrefix(item, fmt.Sprintf("%d. ", i+1)) {
					t.Fatalf("item %d out of order", i+1)
				}
			}
			checks := map[int]string{3: "--unattended", 4: "memory_message_pop", 6: "--strict", 7: "roundfix deliver", 8: "memory_message_send"}
			for itemNumber, want := range checks {
				if !strings.Contains(items[itemNumber-1], want) {
					t.Errorf("item %d omits %q", itemNumber, want)
				}
			}
			for _, forbidden := range []string{"resolved_to", "_triaged", "Inbox Entr", "secondbrain.root"} {
				if strings.Contains(text, forbidden) {
					t.Errorf("reference contains forbidden %q", forbidden)
				}
			}
		})
	}
}

func TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage(t *testing.T) {
	t.Parallel()
	for name, source := range roundfixLayoutSources() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			referenceData, err := fs.ReadFile(source, "roundfix/references/modes.md")
			if err != nil {
				t.Fatal(err)
			}
			reference := string(referenceData)
			followUp := roundfixSection(reference, "--follow-up")
			step1 := roundfixWhitespace(roundfixNumberedItem(followUp, 1))
			step4 := roundfixWhitespace(roundfixNumberedItem(followUp, 4))
			step8 := roundfixWhitespace(roundfixNumberedItem(followUp, 8))
			for _, phrase := range []string{
				"take the project's Session Handoff first",
				"use that block and claim no other; an accept that answers `consumed_by_hook` means the same",
				"list the open handoffs with `memory_handoff_list`",
				"when it names a `newer_unended_session`, read that session's newest observations with `memory_read_session_observations` before anything else, then accept the project's most recent open handoff with `memory_handoff_accept` and its `handoff_id`",
				"one message entry with id `handoff:<handoff id>:<n>`, `from` this project's name, `created_at` the handoff's creation date and no `type`",
				"name them in the plan and the report",
				"settle the handed-in messages as step 4 says",
			} {
				if !strings.Contains(step1, roundfixWhitespace(phrase)) {
					t.Errorf("step 1 omits %q", phrase)
				}
			}
			for _, phrase := range []string{
				"No message the mode handed in stays pending when the turn ends, also when the mode stops early: at step 1, when the user stops it, or when step 3 brings no answer.",
				"A message the plan records becomes a Spec, a Backlog Entry or a Finding and is then consumed.",
				"An informational message is consumed and the report names it.",
				"A message not yet handled becomes a Backlog Entry that records its sender, date and text, and is then consumed.",
				"Consume each one with `memory_message_pop` and its `message_id` only after its record is written",
				"A `handoff:` entry is not an ai-memory message",
			} {
				if !strings.Contains(step4, roundfixWhitespace(phrase)) {
					t.Errorf("step 4 omits %q", phrase)
				}
			}
			if !strings.Contains(step8, roundfixWhitespace("report the accepted handoff and what became of each of its items, each consumed message and where it was recorded")) {
				t.Error("step 8 omits the handoff and consumed-message report")
			}

			optionalSources := roundfixWhitespace(roundfixSection(reference, "Optional sources"))
			for _, phrase := range []string{
				"it brings the project's Session Handoff and its messages",
				"a session-aware client omits `workspace` and `project`, and a static client passes both",
				"`--status` and `--unattended` accept no handoff and consume no message",
			} {
				if !strings.Contains(optionalSources, roundfixWhitespace(phrase)) {
					t.Errorf("Optional sources omits %q", phrase)
				}
			}
			if !strings.Contains(optionalSources, roundfixWhitespace("List this project's messages with `memory_message_list`")) ||
				strings.Contains(strings.ToLower(optionalSources), "without consuming them") {
				t.Error("Optional sources still lists messages only to leave them pending")
			}
			status := roundfixWhitespace(roundfixSection(reference, "--status"))
			if !strings.Contains(status, roundfixWhitespace("accepts no handoff, consumes no ai-memory message and arms nothing")) {
				t.Error("--status section does not say it accepts no handoff")
			}

			skillData, err := fs.ReadFile(source, "roundfix/SKILL.md")
			if err != nil {
				t.Fatal(err)
			}
			modes := roundfixWhitespace(roundfixSection(string(skillData), "Modes"))
			for _, bullet := range []string{
				"- `--follow-up`: take the project's handoff, inventory pending work, plan, gather decisions, author and deliver the resulting work, and leave no handed-in message pending.",
				"- `--status`: run a short read-only status report that accepts no handoff and consumes no message.",
			} {
				if !strings.Contains(modes, roundfixWhitespace(bullet)) {
					t.Errorf("Modes omits %q", bullet)
				}
			}
		})
	}
}

func TestTheUnattendedModeAsksWithinTheBaselineLimit(t *testing.T) {
	t.Parallel()
	for name, source := range roundfixLayoutSources() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			core, err := os.ReadFile("../internal/baseline/assets/modules/core.json")
			if err != nil {
				t.Fatal(err)
			}
			section, err := fs.ReadFile(source, "roundfix/references/modes.md")
			if err != nil {
				t.Fatal(err)
			}
			for label, text := range map[string]string{"core module": string(core), "unattended section": roundfixSection(string(section), "--unattended")} {
				if !strings.Contains(strings.ToLower(text), "at most three questions per call") {
					t.Errorf("%s omits the Baseline limit", label)
				}
			}
			for _, sentence := range []string{
				"Ask at most three questions per call, never more than the tool accepts, as the stop-and-ask clause of `docs/agents/agent-instructions.md` allows, and put only independent questions in one call.",
				"Answers and authorizations about an item that becomes a Spec go into that Spec's `_authorization.md`, and its decisions also go into the PRD's `## Decisions`.",
				"An item that does not become a Spec now becomes a Backlog Entry or a Finding that records the information and decisions and never an authorization.",
			} {
				if !strings.Contains(string(section), sentence) {
					t.Errorf("unattended section omits %q", sentence)
				}
			}
		})
	}
}

func TestTheRoundfixSkillNamesNoPersonalAssistant(t *testing.T) {
	t.Parallel()
	for name, source := range roundfixLayoutSources() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := fs.WalkDir(source, "roundfix", func(path string, entry fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if entry.IsDir() {
					return nil
				}
				data, err := fs.ReadFile(source, path)
				if err != nil {
					return err
				}
				lower := strings.ToLower(string(data))
				if strings.Contains(lower, "jarvis") || strings.Contains(lower, "gt3") {
					t.Errorf("%s names a personal assistant or channel", path)
				}
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestTheRoundfixSkillReadsItsModeFromTheRequest(t *testing.T) {
	t.Parallel()
	for name, source := range roundfixLayoutSources() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			data, err := fs.ReadFile(source, "roundfix/SKILL.md")
			if err != nil {
				t.Fatal(err)
			}
			var frontmatter struct {
				Metadata map[string]any `yaml:"metadata"`
			}
			parts := strings.SplitN(string(data), "---", 3)
			if len(parts) != 3 {
				t.Fatal("skill entry has no frontmatter")
			}
			if err := yaml.Unmarshal([]byte(parts[1]), &frontmatter); err != nil {
				t.Fatal(err)
			}
			if got := frontmatter.Metadata["opencode/slash"]; got != "true" {
				t.Fatalf("metadata opencode/slash = %#v, want %q", got, "true")
			}
			text := string(data)
			modes := roundfixSection(text, "Modes")
			want := "A request names a Skill Mode with the first word after the skill's name: `follow-up`, `status` or `unattended`, with or without a leading `--`, in every harness: `/roundfix follow-up` or `/roundfix --follow-up` in Claude Code, `$roundfix status` or `$roundfix --status` in Codex, and `@roundfix unattended`, `@roundfix --unattended` or `/roundfix unattended` in OpenCode. A Skill Mode is not the CLI command `roundfix follow-up`, which the follow-up mode runs as one of its steps. Read the mode from the request text; a request whose first word after the skill's name is anything else, such as `/roundfix watch 123`, uses the skill as the rest of this file describes."
			if !strings.HasPrefix(modes, want) {
				t.Fatalf("Modes section does not start with required paragraph")
			}
			for _, mode := range []string{"--follow-up", "--status", "--unattended"} {
				if !strings.Contains(modes, mode) {
					t.Errorf("Modes section omits %s", mode)
				}
			}
			if !strings.Contains(modes, "references/modes.md") {
				t.Error("Modes section does not link its reference")
			}
			if regexp.MustCompile(`\$(?:[0-9]+|ARGUMENTS)`).MatchString(parts[2]) {
				t.Error("skill body contains an argument placeholder")
			}
			reference, err := fs.ReadFile(source, "roundfix/references/modes.md")
			if err != nil {
				t.Fatal(err)
			}
			referenceText := string(reference)
			for _, heading := range []string{"How a mode is chosen", "Optional sources", "--status"} {
				if roundfixSection(referenceText, heading) == "" {
					t.Errorf("modes reference has no non-empty %q section", heading)
				}
			}
			chosen := strings.Join(strings.Fields(roundfixSection(referenceText, "How a mode is chosen")), " ")
			for _, phrase := range []string{
				"with or without a leading `--`",
				"`/roundfix follow-up`, `/roundfix --follow-up`, `$roundfix status` and `@roundfix --unattended` each select a mode",
				"This is the Skill Mode, not the CLI command `roundfix follow-up`",
			} {
				if !strings.Contains(chosen, phrase) {
					t.Errorf("How a mode is chosen omits %q", phrase)
				}
			}
		})
	}
}
