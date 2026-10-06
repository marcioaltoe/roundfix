---
task: task_02
spec: 0233-a-light-tier-on-open-models
status: completed
type: backend
complexity: medium
---

# Task 02: User Config carries the light models and the light ceiling, and a tier rule and a Light Spend Log decide and record light Tasks

## Overview

Dispatch cannot choose a light model until three facts exist: the Task's
authored `complexity`, which the Task parser decodes but drops; the machine's
light model list and monthly ceiling, which User Config does not know; and a
record of what light Tasks have spent this month. This Task adds them, plus
the one helper that names the OpenRouter key variable, as pure configuration
and a new small `internal/lighttier` package, without changing any dispatch.
It is verifiable on its own through configuration loading, the tier rule over
Task files and the log in a temporary Home.

## Requirements

1. MUST carry the Task front matter's `complexity`, trimmed, on the parsed
   Task, both when a graph is loaded and when a Task is reloaded, and MUST
   leave every other Task field and every parse error unchanged.
2. MUST add User Config `openrouter.light_models` and
   `openrouter.implement_monthly_ceiling_usd` with the defaults, refusals and
   messages of API Contracts 1 and 2: an unset list means
   `deepseek/deepseek-v4.1-flash`, an explicit empty list turns the tier off,
   an entry must be `<author>/<slug>` and pass the existing subscription
   predicate with runtime `opencode` and model `openrouter/` plus the entry
   (TechSpec Invariant 3), and the ceiling must be a finite number greater
   than zero, 10 when unset.
3. MUST remove either key from Project Config before validation with the
   existing User Config-only warning naming the key (API Contract 3), and MUST
   keep every existing key, warning and deprecated-key behavior unchanged.
4. MUST add `OpenRouterImplementKey`, the only function that names the
   implementation key's variable: it reads `ROUNDFIX_OPENROUTER_API_KEY`
   through the supplied lookup and returns that name and whether it holds a
   non-empty value; it MUST NOT read `OPENROUTER_API_KEY`, and MUST NOT add
   any other variable name, which Spec 0234 owns.
5. MUST add `internal/lighttier` with `TierFor` (TechSpec Invariant 1: a Task
   is `light` only when its complexity is `low`, its type is not `qa`, and no
   `creates`, `interface` or `deletes` Context reference names a Governed
   Path; an `instruction` reference never counts), `Plan` with `Enabled`
   (Invariant 2), and the Light Spend Log of the TechSpec's Data Models:
   `AppendSpend` writes one JSON line to
   `<home>/.roundfix/openrouter/implement/<YYYY-MM>.jsonl` for the UTC month
   with directory mode 0700 and file mode 0600, and `ReadMonth` returns the
   month's summed `cost_usd`, zero for a missing file, and an error for a line
   that does not parse or carries a negative cost.
6. MUST add tests: through `Load` with disposable Homes and repositories, the
   defaults, a custom list, an empty list, each refused entry
   (`openai/gpt-6.1-sol`, `anthropic/claude-opus-5.5`, `openrouter/auto`,
   `deepseek`), each refused ceiling (`0`, `-1`, `.inf`, a string), and the
   Project Config warning for both keys; `OpenRouterImplementKey` over a fake
   lookup with the variable set, empty and absent, and with only
   `OPENROUTER_API_KEY` set; `TierFor` over Tasks parsed from disposable Spec
   files covering each clause of Invariant 1; and the Light Spend Log's
   append, sum, missing file, malformed line, negative cost, month boundary
   and file modes.
7. MUST NOT change dispatch, the runner, the daemon, the judge, any guide or
   skill, or any built-in or Recommended Profile; no test may reach the
   network or the real `~/.roundfix`.

## Subtasks

- [ ] Carry `complexity` on the parsed Task.
- [ ] Add the two User Config keys and the key helper.
- [ ] Add `internal/lighttier` with the tier rule and the Light Spend Log.
- [ ] Add the configuration, helper, tier and log tests.

## Acceptance Criteria

- [ ] User Config loads both keys with their defaults and refuses every
      invalid value with API Contract 1 or 2's message.
- [ ] Project Config ignores both keys with one warning each.
- [ ] `TierFor` returns `light` only for a non-QA `low` Task that declares no
      Governed Path to write.
- [ ] The Light Spend Log appends private lines and sums the UTC month.

## Context

- instruction: `docs/adr/0238-a-light-tier-runs-low-complexity-tasks-on-an-open-model.md`
- interface: `internal/spec/spec.go`
- interface: `internal/spec/task.go`
- interface: `internal/config/config.go`
- creates: `internal/config/light_tier.go`
- creates: `internal/config/light_tier_test.go`
- creates: `internal/spec/task_complexity_test.go`
- creates: `internal/lighttier/tier.go`
- creates: `internal/lighttier/spend_log.go`
- creates: `internal/lighttier/tier_test.go`
- creates: `internal/lighttier/spend_log_test.go`

## Verification

