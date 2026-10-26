// Suite: Python profile composition and skill dispatch
// Invariant: Python guidance governs its toolchain with force and conditionally dispatches only chosen application tools.
// Boundary IN: embedded module dispatch and greenfield plan postimages for a temporary adopter.
// Boundary OUT: Python tool execution, skill installation, and Daemon Verification.

package baseline

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

func TestThePythonModuleDispatchesItsSkills(t *testing.T) {
	t.Parallel()
	asset, ok := mustEmbeddedCatalog(t).Module("python")
	if !ok {
		t.Fatal("missing Python module")
	}
	var module struct {
		RequiredSkills []string `json:"requiredSkills"`
		SkillDispatch  []struct {
			Skill    string `json:"skill"`
			Triggers []struct {
				ID   string `json:"id"`
				When string `json:"when"`
			} `json:"triggers"`
		} `json:"skillDispatch"`
	}
	if err := json.Unmarshal(asset.Data, &module); err != nil {
		t.Fatal(err)
	}
	want := []string{"uv", "ruff", "ty", "pytest", "python-services", "typer", "fastapi", "python-file-pipelines", "developing-with-streamlit"}
	if !slices.Equal(module.RequiredSkills, want) {
		t.Fatalf("required skills = %v, want %v", module.RequiredSkills, want)
	}
	if len(module.SkillDispatch) != len(want) {
		t.Fatalf("dispatch count = %d, want %d", len(module.SkillDispatch), len(want))
	}
	conditional := []string{"typer", "fastapi", "python-file-pipelines", "developing-with-streamlit"}
	for index, dispatch := range module.SkillDispatch {
		if dispatch.Skill != want[index] {
			t.Errorf("dispatch %d = %q, want %q", index, dispatch.Skill, want[index])
		}
		if len(dispatch.Triggers) != 1 {
			t.Errorf("%s has %d triggers, want one", dispatch.Skill, len(dispatch.Triggers))
			continue
		}
		trigger := dispatch.Triggers[0]
		if trigger.ID != "trigger.python."+dispatch.Skill {
			t.Errorf("%s trigger ID = %q", dispatch.Skill, trigger.ID)
		}
		if got, want := strings.HasPrefix(trigger.When, "Only when the repository"), slices.Contains(conditional, dispatch.Skill); got != want {
			t.Errorf("%s conditional = %t, want %t: %s", dispatch.Skill, got, want, trigger.When)
		}
	}
}

func TestThePythonProfileRendersThePythonGuide(t *testing.T) {
	t.Parallel()
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository: newPlanRepository(t), ProfileID: "python", Decisions: planTestDecisions(),
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan == nil || outcome.Result.State != "ready" {
		t.Fatalf("Python outcome = %+v, want ready plan", outcome.Result)
	}
	guide := strings.Join(strings.Fields(string(planPostimage(t, *outcome.Plan, "docs/agents/python.md").Content)), " ")
	want := []string{
		"These rules govern the repository's Python project: its package, its entry points and its tests, managed with uv. Code in another language follows its own guide.",
		"- **mandatory**: Change dependencies only through uv: `uv add`, `uv add --group dev`, `uv remove` and `uv lock --upgrade-package`. Commit `pyproject.toml` and `uv.lock` together, and never install into the project environment with `pip`.",
		"- **mandatory**: Run Python, pytest, ruff and every other project tool through `uv run`, so each uses the locked project environment and the Python version the project pins, never a globally installed interpreter or tool.",
		"- **mandatory**: Create a worktree's or CI job's environment with `uv sync --frozen` and the development group the gate needs, so it installs exactly what `uv.lock` records; a step that must also refuse a lockfile that no longer matches `pyproject.toml` uses `--locked` instead.",
		"- **mandatory**: The repository Verification checks the lockfile with `uv lock --check --offline`: it compares `uv.lock` with `pyproject.toml` from local files alone, fails when they no longer match, and needs no network or cache, so it runs inside a Run as in CI.",
		"- **mandatory**: Use focused commands such as `uv run pytest <node id>`, `uv run ruff check <path>` and `uv run ty check <path>` while iterating, then run the selected repository Verification.",
		"- **prohibited**: Do not hide a ruff finding: never pass `--exit-zero`, add a blanket `# noqa`, or widen the ruff configuration to make a finding disappear. A `# noqa` names the rule code it silences and the reason.",
		"- **mandatory**: The repository Verification type-checks the package and its tests with `uv run ty check --error-on-warning`, ty being a development dependency. Do not hide a ty diagnostic: never pass `--exit-zero` or `--exit-zero-on-warning`, add a blanket `# ty: ignore` or `# type: ignore`, or widen the ty configuration to make a diagnostic disappear. A suppression names the rule it silences and the reason.",
		"- **mandatory**: A Task Verification selects a new test by its pytest node id, such as `tests/test_report.py::test_totals`, and never passes that test's file or directory in the same command. pytest collects a node id once when another argument already contains it and exits 0 even before the test exists, so the command could never fail before the work; a missing node id alone exits 4 and an empty selection exits 5. `roundfix spec check` reports a command that names a test both ways.",
		"- **mandatory**: Keep tests hermetic: set or clear with `monkeypatch` every environment variable the code under test reads, write only under `tmp_path`, and never read the host's credentials, home directory, network or tool state.",
		"- **mandatory**: A test that needs data from outside the repository, such as a mounted drive or a credentialed service, calls `pytest.skip` with a reason that names the environment variable or path it lacks, so the selected Verification passes on a machine without that data. A Task Verification never depends on such a test.",
		"- **mandatory**: When a change touches code that branches on the operating system, such as on `sys.platform`, `os.name` or `platform.system()`, exercise every branch it names: test each one with `monkeypatch` over the value the code reads, or run `uv run ty check --python-platform <platform>` for each platform it names, because a test run on the host executes only the host's branch.",
		"- **mandatory**: Keep each entry point thin, whether a console script, a web application factory or a worker's start-up: it reads its input and calls behavior that lives in the package. Test behavior through the package and the entry point's observable contract: for a console script, standard output, standard error and the exit code it raises with `SystemExit`; for a web application or a worker, the responses or messages it produces.",
	}
	for _, text := range want {
		if !strings.Contains(guide, strings.Join(strings.Fields(text), " ")) {
			t.Errorf("Python guide lacks %q", text)
		}
	}
	dispatch := string(planPostimage(t, *outcome.Plan, "docs/agents/skill-dispatch.md").Content)
	for _, skill := range []string{"uv", "ruff", "ty", "pytest", "python-services", "typer", "fastapi", "python-file-pipelines", "developing-with-streamlit"} {
		if !strings.Contains(dispatch, "`"+skill+"`") {
			t.Errorf("dispatch guide lacks %s", skill)
		}
	}
	root := string(planPostimage(t, *outcome.Plan, "AGENTS.md").Content)
	if !strings.Contains(root, "Python implementation rules are mandatory: `docs/agents/python.md`") {
		t.Errorf("root lacks Python guide pointer: %s", root)
	}
	for _, postimage := range outcome.Plan.Postimages {
		if slices.Contains([]string{"docs/agents/go.md", "docs/agents/rust.md", "docs/agents/tui.md", "docs/agents/cli.md"}, postimage.Path) {
			t.Errorf("python plan renders %s", postimage.Path)
		}
	}
}
