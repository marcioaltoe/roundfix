---
task: task_01
spec: 0230-retire-the-jev-router
status: pending
type: backend
complexity: high
---

# Task 01: Agent sessions no longer reach OpenRouter, and the runner refuses the models the subscription rule reserves

## Overview

The Jev Router lives in the runner, the daemon's Agent Session owner, the CLI
engine wiring and the `internal/jevrouter` package. This Task deletes all of
it and adds the subscription rule's predicate in the config package, which
the runner checks before it prepares a session or sends a prompt, so no
Agent Session can reach an OpenAI or Anthropic model through OpenRouter
whatever configuration produced its selection. The configuration rules that
still name the router stay until task_02. The slice is verifiable on its own:
a refused selection never starts the ACP adapter and falls back with the
reason `subscription_only`, and outside the judge no source names
OpenRouter's host.

## Requirements

1. MUST delete the `internal/jevrouter` package with its tests, the agent
   package's router file and its tests, and the daemon's router gate tests.
2. MUST remove from the ACPX runner the inline `roundfix-openrouter`
   provider, the relay and its per-session tokens, the router observation on
   `ExecuteResult`, the router branch of the prompt failure classification,
   the endpoint override field and the fixture hooks the deleted tests used;
   no non-test source sets `OPENCODE_CONFIG_CONTENT` any more, so an
   inherited value keeps passing through unchanged, and a runner test SHOULD
   prove that pass-through for an allowed `opencode` model, as the deleted
   router test did.
3. MUST remove from the daemon the Jev Router gate interface and type, the
   dependencies `JevRouter`, `JevMonthlyCeilingUSD` and
   `JevRouterMinCreditUSD`, the default gate built in `NewEngine`, and the
   routed branch of the prepared-prompt path that read spend and appended
   `router-prompt` lines; and MUST remove the router refusal codes from the
   selection reason classification, which returns `subscription_only` for
   the new refusal (TechSpec Invariant 8).
4. MUST remove the ceiling and floor parameters that only fed the gate from
   the resolve and Implement Run engine constructors in the CLI.
5. MUST add to the config package's profile code `SubscriptionRule`,
   `SubscriptionOnlyReason` and `CheckSubscriptionRule` with the shape and
   Invariants 1 to 5 of the TechSpec's Interfaces section and the messages of
   API Contracts 1 and 2, without yet calling it from configuration
   validation.
6. MUST call it in the runner's runtime selection validation and at the
   start of `RunPrompt`, before any acpx process starts, wrapping a refusal as
   a selection failure with reason `subscription_only: <message>` and field
   `agent model` (Invariant 7, API Contract 4).
7. MUST add tests: a table test of `CheckSubscriptionRule` covering every
   case the TechSpec's Testing Approach lists; a runner test that
   `RunPrompt` and `PrepareSession` refuse `openrouter/anthropic/claude-opus-5.5`,
   `openrouter/openai/gpt-6.1-sol` and `roundfix-openrouter/typesafe/jev-router`
   on `opencode` with the reason and leave no acpx invocation, while
   `openrouter/deepseek/deepseek-v4-pro` still runs; and a daemon test that a
   `subscription_only` selection failure before Agent work activates the
   fallback and the fallback Run Event's `reason_code` is `subscription_only`.
8. MUST NOT change the judge package or the `spec judge` command, any
   configuration validation rule, any guide or skill, or any built-in or
   Recommended Profile; no test may reach the network or the real
   `~/.roundfix`.

## Subtasks

- [ ] Delete the router package, the agent router file and the router tests.
- [ ] Remove the relay, provider and observation wiring from the runner.
- [ ] Remove the gate and its dependencies from the daemon and the CLI engines.
- [ ] Add the predicate and call it in the runner.
- [ ] Map `subscription_only` in the daemon's reason classification.
- [ ] Add the predicate, runner and daemon tests.

## Acceptance Criteria

