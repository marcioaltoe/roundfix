package skills

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"
)

func TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome(t *testing.T) {
	t.Parallel()
	cases := []struct {
		skill, heading string
		checks         [][2]string
	}{
		{"qa-gate", "### A fresh home keeps the host toolchains", [][2]string{{"suites_keep_the_session_home", "A test-suite run, the repository Verification and every other command that does not run the built product keep the session's own `HOME` and `TMPDIR`, unless the `qa` Task requires otherwise."},
			{"a_product_row_may_get_a_fresh_home", "A row that runs the built product gives it a fresh temporary `HOME` when the Task asks for one, so that the product never reads the user's own configuration."},
			{"locations_read_before_home_changes", "Whenever a command runs under a fresh `HOME`, read the locations of the toolchains and caches the repository's stacks use before `HOME` changes, and pass them to that command unchanged: the fresh home hides the user's configuration, never the installed toolchains."},
			{"go_locations", "For Go, pass `GOMODCACHE` and `GOCACHE` as `go env` prints them."},
			{"rust_locations", "For Rust, pass `RUSTUP_HOME` and `CARGO_HOME`, their values or `$HOME/.rustup` and `$HOME/.cargo` when they are unset."},
			{"another_stack", "For another stack, pass the variables that locate its installed toolchain."},
			{"never_installs_a_toolchain", "The command uses those toolchains and never installs, updates or removes one."},
			{"a_home_prohibition_lets_them_pass", "A requirement that forbids reading the real home directory still lets the gate pass these locations, unless it names them."},
			{"created_outside_the_checkout", "Create the fresh home outside the checkout and remove it afterwards."},
			{"a_hidden_toolchain_is_no_environment_block", "A command that fails only because a fresh `HOME` hid a toolchain is a gate error, never an environment block: run it again with the toolchain locations and record that result."}}},
		{"write-tasks", "## QA environment", [][2]string{{"the_gate_keeps_the_toolchains", "Write a `qa` Task so that the gate keeps the toolchains the repository's stacks install under the home directory."},
			{"suites_keep_the_session_home", "Test-suite runs and the repository Verification keep the session's own `HOME` and `TMPDIR`; only a row that runs the built product needs a fresh temporary `HOME`, so that the product never reads the user's own configuration."},
			{"the_locations_are_passed", "Every command under a fresh `HOME` receives the locations of those toolchains and caches, read before `HOME` changes and passed unchanged, as the qa-gate skill's section \"A fresh home keeps the host toolchains\" lists them: for Go, `GOMODCACHE` and `GOCACHE`; for Rust, `RUSTUP_HOME` and `CARGO_HOME`."},
			{"a_home_prohibition_names_its_exception", "A requirement that forbids reading the real home directory names these locations as its exception."},
			{"a_suite_states_why", "A suite runs under a fresh `HOME` only when the Task states why, and then it keeps the same locations."}}},
	}
	for copyName, source := range roundfixLayoutSources() {
		t.Run(copyName, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(tc.skill, func(t *testing.T) {
					data, err := fs.ReadFile(source, tc.skill+"/SKILL.md")
					if err != nil {
						t.Fatal(err)
					}
					text := string(data)
					heading := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(tc.heading) + `\s*$`)
					locations := heading.FindAllStringIndex(text, -1)
					if len(locations) != 1 {
						t.Fatalf("heading count = %d, want 1", len(locations))
					}
					start := locations[0][0]
					section := text[locations[0][1]:]
					level := strings.SplitN(tc.heading, " ", 2)[0]
					next := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(level) + ` [^\n]+$`)
					if loc := next.FindStringIndex(section); loc != nil {
						section = section[:loc[0]]
					}
					normalized := strings.Join(strings.Fields(section), " ")
					for _, check := range tc.checks {
						t.Run(check[0], func(t *testing.T) {
							phrase := strings.Join(strings.Fields(check[1]), " ")
							if count := strings.Count(normalized, phrase); count != 1 {
								t.Errorf("phrase count = %d, want 1: %s", count, phrase)
							}
						})
					}
					if tc.skill == "qa-gate" {
						t.Run("before_the_qa_settlement_section", func(t *testing.T) {
							settlement := strings.Index(text, "### QA settlement")
							if settlement < 0 || start >= settlement {
								t.Error("fresh-home section must precede QA settlement")
							}
						})
					}
				})
			}
		})
	}
}
