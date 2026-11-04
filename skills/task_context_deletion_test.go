package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTheWriteTasksGuidesListARemovedFileOnlyUnderDeletes(t *testing.T) {
	t.Parallel()

	repoRoot := filepath.Clean(filepath.Join(".."))
	declaredPathRule := normalizeAcceptanceText("Every path a Task edits is declared under `interface:` or `creates:`; every path it removes is declared only under `deletes:`, never `interface:` and never `instruction:`. A Task that reads a file before removing it still lists it once, under `deletes:`.")
	templateRule := normalizeAcceptanceText("A file the Task removes is listed only as deletes:, never as interface: or instruction:.")
	for _, root := range []string{
		filepath.Join(repoRoot, ".agents", "skills", "write-tasks"),
		filepath.Join(repoRoot, "skills", "write-tasks"),
	} {
		t.Run(filepath.Base(filepath.Dir(root))+"/"+filepath.Base(root), func(t *testing.T) {
			skill, err := os.ReadFile(filepath.Join(root, "SKILL.md"))
			if err != nil {
				t.Fatalf("read write-tasks skill %s: %v", root, err)
			}
			normalizedSkill := normalizeAcceptanceText(string(skill))
			if !strings.Contains(normalizedSkill, declaredPathRule) {
				t.Errorf("%s is missing the removed-file declaration rule", root)
			}
			if strings.Contains(string(skill), "declared under `deletes:`, never `instruction:`") {
				t.Errorf("%s retains the obsolete removed-file declaration rule", root)
			}

			template, err := os.ReadFile(filepath.Join(root, "references", "task-template.md"))
			if err != nil {
				t.Fatalf("read Task template %s: %v", root, err)
			}
			if !strings.Contains(normalizeAcceptanceText(string(template)), templateRule) {
				t.Errorf("%s is missing the Task template removed-file declaration rule", root)
			}
		})
	}
}
