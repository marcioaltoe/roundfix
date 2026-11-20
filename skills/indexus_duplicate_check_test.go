package skills

import (
	"regexp"
	"testing"
)

func TestWriteIdeaRunsTheDuplicateCheckThroughTheIndexus(t *testing.T) {
	t.Parallel()
	for _, phrase := range []struct{ name, text string }{
		{"duplicate_check_with_the_tool", "Run the advisory duplicate check first with the Indexus tool `pesquisar_acervos`"},
		{"each_returned_page_opens", "verify that each returned page opens"},
		{"no_label_or_rank_decides", "a relevance label or a rank alone never decides"},
		{"absent_tool_substitutes_nothing", "say so in the message and substitute no other search tool"},
	} {
		t.Run(phrase.name, func(t *testing.T) {
			assertIdeaStageSkillPhrases(t, "write-idea", "SKILL.md", phrase.text)
		})
	}
}

func TestNoOwnedSkillNamesTheRetiredSearchCommand(t *testing.T) {
	t.Parallel()
	files, err := Files()
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("owned-skill scan read no files")
	}
	retired := regexp.MustCompile(`(?i)\bqmd\b`)
	for _, file := range files {
		if retired.Match(file.Data) {
			t.Errorf("%s names the retired search command", file.Path)
		}
	}
}
