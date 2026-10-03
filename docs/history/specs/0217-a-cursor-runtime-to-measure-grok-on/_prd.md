---
spec: 0217-a-cursor-runtime-to-measure-grok-on
status: archived
created: 2026-10-01
surfaces: [backend, cli, docs]
archived: "2026-10-03"
source_slug: 0217-a-cursor-runtime-to-measure-grok-on
---


# A Cursor runtime to measure Grok on

Roundfix selects an Agent for each work category from three ACP Runtimes:
`codex`, `claude` and `opencode`. The maintainer also has Cursor, whose plan
includes Grok models, and its agent CLI, `cursor-agent`, is installed on the
maintainer's machine. No Agent Selection can name it, so Grok cannot be tried
on a real Task, and a fourth provider that fails independently of the other
three cannot join a Fallback Chain.

The adopted Backlog Entry
([references/2026-09-30-run-grok-through-the-cursor-runtime.md](references/2026-09-30-run-grok-through-the-cursor-runtime.md))
asked for measurement before a Spec. The authoring session measured what can
be measured without the maintainer's Cursor login and without sending a
prompt, on 2026-10-01:

- `cursor-agent` is version `2026.07.09-a3815c0`. Its hidden `acp`
  subcommand prints `Start the Cursor Agent as an ACP (Agent Client Protocol)
  server`. Answering an ACP `initialize`, it advertises `loadSession: true`,
  session listing, image prompts, and one authentication method,
  `cursor_login`.
- The installed acpx `0.19.4` maps its built-in `cursor` agent to
  `cursor-agent acp`, and its `grok-build` agent to `grok agent stdio`. No
  `grok` binary is installed.
- The CLI is not logged in. `cursor-agent status` prints `Not logged in`;
  `cursor-agent models` and `--list-models` refuse with `Authentication
  required`; and an ACP `session/new` is refused with JSON-RPC error `-32000`
  `Authentication required`. So the Grok model IDs Cursor offers this account
  could not be listed.
- Cursor's published ACP sessions advertise a `model` session config option
  whose values carry their parameters in brackets, such as
  `grok-4-20[thinking=true]`,
  `gpt-5.4[reasoning=medium,context=272k,fast=false]` and `default[]`.
  Roundfix's capability projection reads a trailing bracket as a reasoning
  effort and refuses an empty one, so today a whole Cursor session would be
  refused as `malformed_model_value`.

This Spec adds `cursor` as an opt-in ACP Runtime that a profile, a fallback
or a one-Run override can name, proves its selections the way the other
three are proved, and gives the Doctor Command a check of the CLI and its
login. It then measures Grok on real work through it. Nothing becomes a
default: whether a Grok selection should enter the Recommended Profile is a
follow-up Backlog Entry, minted only if the measurement supports it.

## Project Constraints