- [ ] For every runtime and model in the predicate table, the result matches
      TechSpec Invariants 1 to 5 and the refusal text matches API Contract 1
      or 2.
- [ ] A refused selection never starts the ACP adapter, through either
      `RunPrompt` or `PrepareSession`; an allowed OpenRouter model still does.
- [ ] Before Agent work, a `subscription_only` refusal hands the Task to the
      fallback with that reason code.
- [ ] Outside the judge, no Go source names OpenRouter's host or writes
      `router-prompt`, no non-test source sets `OPENCODE_CONFIG_CONTENT`, and
      the judge package is byte-identical.

## Context

- instruction: `docs/adr/0235-openai-and-anthropic-models-run-only-through-the-codex-and-claude-subscriptions.md`
- deletes: `internal/jevrouter/credits.go`
- deletes: `internal/jevrouter/credits_test.go`
- deletes: `internal/jevrouter/key.go`
- deletes: `internal/jevrouter/key_test.go`
- deletes: `internal/jevrouter/key_limit_test.go`
- deletes: `internal/jevrouter/ledger.go`
- deletes: `internal/jevrouter/ledger_test.go`
- deletes: `internal/jevrouter/relay.go`
- deletes: `internal/jevrouter/relay_test.go`
- deletes: `internal/jevrouter/spend.go`
- deletes: `internal/jevrouter/spend_test.go`
- deletes: `internal/agent/jev_router.go`
- deletes: `internal/agent/jev_router_test.go`
- deletes: `internal/daemon/jev_router_gate_test.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/agent.go`
- interface: `internal/agent/acpx_runner_test.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/agent_session_owner.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/implement.go`
- interface: `internal/config/profiles.go`
- creates: `internal/config/subscription_rule_test.go`
- creates: `internal/agent/subscription_rule_test.go`
- creates: `internal/daemon/subscription_rule_test.go`

## Verification

- `test ! -e internal/jevrouter && test ! -e internal/agent/jev_router.go && test ! -e internal/agent/jev_router_test.go && test ! -e internal/daemon/jev_router_gate_test.go && ! grep -rqE 'jevrouter|JevRouter|router-prompt|openrouter_credit_|jev_router_key|jev_ceiling_reached|jev_spend_unreadable' internal/agent internal/daemon internal/cli cmd && ! grep -rq --include='*.go' --exclude='*_test.go' OPENCODE_CONFIG_CONTENT internal cmd && hosts="$(grep -rl 'openrouter[.]ai' internal cmd | sort | tr '\n' ' ')" && test "$hosts" = "internal/cli/spec_judge_test.go internal/judge/questions.json " && git diff --quiet HEAD -- internal/judge internal/cli/spec_judge.go && go build -buildvcs=false ./... && go vet ./internal/agent ./internal/daemon ./internal/cli ./internal/config` — expected: exit 0; before this Task the router package, files and terms exist, so the command fails.
- `out="$(go test -count=1 -v -run '^(TestCheckSubscriptionRule|TestRunPromptRefusesSubscriptionModelsBeforeTheAdapter|TestPrepareSessionRefusesSubscriptionModels|TestSubscriptionOnlyRefusalActivatesTheFallback)$' ./internal/config ./internal/agent ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCheckSubscriptionRule TestRunPromptRefusesSubscriptionModelsBeforeTheAdapter TestPrepareSessionRefusesSubscriptionModels TestSubscriptionOnlyRefusalActivatesTheFallback; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the four tests exists, so the command fails.

## References

- `_prd.md` → Goals; User Stories 1, 2 and 4; Core Features 1, 3 and 5; Success Metrics 2 and 3
- `_techspec.md` → Interfaces; Invariants 1 to 9; API Contract 1; API Contract 2; API Contract 4; API Contract 6; Testing Approach; Build Order 2
- ADR-0235; ADR-0050; ADR-0114; ADR-0201
