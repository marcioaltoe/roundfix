package skillcoverage_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/skillcoverage"
)

func TestParseMapRefusesEveryMalformedEntry(t *testing.T) {
	mapJSON := func(entries string) []byte {
		return []byte(`{"schemaVersion":"` + skillcoverage.MapSchema + `","surfaces":[` + entries + `]}`)
	}
	validEntry := `{"id":"command run","skills":[".agents/skills/roundfix/references/run.md"],"sources":["internal/cli/","docs/*.md","internal/cli/run.go"]}`
	cases := []struct{ name, entries, field string }{
		{"empty id", `{"id":"","skills":["skill.md"]}`, "id"},
		{"blank id", `{"id":"  ","uncovered":"reason"}`, "id"},
		{"duplicate id", validEntry + `,` + validEntry, "command run"},
		{"unsorted ids", `{"id":"z","uncovered":"reason"},{"id":"a","uncovered":"reason"}`, "a"},
		{"unknown entry field", `{"id":"command run","skills":["skill.md"],"extra":true}`, "command run"},
		{"neither coverage", `{"id":"command run"}`, "command run"},
		{"both coverage", `{"id":"command run","skills":["skill.md"],"uncovered":"reason"}`, "command run"},
		{"both fields with empty skills", `{"id":"command run","skills":[],"uncovered":"reason"}`, "command run"},
		{"empty skills", `{"id":"command run","skills":[]}`, "command run"},
		{"null skills", `{"id":"command run","skills":null}`, "command run"},
		{"blank uncovered", `{"id":"command run","uncovered":"  "}`, "command run"},
		{"wrong skills type", `{"id":"command run","skills":"skill.md"}`, "command run"},
		{"blank review", `{"id":"command run","uncovered":"reason","review":"   "}`, "review"},
		{"undated review", `{"id":"command run","uncovered":"reason","review":"looked at it"}`, "review"},
		{"review without reason", `{"id":"command run","uncovered":"reason","review":"2026-10-08 — "}`, "review"},
		{"review with impossible date", `{"id":"command run","uncovered":"reason","review":"2026-13-40 — renamed only"}`, "review"},
		{"review with hyphen separator", `{"id":"command run","uncovered":"reason","review":"2026-10-08 - renamed only"}`, "review"},
	}
	for _, field := range []string{"skills", "sources"} {
		for _, value := range []string{"", "/absolute.md", "../escape.md", "a/../b.md", "a/./b.md", "a//b.md", `a\b.md`, ".", "C:/absolute.md"} {
			encoded, err := json.Marshal(value)
			if err != nil {
				t.Fatal(err)
			}
			entry := `{"id":"command run","uncovered":"reason","sources":[` + string(encoded) + `]}`
			if field == "skills" {
				entry = `{"id":"command run","skills":[` + string(encoded) + `]}`
			}
			cases = append(cases, struct{ name, entries, field string }{field + " invalid " + value, entry, "command run"})
		}
	}
	cases = append(cases, struct{ name, entries, field string }{"malformed source pattern", `{"id":"command run","uncovered":"reason","sources":["internal/["]}`, "command run"})
	cases = append(cases, struct{ name, entries, field string }{"skill directory", `{"id":"command run","skills":["skills/"]}`, "command run"})
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := skillcoverage.ParseMap(mapJSON(tc.entries))
			if err == nil || !strings.Contains(err.Error(), tc.field) {
				t.Fatalf("error = %v, want refusal naming %q", err, tc.field)
			}
		})
	}
	for _, tc := range []struct{ name, data, field string }{
		{"wrong map schema", `{"schemaVersion":"wrong","surfaces":[]}`, "schemaVersion"},
		{"unknown map field", `{"schemaVersion":"` + skillcoverage.MapSchema + `","extra":1}`, "extra"},
		{"missing map surfaces", `{"schemaVersion":"` + skillcoverage.MapSchema + `"}`, "surfaces"},
		{"null map surfaces", `{"schemaVersion":"` + skillcoverage.MapSchema + `","surfaces":null}`, "surfaces"},
		{"trailing map JSON", string(mapJSON(validEntry)) + ` {}`, "trailing JSON"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := skillcoverage.ParseMap([]byte(tc.data))
			if err == nil || !strings.Contains(err.Error(), tc.field) {
				t.Fatalf("error = %v, want %q", err, tc.field)
			}
		})
	}
	m, err := skillcoverage.ParseMap(mapJSON(validEntry + `,{"id":"exit codes","uncovered":"no skill covers this"}`))
	if err != nil || len(m.Surfaces) != 2 {
		t.Fatalf("valid map = %+v, error = %v", m, err)
	}

	recordJSON := func(entries string) []byte {
		return []byte(`{"schemaVersion":"` + skillcoverage.RecordSchema + `","surfaces":{` + entries + `}}`)
	}
	validFP := skillcoverage.Fingerprint([]byte("content"))
	for _, tc := range []struct{ name, entries, field string }{
		{"empty record id", `"":"` + validFP + `"`, "id"},
		{"blank record id", `"  ":"` + validFP + `"`, "id"},
		{"duplicate record id", `"run":"` + validFP + `","run":"` + validFP + `"`, "run"},
		{"short fingerprint", `"run":"sha256:abc"`, "run"},
		{"wrong prefix", `"run":"md5:` + strings.Repeat("a", 64) + `"`, "run"},
		{"uppercase hex", `"run":"sha256:` + strings.Repeat("A", 64) + `"`, "run"},
		{"non hex", `"run":"sha256:` + strings.Repeat("g", 64) + `"`, "run"},
		{"long fingerprint", `"run":"` + validFP + `0"`, "run"},
		{"wrong fingerprint type", `"run":1`, "run"},
		{"unknown fingerprint field", `"run":{"fingerprint":"` + validFP + `"}`, "run"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := skillcoverage.ParseRecord(recordJSON(tc.entries))
			if err == nil || !strings.Contains(err.Error(), tc.field) {
				t.Fatalf("error = %v, want %q", err, tc.field)
			}
		})
	}
	for _, tc := range []struct{ name, data, field string }{
		{"wrong record schema", `{"schemaVersion":"wrong","surfaces":{}}`, "schemaVersion"},
		{"unknown record field", `{"schemaVersion":"` + skillcoverage.RecordSchema + `","extra":1}`, "extra"},
		{"missing record surfaces", `{"schemaVersion":"` + skillcoverage.RecordSchema + `"}`, "surfaces"},
		{"null record surfaces", `{"schemaVersion":"` + skillcoverage.RecordSchema + `","surfaces":null}`, "surfaces"},
		{"trailing record JSON", string(recordJSON(`"run":"`+validFP+`"`)) + ` {}`, "trailing JSON"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := skillcoverage.ParseRecord([]byte(tc.data))
			if err == nil || !strings.Contains(err.Error(), tc.field) {
				t.Fatalf("error = %v, want %q", err, tc.field)
			}
		})
	}
	if _, err := skillcoverage.ParseRecord(recordJSON(`"run":"` + validFP + `"`)); err != nil {
		t.Fatalf("valid record: %v", err)
	}
}

