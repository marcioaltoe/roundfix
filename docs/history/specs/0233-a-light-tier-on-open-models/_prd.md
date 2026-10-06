---
spec: 0233-a-light-tier-on-open-models
status: archived
created: 2026-10-05
surfaces: [backend, cli, docs]
archived: "2026-10-05"
source_slug: 0233-a-light-tier-on-open-models
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial: the Run sandbox denied the OpenCode docs, the GitHub issue and the OpenRouter key endpoint (R11-R13), and R16 has no Pull Request yet; the queue opens it. Every behavior row passed. The live light-tier check happens on Spec 0234''s own delivery, whose low Tasks run on the light tier by default; the operator records its cost.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 0a9172b6342f30b2c1ca0d0238a25b648687cea5
---


# A light tier on open models

Every Task in a Run runs on its Agent Work Category's Preferred Selection, so a
one-line documentation Task spends the same codex or claude subscription quota
as a new subsystem. Since ADR-0235 OpenAI and Anthropic models run only
through those subscriptions, which leaves open models on OpenRouter as the
cheap place for small work. The adopted Backlog Entry,
[a judge-assigned model tier per Task](references/2026-10-05-a-judge-assigned-model-tier-per-task.md),
proposed a tier per Task and measured it on 2026-10-05: six archived
`complexity: low` Tasks replayed on OpenCode with
`deepseek/deepseek-v4.1-flash` all passed Verification on the first attempt,
in 0.41 times Codex's agent time, for US$0.32 together, and the Jev judge's
tier was no better than the authored `complexity`. This Spec makes Roundfix
run each light Task on an open model by default, within a monthly ceiling the
maintainer sets, escalating once to the subscription when Verification fails,
and lets the advisory judge suggest a tier while the Tasks are written.

## Prerequisites

None of this Spec's Verification depends on another Spec's artifacts. Spec
0234, authored in the same cycle, changes the helper that names the light
tier's key variable and may raise the Roundfix Skill's version; the operator
orders the queue so that the later Spec raises the version from the earlier
one's value.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the new
  keys `openrouter.implement_monthly_ceiling_usd` and `openrouter.light_models`
  follow the existing dotted configuration names, and the tier values `light`
  and `standard` are plain words. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — a light Agent Session reaches
  OpenRouter through OpenCode with the operator's OpenRouter key. Roundfix
  passes the key's variable name, never its value, removes the generic
  OpenRouter variable from that session, and makes no request of its own; the
  advisory judge's recipients and keys are unchanged. No test or Verification
  reaches the network. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0238 (this Spec) decides the tier:
  "A Task is `light` when its `complexity` is `low`, its `type` is not `qa`,
  and it declares no Governed Path to create, change or delete". ADR-0235:
  "Other OpenRouter models stay selectable through OpenCode"; every light
  model passes its refusal. ADR-0108: "Roundfix therefore warms an OpenCode
  session" before applying an effort; a light selection has none. ADR-0050:
  "once Agent work begins, Roundfix fails the Work Item instead of switching
  models over potentially modified state", which stays for the Fallback
  Chain, and ADR-0114: "A Fallback Selection may switch ACP Runtime
  automatically only while Agent work has not begun", which lets a light
  session that fails to start fall to the Preferred Selection. ADR-0049: "each
  present profile replaces the complete lower-precedence profile"; the light
  tier is not a profile and leaves profiles unchanged. ADR-0231: "It must be a
  finite number greater than zero", which the new ceiling follows as a User
  Config value. ADR-0002: "Roundfix uses YAML for User Config at", the file the
  operator edits, and ADR-0027: "truly unknown keys keep failing strict
  validation", so an older binary refuses the new keys. ADR-0200: "Spec
  authoring gets an advisory judge that never gates", and ADR-0201: "Each
  request appends one line to a Judge Log in Roundfix Home"; the tier
  suggestion is one more advisory judgment. ADR-0187 splits the Roundfix
  Skill by command and ADR-0189 ties an owned skill's version to its content,
  so the skill edits raise their versions, and ADR-0233 regenerates the
  raised versions at merge. ADR-0184: "A TechSpec now declares numbered
  Surface Transcripts"; this Spec declares the transcripts of its warnings and
  its judge output. The authored QA gate follows ADR-0080: "QA verdicts
  distinguish environment-blocked rows", and ADR-0091: "required to be
  terminal and to depend on every leaf"; ADR-0096, ADR-0097 and ADR-0167 bind
  its machine stage, its row carry and its pre-PR Pull Request row, ADR-0104:
  "Every Spec therefore rests at least one named" acceptance row on outside
  evidence, ADR-0194, ADR-0195 and ADR-0210 bind what a QA row records, when
  it is observed again and its evidence snapshot, and ADR-0182 settles each
  Task on the facts its gate checks. ADR-0093, ADR-0117, ADR-0156, ADR-0168,
  ADR-0176 and ADR-0183 check this Spec's consistency by citation and receipt. ADR-0069 cites ADR-0050 but decides the Baseline semantic analysis, ADR-0180 cites ADR-0069 but decides the dated Recommended Profile, which this Spec leaves unchanged, ADR-0181 cites ADR-0180 but decides where a configuration is compared with that profile, ADR-0209 cites ADR-0200 but decides how the judge suggests sources and ADR-0208 cites ADR-0209 but decides how sources are grouped, ADR-0217 cites ADR-0049 but decides how a Cursor selection is named, and ADR-0229 cites ADR-0167 but decides how an operator archive resumes a park and ADR-0237 cites ADR-0229 but decides how a Delivery Retry records a merge made outside the queue; this Spec changes none of them, so none applies.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill's canonical files and
  its `SKILL.md` mirror, and the write-tasks skill and its mirror, are
  Governed Paths. The maintainer authorized skill edits on 2026-09-30
  ("considere autorizado a ajustar todas as skills se necessário") and granted
  the governed paths this Spec declares on 2026-10-05 ("Concedo"). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0233-a-light-tier-on-open-models/_authorization.md`; bounded
  files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-tasks/SKILL.md`, `skills/roundfix/SKILL.md`,
  `skills/write-tasks/SKILL.md`.

