---
task: task_02
spec: 0233-a-light-tier-on-open-models
status: pending
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
