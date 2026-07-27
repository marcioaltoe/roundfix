---
status: pending
created_at: 2026-07-27
updated_at: 2026-07-27
---

# Agent Selection — standardize on the official Claude adapter and stop reading `[...]` as reasoning effort (2026-07-27)

The maintainer directs Roundfix to drop every reference to
`@zed-industries/claude-code-acp` and support only
`@agentclientprotocol/claude-agent-acp`, including Doctor and Setup. This
finding records the empirical verification that the replacement works, the
exact code sites that still name the deprecated package, and one **new
blocking defect** that only becomes visible after the migration: Roundfix
misreads the replacement's advertised model identifiers.

This continues
[2026-07-26 — the configured Claude ACP adapter is deprecated](2026-07-26-claude-adapter-configoptions-migration.md),
which proposed the migration while `claude-agent-acp` was still uninstalled.
That finding's premise that `claude-opus-5` "is advertised and gets selected
successfully" describes the deprecated adapter only, and no longer holds — see
section 3.

No repository change was made: an Implement Run for Spec
`0038-terminal-run-worktree-reconciliation` was active, and editing tracked
files during a Run strands it in Integration Pending.

## Session evidence

- Repository `/Users/marcio/dev/roundfix` at `8ec92ad`, macOS 25.5.0.
- `acpx 0.12.1`. `codex-acp` resolves to `@agentclientprotocol/codex-acp@1.1.7`.
- `@agentclientprotocol/claude-agent-acp@0.63.0` installed globally during this
  session; `~/.acpx/config.json` `agents.claude.command` repointed from
  `claude-code-acp` to `claude-agent-acp`. Prior config backed up.
- Every proof below ran through
  `roundfix profiles configure --scope project --file <fragment> --dry-run --json`,
  which proves an exact tuple without writing config or creating a Run. From
  inside a Claude Code session the probe needs
  `env -u CLAUDECODE -u CLAUDE_CODE_ENTRYPOINT -u CLAUDE_CODE_SSE_PORT`.

## 1. The replacement adapter resolves the capability failure

With the deprecated adapter, every Claude tuple failed identically —
including `reasoning_effort: ""`, which rules out model and effort as causes:

```text
classification: capability_evidence_invalid
adapter error: capability evidence invalid: missing_config_options
```

With `@agentclientprotocol/claude-agent-acp@0.63.0`, the same probe reaches
model evaluation and reports the controls it advertises:

```text
advertised Agent Models: default, opus[1m], claude-fable-5[1m], sonnet, haiku
advertised reasoning efforts: default, low, medium, high, xhigh, max
```

`claude / sonnet / xhigh` proves clean. Claude is a usable runtime again, and
the migration is confirmed as the correct fix.

## 2. Roundfix reads `[...]` as reasoning effort, so `opus[1m]` is unselectable

`parseModelCapability` (`internal/agent/selection_capabilities.go:469`) treats
any advertised identifier shaped `name[suffix]` as Roundfix's model-variant
encoding `canonical[effort]`. The replacement adapter uses that bracket for
the **1M context window**, not for effort. The collision is proven:

| Requested selection | Result |
| --- | --- |
| `claude / opus[1m] / xhigh` | fails `model_not_advertised` — although `opus[1m]` is in the advertised list |
| `claude / opus / 1m` | **passes** — Roundfix parsed canonical `opus` + effort `1m` |
| `claude / sonnet / xhigh` | passes — no brackets, so the independent reasoning control is used |

Two consequences. A user cannot select a model by the identifier the adapter
actually advertises. And the only route to Opus today is to write a context
window into `reasoning_effort`, which forces the model-variant encoding and
therefore **discards explicit effort control** — the Task runs at the adapter's
default effort, not the configured one.

The adapter advertises an independent reasoning select *and* bracketed model
identifiers at the same time. That combination is what Roundfix does not model:
when an independent reasoning control exists, the model identifier must be
treated as opaque instead of being parsed for an embedded effort.

## 3. `claude-opus-5` is not advertised by the replacement adapter

The advertised set is `default, opus[1m], claude-fable-5[1m], sonnet, haiku`.
Both `.roundfixrc.yml` files that pin `frontend` — this repository's and
`fluxus`'s — request `claude / claude-opus-5 / xhigh`, which no adapter
lineage advertises. Those profiles stay broken after the migration until they
are repointed, so any Spec with a `frontend` Task is blocked in both
repositories. Repointing is a supervisor-owned boundary commit, not Run work.