## Goals

- A `complexity: low` Task that is not QA and declares no Governed Path runs
  on an open OpenRouter model by default, in every repository on the machine,
  with no configuration.
- Light spend never passes the month's ceiling by more than the Tasks already
  running when it was reached, and the maintainer can read what every light
  Task cost.
- A machine without an OpenRouter key, or past its ceiling, runs every Task
  exactly as today and says why.
- A light Task that fails Verification gets one repair on the subscription
  before it fails.
- The advisory judge can suggest a tier for each Task while the Tasks are
  written, without changing where any Task runs.

## User Stories

1. As the maintainer, I want small Tasks to run on a cheap open model by
   default, so that they stop spending my codex and claude quota.
2. As the maintainer, I want a monthly ceiling for that spend, separate from
   the judge's, summed across my repositories, so that the tier cannot cost
   more than I decided.
3. As the maintainer, I want to change or empty the list of open models on my
   machine, so that I can widen the tier after the first escalation counts or
   turn it off.
4. As an operator without an OpenRouter key, I want my Runs to behave as
   before with a warning, so that the new default never breaks a Run.
5. As the maintainer, I want a light Task that fails Verification to be
   repaired once on the subscription, so that a weak answer costs one repair,
   not a failed Task.
6. As a Task author, I want the judge to suggest a tier for each Task, so that
   I can reconsider a `complexity` before the planning Pull Request merges.

## Core Features

1. **The tier.** At dispatch a Task is `light` when its `complexity` is
   `low`, its `type` is not `qa`, and none of the files it declares to create,
   change or delete is a Governed Path; otherwise it is `standard`. A
   `standard` Task, every QA gate and every review Batch dispatch exactly as
   today. A one-Run Agent Selection override turns the tier off for that Run.
