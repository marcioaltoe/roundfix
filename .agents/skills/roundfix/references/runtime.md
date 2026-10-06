## acpx dependency

Roundfix drives ACP Runtimes through acpx `0.12.0` or newer. Node.js 22.13 or
newer with npm/npx is a prerequisite. Prefer the Setup Command after
installing Roundfix; it verifies Node, installs the minimum tested acpx only
when acpx is missing or older, accepts newer versions without downgrading,
proves the effective adapter identities, proves the
generated Agent Selection Profiles, offers authorized local adapter migration,
and offers User Config and Project Config creation. Review work uses owned review Agent Sessions, Spec
Tasks use per-Task Agent Sessions named `roundfix-<run-id>-<task_id>` in their
Task Worktrees, and QA uses its own Agent Session after Tasks settle.

Use the Doctor Command, `roundfix doctor`, to diagnose Run readiness without
installing dependencies, writing config, or changing files. The `environment:`
line lists `spec judge keys` by name in preference order:
`ROUNDFIX_OPENROUTER_JUDGE_API_KEY`, `ROUNDFIX_OPENROUTER_API_KEY`,
`ROUNDFIX_TYPESAFE_API_KEY`; then `implementation keys`:
`ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY`, `ROUNDFIX_OPENROUTER_API_KEY`.
Each entry is `set` or `not set`, never a key value; both stages fall back
to the shared key when their stage key is unset. Doctor runs the
shared Node.js, minimum-supported acpx, effective adapters, configured Agent
Selection Profiles, Repository Skill Set, process residue, storage check, and
codex runtime hygiene checks and prints one line per check with status `ok`,
`failed`, or `skipped`; residue and storage can also report `found` or
`partial`. Adapter
Readiness requires the effective Codex command to prove official
`@agentclientprotocol/codex-acp` lineage at version `2.0.1` or newer and the
effective Claude command to prove official
`@agentclientprotocol/claude-agent-acp` lineage at version `0.84.0` or newer;
executable presence and a matching name are not proof. The `profiles:` line is
the selection authority: it exact-proves every distinct Preferred Selection
and fallback through disposable ACP Sessions and reports affected category
references plus one deterministic next action. A proof whose setup times out
is retried once. A second timeout is classified `temporary`; rerun the command
when load drops because the configured profile was not shown to be wrong.
Doctor has no separate legacy `agent:` or `model:` authority. Failed checks
include `next: <action>` when Roundfix knows the remediation.

When `NODE_OPTIONS` names a preload whose file no longer exists, Roundfix drops
that preload from the environment it gives an agent process. It checks absolute
preload paths from the last `NODE_OPTIONS` entry, keeps existing paths and
package names, and reports each dropped path once per process on standard error:

```text
roundfix: notice: NODE_OPTIONS preload "<path>" does not exist; Roundfix left it out of the agent environment
```

The `pre-pr-review:` line reports the resolved pre-Pull-Request review provider
and the configuration layer that supplied it. An explicit `none` reports that
review is disabled by configuration. This check reads policy only: it invokes
no provider and mutates nothing. Use `roundfix review` to run the configured
policy over the current candidate; that command is the enforcement point for
the review policy, while publication and merge gating remain separate.

Profile readiness covers every Agent Work Category the effective configuration
defines — the five required categories plus each optional category
(`data`, `infra`, `docs`, `test`, `chore`) a profile actually declares. A
category that resolves only by inheriting `general` adds no distinct tuple and
is not enumerated. The `adapter:` line follows the same scope, so it names every
ACP Runtime the configured tuples reference, including one only an optional
category selects. A configured profile that fails therefore fails the
`profiles:` line instead of leaving it `ok`.

The `opencode` runtime accepts a non-empty reasoning effort with the
`runtime_deferred` encoding. OpenCode advertises effort per model only after an
Agent Session's first prompt, so token-free Preflight cannot apply the effort:
it proves that the model is advertised and current and that the requested
effort is among the values that model advertises. Before any work turn, the Run
ensures the Agent Session with its model, sends one minimal warm-up prompt to
raise the queue owner, applies the requested effort, and observes the effective
value. The Run therefore proves the effort applied. An empty effort remains
`runtime_managed`: Roundfix declines to assign the advertised control and the
Agent Model opens at its own value.

### Cursor

`cursor` is opt-in and reached as `cursor-agent acp` through acpx. A Cursor
selection names the advertised model value verbatim with
`reasoning_effort: ""`; a non-empty effort is refused. The login is the
maintainer's own: Roundfix checks it with `cursor-agent status` and never
performs it. A missing login is reported as `cursor_login_required`.

The blocking `skills:` line runs after, and independently from, `profiles:`.
For each Roundfix-owned skill, the minimum version is the version of that
skill the running binary carries, and Doctor compares it with the version
declared by the installed `SKILL.md`.
Readiness is this version comparison, not a content match. The three states are:

- `satisfies` — the declared version is at or above the declared minimum, so
  readiness passes.
- `below minimum` — the declared version is below the minimum, so readiness
  blocks. The failure names the skill, the minimum, the version found, and the
  `roundfix skills install --target project` upgrade path.
- `unversioned or unresolvable` — the skill has no comparable declaration or
  its version source is unreachable. Doctor renders this state as
  `unversioned`, distinct from both `satisfies` and `below minimum`; an
  unreachable source is never reported as a missing skill.