## 4. Code sites that still name the deprecated package

Codex has a complete lineage contract; Claude has none — only a bare command
name, which is why a same-named binary of either lineage passes today.

| Location | Current value | Problem |
| --- | --- | --- |
| `internal/agent/acpx_runner.go:54` | `"claude": "claude-code-acp"` | default adapter command is the deprecated lineage |
| `internal/agent/acpx_runner.go:59` | `"claude-code-acp": "npm install -g @zed-industries/claude-code-acp"` | install hint installs the deprecated package |
| `internal/agent/acpx_runner.go:60` | `"claude-agent-acp": "npm install -g @zed-industries/claude-agent-acp"` | **wrong scope** — resolves to `@zed-industries/claude-agent-acp@0.23.1`, not the official `@agentclientprotocol/claude-agent-acp@0.63.0` |
| `internal/agent/acpx_runner.go:27-30` | `CodexAdapterPackage`, `PinnedCodexAdapterVersion`, `legacyCodexAdapterPackage`, `defaultCodexAdapterCommand` | Codex-only; no Claude equivalent exists |
| `docs/adr/0040-reasoning-effort-is-assigned-only-when-configured.md:7` | "claude-code-acp does not implement the reasoning config option at all" | premise is now false for the official adapter |

Test fixtures naming the deprecated command:
`internal/agent/acpx_runner_test.go:175,2693` and
`internal/cli/implement_test.go:4163`.

Both packages named `claude-agent-acp` exist on npm, so a lineage check cannot
rely on the command name alone: `@zed-industries/claude-agent-acp@0.23.1` and
`@agentclientprotocol/claude-agent-acp@0.63.0` are different lineages.

## Implementation instructions

1. Introduce the Claude counterparts of the Codex lineage constants —
   official package `@agentclientprotocol/claude-agent-acp`, a pinned minimum
   version, and `@zed-industries/claude-code-acp` plus
   `@zed-industries/claude-agent-acp` as recognized legacy lineages. Make the
   default Claude adapter command the official pinned form, matching how Codex
   resolves `npx -y @agentclientprotocol/codex-acp@<pin>`.
2. Extend Adapter Readiness so Doctor and Setup prove Claude lineage the same
   way they prove Codex: report the effective command, resolved package, and
   version; classify a legacy lineage as a migration with the official install
   action; and never accept a matching executable name as proof.
3. Fix the wrong-scope install command at `acpx_runner.go:60`.
4. Stop parsing `[...]` as an embedded reasoning effort when the adapter
   advertises an independent reasoning control; treat the model identifier as
   opaque so `opus[1m]` selects with an explicit effort. Keep the existing
   variant encoding for adapters that genuinely use it, and cover both shapes
   with tests.
5. Supersede ADR-0040's claim about Claude reasoning support, recording that
   the official adapter advertises the reasoning control.
6. Update the operator docs and the Roundfix Skill pair so the Claude adapter
   contract reads like the Codex one. The Skill pair edit changes its
   `contentDigest` and requires the derived baseline/catalog digest pins to be
   propagated under express authorization — see the Tooling authority entries
   of archived Spec `0037-terminal-outcome-integrity`.

## Suggested acceptance checks

- Doctor on a machine whose `claude` command resolves to either
  `@zed-industries` lineage fails with the official install action and does not
  accept the name as proof.
- Doctor with `@agentclientprotocol/claude-agent-acp` at or above the pin
  reports `adapter: ok` naming package and version for Claude as it does for
  Codex.
- `claude / opus[1m] / xhigh` proves clean and applies `xhigh` through the
  independent reasoning control.
- `claude / opus / 1m` no longer resolves an effort of `1m`.
- An adapter that really uses the `canonical[effort]` variant encoding keeps
  working.
- No source, test, or documentation path outside `docs/findings/` and archived
  Specs mentions `@zed-industries/claude-code-acp`.

## What worked — keep

- The dry-run tuple proof
  (`profiles configure --dry-run`) diagnoses selection problems without writing
  config or starting a Run, and its classifications (`missing_config_options`
  versus `model_not_advertised`) were precise enough to separate an adapter
  capability gap from a parsing defect.
- Roundfix's error message listing advertised models and efforts is what made
  the bracket collision visible; keep that detail in selection failures.
