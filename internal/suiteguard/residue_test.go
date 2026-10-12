package suiteguard

import (
	"testing"
)

func TestResidueMarkRequiresTheExactEnvironmentAssignment(t *testing.T) {
	t.Parallel()
	const mark = "123.0123456789abcdef"
	for _, tc := range []struct {
		name, environment string
		want              bool
	}{
		{"exact", testBinaryMarkEnv + "=" + mark + "\x00OTHER=value\x00", true},
		{"another mark", testBinaryMarkEnv + "=" + mark + "0\x00", false},
		{"another key", "OTHER_" + testBinaryMarkEnv + "=" + mark + "\x00", false},
		{"value containing assignment", "OTHER=" + testBinaryMarkEnv + "=" + mark + "\x00", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := environmentHasMark([]byte(tc.environment), mark); got != tc.want {
				t.Fatalf("mark match = %v, want %v", got, tc.want)
			}
		})
	}
	if environmentHasMark([]byte(testBinaryMarkEnv+"=\x00"), "") {
		t.Fatal("empty mark matched")
	}
}
