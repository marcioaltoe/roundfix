package config

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"
)

const jevCeilingWarning = "config: jev.monthly_ceiling_usd in Project Config is ignored; set jev.monthly_ceiling_usd in User Config\n"

func jevCeilingFixture(t *testing.T, user, project string) (LoadOptions, *bytes.Buffer) {
	t.Helper()
	home, repo := t.TempDir(), t.TempDir()
	mustMkdir(t, filepath.Join(home, ".roundfix"))
	mustMkdir(t, filepath.Join(repo, ".git"))
	mustWrite(t, filepath.Join(home, ".roundfix/config.yml"), user)
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), project)
	stderr := &bytes.Buffer{}
	return LoadOptions{HomeDir: home, WorkDir: repo, Stderr: stderr}, stderr
}

func TestJevMonthlyCeilingIsReadFromUserConfig(t *testing.T) {
	opts, stderr := jevCeilingFixture(t, "jev:\n  monthly_ceiling_usd: 50\n", "")
	loaded, err := Load(opts)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Config.Jev.MonthlyCeilingUSD != 50 || stderr.Len() != 0 {
		t.Fatalf("ceiling=%v stderr=%s", loaded.Config.Jev.MonthlyCeilingUSD, stderr)
	}
}

func TestJevMonthlyCeilingIsUnsetByDefault(t *testing.T) {
	opts, _ := jevCeilingFixture(t, "", "")
	loaded, err := Load(opts)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Config.Jev.MonthlyCeilingUSD != 0 || Builtin().Jev.MonthlyCeilingUSD != 0 {
		t.Fatal("unset ceiling must stay zero")
	}
}

func TestProjectConfigCannotSetTheJevCeiling(t *testing.T) {
	for _, user := range []string{"", "jev:\n  monthly_ceiling_usd: 25\n"} {
		for _, project := range []string{"50", "0", "invalid"} {
			t.Run(fmt.Sprintf("%q/%s", user, project), func(t *testing.T) {
				opts, stderr := jevCeilingFixture(t, user, "jev:\n  monthly_ceiling_usd: "+project+"\n")
				loaded, err := Load(opts)
				if err != nil {
					t.Fatal(err)
				}
				want := 0.0
				if user != "" {
					want = 25
				}
				if loaded.Config.Jev.MonthlyCeilingUSD != want || stderr.String() != jevCeilingWarning {
					t.Fatalf("ceiling=%v stderr=%q", loaded.Config.Jev.MonthlyCeilingUSD, stderr.String())
				}
			})
		}
	}
}

func TestJevMonthlyCeilingRefusesANonPositiveOrInfiniteValue(t *testing.T) {
	for _, value := range []string{"0", "-1", ".inf", "-.inf", ".nan", "null", "nope", "true", "[]", "'50'"} {
		t.Run(value, func(t *testing.T) {
			opts, _ := jevCeilingFixture(t, "jev:\n  monthly_ceiling_usd: "+value+"\n", "")
			_, err := Load(opts)
			want := fmt.Sprintf("parse config %q: jev.monthly_ceiling_usd must be a finite number greater than 0", filepath.Join(opts.HomeDir, ".roundfix/config.yml"))
			if err == nil || err.Error() != want {
				t.Fatalf("err=%v want=%s", err, want)
			}
		})
	}
}
