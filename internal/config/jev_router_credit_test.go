package config

import (
	"fmt"
	"path/filepath"
	"testing"
)

const jevRouterCreditWarning = "config: jev.router_min_credit_usd in Project Config is ignored; set jev.router_min_credit_usd in User Config\n"

func TestJevRouterMinCreditIsReadFromUserConfig(t *testing.T) {
	opts, stderr := jevCeilingFixture(t, "jev:\n  router_min_credit_usd: 20\n", "")
	loaded, err := Load(opts)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Config.Jev.RouterMinCreditUSD != 20 || stderr.Len() != 0 {
		t.Fatalf("floor=%v stderr=%s", loaded.Config.Jev.RouterMinCreditUSD, stderr)
	}
}

func TestJevRouterMinCreditIsUnsetByDefault(t *testing.T) {
	opts, _ := jevCeilingFixture(t, "", "")
	loaded, err := Load(opts)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Config.Jev.RouterMinCreditUSD != 0 || Builtin().Jev.RouterMinCreditUSD != 0 {
		t.Fatal("unset floor must stay zero")
	}
}

func TestProjectConfigCannotSetTheJevRouterMinCredit(t *testing.T) {
	for _, user := range []string{"", "jev:\n  router_min_credit_usd: 25\n"} {
		for _, project := range []string{"50", "0", "invalid"} {
			t.Run(fmt.Sprintf("%q/%s", user, project), func(t *testing.T) {
				opts, stderr := jevCeilingFixture(t, user, "jev:\n  router_min_credit_usd: "+project+"\n")
				loaded, err := Load(opts)
				if err != nil {
					t.Fatal(err)
				}
				want := 0.0
				if user != "" {
					want = 25
				}
				if loaded.Config.Jev.RouterMinCreditUSD != want || stderr.String() != jevRouterCreditWarning {
					t.Fatalf("floor=%v stderr=%q", loaded.Config.Jev.RouterMinCreditUSD, stderr.String())
				}
			})
		}
	}
}

func TestJevRouterMinCreditRefusesANonPositiveOrInfiniteValue(t *testing.T) {
	for _, value := range []string{"0", "-1", ".inf", "-.inf", ".nan", "null", "nope", "true", "[]", "'50'"} {
		t.Run(value, func(t *testing.T) {
			opts, _ := jevCeilingFixture(t, "jev:\n  router_min_credit_usd: "+value+"\n", "")
			_, err := Load(opts)
			want := fmt.Sprintf("parse config %q: jev.router_min_credit_usd must be a finite number greater than 0", filepath.Join(opts.HomeDir, ".roundfix/config.yml"))
			if err == nil || err.Error() != want {
				t.Fatalf("err=%v want=%s", err, want)
			}
		})
	}
}