- Identifier strategy: not applicable — `cursor` is a runtime name in the
  existing Agent Selection vocabulary, and a Cursor model is named by the
  value the adapter advertises; no identifier scheme is added. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — Roundfix makes no HTTP request of
  its own. The `cursor-agent` CLI authenticates with the maintainer's own
  Cursor login. Roundfix only reads whether that login exists, through
  `cursor-agent status`. It never runs `cursor-agent login` or `logout`,
  never passes `--api-key` or `--auth-token`, never reads `CURSOR_API_KEY` or
  `CURSOR_AUTH_TOKEN`, and never prints or stores the account it finds. The
  prompts a Cursor Agent Session sends go to Cursor under the maintainer's
  plan, as prompts on the other runtimes go to their providers. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0217 (this Spec) decides that a
  Cursor selection names the advertised model value verbatim with a
  model-managed reasoning effort, and that the Cursor login is the
  maintainer's own, checked and never performed. Under ADR-0017 Roundfix
  drives ACP Runtimes through acpx, so `cursor` is reached through acpx's
  built-in `cursor` agent. ADR-0037: "passes both explicitly for every Agent
  Session, and fails Preflight Validation when the runtime does not support
  them". ADR-0049: "every configured tuple must be proven through the
  installed ACP adapter rather than accepted from a static compatibility
  assumption". A Cursor selection is one more tuple held to both. ADR-0050
  activates a configured fallback after notification, which is how a refused
  Cursor selection hands over before work begins. ADR-0105 bounds what the
  capability projection keeps, and a Cursor session is read by the same
  projection. ADR-0107: "requiring Exact Agent Selection Proof for every
  distinct tuple", so a Cursor tuple in an optional category is proved too.
  ADR-0108 warms an OpenCode session to apply an effort and stays unchanged,
  because a Cursor effort is never applied separately. ADR-0180 derives
  built-in selections from one dated Recommended Profile, and ADR-0181
  compares a configuration with it only where a person asked; this Spec
  leaves the Recommended Profile unchanged. ADR-0125: "Fixtures are therefore
  compiled once", which binds the fake `cursor-agent` of every test.
  ADR-0187 splits the Roundfix Skill by command, and ADR-0189 ties an owned
  skill's version to its content, so the runtime reference that names
  `cursor` raises the skill's version. ADR-0198 counts the tokens an adapter
  reports by its own scope, so a Cursor prompt's usage is recorded the same
  way. ADR-0208 groups work items that share a context; this Spec and Spec
  0218 were split because a fourth runtime and a routed model on an existing
  runtime do not share one. ADR-0211 drops a missing Node preload from the
  agent environment and holds for the measurement Runs. This Spec's gate is
  bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155, ADR-0156 and
  ADR-0167, and ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check
  its consistency by citation and receipt. ADR-0182 runs Settlement Checks
  before each Task commit, ADR-0178 authorizes a Task commit by its grant,
  and ADR-0166 records undeclared paths; every Task here declares its paths.
  ADR-0114: "A Fallback Selection may switch ACP Runtime automatically only
  while Agent work has not begun", and a Cursor selection refused at the
  login check has begun none. ADR-0184: "A TechSpec described a command's new
  behavior in prose, and each reader rebuilt the exact output from it", so
  the changed command surfaces are stated as Surface Transcripts. ADR-0069 cites ADR-0050 but decides the
  Baseline semantic analysis, ADR-0096 and ADR-0097 cite ADR-0080 but decide
  the gate's machine stage and row carry, ADR-0194, ADR-0195 and ADR-0210
  cite ADR-0097 but decide what a QA row records, when it is observed again
  and its evidence snapshot, ADR-0192 cites ADR-0178 but decides
  conflicts on derived paths, ADR-0199 cites ADR-0198 but decides a queue's
  token ceiling, and ADR-0209 cites ADR-0208 but decides how the judge
  suggests grouping; this Spec changes none of them, so none applies.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 ("considere autorizado a ajustar todas as skills se
  necessário") and the cycle authorization of 2026-10-01 ("Tudo, de A a F"),
  recorded in [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runtime.md`. Sanctioned regeneration:
  `make skills-sync`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`.

## Goals

- A profile, a fallback or a one-Run override can name `cursor` with a model
  value Cursor advertises, and the selection is proved before any Run
  depends on it.
- A Cursor selection states its speed, effort and context in the model value
  it names, so Fast billing or a different effort never happens through a
  setting Roundfix chose silently.
- Roundfix tells the maintainer, before a Run, whether `cursor-agent` is
  installed and logged in, and never logs in itself.
- The Grok model IDs Cursor offers, the behavior of an unattended Cursor
  session, and Grok's result on real Tasks against the current default are on
  record, so a later decision about defaults rests on measurement.

## User Stories

1. As a maintainer, I want to name `cursor` and a Grok model in a profile or
   a fallback, so that a Task can run on Grok without changing any default.
2. As a maintainer, I want `roundfix profiles validate` and the Doctor
   Command to prove a Cursor selection the way they prove the other runtimes,
   so that a Run never discovers a broken Cursor selection after it started.
3. As a maintainer, I want the Doctor Command to tell me when `cursor-agent`
   is not logged in and that I must log in myself, so that I never hand a
   credential to Roundfix.
4. As a maintainer deciding about defaults, I want a record of Grok on real
   Tasks through Cursor next to the current default, so that the decision
   rests on time, retries and outcome, not on published benchmarks alone.

## Core Features

1. **`cursor` is the fourth ACP Runtime.** Profiles, fallbacks and the
   `--agent` one-Run override accept `cursor`, and every message and help
   text that lists the supported runtimes names `codex`, `claude`, `cursor`
   and `opencode`, from one list. acpx runs it through its built-in `cursor`
   agent, `cursor-agent acp`. The legacy `runtimes:` section and the
   Recommended Profile do not gain it.
2. **A Cursor model is the advertised value, verbatim.** A Cursor selection's
   model is exactly a value the session's `model` option advertises, brackets
   and parameters included, such as `grok-4-20[thinking=true]`, and its
   `reasoning_effort` is `""`. A non-empty effort on a Cursor selection is
   refused when the configuration loads, naming the rule. The capability
   projection reads Cursor's bracketed parameters, including the empty
   `default[]`, as part of one model identity, never as an effort. Proof sets
   that exact value and requires the session to report it as current.
3. **The login is checked, never performed.** When an effective profile
   names `cursor`, the Doctor Command's `adapter:` line checks that
   `cursor-agent` resolves and that `cursor-agent status` reports a login.
   Without one it fails with `cursor_login_required` and the next action
   `run cursor-agent login in a terminal yourself`. A Run's preflight refuses
   a Cursor selection with the same classification, so the Fallback Chain
   moves on after notification. No Roundfix code path runs `login`, passes a
   Cursor credential, or prints the account.
4. **The runtime is documented.** The configuration and command guides and
   the Roundfix Skill's runtime reference name `cursor`, its model-value rule,
   its login rule, and that it is opt-in.
5. **Grok through Cursor, measured.** With the maintainer logged in, the
   measurement records the model values a real Cursor session advertises,
   proves one Grok selection, and replays two archived Tasks on Grok and on
   the current default under the protocol in the TechSpec. It records time,
   prompts, Verification repairs, outcome, whether the session ran without a
   person, and the billing mode the model value states. Without a login it
   stops with the named blocker `cursor_login_required`.

## User Experience

A maintainer adds a Cursor fallback to a profile fragment, for example
`{runtime: cursor, model: "grok-4-20[thinking=true]", reasoning_effort: ""}`,
and applies it with `roundfix profiles configure`, which proves it before
writing. When the CLI is not logged in, the proof and `roundfix doctor` both
say so and name the command the maintainer runs: Roundfix never opens a
browser or asks for a key. Help text and errors list four runtimes.

## Declared breaks

- Messages and help texts that list the supported runtimes change from
  `codex, claude, opencode` to `codex, claude, cursor, opencode`, including
  the `profiles.<category>.preferred.runtime` refusal that
  `internal/config/config_test.go` pins.
- A capability projection for the `cursor` runtime no longer refuses an empty
  bracket such as `default[]`; projections for the other runtimes keep
  today's parsing.

## Non-Goals / Out of Scope

- Any change to the Recommended Profile, the built-in profiles or this
  repository's Project Config. A default change is a follow-up Backlog Entry.
- The `grok-build` route to xAI directly, and any API key: no `grok` binary is
  installed, and its key would be a credential Roundfix would hold.
- Logging in, refreshing a login, or reading a Cursor credential.
- Cursor in the legacy `runtimes:` section, the interactive model picker
  catalog, or Setup's acpx overrides.
- Cursor's interactive extension methods (`cursor/ask_question`,
  `cursor/create_plan`) beyond recording whether a session sent one.
- A Grok comparison over many Specs; the measurement is two replayed Tasks,
  recorded as a small sample.

## Success Metrics

1. Success Metric: a profile naming `cursor` and `grok-4-20[thinking=true]`
   with an empty effort loads, and one naming a non-empty effort is refused
   with the model-value rule.
2. Success Metric: the published Cursor `session/new` shape projects without
   an issue, `default[]` included, and a Grok selection assigns its exact
   value as `model_managed`; the same shape under another runtime keeps
   today's result.
3. Success Metric: with a fake `cursor-agent` that reports no login, Doctor's
   `adapter:` line fails with `cursor_login_required` and the login next
   action, and preflight refuses the selection; with one that reports a
   login, both pass; no production source names a Cursor login, key or
   token argument.
4. Success Metric: `roundfix implement --help` lists `codex, claude, cursor,
   opencode`.
5. Success Metric: the measurement record names the Grok model values a real
   session advertised, one proved Grok selection, and for each replayed Task
   and selection the time, prompts, repairs and outcome, or the named login
   blocker.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- **Cursor's ACP documentation** (<https://cursor.com/docs/cli/acp>, read
  2026-10-01): `agent acp` is a hidden ACP server over stdio; a client calls
  `initialize`, `authenticate` with `methodId: "cursor_login"`, then
  `session/new` or `session/load`; a client can pre-authenticate with
  `agent login`; tool approval arrives as `session/request_permission`; and
  `cursor/ask_question` and `cursor/create_plan` are blocking extension
  methods.