2. **The light selection.** A light Task's candidates are, in order, each
   model of the machine's light model list on the `opencode` runtime with no
   reasoning effort, then its category's Preferred Selection and Fallback
   Chain. The list is a User Config value of OpenRouter model ids, defaults to
   `deepseek/deepseek-v4.1-flash`, refuses any id the subscription rule
   refuses or that is not `<author>/<slug>`, and turns the tier off when empty.
   A light selection sends no warm-up prompt.
3. **The key.** A light session reads the OpenRouter key from the one variable
   Roundfix's key helper names, `ROUNDFIX_OPENROUTER_API_KEY` today. The
   session receives a reference to that variable, never its value, and does
   not receive the generic `OPENROUTER_API_KEY`. With no key in the Run's
   environment a light Task runs on its Preferred Selection and the Run prints
   a warning naming the variable.
4. **The ceiling.** `openrouter.implement_monthly_ceiling_usd` is a User
   Config value, a finite number greater than zero, US$10 when unset, ignored
   with a warning in Project Config. Before each light Task Roundfix sums the
   month's Light Spend Log; at or above the ceiling, or when the log cannot be
   read, the Task runs on its Preferred Selection and the Run prints a warning.
5. **The Light Spend Log.** After every prompt of a light session Roundfix
   appends one line to a per-month file in Roundfix Home with the Run, the
   repository, the Spec, the Task, the model and the increase of the session
   cost OpenCode reported, or zero with an `unreported` source when it reported
   none. The file never holds the key, a prompt or a diagnostic.
6. **One escalation.** When a light Task's first Verification fails, its one
   Verification Feedback turn runs on the Preferred Selection, or its Fallback
   Chain before Agent work, in a new Agent Session that receives the Task
   prompt, the diagnostics and the statement that the working tree holds
   another model's attempt. A light session that fails to start falls to the
   Preferred Selection the same way. Both are recorded as Run Events.
7. **The advisory tier.** `roundfix spec judge <slug> --stage tasks` asks the
   judge for a tier, `light`, `standard` or `heavy`, for each non-QA Task and
   prints it as a suggestion; it changes no file, is never read at dispatch,
   and writes the usual Judge Log line for each request.
8. **The record says so.** The configuration guide, the Roundfix Skill's
   `runtime` and `spec` references, the `spec` command reference, the
   write-tasks skill and the model reference describe the tier, its keys, its
   warnings, the Light Spend Log, the escalation and the judge's stage.

## User Experience

A machine with the key and the defaults sees nothing new except light Tasks
finishing on `opencode` with `openrouter/deepseek/deepseek-v4.1-flash` in their
Agent Selection Run Events. A machine without the key, or past its ceiling,
sees one warning line per light Task on standard error, naming the reason, and
the Task runs where it ran before. An escalated Task shows a Run Event and a
progress line naming the Preferred Selection. `roundfix spec judge --stage
tasks` prints one `suggested model-tier` line per non-QA Task.

## Non-Goals / Out of Scope

- A `heavy` tier, or any dispatch change for `medium` and `high` Tasks.
- Reading the tier from the judge or from a Task front-matter field at
  dispatch; the judge's tier stays advisory until it is measured against
  light outcomes.
- Separate OpenRouter keys for the judge and the implementation; Spec 0234
  owns the key names and changes the helper this Spec introduces.
- A Project Config switch for the tier, or per-repository ceilings.
- Proving the light models in Preflight.
- Making ADR-0108's warm-up inert for OpenCode selections that carry an
  effort.
- Escalating a light Task whose Agent turn fails after work began; ADR-0050
  still fails it.
- Changing any built-in or Recommended Profile, `.roundfixrc.yml` or the
  changelog.

## Success Metrics

1. Success Metric: with the key present and an empty Light Spend Log, a
   `low` `docs` Task runs its first prompt on `opencode` with
   `openrouter/deepseek/deepseek-v4.1-flash` and no effort, with no
   `Session setup.` prompt, while `medium`, `qa`, and Governed-Path Tasks run
   on their Preferred Selection.
2. Success Metric: with no key, with a month at the ceiling, and with an
   unreadable log, a light Task runs on its Preferred Selection and the Run
   prints one warning naming the reason, and a light session that fails to
   start is taken by the Preferred Selection before Agent work.
