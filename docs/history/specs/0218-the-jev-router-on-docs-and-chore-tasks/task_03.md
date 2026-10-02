---
task: task_03
spec: 0218-the-jev-router-on-docs-and-chore-tasks
status: completed
type: backend
complexity: high
---

# Task 03: Every routed prompt runs under the shared ceiling and leaves a Judge Log line

## Overview

Tasks 01 and 02 give Roundfix a routed selection and a way to read the
month's Jev spend. Nothing yet stops a routed prompt at the ceiling or
records what it cost. This Task adds a `JevRouterGate` to the Daemon's
dependencies and calls it around each routed prompt in
`agentSessionOwner.runPrepared`: a refusal before Agent work begins moves the
Fallback Chain on, a refusal after it fails the Work Item (ADR-0050,
ADR-0114), and each routed prompt appends one `router-prompt` line
(ADR-0218). It also describes the router in the Roundfix Skill.

## Requirements

1. MUST add `JevRouterGate` and `Dependencies.JevRouter` in
   `internal/daemon/engine.go`, defaulting when nil to a gate built from the
   process environment, Roundfix Home, the `https://openrouter.ai/api/v1`
   endpoint and the judge's loaded ceiling.
2. MUST call `Before` in `agentSessionOwner.runPrepared` only when the active
   runtime and model satisfy `agent.IsJevRouterSelection`, and MUST turn its
   refusal into the `SelectionFailureError` reasons the TechSpec "The gate
   and the record" states (`jev_spend_unreadable: <error>`,
   `jev_ceiling_reached: month's Jev spend US$<total> of US$<ceiling>`), so
   the existing `Run` loop falls back before work and fails the Work Item
   after it. A refused prompt MUST NOT reach the runner.
3. MUST call `After` once per routed prompt that ran, with the Run, Spec,
   scope, category, attempt, latency, tokens, the usage `Before` returned and
   whether the prompt failed. A failed `After` MUST be written to the Run's
   progress stream and MUST NOT change the prompt's result.
4. MUST leave non-routed prompts byte-for-byte as today, with no key
   endpoint call.
5. MUST add a `### Jev Router` section to
   `.agents/skills/roundfix/references/runtime.md` stating the routed
   selection, the Project Config rule, the key variable, the shared US$5
   ceiling, and the refusal classifications `jev_router_key_missing`,
   `jev_spend_unreadable` and `jev_ceiling_reached`; raise both version
   fields of `.agents/skills/roundfix/SKILL.md`; run `make skills-sync`; and
   re-record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
6. MUST test with a fake gate and a fake runner only, with a sentinel key
   that appears in no line, and MUST NOT change the `### QA settlement`
   section of any skill.

## Subtasks

- [ ] Add the gate dependency and its default.
- [ ] Call it before and after routed prompts in the session owner.
- [ ] Prove fallback before work, failure after work, and one line per prompt.
- [ ] Describe the router in the Roundfix Skill and raise its version.

## Acceptance Criteria

- [ ] At the ceiling, the first routed prompt never runs and the fallback
      selection starts after its notification; after work began, the Work
      Item fails with `jev_ceiling_reached`.
- [ ] An unreadable spend falls back with `jev_spend_unreadable`.
- [ ] Below the ceiling, one routed prompt appends one `router-prompt` line
      whose cost is the usage change, and a non-routed prompt calls no gate.

## Context