func TestRecordRoundTripsByteForByte(t *testing.T) {
	fp := skillcoverage.Fingerprint([]byte("abc"))
	if fp != "sha256:ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
		t.Fatalf("fingerprint = %q", fp)
	}
	for _, surfaces := range []map[string]string{{"z": fp, "a": fp}, {}} {
		t.Run("canonical record", func(t *testing.T) {
			original := skillcoverage.Record{SchemaVersion: skillcoverage.RecordSchema, Surfaces: surfaces}
			encoded := skillcoverage.EncodeRecord(original)
			parsed, err := skillcoverage.ParseRecord(encoded)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(encoded, skillcoverage.EncodeRecord(parsed)) {
				t.Fatal("re-encoded bytes differ")
			}
			if !bytes.HasSuffix(encoded, []byte("\n")) || bytes.HasSuffix(encoded, []byte("\n\n")) {
				t.Fatalf("expected one final newline: %q", encoded)
			}
			if len(surfaces) > 0 {
				expected := "{\n  \"schemaVersion\": \"" + skillcoverage.RecordSchema + "\",\n  \"surfaces\": {\n    \"a\": \"" + fp + "\",\n    \"z\": \"" + fp + "\"\n  }\n}\n"
				if string(encoded) != expected {
					t.Fatalf("encoded = %s, want %s", encoded, expected)
				}
			}
		})
	}
}

