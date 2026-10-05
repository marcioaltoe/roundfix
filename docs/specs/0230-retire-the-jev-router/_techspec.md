---
spec: 0230-retire-the-jev-router
prd: _prd.md
created: 2026-10-05
---

# Retire the Jev Router — Technical Spec

## Executive Summary

The Jev Router is wired through four layers: configuration validation that
lets only Project Config select `roundfix-openrouter/typesafe/jev-router`;
the ACPX runner, which gives OpenCode an inline `roundfix-openrouter`
provider through `OPENCODE_CONFIG_CONTENT` and runs a loopback relay to
OpenRouter; the daemon's Agent Session owner, which gates each routed prompt
on OpenRouter spend and credit and appends a `router-prompt` line to the
Judge Log; and the `internal/jevrouter` package that holds the reads, the
relay and the ledger. This design deletes the package, the runner and daemon
wiring and the configuration rules, and adds one predicate in the config
package, `CheckSubscriptionRule`, called by configuration validation for every
Agent Selection and again by the runner before a session or prompt. The
primary trade-off is the refusal's shape: it names OpenRouter authors
(`openai`, `anthropic`, and the router authors `openrouter`, `typesafe` and
`@` presets) instead of refusing OpenRouter as a whole, which keeps the other
vendors selectable at the cost of not recognizing a router under an unknown
author (ADR-0235). The judge package is left byte-identical.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the
  refusal code `subscription_only` follows the existing selection refusal
  codes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the change deletes outbound HTTP: the
  key read and the credits read of the router gate, and the relay that
  forwarded OpenCode's requests to `https://openrouter.ai/api/v1`. After it,
  outside the judge no Go source names OpenRouter's host. The judge's
  requests, recipients and keys are unchanged. No test or Verification
  reaches the network. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0235 governs the whole design:
  "Configuration validation refuses an `opencode` or `opencode-custom` Agent
  Selection whose model is under the retired `roundfix-openrouter` provider",
  and "the ACPX runner applies the same check before it prepares a session".
  ADR-0027: "recognized deprecated keys are ignored with one stderr warning
  naming the replacement"; the retired key has no replacement, so its warning
  says to remove it. ADR-0201: "Each request appends one line to a Judge Log
  in Roundfix Home"; the judge keeps writing and summing it. ADR-0231: "It
  must be a finite number greater than zero"; the ceiling stays. ADR-0050:
  "once Agent work begins, Roundfix fails the Work Item instead of switching
  models over potentially modified state", and ADR-0114: "A Fallback Selection
  may switch ACP Runtime automatically only while Agent work has not begun";
  the runner's refusal is a selection failure before Agent work. ADR-0049:
  "each present profile replaces the complete lower-precedence profile"; a
  refused selection fails its profile's validation whole. ADR-0187 and
  ADR-0189 govern the skill edit and its version. ADR-0184: "A TechSpec now
  declares numbered Surface Transcripts"; two are declared below. The authored
  QA gate follows ADR-0080: "QA verdicts distinguish environment-blocked
  rows", and ADR-0091: "required to be terminal and to depend on every leaf";
  ADR-0096, ADR-0097 and ADR-0167 bind its machine stage, its row carry and
  its pre-PR Pull Request row, ADR-0104: "Every Spec therefore rests at least
  one named" acceptance row on outside evidence, ADR-0194, ADR-0195 and
  ADR-0210 bind what a QA row records, when it is observed again and its
  evidence snapshot, and ADR-0182 settles each Task on the facts its gate
  checks. ADR-0093, ADR-0117, ADR-0156, ADR-0168, ADR-0176 and ADR-0183 check
  this Spec's consistency by citation and receipt, and ADR-0233 regenerates
  the raised skill version at merge. ADR-0069 cites ADR-0050 but decides the
  Baseline semantic analysis, ADR-0180 cites ADR-0069 but decides the dated
  Recommended Profile, which this Spec leaves unchanged, ADR-0181 cites
  ADR-0180 but decides where a configuration is compared with that profile,
  ADR-0209 cites ADR-0200 but decides how the judge suggests sources and
  ADR-0208 cites ADR-0209 but decides how sources are grouped, ADR-0217 cites
  ADR-0049 but decides how a Cursor selection is named, and ADR-0229 cites
  ADR-0167 but decides how an operator archive resumes a park; this Spec
  changes none of them, so none applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — task_03 edits the Roundfix Skill, whose
  canonical files and `SKILL.md` mirror are Governed Paths, and the comment
  of `.roundfixrc.yml`; express maintainer authorization: "considere
  autorizado a ajustar todas as skills se necessário" and "Autorizar os dois"
  (2026-09-30), with the retirement decided on 2026-10-05; bounded files:
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `skills/roundfix/SKILL.md`, `.roundfixrc.yml`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0230-retire-the-jev-router/_authorization.md`.

## System Architecture

No new package, command or flag.

| Component | Where | Change |
| --- | --- | --- |
| Subscription rule | `internal/config/profiles.go`, beside `normalizeSelection` | New `CheckSubscriptionRule`, `SubscriptionRule`, `SubscriptionOnlyReason` |
| Selection validation | `normalizeSelection`, `ResolveProfile`, `validateProfiles`, `validateJevRouterProfileSource` in `internal/config/profiles.go`; `applyConfigContent` in `internal/config/config.go` | Calls the rule; the Project-Config-only router rules, the router's empty-effort rule and the legacy-runtimes router check are removed |
| Retired key | `deprecatedConfigKeys`, `warn`, `Jev`, `jevOverlay`, `applyOverlay` in `internal/config/config.go` | `jev.router_min_credit_usd` becomes a deprecated key; the field and overlay go |
| Router selection | `internal/config/jev_router.go`, `internal/agent/jev_router.go` | Deleted |
| Runner | `ACPXRunner`, `routerSession`, `RunPrompt`, `codexEnvForSession`, `clearSessionState`, `validateRuntimeSelection` in `internal/agent/acpx_runner.go`; `ExecuteResult` in `internal/agent/agent.go` | Relay, inline provider and router observation removed; the rule is checked before a session or prompt |
| Daemon | `JevRouterGate`, `jevRouterGate`, `Dependencies`, `NewEngine` in `internal/daemon/engine.go`; `runPrepared`, `selectionReasonCode` in `internal/daemon/agent_session_owner.go` | Gate and Judge Log append removed; `subscription_only` is a named reason |
| Engine wiring | `newResolveEngine` in `internal/cli/cli.go`; `executeImplementCycle` in `internal/cli/implement.go` | The ceiling and floor parameters the gate needed are removed |
| Router package | `internal/jevrouter/` | Deleted |
| Record | configuration guide, Roundfix Skill `runtime` reference, model reference, `.roundfixrc.yml` comment | The router passages give way to the rule |

The judge (`internal/judge`, `internal/cli/spec_judge.go`) is not touched.
Its exported `ReadMonth`, `AppendLogLine` and `LogLine` lose their only
callers outside the package; they stay, so the judge package stays
byte-identical.

## Implementation Design

### Interfaces

```go
// internal/config/profiles.go
const SubscriptionRule = "OpenAI and Anthropic models run only through the codex and claude subscriptions"
const SubscriptionOnlyReason = "subscription_only"