Roundfix never applies this owned-skill version comparison to third-party
skills or holds them to a version Roundfix invented for them. The required
external set comes from the repository's Setup Manifest at
`docs/agents/setup-context.json`: Doctor unions the `requiredSkills` declared
by its selected modules, removes Roundfix-owned names, and checks each
remaining skill against its `computedHash` in `skills-lock.json`. The external
count therefore follows the repository's selected modules instead of a fixed
Roundfix development list.

Outside a Git repository, Doctor does not inspect the Repository Skill Set and
prints
`skills: failed (Repository Skill Set readiness requires a Git repository; next: run roundfix doctor from a Git repository)`.
When the Setup Manifest is absent or unreadable, Doctor still checks the
running binary's owned set, requires zero external skills, fails the `skills:`
line with `Setup Manifest is absent or unreadable`, and includes
`roundfix baseline` as the Baseline-adoption next action.

Surface a failed `skills:` line and its printed `next:` remediation before
work continues. Owned failures print
`roundfix skills install --target project`. Each named missing or outdated
external skill prints its own
`bunx skills add marcioaltoe/skills@<skill>` command. When an external failure
does not identify specific skill names, Doctor instead prints
`bunx skills experimental_install && bunx skills update -p -y`; a mixed named
failure prints the owned action followed by the skill-scoped external actions.
Doctor is diagnosis-only: it never runs these commands, accesses the network,
installs or updates skills, or writes repository state. Apply remediation only
after explicit workflow authorization, then rerun Doctor.
On macOS, the codex hygiene check resolves `CODEX_PATH` first and then `codex`
on `PATH`, inspects the `com.apple.quarantine` attribute (the real XProtect
trigger), and verifies the binary's code signature (not `spctl --assess`, which
rejects any signed CLI that is not a notarized app). A quarantined or
improperly-signed codex fails with the next action to reinstall codex with the official curl installer
into `~/.local/bin`, then set `CODEX_PATH` to that binary. On non-Darwin
platforms the codex check is `skipped` and never fails the command by itself.

When Roundfix launches codex through `codex-acp` on macOS, it uses the same
configured-path-then-`PATH` resolution and passes a verified-clean codex to
acpx through `CODEX_PATH`. If no clean codex is available, Roundfix surfaces
the hygiene failure instead of silently spawning a known unsafe binary.

Known constraint: acpx `0.12.0` has a hard 10 MiB queue-owner per-message
buffer in `src/cli/queue/ipc.ts`, bundled in the installed package at
`dist/output-CjdF5rHk.js`, with no CLI, config, or environment override found.
Large docs-task payloads, especially turns that print or return large
skill/docs file content, can trigger `-32603 Message buffer exceeded 10485760
bytes`. Treat this as an upstream acpx limit: keep payloads smaller when
practical, and rely on the ADR-0020 classification and the Settle Command for
completed Task work preserved in a spec Run Worktree.

ADR-0020 classification: when acpx has delivered a valid
`session/prompt` result for a Batch before a later nonzero acpx exit, Roundfix
journals the anomaly with the stderr tail and proceeds to the Daemon's
verification. Without that parsed result, the nonzero exit remains a Batch
failure. If acpx rejects the selected Agent Model with its not-advertised
stderr, Roundfix reports the terminal reason as
`Agent Model "<model>" not advertised by runtime "<runtime>"; advertised: <list>`
in Work Item reasons, Run Events, and final report reason lines instead of a
generic `agent/protocol error`. Verification remains the only gate for settling
and committing.

### OpenAI and Anthropic subscription rule

OpenAI and Anthropic models run only through the codex and claude
subscriptions. For `opencode` and `opencode-custom`, Roundfix refuses
`roundfix-openrouter/typesafe/jev-router` and other selections under the
retired `roundfix-openrouter` provider. Under the `openrouter` provider it
refuses authors `openai` and `anthropic`, router authors `openrouter` and
`typesafe`, and `@` presets. The refusal reason is `subscription_only`; a
refusal before Agent work activates the configured fallback. Other OpenRouter
models stay selectable. An open model through OpenCode reads
`ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY` first and the shared
`ROUNDFIX_OPENROUTER_API_KEY` second; empty values count as unset and the
generic `OPENROUTER_API_KEY` is never read.
`jev.router_min_credit_usd` is a deprecated key.

### Light implementation tier

At dispatch, a Task is `light` when `complexity` is `low`, `type` is not `qa`,
and it declares no Governed Path to create, change, or delete. All other Tasks,
QA gates, and review Batches are `standard`. A one-Run Agent Selection override
turns the tier off for that Run.

A light Task derives a profile from the User Config `openrouter.light_models`
list on the `opencode` runtime, with no reasoning effort and
`profile_source` `light-tier`; its category Preferred Selection and Fallback
Chain follow the light models. An empty list disables the tier. The default
model is `deepseek/deepseek-v4.1-flash`, and each configured light model must
pass the subscription rule. The light session receives the key variable name
from Roundfix's key helper and never the key value or generic
`OPENROUTER_API_KEY`.

A skipped light Task emits phase `light_tier_skipped` with reason code
`key_missing`, `spend_unreadable`, or `ceiling_reached`. A light Task that
fails Verification escalates once. The escalation emits phase
`light_tier_escalated`, naming the light model and failed commands, and runs a
new session on the category's Preferred Selection or its pre-work fallback.
The light prompt, read files, and diagnostics are sent to OpenRouter and its
provider; the Light Spend Log records the reported cost without credentials.
