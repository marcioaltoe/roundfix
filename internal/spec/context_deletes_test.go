package spec

import (
	"strings"
	"testing"
)

func TestContextAcceptsADeletesEntry(t *testing.T) {
	refs, err := parseTaskContextRefs([]byte("## Context\n\n- deletes: `old/file.txt`\n"))
	if err != nil || len(refs) != 1 || refs[0].Kind != ContextKindDeletes || refs[0].Path != "old/file.txt" {
		t.Fatalf("Context = %#v, %v", refs, err)
	}
	for _, path := range []string{"../old.txt", "/old.txt", "old/../file.txt", ""} {
		t.Run(path, func(t *testing.T) {
			if _, err := parseTaskContextRefs([]byte("## Context\n- deletes: " + path)); err == nil {
				t.Fatal("accepted invalid path")
			}
		})
	}
}

func TestContextRefusesAPathUnderTwoKinds(t *testing.T) {
	for _, pair := range [][2]string{{"interface", "deletes"}, {"deletes", "interface"}, {"creates", "deletes"}, {"deletes", "creates"}, {"deletes", "deletes"}} {
		t.Run(pair[0]+"/"+pair[1], func(t *testing.T) {
			_, err := parseTaskContextRefs([]byte("## Context\n- " + pair[0] + ": old/file.txt\n- " + pair[1] + ": old/file.txt\n"))
			if err == nil || !strings.Contains(err.Error(), pair[0]) || !strings.Contains(err.Error(), pair[1]) {
				t.Fatalf("error = %v, want both kinds", err)
			}
		})
	}
}

func TestUndeclaredTaskPathsCountsADeletesPath(t *testing.T) {
	task := Task{Context: []TaskContextRef{{Kind: ContextKindDeletes, Path: "old/file.txt"}}}
	paths := UndeclaredTaskPaths(task, "task_01.md", []string{"old/file.txt", "other.txt"}, nil)
	if len(paths) != 1 || paths[0] != "other.txt" {
		t.Fatalf("undeclared = %v", paths)
	}
}

func TestWaveCollisionCountsADeletesPath(t *testing.T) {
	root := t.TempDir()
	writeCollisionFile(t, root, "old/file.txt", "old")
	graph := &Graph{Tasks: []Task{
		{ID: "task_01", Context: []TaskContextRef{{Kind: ContextKindInterface, Path: "old/file.txt"}}},
		{ID: "task_02", Context: []TaskContextRef{{Kind: ContextKindDeletes, Path: "old/file.txt"}}},
	}}
	collisions, err := Collisions(root, graph)
	if err != nil || len(collisions) != 1 || collisions[0].Paths["old/file.txt"] != TouchFromContext {
		t.Fatalf("collisions = %#v, %v", collisions, err)
	}
	graph.Tasks[1].Needs = []string{"task_01"}
	collisions, err = Collisions(root, graph)
	if err != nil || len(collisions) != 0 {
		t.Fatalf("ordered collisions = %#v, %v", collisions, err)
	}
}
