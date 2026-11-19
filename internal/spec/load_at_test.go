package spec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadAtReadsThroughTheGivenReader(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"same_graph_as_load", "an_absent_prd_is_not_found", "an_absent_task_file_is_named"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root := defaultSpecsRoot(t.TempDir())
			files := diamondSpecFiles()
			writeSpecDir(t, root, "demo", files)
			paths := make(map[string][]byte, len(files))
			for file, content := range files {
				paths[filepath.Join(root, "demo", file)] = []byte(content)
			}
			switch name {
			case "an_absent_prd_is_not_found":
				delete(paths, filepath.Join(root, "demo", "_prd.md"))
			case "an_absent_task_file_is_named":
				delete(paths, filepath.Join(root, "demo", "task_01.md"))
			}
			read := func(path string) ([]byte, error) {
				content, ok := paths[path]
				if !ok {
					return nil, fmt.Errorf("read map path %q: %w", path, os.ErrNotExist)
				}
				return content, nil
			}
			var want *Graph
			if name == "same_graph_as_load" {
				var err error
				want, err = Load(root, "demo")
				if err != nil {
					t.Fatal(err)
				}
				if err := os.RemoveAll(filepath.Join(root, "demo")); err != nil {
					t.Fatal(err)
				}
			}
			got, err := LoadAt(root, "demo", read)
			switch name {
			case "same_graph_as_load":
				if err != nil {
					t.Fatal(err)
				}
				if len(got.Tasks) != len(want.Tasks) {
					t.Fatalf("LoadAt tasks=%d, Load tasks=%d", len(got.Tasks), len(want.Tasks))
				}
				for i, task := range got.Tasks {
					other := want.Tasks[i]
					if task.ID != other.ID || task.File != other.File || task.Status != other.Status || task.Type != other.Type || !reflect.DeepEqual(task.Needs, other.Needs) {
						t.Errorf("LoadAt task=%+v, Load task=%+v", task, other)
					}
				}
			case "an_absent_prd_is_not_found":
				var missing SpecNotFoundError
				if !errors.As(err, &missing) || missing.Slug != "demo" || got != nil {
					t.Fatalf("LoadAt graph=%+v error=%v, want SpecNotFoundError for demo", got, err)
				}
			case "an_absent_task_file_is_named":
				var missing MissingTaskFileError
				if !errors.As(err, &missing) || missing.TaskID != "task_01" || missing.Path != filepath.Join(root, "demo", "task_01.md") || got != nil {
					t.Fatalf("LoadAt graph=%+v error=%v, want MissingTaskFileError for task_01", got, err)
				}
			}
		})
	}
}
