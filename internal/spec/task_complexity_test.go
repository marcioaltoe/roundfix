package spec

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestTaskCarriesComplexity(t *testing.T) {
	for _, raw := range []string{"low", "  medium  ", "", "custom"} {
		t.Run(raw, func(t *testing.T) {
			root := t.TempDir()
			content := strings.Replace(taskFixture("task_01", "Fixture", "pending", "backend", defaultVerificationSection), "complexity: low", "complexity: '"+raw+"'", 1)
			writeSpecDir(t, root, "demo", map[string]string{
				"_prd.md":    prdFixture("active"),
				"_tasks.md":  manifestFixture("spec-tasks/v1", "    - id: task_01\n      file: task_01.md\n      needs: []\n"),
				"task_01.md": content,
			})
			graph, err := Load(root, "demo")
			if err != nil {
				t.Fatal(err)
			}
			task := graph.Tasks[0]
			if task.Complexity != strings.TrimSpace(raw) {
				t.Fatalf("Complexity=%q", task.Complexity)
			}
			before := task
			writeFile(t, filepath.Join(root, "demo/task_01.md"), strings.Replace(content, "complexity: '"+raw+"'", "complexity: ' high '", 1))
			if err := ReloadTask(root, &task); err != nil {
				t.Fatal(err)
			}
			if task.Complexity != "high" {
				t.Fatalf("reloaded Complexity=%q", task.Complexity)
			}
			before.Complexity = "high"
			if !reflect.DeepEqual(task, before) {
				t.Fatalf("reload changed another field: got %+v want %+v", task, before)
			}
			writeFile(t, filepath.Join(root, "demo/task_01.md"), strings.Replace(content, "complexity: '"+raw+"'\n", "", 1))
			if err := ReloadTask(root, &task); err != nil {
				t.Fatal(err)
			}
			if task.Complexity != "" {
				t.Fatalf("missing Complexity=%q", task.Complexity)
			}
			graph, err = Load(root, "demo")
			if err != nil {
				t.Fatal(err)
			}
			if graph.Tasks[0].Complexity != "" {
				t.Fatalf("loaded missing Complexity=%q", graph.Tasks[0].Complexity)
			}
		})
	}
}