- **Cursor's forum thread on ACP model switching**
  (<https://forum.cursor.com/t/bug-agent-acp-model-switching-updates-session-metadata-but-does-not-change-the-inference-backend/157312>,
  read 2026-10-01): the published `session/set_config_option` response lists
  model values such as `grok-4-20[thinking=true]` and `default[]`, and Cursor
  staff report the switching fix shipped on 2026-04-30. The adopted entry's
  claim that model switching is unsupported under ACP rests on an older
  report, so the measurement decides it.
- **The installed acpx and `cursor-agent`**, programs this Spec did not
  build: acpx `0.19.4`'s agent registry maps `cursor` to `cursor-agent acp`,
  and `cursor-agent` answered the 2026-10-01 probes quoted above.
- **The live measurement** of Core Feature 5, against Cursor's service, is a
  measurement of a runtime this Spec did not build. Until the maintainer logs
  in, it is blocked by `cursor_login_required`.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and the query
`qmd query "Grok Cursor ACP runtime cursor-agent modelos de codificação"
--all --files --min-score 0.3`. The digest
`wiki/sources/digest-roundfix-modelos-de-codificacao-ago-set-2026-e-runtimes-cursor-grok-2026-09-30.md`
records Grok 4.7 below the current defaults on Terminal-Bench 4.0 (25.8
against 59.6 for Opus 5.5) with about 81k output tokens per task, and the
four forum-reported ACP problems; that is why this Spec measures before any
default moves, and why the measurement records billing mode and model
switching. Exa found Cursor's ACP page, its parameters reference
(<https://cursor.com/docs/cli/reference/parameters>), the forum thread above
and the ACP session config options specification
(<https://agentclientprotocol.com/protocol/v2/session-config-options>), which
lets an agent change an option and report the complete state, the shape the
proof relies on.

## Decisions

- **Name the advertised value verbatim.** See ADR-0217. Parsing Cursor's
  bracket into a reasoning effort would invent a mapping for `thinking`,
  `context` and `fast` that the adapter does not state, and composing a
  value would let Roundfix pick Fast billing on the maintainer's behalf.
- **Check the login, never perform it.** See ADR-0217. A login is a
  credential flow that needs a person; the maintainer decided that a missing
  login stops the work with a named blocker.
- **Opt-in only.** The maintainer decided on 2026-10-01 that both
  experiments deliver measurement and an opt-in entry, never a default
  change.
- **Split from the Jev Router experiment.** The two share no file that would
  make one Spec smaller, and together they exceed four implementation Tasks
  (ADR-0208).

## Open Questions

- Whether acpx `0.19.4` calls `authenticate` with `cursor_login` before
  `session/new`, as Cursor's documentation lists, is unknown until a logged-in
  session is observed. Default until measured: the runtime relies on acpx; if
  a logged-in session is still refused, the measurement records it as a
  blocker and a follow-up.

## Technical candidate

The [_techspec.md](_techspec.md) records the parsing rule, the readiness
check, the measurement protocol, coverage and build order.
