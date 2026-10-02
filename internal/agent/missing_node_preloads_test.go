package agent

import (
	"reflect"
	"testing"
)

func TestMissingNodePreloadsReusesTheAgentSplitter(t *testing.T) {
	for _, value := range []string{
		`--require="/missing path.cjs" --import=file:///missing.mjs -r /existing.cjs --require package-name`,
		`--require '/single quoted.cjs'`,
		`--require "/unbalanced`,
		`--require=/missing.cjs --require=/missing.cjs`,
		`--require=relative.cjs`,
	} {
		t.Run(value, func(t *testing.T) {
			exists := func(path string) bool { return path == "/existing.cjs" }
			_, want, _ := agentNodeOptions(value, exists)
			got := MissingNodePreloads(value, exists)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("missing paths = %v, want %v", got, want)
			}
		})
	}
	got := MissingNodePreloads(`--require="/missing path.cjs" --import=file:///missing.mjs -r /existing.cjs --require package-name`, func(path string) bool { return path == "/existing.cjs" })
	if !reflect.DeepEqual(got, []string{"/missing path.cjs", "/missing.mjs"}) {
		t.Fatalf("missing paths = %v", got)
	}
}