// CheckSubscriptionRule returns nil, or an error naming field and model
// and ending with SubscriptionRule (ADR-0235).
func CheckSubscriptionRule(field, runtime, model string) error
```

```text
1. The rule applies only when the runtime, trimmed and without a "-custom" suffix, is "opencode".
2. The model is trimmed and compared lowercased; its provider is the text before the first "/".
3. Provider "roundfix-openrouter": refused with the retired-router message, whatever follows.
4. Provider "openrouter": the author is the next "/"-separated segment, with a leading "~" and any ":" suffix removed.
   Refused with the OpenRouter message when the author is "openai", "anthropic", "openrouter" or "typesafe", or starts with "@".
5. Every other runtime, provider and author is allowed.
6. normalizeSelection calls it with field "<path>.model" right after the model is known non-empty, before any effort rule.
7. The runner calls it with field "agent model" in validateRuntimeSelection and at the start of RunPrompt, before any acpx process starts,
   and wraps a refusal as SelectionFailureError{Runtime: <runtime id>, Reason: "subscription_only: <message>"}.
8. selectionReasonCode returns "subscription_only" for that reason and no longer knows the router reasons.
9. No package other than internal/judge and the judge's CLI test names OpenRouter's API host.
```

### Data Models

`Jev` keeps `MonthlyCeilingUSD` and loses `RouterMinCreditUSD`; `jevOverlay`
loses `router_min_credit_usd`. `deprecatedConfigKey` gains a way to state
that a key has no replacement. `ExecuteResult` loses `Router`; `daemon.
Dependencies` loses `JevRouter`, `JevMonthlyCeilingUSD` and
`JevRouterMinCreditUSD`. The Judge Log schema is unchanged; no new
`router-prompt` line is written.

### API Contracts

1. API Contract: the retired-router refusal is
   `<field> "<model>" names the retired Jev Router; OpenAI and Anthropic models run only through the codex and claude subscriptions`.
2. API Contract: the OpenRouter refusal is
   `<field> "<model>" can reach OpenAI or Anthropic models through OpenRouter; OpenAI and Anthropic models run only through the codex and claude subscriptions`.
   Both messages print the model as configured, trimmed.
3. API Contract: configuration loading fails with API Contract 1 or 2 for a
   refused selection in Project Config, User Config (profiles or the legacy
   `runtimes` defaults once they become profiles) or a one-Run override;
   `roundfix profiles configure` refuses the fragment with exit 2 and writes
   nothing.
4. API Contract: a refused selection that reaches the runner fails before
   any acpx process starts with reason code `subscription_only`; before Agent
   work the Fallback Chain takes the Task and the fallback Run Event's
   `reason_code` is `subscription_only`.
5. API Contract: `jev.router_min_credit_usd` in User Config or Project
   Config, with any value, is removed before validation with the one warning
   `config: jev.router_min_credit_usd is deprecated and ignored; the Jev Router was retired, so remove it`.
   The existing deprecated-key warnings keep their text.
6. API Contract: Roundfix no longer sets `OPENCODE_CONFIG_CONTENT`, opens a
   relay, reads OpenRouter's key or credits endpoints, or appends a
   `router-prompt` line; `roundfix spec judge` output and its Judge Log lines
   are unchanged.

### Surface Transcripts

1. Surface Transcript: `roundfix profiles configure` refuses a fragment whose
   `docs` preferred selection is `opencode` with
   `openrouter/anthropic/claude-opus-5.5`, in a disposable repository with a
   disposable Home.

   ```transcript
   $ roundfix profiles configure --scope project --file <fragment>
   stdout:
   stderr:
   roundfix: profiles failed: read profiles file "<fragment>": profiles.docs.preferred.model "openrouter/anthropic/claude-opus-5.5" can reach OpenAI or Anthropic models through OpenRouter; OpenAI and Anthropic models run only through the codex and claude subscriptions
   Run 'roundfix profiles --help' for usage.
   exit: 2
   ```

2. Surface Transcript: `roundfix profiles show` with a disposable Home whose
   User Config holds only `jev.router_min_credit_usd: 0`.

   ```transcript
   $ roundfix profiles show
   stdout:
   ...
   stderr:
   config: jev.router_min_credit_usd is deprecated and ignored; the Jev Router was retired, so remove it
   exit: 0
   ```

## Coverage Map

- Goal 1 → `CheckSubscriptionRule`, `normalizeSelection`, `validateRuntimeSelection`, `RunPrompt`
- Goal 2 → deletion of `internal/jevrouter`, the runner and daemon wiring, the config rules; the record
- Goal 3 → `deprecatedConfigKeys`, `warn`
- Goal 4 → judge package untouched; Invariant 9
- Story 1 → `CheckSubscriptionRule`, `normalizeSelection`
- Story 2 → API Contracts 1 and 2
- Story 3 → API Contract 5
- Story 4 → API Contract 6, judge package untouched
- Story 5 → configuration guide, `runtime` reference, model reference, `.roundfixrc.yml` comment
- Core Feature 1 → runner, daemon, engine wiring, router package deletion
- Core Feature 2 → `CheckSubscriptionRule`, `normalizeSelection`, `profiles configure`
- Core Feature 3 → `validateRuntimeSelection`, `RunPrompt`, `selectionReasonCode`
- Core Feature 4 → `deprecatedConfigKeys`, `warn`
- Core Feature 5 → judge package untouched
- Core Feature 6 → the record
- Success Metric 1 → `normalizeSelection`, `profiles configure`
- Success Metric 2 → `validateRuntimeSelection`, `RunPrompt`, `selectionReasonCode`
- Success Metric 3 → runner and daemon deletion, Invariant 9
- Success Metric 4 → `deprecatedConfigKeys`, `warn`
- Success Metric 5 → the record

## Integration Points

- OpenCode: Roundfix no longer passes it a provider definition; an inherited
  `OPENCODE_CONFIG_CONTENT` in the operator's environment passes through
  unchanged, as it already does for every non-router session.
- OpenRouter: no Roundfix call outside the judge.

## Testing Approach

- `CheckSubscriptionRule` is a table test in the config package covering
  Invariants 1 to 5, including `opencode-custom`, upper-case and `~` and `:`
  forms, `openrouter/openrouter/auto`, `openrouter/typesafe/jev-router`,
  `openrouter/@preset/x`, other OpenRouter vendors, OpenCode's other
  providers and the other runtimes.
- Configuration scopes are tested through `config.Load` with disposable
  Homes and repositories, `ResolveProfile` with an override, and the
  `profiles configure` command through the CLI test harness (`runCLI`), which
  asserts the file is unchanged.
- The runner is tested through the existing fake acpx harness: a refused
  selection fails `RunPrompt` and `PrepareSession` with the reason and leaves
  no acpx invocation; an allowed OpenRouter model still runs.
- The daemon is tested through the existing fallback boundary owner: a
  `subscription_only` selection failure before work activates the fallback
  and the Run Event names the reason.
- The deleted tests are the router's own tests; no other test changes
  expectations except the call sites that remove the router fixture hooks.

## Build Order

1. State the rule in the guides, the Roundfix Skill, the model reference and
   `.roundfixrc.yml`, raising the skill version, so the CLI change ships with
   its guide.
2. Retire the routed runtime and add the subscription rule's predicate with
   the runner check and the daemon reason (`internal/jevrouter`,
   `internal/agent`, `internal/daemon`, `internal/cli` engine wiring, and
   the predicate in `internal/config/profiles.go`) (depends on: 1).
3. Apply the rule in configuration validation, remove the router's
   configuration rules and field, and deprecate the floor key (depends on: 2).
4. Final QA gate (depends on: 1, 2, 3).

## Risks & Considerations

- A Run started before the upgrade that resumes with the router selection is
  refused before its next prompt (Invariant 7) and falls back when no Agent
  work had begun; otherwise its Work Item fails, as ADR-0050 requires.
- Third-party routers under other authors and operator-defined OpenCode
  providers are not recognized (ADR-0235); the operator can restrict models
  on the OpenRouter account.
- Past `router-prompt` lines stay in the Judge Log; the judge reads every
  line, so October 2026's ceiling still counts them.
- The Roundfix Skill version is raised by one patch from the tree the Task
  starts from (rule of ADR-0189).

## Vocabulary Contract

- emits: `internal/config/profiles.go`
  pattern: `subscription_only`
  documented-in: `.agents/skills/roundfix/references/runtime.md`
- emits: `internal/config/profiles.go`
  pattern: `run only through the codex and claude subscriptions`
  documented-in: `docs/user-guide/configuration.md`
- emits: `internal/config/config.go`
  pattern: `the Jev Router was retired, so remove it`
  documented-in: `docs/user-guide/configuration.md`

No new glossary term is adopted. "Subscription rule" is used in its plain
sense beside **Agent Selection**, **Fallback Chain**, **Judge Log** and
**Agent work**.

## Decisions

- Retire, do not restrict, the router. See ADR-0235.
- Refuse by OpenRouter author, check at configuration and at the runner.
  See ADR-0235.
- The retired floor key is a deprecated key without a replacement (ADR-0027).
- Keep the judge package byte-identical, including helpers only the router
  called.
- ADR-0218 and ADR-0234 were retired to `docs/history/adr/`, and ADR-0231
  noted, at authoring, because ADR-0235 is the decision this Spec carries out.
