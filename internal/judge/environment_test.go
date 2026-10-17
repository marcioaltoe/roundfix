package judge

import (
	"slices"
	"strings"
	"testing"
)

func TestWithoutKeyVariablesRemovesEveryJudgeVariable(t *testing.T) {
	t.Parallel()
	questions, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	variables := append(questions.KeyVariables(), "TYPESAFE_API_KEY")
	if len(variables) == 0 {
		t.Fatal("embedded questions declare no Judge Key Variables")
	}
	environ := []string{"PATH=/fixture/bin", "HOME=/fixture/home", "GOCACHE=/fixture/cache", "CONTROL=first=part"}
	want := slices.Clone(environ)
	for _, variable := range variables {
		environ = append(environ, variable+"=fixture", variable+"=", variable+"=duplicate", variable)
		near := variable + "_CONTROL=kept"
		environ = append(environ, near)
		want = append(want, near)
	}
	environ = append(environ, "CONTROL=last", "EMPTY=", "bare-entry")
	want = append(want, "CONTROL=last", "EMPTY=", "bare-entry")
	before := slices.Clone(environ)
	got := WithoutKeyVariables(environ)
	if !slices.Equal(got, want) {
		t.Fatal("unrelated entries changed or were reordered")
	}
	for _, entry := range got {
		name, _, _ := strings.Cut(entry, "=")
		if slices.Contains(variables, name) {
			t.Errorf("judge variable %s remains", name)
		}
	}
	if !slices.Equal(environ, before) {
		t.Fatal("input environment mutated")
	}
	got[0] = "changed"
	if !slices.Equal(environ, before) {
		t.Fatal("result aliases input environment")
	}
}
