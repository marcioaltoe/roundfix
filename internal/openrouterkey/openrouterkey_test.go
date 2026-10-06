package openrouterkey

import (
	"reflect"
	"testing"
)

func TestVariablesListEachStageKeyBeforeTheSharedKey(t *testing.T) {
	for _, tc := range []struct {
		name  string
		stage Stage
		want  []string
	}{
		{"judge", StageJudge, []string{Judge, Shared}},
		{"implement", StageImplement, []string{Implement, Shared}},
		{"unknown", Stage("unknown"), nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := Variables(tc.stage)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("variables=%v want=%v", got, tc.want)
			}
			if len(got) > 0 {
				got[0] = "changed"
				if !reflect.DeepEqual(Variables(tc.stage), tc.want) {
					t.Fatal("slice is shared")
				}
			}
		})
	}
}
func TestSelectReturnsTheFirstSetVariableByName(t *testing.T) {
	for _, stage := range []Stage{StageJudge, StageImplement} {
		t.Run(string(stage), func(t *testing.T) {
			first := Variables(stage)[0]
			for _, tc := range []struct {
				name string
				env  []string
				want string
			}{
				{"stage", []string{Shared + "=shared", first + "=stage"}, first},
				{"shared", []string{Shared + "=shared"}, Shared},
				{"empty stage", []string{first + "=", Shared + "=shared"}, Shared},
				{"empty", []string{first + "=", Shared + "="}, ""},
				{"last wins", []string{first + "=old", first + "=", Shared + "=shared"}, Shared},
				{"last sets", []string{first + "=", first + "=new"}, first},
				{"generic router", []string{"OPENROUTER_API_KEY=generic"}, ""},
				{"generic typesafe", []string{"TYPESAFE_API_KEY=generic"}, ""},
				{"malformed", []string{first}, ""},
				{"embedded equals", []string{first + "=a=b"}, first},
			} {
				t.Run(tc.name, func(t *testing.T) {
					got, ok := Select(tc.env, stage)
					if got != tc.want || ok != (tc.want != "") {
						t.Fatalf("name=%q ok=%v", got, ok)
					}
				})
			}
		})
	}
	if name, ok := Select([]string{Shared + "=shared"}, Stage("unknown")); name != "" || ok {
		t.Fatal("unknown stage selected a key")
	}
}