- creates: `internal/daemon/jev_router_gate_test.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/agent_session_owner.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `docs/adr/0218-the-jev-router-is-a-project-selected-opencode-model-under-the-shared-jev-ceiling.md`
- instruction: `docs/adr/0050-configured-fallbacks-activate-after-notification.md`
- instruction: `docs/adr/0114-opening-an-agent-session-is-not-agent-work.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestJevRouterCeilingFallsBackBeforeWork|TestJevRouterCeilingFailsTheTaskAfterWork|TestJevRouterUnreadableSpendFallsBack|TestJevRouterPromptAppendsOneLine|TestNonRoutedPromptCallsNoGate)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestJevRouterCeilingFallsBackBeforeWork TestJevRouterCeilingFailsTheTaskAfterWork TestJevRouterUnreadableSpendFallsBack TestJevRouterPromptAppendsOneLine TestNonRoutedPromptCallsNoGate; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for phrase in '### Jev Router' 'roundfix-openrouter/typesafe/jev-router' 'jev_ceiling_reached' 'jev_router_key_missing' 'ROUNDFIX_OPENROUTER_API_KEY'; do tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/runtime.md | grep -qF -- "$phrase" || { printf 'missing phrase in the runtime reference: %s\n' "$phrase" >&2; exit 1; }; done; diff -r .agents/skills/roundfix skills/roundfix >/dev/null || { printf 'mirror differs: skills/roundfix\n' >&2; exit 1; }; vout="$(go test -count=1 -v -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills 2>&1)" || { printf '%s\n' "$vout"; exit 1; }; printf '%s\n' "$vout" | grep -q -- '--- PASS: TestEveryOwnedSkillVersionIsRecorded' || { printf '%s\n' "$vout"; exit 1; }; make skills-sync-check` — expected: exit 0; before this Task the named Daemon tests do not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 3; User Story 3; Core Features 4, 5 and 7; Success Metric 3
- [_techspec.md](_techspec.md) — Interfaces; The gate and the record; API Contracts 3-4; Testing Approach 3-4; Build Order 3
- ADR-0218; ADR-0050; ADR-0114; ADR-0189

## Result

Implemented the Task 03 slice for Daemon Verification. Status and the
declared Verification remain Daemon-owned.

- Added `Dependencies.JevRouter` and its `Before`/`After` interface. The nil
  default captures the process environment and Home, the OpenRouter key
  endpoint, the judge's loaded ceiling, and the engine clock. `Before`
  refuses unreadable spend or spend at/above the ceiling with the specified
  reasons and four-decimal amounts. `After` reads key usage and appends the
  existing Judge Log schema through `jevrouter.Ledger`.
- Routed session prompts call `Before` before runner dispatch and `After`
  once after dispatch, preserving existing token usage recording. The record
  carries Run, Spec, scope, category, repository, attempt, prompt latency,
  input/output tokens, initial usage, and failure state. Recording survives
  prompt cancellation; append errors go to progress without replacing the
  prompt result. Non-routed prompts retain their existing dispatch and never
  call the gate. Fallback notifications preserve all three Jev reason codes.
- Added the Roundfix runtime reference's `### Jev Router` section, raised
  both skill version fields to `0.1.16`, synchronized the embedded mirror,
  and recorded the new version. `### QA settlement` remains byte-identical.

Acceptance evidence from focused checks:

| Acceptance criterion | Evidence |
| --- | --- |
| Ceiling refusal before work falls back after notification; after work it fails the Work Item | `TestJevRouterCeilingFallsBackBeforeWork` executes TaskCycle, observes both notifications before fallback preparation, and proves the refused prompt never reaches the runner. `TestJevRouterCeilingFailsTheTaskAfterWork` refuses the Verification Feedback prompt after the first prompt, proves no second runner call or fallback, and observes the Task's failed status and `jev_ceiling_reached` outcome reason. |
| Unreadable spend falls back with `jev_spend_unreadable` | `TestJevRouterUnreadableSpendFallsBack` observes the classified fallback event, notifications before fallback preparation, one fallback prompt, and no record for the refused prompt. |
| One below-ceiling routed prompt appends one usage-delta line; non-routed prompts call no gate | `TestJevRouterPromptAppendsOneLine` reads one real temporary Judge Log line with US$0.15 cost from usage 1.20 to 1.35, correct identity, attempt, latency and tokens, and no key, prompt or answer. `TestNonRoutedPromptCallsNoGate` covers another OpenCode model and another runtime and preserves the runner result. |

Additional focused coverage proves append failure preserves successful and
failed prompt results, canceled prompts still append a skipped line,
Runner-only implementations record every prompt with unreported tokens,
key-missing notifications retain their classification, and the default gate
uses process Home/environment and the judge's ceiling without invoking it.
Prompt execution tests use fake gates and fake runners; the key is a sentinel
checked against progress, Run Events and the temporary accounting line. No
live key endpoint or Agent runtime was invoked.

Focused commands and outcomes:

- Before implementation, `GOCACHE=/private/tmp/roundfix-task03-gocache go test
  ./internal/daemon -run TestJevRouterPromptAppendsOneLine -count=1` failed to
  compile because `Dependencies.JevRouter` and scope Spec metadata were absent.
  The initial test also used a nonexistent `spec.ParseTask`; that fixture
  mistake was corrected to the existing `spec.Load` API.
- Final `GOCACHE=/private/tmp/roundfix-task03-gocache go test -race -count=1
  ./internal/daemon -run
  'Test(JevRouter|NonRoutedPrompt|FallbackEligibility|NoWorkStarted|NoFallbackAfterAgentWorkStarted|EachPromptRecordsItsUsageWithTheOwnerScope)'`:
  exit 0 (`ok roundfix/internal/daemon`, 3.644s).
- `make skills-sync`: exit 0.
- `GOCACHE=/private/tmp/roundfix-task03-gocache go test ./skills -run
  '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`: exit 0
  after the final documentation edit.
- `GOCACHE=/private/tmp/roundfix-task03-gocache make baseline-digests`:
  exit 0; `changed:false`, no derived artifacts changed.
- Python byte comparisons confirmed the complete Roundfix skill mirror and
  unchanged QA settlement sections; `git diff --check`: exit 0.

No follow-up work was added. The Task Graph, other Task files, and unrelated
paths were not edited; no commit, push, or Pull Request was made. The Task's
declared Verification command was not run.

## Carry-forward provenance

- Source Run: `run_20261002T194325Z_56cb22f8a792bdff`
- Source commit: `399b346720c7783b1e08ccc150adfa344dcb59e1`
