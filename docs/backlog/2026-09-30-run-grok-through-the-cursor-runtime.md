---
type: feat
status: open
created: 2026-09-30
spec: null
reason: null
---

# Run Grok through a Cursor runtime

## Opportunity

Roundfix selects an Agent per work category from three ACP Runtimes: `codex`, `claude` and `opencode`. The maintainer also has Cursor, whose agent CLI can run Grok models. Roundfix cannot route a Task to it, so a fourth provider that might suit some categories is unreachable from a Run.

## Value

A fourth Runtime widens every Fallback Chain with a provider that fails independently of the other three, and lets Grok be measured on real Tasks against the current defaults. The hypothesis is that Grok through Cursor is a useful `preferred` or fallback selection for at least one category, at a cost the Cursor plan already covers.

## Shape

This is intent only. It is sequenced after the work already in flight (Wave 6, the Jev integration, and the model catalog refresh).

- Add `cursor` as an ACP Runtime next to the three in `internal/agent/agent.go`. Its adapter command goes in `internal/agent/acpx_runner.go` and its models in the picker catalog. The pinned `acpx` (0.19.3) already lists `cursor` and `grok-build` agents, and `cursor-agent` 2026.07.09 is installed on the maintainer's machine.
- Prove each configured `cursor` selection the way the other Runtimes are proved, before any Run depends on it, and give Doctor a readiness check for the CLI and its login.
- Before a Spec, measure:
  - which Grok model IDs the adapter advertises;
  - whether a Daemon-run session can approve tools without a person;
  - how the plan bills a long Run;
  - one recent Spec's Tasks run on Grok, compared with the current default for time, retries and QA outcome.
- Decide whether `grok-build`, going to xAI directly, is a better route than Cursor for the same models.

Blockers to disprove first, reported on Cursor's forum for its ACP mode (read on 2026-09-30):

- Runtime model switching from a third-party ACP client is unsupported, so the model must be pinned with `cursor-agent --model <id> acp`. That conflicts with how Roundfix proves a selection before a Run.
- Reasoning effort and the non-Fast option are not honored under ACP, so Fast billing can happen silently.
- `session/load` fails right after `session/new`, which breaks the named sessions Roundfix uses.
- Web search prompts even in unrestricted mode.

Other facts read the same day:

- The installed `cursor-agent` is not logged in and predates Cursor's July fixes.
- Cursor encodes effort and speed in the model slug (for example `cursor-grok-4.6-high-fast`); no Grok 4.7 slug was found.
- xAI's own CLI, which `acpx` exposes as `grok-build`, takes an API key and an always-approve flag. No real `grok` binary was found on this machine.
- A third candidate Runtime that speaks ACP and logs in with a Grok subscription exists, but it is experimental.
- Published results for Grok 4.7 place it below the current defaults on a neutral harness (Terminal-Bench 4.0: 25.8, against 59.6 for Opus 5.5), and it used about three times the output tokens per task.

Evidence: the local facts were read on 2026-09-30 with `cursor-agent --version` and `acpx --help`. The external facts are in the Secondbrain research entry `inbox/secondbrain/2026-09-30-modelos-de-codificacao-ago-set-2026-e-runtimes-cursor-grok.md`.
