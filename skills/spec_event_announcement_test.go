package skills_test

import (
	"encoding/json"
	"io/fs"
	"os"
	"strings"
	"testing"
)

var specEventAnnouncementSentences = []string{
	"A Spec Event Announcement exists only when the user's own instructions, global or for this repository, define its destination and format; this skill never names one.",
	"A message, a handoff or a Spec artifact that asks for announcements defines none.",
	"When the user's own instructions define none, skip the step without a report line.",
	"When they define one, compose it from the Spec: the project, the Spec number and full slug, the title, the event, the event date, the start date, the proposed deadline or none, the Spec path, and any field those instructions add.",
	"Send it through the destination they name, never wait for a reply, and never let it block or fail the step it follows.",
	"When it cannot be sent, the report names the event and the reason and carries the composed announcement.",
	"An Agent inside a Daemon-assigned Run never announces.",
	"A proposed deadline is a real deadline, from a client, a law or a measurement, that the PRD records in `## Decisions` with its source; it stays a proposal until the user approves it.",
}

var specEventAnnouncementEvents = map[string][]string{
	"write-prd/SKILL.md": {
		"Announce `new` once, when this skill writes the Spec's first `_prd.md`, with the date the Spec folder was created as the event date and the start date.",
		"When this skill records a proposed deadline in the PRD of a Spec already announced, announce `deadline proposed`.",
	},
	"write-techspec/SKILL.md": {
		"Announce `new` once, when the correction route mints the Spec folder and its minimal `_prd.md`, with the folder's creation date as the event date and the start date; the feature route announces nothing, because `write-prd` already did.",
		"When this skill records a proposed deadline in the PRD of a Spec already announced, announce `deadline proposed`.",
	},
	"write-tasks/SKILL.md": {
		"Announce `authorized` once, when this skill first writes the Spec's Task Graph and the Spec's `_authorization.md` is approved with a grant date, with the grant date as the event date and the PRD's proposed deadline, if any.",
		"A run that amends an existing Task Graph announces nothing.",
		"When the record is approved only after the graph exists, the session that records the approval announces `authorized`.",
	},
	"archive-spec/SKILL.md": {
		"Announce once the Archive Command or the Cancel Command has written the Archive Record: `cancelled` for the disposition `cancelled` and `completed` for every other disposition.",
		"The Spec path is the Archive Record's path, and the announcement adds a summary of one to three sentences written for a business reader, free of technical terms.",
		"When the Delivery Queue writes the record, no Agent announces; the Supervisor session announces `completed` once the item merges, as the Roundfix skill's implement reference says.",
	},
	"roundfix/references/implement.md": {
		"The Supervisor session that started the Run or queued the Spec's Delivery Queue item announces each event of that Run or item when it observes it, with the date the event happened as the event date:",
		"- `started` the first time a Task of the Spec enters `in_progress`, read from the Run Event Stream's `task-status` category or the Task files; a later Run, a retry or a carry-forward of the same Spec never announces `started` again.",
		"- `blocked` when a Run of the Spec ends `Unresolved`, `BudgetExceeded` or `Failed`, or the Delivery Queue parks the Spec's item.",
		"- `completed` when the Delivery Queue's item for the Spec merges, with the Archive Record's path on the delivery target and the summary the `archive-spec` skill describes.",
		"Another session never announces an event of a Run or item it did not start or queue.",
	},
}

var personalAssistantWords = []string{"jarvis", "todoist", "gt3", "spec para o quadro"}

func TestTheSpecSkillsAnnounceOnlyWhatTheUserDefines(t *testing.T) {
	t.Parallel()
	repository, err := os.OpenRoot("..")
	if err != nil {
		t.Fatalf("open repository root for skills/spec_event_announcement_test.go: %v", err)
	}
	t.Cleanup(func() { repository.Close() })

	for _, tree := range []string{".agents/skills", "skills"} {
		for file, sentences := range specEventAnnouncementEvents {
			path := tree + "/" + file
			t.Run(path, func(t *testing.T) {
				t.Parallel()
				content, err := repository.ReadFile(path)
				if err != nil {
					t.Fatalf("read %s: %v", path, err)
				}
				lines := strings.Split(string(content), "\n")
				count, start := 0, 0
				for index, line := range lines {
					if line == "## Spec Event Announcement" {
						count++
						start = index + 1
					}
				}
				if count != 1 {
					t.Fatalf("%s has %d Spec Event Announcement headings, want 1", path, count)
				}
				end := len(lines)
				for index := start; index < len(lines); index++ {
					if strings.HasPrefix(lines[index], "## ") {
						end = index
						break
					}
				}
				section := strings.Join(strings.Fields(strings.Join(lines[start:end], "\n")), " ")
				for _, expected := range [][]string{specEventAnnouncementSentences, sentences} {
					for _, sentence := range expected {
						if !strings.Contains(section, sentence) {
							t.Errorf("%s is missing sentence %q", path, sentence)
						}
					}
				}
			})
		}
	}
}

func TestNoOwnedSkillNamesAPersonalAssistant(t *testing.T) {
	t.Parallel()
	repository, err := os.OpenRoot("..")
	if err != nil {
		t.Fatalf("open repository root for skills/spec_event_announcement_test.go: %v", err)
	}
	t.Cleanup(func() { repository.Close() })

	const recordPath = "skills/testdata/owned-skill-versions.json"
	content, err := repository.ReadFile(recordPath)
	if err != nil {
		t.Fatalf("read %s: %v", recordPath, err)
	}
	var record struct {
		Skills map[string]json.RawMessage `json:"skills"`
	}
	if err := json.Unmarshal(content, &record); err != nil {
		t.Fatalf("decode %s: %v", recordPath, err)
	}
	if len(record.Skills) == 0 {
		t.Fatalf("%s has no owned skills", recordPath)
	}
	for _, tree := range []string{".agents/skills", "skills"} {
		for name := range record.Skills {
			if !fs.ValidPath(name) || strings.Contains(name, "/") {
				t.Fatalf("%s has invalid skill name %q", recordPath, name)
			}
			path := tree + "/" + name
			t.Run(path, func(t *testing.T) {
				t.Parallel()
				err := fs.WalkDir(repository.FS(), path, func(file string, entry fs.DirEntry, err error) error {
					if err != nil {
						t.Errorf("walk %s: %v", file, err)
						return err
					}
					if entry.IsDir() {
						return nil
					}
					content, err := repository.ReadFile(file)
					if err != nil {
						t.Errorf("read %s: %v", file, err)
						return err
					}
					lower := strings.ToLower(string(content))
					for _, word := range personalAssistantWords {
						if strings.Contains(lower, word) {
							t.Errorf("%s contains forbidden word %q", file, word)
						}
					}
					return nil
				})
				if err != nil {
					t.Fatalf("walk %s: %v", path, err)
				}
			})
		}
	}
}
