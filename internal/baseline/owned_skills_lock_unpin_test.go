package baseline

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRecordInstalledOwnedSkillsRemovesTheRefOfAnInstalledOwnedSkill(t *testing.T) {
	t.Parallel()
	stale := strings.Repeat("0", 64)
	for _, tt := range []struct {
		name   string
		before string
		after  string
		hashed bool
	}{
		{
			name: "skills tool format",
			before: `{
  "version": 1,
  "skills": {
    "roundfix": {
      "source": "marcioaltoe/skills",
      "ref": "8d71fade68c54693bedbd917371f0514311af157",
      "sourceType": "github",
      "skillPath": "skills/roundfix",
      "computedHash": "` + stale + `"
    },
    "testing-boss": {"source":"external/skills","ref":"keep","computedHash":"` + stale + `"},
    "write-prd": {"ref":"not installed","computedHash":"` + stale + `"}
  }
}
`,
			after: `{
  "version": 1,
  "skills": {
    "roundfix": {
      "source": "marcioaltoe/skills",
      "sourceType": "github",
      "skillPath": "skills/roundfix",
      "computedHash": "<hash>"
    },
    "testing-boss": {"source":"external/skills","ref":"keep","computedHash":"` + stale + `"},
    "write-prd": {"ref":"not installed","computedHash":"` + stale + `"}
  }
}
`,
			hashed: true,
		},
		{
			name:   "ref first",
			before: `{"version":1,"skills":{"roundfix":{ "ref" : "x" , "computedHash":"` + stale + `"}}}`,
			after:  `{"version":1,"skills":{"roundfix":{ "computedHash":"<hash>"}}}`,
			hashed: true,
		},
		{
			name:   "ref last with a current hash",
			before: `{"version":1,"skills":{"roundfix":{"computedHash":"<hash>", "ref":""}}}`,
			after:  `{"version":1,"skills":{"roundfix":{"computedHash":"<hash>"}}}`,
		},
		{
			name:   "ref alone",
			before: `{"version":1,"skills":{"roundfix":{"ref":null}}}`,
			after:  `{"version":1,"skills":{"roundfix":{"computedHash":"<hash>"}}}`,
			hashed: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root, hash := ownedLockFixture(t)
			before := strings.ReplaceAll(tt.before, "<hash>", hash)
			want := strings.ReplaceAll(tt.after, "<hash>", hash)
			filename := filepath.Join(root, "skills-lock.json")
			if err := os.WriteFile(filename, []byte(before), 0o640); err != nil {
				t.Fatal(err)
			}
			update, err := RecordInstalledOwnedSkills(t.Context(), root, []string{"write-prd", "roundfix", "roundfix"})
			if err != nil {
				t.Fatal(err)
			}
			after, err := os.ReadFile(filename)
			if err != nil {
				t.Fatal(err)
			}
			if string(after) != want {
				t.Fatalf("lock=%s\nwant=%s", after, want)
			}
			var hashed []string
			if tt.hashed {
				hashed = []string{"roundfix"}
			}
			if !reflect.DeepEqual(update.Hashed, hashed) || !reflect.DeepEqual(update.Unpinned, []string{"roundfix"}) || !reflect.DeepEqual(update.Names(), []string{"roundfix"}) {
				t.Fatalf("update=%+v want hashed=%v unpinned=[roundfix]", update, hashed)
			}
			info, err := os.Stat(filename)
			if err != nil {
				t.Fatal(err)
			}
			if info.Mode().Perm() != 0o640 {
				t.Fatalf("mode=%v", info.Mode())
			}
			second, err := RecordInstalledOwnedSkills(t.Context(), root, []string{"write-prd", "roundfix"})
			if err != nil || len(second.Names()) != 0 {
				t.Fatalf("second=%+v err=%v", second, err)
			}
			afterInfo, err := os.Stat(filename)
			if err != nil {
				t.Fatal(err)
			}
			after, err = os.ReadFile(filename)
			if err != nil {
				t.Fatal(err)
			}
			if !os.SameFile(info, afterInfo) || !info.ModTime().Equal(afterInfo.ModTime()) || string(after) != want {
				t.Fatal("second run wrote the lock")
			}
		})
	}
}