3. Success Metric: the light session's environment carries a reference to
   the key variable and neither the key value nor `OPENROUTER_API_KEY`, and
   the Light Spend Log gains one line per light prompt whose costs add up to
   OpenCode's last reading for the session.
4. Success Metric: a light Task whose first Verification fails has its one
   repair on the Preferred Selection in a new Agent Session, and fails after a
   second failure without another escalation.
5. Success Metric: User Config loads the two new keys with their defaults,
   refuses an invalid ceiling and any model id the subscription rule refuses,
   and Project Config ignores both with a warning.
6. Success Metric: `roundfix spec judge --stage tasks` prints one suggestion
   per non-QA Task with a recorded Judge Log line and changes no file.
7. Success Metric: a real light Task run by the QA gate settles on an open
   model, and its Light Spend Log cost is within 5 % of the movement of the
   OpenRouter key's usage.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not produce:

- The "Measurement, 2026-10-05" section of the adopted Backlog Entry, an
  operator measurement this Spec did not design: six DeepSeek replays, all
  first-attempt passes, agent time 0.41 times Codex's, US$0.3166 in OpenCode's
  figures within about 2 % of the key's movement, the warm-up observation, and
  the judge's tier compared with `complexity` on 50 archived Tasks.
- acpx 0.19.4 as installed on the maintainer's machine, read on 2026-10-05:
  its session option for allowed tools is forwarded only into the Claude Code
  session metadata and the Qoder command line, so the inert warm-up's tool
  restriction does not reach OpenCode.
- OpenCode's configuration and providers documentation
  (<https://opencode.ai/docs/config/>, <https://opencode.ai/docs/providers/>),
  read 2026-10-05: inline configuration through `OPENCODE_CONFIG_CONTENT`
  overrides the global and project files, `{env:NAME}` substitutes an
  environment variable, a provider's `options.apiKey` sets its key, and
  OpenRouter models are named with the full OpenRouter id after the
  `openrouter/` prefix.
- OpenCode issue 13219
  (<https://github.com/anomalyco/opencode/issues/13219>), read 2026-10-05:
  inline configuration was reported to skip `{env:}` substitution. The
  measurement's six replays, which passed the key that way on OpenCode
  1.18.34, contradict it for that version; the QA gate's live light Task
  observes it again through Success Metric 7.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and
`qmd query "model tier per task cheap open model OpenRouter DeepSeek" --all
--files --min-score 0.3`; its top results are this repository's own mirrored
model reference, Spec 0218's router measurement and the adopted Backlog Entry,
and its subscription-authentication synthesis records that an OpenAI or
Anthropic key bills pay-per-use while the codex and claude CLIs consume the
subscriptions, the boundary ADR-0235 set. Exa found OpenCode's configuration
and providers pages and issue 13219 cited above. The brain-side Inbox for this
repository holds two Baseline entries that share no context with this Spec.
The repository has no other open Backlog Entry and no unresolved Finding.

## Decisions

- Map `light` from `complexity: low`, excluding QA and Governed-Path Tasks;
  the judge's tier is advisory only. See ADR-0238.
- The allow-list and the ceiling are User Config values, because the key and
  its spend are the machine's; an empty list turns the tier off. See ADR-0238.
- Cost is OpenCode's reported session cost, appended to a Light Spend Log in
  Roundfix Home; no OpenRouter usage request. See ADR-0238.
- No reasoning effort on a light selection, so no warm-up prompt. See
  ADR-0238.
- Escalation keeps the light attempt's changes and runs the one repair on the
  Preferred Selection in a new Agent Session. See ADR-0238.
- A one-Run override turns the tier off for that Run.

## Open Questions

- Should a repository be able to turn the tier off in Project Config? Default:
  not in this Spec; the operator empties the light model list.
- Should a light Agent turn that fails after work began escalate too? Default:
  no; ADR-0050 fails the Task, and the first real Runs' escalation counts
  decide.
- Should ADR-0108's warm-up become inert for OpenCode selections with an
  effort? Default: left to a later decision; the light tier does not depend on
  it.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