func TestSurfacesForPathFollowTheSourceGrammar(t *testing.T) {
	m := skillcoverage.Map{Surfaces: []skillcoverage.Surface{
		{ID: "directory", Sources: []string{"internal/cli/", "internal/cli/run.go"}},
		{ID: "exact", Sources: []string{"internal/cli/run.go"}},
		{ID: "wildcard", Sources: []string{"docs/*.md"}},
	}}
	for _, tc := range []struct {
		path string
		ids  []string
	}{
		{"internal/cli/run.go", []string{"directory", "exact"}},
		{"internal/cli/deep/file.go", []string{"directory"}},
		{"internal/cli", nil}, {"internal/cliff/file.go", nil},
		{"internal/cli/run.go.bak", []string{"directory"}},
		{"docs/guide.md", []string{"wildcard"}}, {"docs/nested/guide.md", nil},
		{"docs/guide.txt", nil}, {"other/docs/guide.md", nil},
	} {
		t.Run(tc.path, func(t *testing.T) {
			var ids []string
			for _, s := range m.SurfacesForPath(tc.path) {
				ids = append(ids, s.ID)
			}
			if !reflect.DeepEqual(ids, tc.ids) {
				t.Fatalf("ids = %v, want %v", ids, tc.ids)
			}
		})
	}
}