- `test -f internal/lighttier/tier.go && test -f internal/lighttier/spend_log.go && test -f internal/config/light_tier.go && go build -buildvcs=false ./... && go vet ./internal/config ./internal/lighttier ./internal/spec` — expected: exit 0; before this Task the package and the configuration file do not exist, so the command fails.
- `cfg="$(go test -count=1 -v -run '^(TestLightTierUserConfig|TestOpenRouterImplementKeyNamesOneVariable)$' ./internal/config 2>&1)" || { printf '%s\n' "$cfg"; exit 1; }; lt="$(go test -count=1 -v -run '^(TestTierFor|TestLightSpendLog)$' ./internal/lighttier 2>&1)" || { printf '%s\n' "$lt"; exit 1; }; sp="$(go test -count=1 -v -run '^TestTaskCarriesComplexity$' ./internal/spec 2>&1)" || { printf '%s\n' "$sp"; exit 1; }; for name in TestLightTierUserConfig TestOpenRouterImplementKeyNamesOneVariable TestTierFor TestLightSpendLog TestTaskCarriesComplexity; do printf '%s\n%s\n%s\n' "$cfg" "$lt" "$sp" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the five tests exists, so the command fails.

## References

- `_prd.md` → Goals; User Stories 2 and 3; Core Features 1, 2, 3, 4 and 5; Success Metric 5
- `_techspec.md` → Interfaces; Invariants; Data Models; API Contract 1; API Contract 2; API Contract 3; Testing Approach; Build Order 2
- ADR-0238; ADR-0235; ADR-0231; ADR-0027; ADR-0002

## Result

Implemented this Task's configuration, parsing, tier rule and spend-log slice.
Dispatch, runner, daemon, judge, guides, skills and selection profiles are
unchanged. The pre-existing `status: in_progress` is preserved; settlement and
declared Verification remain with the Daemon.

- User Config: `Config.OpenRouter` loads the default model list and US$10
  ceiling, preserves an explicit empty list as off, and accepts custom models
  and a positive finite ceiling. Model ids are checked for `<author>/<slug>`
  without whitespace, then passed to the existing subscription predicate as
  `opencode` / `openrouter/<id>`. `TestLightTierUserConfig` exercises defaults,
  custom values, off, all four required refused entries, malformed ids and
  list shapes, and every required refused ceiling with its contract message.
- Project Config: both keys are removed before validation, preserving the
  User Config values and emitting one existing-format warning per key.
  `TestLightTierUserConfig` proves that invalid project values are ignored,
  with exact warning text, both with defaults and with explicit user values.
- Tier rule: `TierFor` returns light only for non-QA `low` Tasks with no
  governed creates/interface/deletes reference; instruction references never
  count. `TestTierFor` parses disposable Task files for every clause, including
  ordinary write references and a governed instruction beside a governed
  write. `TestPlanEnabled` covers zero/empty plans and keeps key availability
  separate from enablement so later dispatch can explain a missing-key skip.
- Light Spend Log: `AppendSpend` records the schema, UTC time and prompt
  metadata as one JSON line under the UTC month's path, enforcing directory
  mode 0700 and file mode 0600, including existing storage. `ReadMonth` sums
  costs, returns zero without creating storage for a missing file, and returns
  an error rather than a partial sum for malformed records or negative costs.
  `TestLightSpendLog` covers append, sum, zero/unreported spend, missing files,
  malformed records (including an invalid timestamp), negative costs,
  September/October separation, a local September time in UTC October, and
  new/existing file modes.
- Supporting contracts: `TestTaskCarriesComplexity` proves graph load and
  reload retain trimmed complexity, tolerate custom/missing values, and leave
  every other reloaded Task field unchanged. The parser's error paths are
  untouched. `TestOpenRouterImplementKeyNamesOneVariable` uses a fake lookup
  for set, empty, absent and generic-only keys, asserting exactly one lookup
  of the implementation variable and no lookup of `OPENROUTER_API_KEY`.

Focused implementation evidence:

- Before implementation, `GOCACHE=/private/tmp/roundfix-task02-go-cache rtk proxy go test -count=1 ./internal/config ./internal/spec` exited 1 on the
  new tests' missing configuration symbols and Task complexity field.
- After the final code and test edits,
  `GOCACHE=/private/tmp/roundfix-task02-go-cache rtk proxy go test -count=1 ./internal/config ./internal/spec ./internal/lighttier` exited 0:
  config 0.999s, spec 19.524s, lighttier 0.358s. These package checks include
  all five required new tests and existing config/spec regression tests.
- `rtk proxy gofmt -w` formatted the changed Go files.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0; changed-path
  inspection found only this Task's declared source/test paths and its Task
  file. The Task status change from pending to in_progress predates this turn.
- The first check using the normal Go cache was sandbox-blocked; all subsequent
  checks used the disposable task cache above. New tests use temporary Homes
  and Spec/repository files, with no network or real Roundfix Home access.

The Task's `## Verification` commands and repository delivery gates were not
run in this Daemon-assigned turn. No commit, push or Pull Request was made.
No follow-up implementation was added to this slice.

## Carry-forward provenance

- Source Run: `run_20261006T004151Z_0e07759cbdbec8e5`
- Source commit: `d9d932f9e78bc4c04526e64ff8db6b012a4636b4`
