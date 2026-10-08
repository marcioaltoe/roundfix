package spec

import (
	"reflect"
	"strings"
	"testing"
)

func TestTaskGraphReadsRequiredSpecs(t *testing.T) {
	t.Parallel()
	for _, declaration := range []string{"", "requires: []\n", "requires: [first, second]\n"} {
		t.Run(declaration, func(t *testing.T) {
			root := t.TempDir()
			files := diamondSpecFiles()
			files["_tasks.md"] = strings.Replace(files["_tasks.md"], "schema:", declaration+"schema:", 1)
			writeSpecDir(t, root, "demo", files)
			graph, err := Load(root, "demo")
			if err != nil {
				t.Fatal(err)
			}
			var want []string
			if strings.Contains(declaration, "first") {
				want = []string{"first", "second"}
			}
			if !reflect.DeepEqual(graph.Requires, want) {
				t.Fatalf("requires = %v, want %v", graph.Requires, want)
			}
		})
	}
}

func TestTaskGraphRefusesAMalformedRequiresList(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"null", "first", "{first: true}", "[42]", "[true]", "[null]", "[{}]", "['']", "['  ']", "[first, first]", "[first, ' first ']", "[demo]"} {
		t.Run(value, func(t *testing.T) {
			root := t.TempDir()
			files := diamondSpecFiles()
			files["_tasks.md"] = strings.Replace(files["_tasks.md"], "schema:", "requires: "+value+"\nschema:", 1)
			writeSpecDir(t, root, "demo", files)
			_, err := Load(root, "demo")
			if err == nil || !strings.Contains(err.Error(), "_tasks.md") {
				t.Fatalf("Load error = %v", err)
			}
		})
	}
}

func TestTaskGraphRefusesARequiresEntryThatIsAPath(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"['../../outside']", "['nested/spec']", "['..']", "['.']", "['back\\\\slash']"} {
		t.Run(value, func(t *testing.T) {
			root := t.TempDir()
			files := diamondSpecFiles()
			files["_tasks.md"] = strings.Replace(files["_tasks.md"], "schema:", "requires: "+value+"\nschema:", 1)
			writeSpecDir(t, root, "demo", files)
			_, err := Load(root, "demo")
			if err == nil || !strings.Contains(err.Error(), "must name a Spec directory") {
				t.Fatalf("Load error = %v", err)
			}
		})
	}
}