func TestCompareJudgesEveryChangedSurface(t *testing.T) {
	oldFP, newFP := skillcoverage.Fingerprint([]byte("old")), skillcoverage.Fingerprint([]byte("new"))
	oldReview, newReview := "2026-10-07 — no skill text needed", "2026-10-08 — still no skill text needed"
	for _, tc := range []struct {
		name, kind             string
		baseEntry, targetEntry *skillcoverage.Surface
		paths                  map[string]bool
		outcome                string
		skills                 []string
	}{
		{"added described", "added", nil, &skillcoverage.Surface{Skills: []string{"new.md"}}, map[string]bool{"new.md": true}, "described", []string{"new.md"}},
		{"changed uses target skill", "changed", &skillcoverage.Surface{Skills: []string{"old.md"}}, &skillcoverage.Surface{Skills: []string{"new.md"}}, map[string]bool{"new.md": true}, "described", []string{"new.md"}},
		{"old skill cannot describe target", "changed", &skillcoverage.Surface{Skills: []string{"old.md"}}, &skillcoverage.Surface{Skills: []string{"new.md"}}, map[string]bool{"old.md": true}, "lagging", []string{"new.md"}},
		{"removed uses base skill", "removed", &skillcoverage.Surface{Skills: []string{"old.md"}}, &skillcoverage.Surface{Skills: []string{"new.md"}}, map[string]bool{"old.md": true}, "described", []string{"old.md"}},
		{"uncovered wins", "changed", nil, &skillcoverage.Surface{Uncovered: "reason", Review: newReview}, nil, "uncovered", nil},
		{"removed uncovered", "removed", &skillcoverage.Surface{Uncovered: "reason"}, nil, nil, "uncovered", nil},
		{"added review", "added", nil, &skillcoverage.Surface{Skills: []string{"skill.md"}, Review: newReview}, nil, "reviewed", []string{"skill.md"}},
		{"new review", "changed", &skillcoverage.Surface{Skills: []string{"skill.md"}}, &skillcoverage.Surface{Skills: []string{"skill.md"}, Review: newReview}, nil, "reviewed", []string{"skill.md"}},
		{"changed review", "changed", &skillcoverage.Surface{Review: oldReview}, &skillcoverage.Surface{Skills: []string{"skill.md"}, Review: newReview}, nil, "reviewed", []string{"skill.md"}},
		{"unchanged review", "changed", &skillcoverage.Surface{Review: newReview}, &skillcoverage.Surface{Skills: []string{"skill.md"}, Review: newReview}, nil, "lagging", []string{"skill.md"}},
		{"removed review cannot count", "removed", &skillcoverage.Surface{Skills: []string{"skill.md"}, Review: oldReview}, &skillcoverage.Surface{Review: newReview}, nil, "lagging", []string{"skill.md"}},
		{"described precedes review", "changed", nil, &skillcoverage.Surface{Skills: []string{"skill.md"}, Review: newReview}, map[string]bool{"skill.md": true}, "described", []string{"skill.md"}},
		{"false path not changed", "changed", nil, &skillcoverage.Surface{Skills: []string{"skill.md"}}, map[string]bool{"skill.md": false}, "lagging", []string{"skill.md"}},
		{"no entry", "changed", nil, nil, nil, "lagging", nil},
		{"empty review", "changed", nil, &skillcoverage.Surface{Skills: []string{"skill.md"}}, nil, "lagging", []string{"skill.md"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			base := skillcoverage.Snapshot{Record: &skillcoverage.Record{Surfaces: map[string]string{"surface": oldFP, "unchanged": oldFP}}}
			target := skillcoverage.Snapshot{Record: &skillcoverage.Record{Surfaces: map[string]string{"surface": newFP, "unchanged": oldFP}}}
			if tc.kind == "added" {
				delete(base.Record.Surfaces, "surface")
			}
			if tc.kind == "removed" {
				delete(target.Record.Surfaces, "surface")
			}
			if tc.baseEntry != nil {
				entry := *tc.baseEntry
				entry.ID = "surface"
				base.Map = &skillcoverage.Map{Surfaces: []skillcoverage.Surface{entry}}
			}
			if tc.targetEntry != nil {
				entry := *tc.targetEntry
				entry.ID = "surface"
				target.Map = &skillcoverage.Map{Surfaces: []skillcoverage.Surface{entry}}
			}
			got := skillcoverage.Compare(base, target, tc.paths)
			want := []skillcoverage.Change{{ID: "surface", Kind: tc.kind, Skills: tc.skills, Outcome: tc.outcome}}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("changes = %+v, want %+v", got, want)
			}
		})
	}
	t.Run("sorted union and absent snapshots", func(t *testing.T) {
		record := &skillcoverage.Record{Surfaces: map[string]string{"z": newFP, "a": oldFP}}
		for _, tc := range []struct {
			base, target skillcoverage.Snapshot
			kind         string
		}{
			{skillcoverage.Snapshot{}, skillcoverage.Snapshot{Record: record}, "added"},
			{skillcoverage.Snapshot{Record: record}, skillcoverage.Snapshot{}, "removed"},
		} {
			got := skillcoverage.Compare(tc.base, tc.target, nil)
			want := []skillcoverage.Change{{ID: "a", Kind: tc.kind, Outcome: "lagging"}, {ID: "z", Kind: tc.kind, Outcome: "lagging"}}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("changes = %+v, want %+v", got, want)
			}
		}
		if got := skillcoverage.Compare(skillcoverage.Snapshot{}, skillcoverage.Snapshot{}, nil); len(got) != 0 {
			t.Fatalf("empty snapshots = %+v", got)
		}
	})
}
