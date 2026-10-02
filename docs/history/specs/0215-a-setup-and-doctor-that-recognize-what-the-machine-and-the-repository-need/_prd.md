---
spec: 0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need
status: archived
created: 2026-10-02
surfaces: [backend, cli, docs]
archived: "2026-10-02"
source_slug: 0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: QA partial with zero findings; the only blocked row is the Pull Request row (11), which needs an open Pull Request
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 97539097ce4062aaa7f71d778e3e088f14887aa9
---


# A setup and doctor that recognize what the machine and the repository need

On 2026-10-01, preparing a second machine, the maintainer asked:

> "o que mais é necessário na outra maquina para o perfeito funcionamento?
> Agora, posso fazer esse setup manualmente mas é importante que a parte do
> roundfix de setup verifique e reconheça os requerimentos para que possa
> rodar sem problemas. SE isso está desatualizado, pode virar uma spec ao
> final do fluxo atual."

`roundfix setup` and `roundfix doctor` prove Node.js, acpx, the adapters, the
Agent Selection Profiles and the Repository Skill Set. They never look at the
GitHub CLI the Delivery Queue publishes through, at Git and the identity every
Task commit uses, at the repository remote, or at the tools the repository's
Verification runs. The runbook the operator sent to every adopter that day
lists those checks for a person to run by hand and says that Doctor does not
make them yet. Doctor is also wrong in two places it already looks. With a
`NODE_OPTIONS` preload whose file is gone, it reports `node: ok` and fails
`acpx` with an install command for a version that is already installed; and it
reports `skills: ok` for upstream skills that match their lock while trailing
the Setup Snapshot the binary pins. This Spec makes Setup and Doctor name each
requirement with a stable code, a status and a next action, makes
`deliver start` refuse up front when it cannot publish, and has this
repository declare what its own delivery needs.

## Prerequisites

Spec 0211 (a delivery queue that finishes without intervention) must be on
the default branch first. It adds the reading of `NODE_OPTIONS` that the
agent environment uses, and this Spec's `environment` line reuses that reading
rather than adding a second one. The Task Graph names it in `requires`.

## Project Constraints

- Identifier strategy: applicable — new public identifiers follow existing
  forms: five Doctor and Setup line names (`gh`, `git`, `remote`, `toolchain`,
  `environment`), one status token `warn`, a family of finding codes with the
  `DR-` prefix, and one Project Config key, `verification.tools`. No
  project-owned Internal Identifier is generated. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — Roundfix reads the forge only through
  the user's own `gh` and `git` and their stored credentials. It never reads,
  prints or stores a token, never passes `--show-token`, and never logs in or
  out. Doctor's three network reads are bounded and read-only, as ADR-0220:
  "It now makes at most three network reads, all against the repository's own
  forge". No test or Verification command reaches the network. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0220 (this Spec) decides the five
  lines, the codes, the `warn` status, the bounded forge reads and the
  `deliver start` refusal, and ADR-0220: "only `failed` changes the exit
  code". ADR-0221 (this Spec) holds a required external skill to its Setup
  Snapshot, ADR-0221: "Trailing alone makes the `skills` line `warn`". Under
  ADR-0087 a required tool is found without running it, ADR-0087: "Discovery
  therefore resolves a bounded symlink chain to its target and judges the
  target, while still never executing it". Under ADR-0211 the `environment`
  line reads `NODE_OPTIONS` as the agent environment does, ADR-0211: "removes,
  from the environment it gives an agent process, each `NODE_OPTIONS` preload
  whose value is a file path that does not exist". Under ADR-0172 Doctor still reports reclaimable Run storage from cheap reads, ADR-0172: "The check never reports `failed`, so
  it never changes Doctor's exit code". ADR-0181 keeps the comparison with the
  recommendation where it is, ADR-0181: "It is read-only and offline". Under ADR-0107 profile readiness still covers every configured Agent Work Category, ADR-0107: "Categories that merely inherit `general` contribute no distinct tuple". ADR-0189 keeps the owned-skill minimum and governs
  the skill edit, ADR-0189: "Every edit to an owned skill now needs a version
  change and a line in the record". Under ADR-0187 the Roundfix Skill and the command reference are still read one command at a time, ADR-0187: "Both documents are now an entry file plus one file
  per command family". ADR-0191 still decides which skills a snapshot lists,
  ADR-0191: "A setup snapshot mirrors one upstream list, and the asset sync is
  its only writer", and this Spec only compares installed trees with it.
  ADR-0192 governs the derived-path declarations this repository adds,
  ADR-0192: "When every conflicted path matches a derived-path declaration in
  Project Config". ADR-0193 governs the `requires` entry naming Spec 0211,
  ADR-0193: "A Spec now names the Specs it needs in the `requires` list".
  Under ADR-0032 Roundfix still detects and avoids an unsafe codex, and the new lines follow its rule, ADR-0032: "Roundfix only diagnoses and instructs". ADR-0184 has the TechSpec
  state the changed output, ADR-0184: "A TechSpec now declares numbered
  Surface Transcripts", and ADR-0183 checks this Spec's receipts, ADR-0183:
  "A claim that attributes behavior to a decision record now carries a Claim
  Receipt". Spec consistency is still checked by citation, never by
  inference, ADR-0093: "It compares citations, declarations, and
  cross-references that are already written down". The QA gate still proves
  machine facts before it spends an agent turn, ADR-0096: "The QA gate runs a
  Daemon-owned mechanical stage before its Agent Session". A QA row still
  carries forward only on declared, unmoved evidence, ADR-0097: "A carried row
  names the report and the head that established it". A composed profile
  still takes a setup composed from upstream setups by name, and the snapshot
  comparison reads that composed setup as the restore command does, ADR-0204:
  "carries the union of their skills and activation bundles". The Baseline
  still takes upstream setup names, ADR-0206: "A setup snapshot takes the
  upstream list's new name", and no snapshot changes here. ADR-0180 does not
  apply: no built-in selection or dated Recommended Profile changes. These do
  not apply either, because this Spec changes none of their subjects: ADR-0117
  (a defect is checked by the stage that can produce it), ADR-0168 (when a
  related-ADR gap opens), ADR-0176 (what citation checks read), ADR-0182 (the
  facts a Task settles on), ADR-0194 (what the Daemon records of a QA row),
  ADR-0195 (rows observed on every pass), ADR-0210 (evidence snapshot
  digests), ADR-0217 (a Cursor selection and its login, a Doctor line its own
  Spec adds), ADR-0218 (the Jev Router's model and ceiling) and ADR-0219 (a
  profile draft's built-in profile). The gate
  is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104,
  ADR-0155, ADR-0156 and ADR-0167.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded
  in [_authorization.md](_authorization.md). The maintainer authorized every
  skill change ("considere autorizado a ajustar todas as skills se
  necessário") and the Baseline source, guides and `.roundfixrc.yml`
  ("Autorizar os dois"). No Makefile, lint, formatter, test-runner or CI file
  and no dependency changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need/_authorization.md`;
  bounded files: `.roundfixrc.yml`, `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- On a machine missing `gh`, its login, Git, a Git identity, the delivery
  remote or a tool the repository's Verification needs, Doctor and Setup name
  the gap with a code and a next action before any Run or queue starts.
- An offline machine gets warnings, never a failure it cannot fix by
  configuration.
- `roundfix deliver start` refuses, before it creates a queue, when this
  machine cannot publish.
- Doctor names a dead `NODE_OPTIONS` preload and a required upstream skill
  that trails its Setup Snapshot, and the managed refresh restores the
  trailing skill.
- This repository declares the tools its Verification needs and its derived
  paths in its Project Config.

## User Stories

1. As a maintainer preparing a new machine, I want `roundfix setup` and
   `roundfix doctor` to tell me that `gh`, Git, the remote or a Verification
   tool is missing or misconfigured, and what to run, so that I do not find
   out from a parked delivery.
2. As a maintainer working offline, I want the forge checks to warn rather
   than fail, so that Doctor still tells me what it could check.
3. As the operator of an unattended queue, I want `deliver start` to refuse
   when `gh` cannot publish to the repository, so that no item is created only
   to park at the first push.
4. As a maintainer whose shell exports a stale preload, I want Doctor to name
   the preload, so that I fix the cause and not reinstall acpx.
5. As a maintainer, I want Doctor to say which upstream skills differ from the
   snapshot my Baseline pins and the managed refresh to restore them, so that
   Agents read the skills the Baseline was built against.
6. As a maintainer of a repository, I want to list the tools its Verification
   needs in Project Config, so that Doctor checks tools no configured command
   names.
7. As a maintainer using optional features, I want Doctor to say which of
   their `ROUNDFIX_` keys this shell has, without showing a value, so that I
   know why the judge skipped.

## Core Features

1. **Coded readiness lines.** Doctor and Setup print five new lines, `gh`,
   `git`, `remote`, `toolchain` and `environment`. Each finding carries a
   stable `DR-` code, a status of `ok`, `warn` or `failed`, and a next action.
   `warn` never changes the exit code, and its next action is printed as for
   `failed`. Existing lines keep their names, order and statuses, except the
   `skills` line of Core Feature 6.
2. **The forge.** The `gh` line fails when `gh` is absent, older than the
   first release whose `auth status` reports JSON, has no account for the
   forge host, holds a token the forge rejects, or the account's permission
   on the repository is below write. It warns when the forge did not answer.
   The `remote` line fails when the delivery remote is absent or points to a
   host that is not a GitHub host, and warns when `git ls-remote` cannot reach
   it. Each network read is bounded, read-only and never waits for input
   (ADR-0220).
3. **Git.** The `git` line fails when Git is absent, older than 2.23.0, or the
   repository has no `user.name` or `user.email`.
4. **The toolchain.** The `toolchain` line fails for each tool missing from
   `PATH` that a configured command starts with, or that Project Config's new
   `verification.tools` list names, and says which command or declaration
   needs it. Tools are found without being run (ADR-0087). A command whose
   first word cannot be read as a tool name warns.
5. **The environment.** The `environment` line warns for each `NODE_OPTIONS`
   preload whose file is missing and reports which `ROUNDFIX_` key variables
   the optional features name are set, without a value. It never reads the
   provider's generic key variables.
6. **Skills that trail the snapshot.** The `skills` line reports a required
   upstream skill that matches its lock but differs from its Setup Snapshot
   tree as trailing, with `warn` when nothing else fails (ADR-0221).
   `roundfix baseline update` restores each trailing skill through the
   existing restore path. Repository Skill Set readiness keeps its meaning.
7. **Delivery refuses early.** `roundfix deliver start` runs the `gh` and
   `remote` checks after the authorization check and refuses with the failing
   finding, its code and next action, before it creates a queue or starts an
   owner. A `warn` does not stop it.
8. **This repository declares what it needs.** Its Project Config lists the
   tools its Verification needs and declares its derived paths with their
   sanctioned regeneration commands, so the queue resolves a conflict
   confined to them.
9. **The guides say so.** The `doctor`, `setup`, `deliver` and `baseline`
   command guides, the configuration guide and the Roundfix Skill's `setup`,
   `deliver` and `baseline` references describe the new lines, codes,
   refusal, key and restore.

## User Experience

A ready machine sees five more `ok` lines between `pre-pr-review` and
`skills` in Doctor and after `acpx` in Setup. A missing login reads
`gh: failed (DR-GH-UNAUTHENTICATED: gh has no account for github.com; next:
gh auth login --hostname github.com)`. An offline machine reads `gh: warn` and
`remote: warn` with the reason and exits zero when nothing else fails. A
refused `deliver start` prints the existing `Preflight failed` block whose
reason names the finding, and exits `2`. No line prints a token, a key value
or an email address.

## Non-Goals / Out of Scope

- Installing, upgrading or logging in to `gh`, Git or any tool; Setup offers
  no install for them.
- Running the repository's Verification or reading a Makefile to infer tools.
- Supporting a forge other than GitHub.
- Checking commit hooks, signing keys or credential helpers.
- Restoring this repository's own trailing upstream skills; the managed
  refresh does it after release, in its own change.
- The Release Plan Command's skills and guides lines (Backlog Entry
  `docs/backlog/2026-09-30-release-plan-does-not-report-skill-and-guide-checks.md`),
  which change that command's frozen JSON contract and stay open.
- The `cut-release` dispatch trigger (Backlog Entry
  `docs/backlog/2026-09-30-a-dispatch-trigger-names-a-skill-the-model-cannot-invoke.md`),
  a Baseline dispatch wording change that belongs with the Baseline wording
  work.
- Changing the agent environment, which Spec 0211 owns.

## Success Metrics

1. Success Metric: on a disposable machine profile with no `gh` account, no
   Git identity, no remote and a missing tool, Doctor prints
   `DR-GH-UNAUTHENTICATED`, `DR-GIT-IDENTITY`, `DR-REMOTE-MISSING` and
   `DR-TOOL-MISSING` lines with next actions and exits `1`; with every fake
   ready it prints five `ok` lines and exits `0`.
2. Success Metric: with the forge unreachable, the `gh` and `remote` lines
   report `warn` and Doctor exits `0` when nothing else fails; no probe
   outlives its bound.
3. Success Metric: `deliver start` with `gh` unauthenticated exits `2` with
   the finding's code in its reason and creates no queue; with an unreachable
   forge it proceeds.
4. Success Metric: with `NODE_OPTIONS` naming a missing preload, the
   `environment` line names that path; with a key variable set to a sentinel,
   no output contains the sentinel.
5. Success Metric: an installed upstream skill that matches its lock and
   differs from its snapshot tree is reported trailing with `warn`, and the
   managed refresh's preview lists it for restore.
6. Success Metric: this repository's Project Config loads with its declared
   tools and derived paths, and Doctor's `toolchain` line in this repository
   names `go` and `make`.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The maintainer's request of 2026-10-01 quoted above, and the adopter runbook
  the operator sent to five repositories that day through the Secondbrain
  inbox (`inbox/conexus/2026-10-01-roteiro-de-atualizacao-do-roundfix-e-validacao-do-ambiente.md`),
  whose step 0 lists Git, `gh auth status` with push, PR and merge, the
  Verification toolchain and the `ROUNDFIX_` keys, and says "O `gh` e as
  ferramentas da Verification ainda não são verificados pelo `roundfix doctor`".
- A measurement on 2026-10-02 at `18ef15eb`: with the shell's `NODE_OPTIONS`
  naming a removed `--require` preload, `roundfix doctor` printed `node: ok`
  and `acpx: failed (run npm install -g acpx@0.12.0 ...)` while acpx 0.19.4
  was installed, and passed with `NODE_OPTIONS` unset; in the same shell the
  Secondbrain's `qmd` died with `Cannot find module`.
- A measurement on 2026-10-02 at `18ef15eb`: eleven of the twenty-nine
  required upstream skills in this repository match their lock entries and
  differ from the `treeDigest` the `go` Setup Snapshot pins at `b3c45a4`,
  while Doctor printed `skills: ok (43 required: 14 Roundfix-owned, 29
  external)`.
- The GitHub CLI's published record: `gh auth status` gained `--json` in
  release 2.81.0 (<https://github.com/cli/cli/releases/tag/v2.81.0>), its
  entries carry a `state` of `success`, `timeout` or `error`
  (`pkg/cmd/auth/status/status.go` in `cli/cli`), and with `--json` it "will
  always exit with zero regardless of any authentication issues". A
  characterization with gh 2.102.0 on 2026-10-02 recorded `{"hosts":{}}` with
  no account, `state` `error` with a `401` for a rejected token, and `state`
  `error` with `connection refused` when the API could not be reached.
- Git's 2.23.0 release notes
  (<https://github.com/git/git/blob/master/Documentation/RelNotes/2.23.0.adoc>):
  "Two new commands "git switch" and "git restore" are introduced", which
  Roundfix runs and prints as next actions.
- The operator's intervention log, entry 66 of 2026-10-01: the tool shell
  lacked the `ROUNDFIX_` keys the interactive shell exported, so the judge
  skipped until the queue was started from an interactive shell.

## Decisions

- Only a definite answer fails and an unanswered network question warns. See
  ADR-0220.
- Doctor contacts the repository's own forge with at most three bounded reads,
  ending the offline guarantee for the `gh` and `remote` lines only. See
  ADR-0220.
- Required tools come from configured commands and an explicit Project Config
  list, never from running the Verification or reading a Makefile. See
  ADR-0220.
- A trailing upstream skill warns and is restored by the managed refresh. See
  ADR-0221.
- Two open Backlog Entries stay apart: the Release Plan lines change another
  command's frozen contract, and the `cut-release` trigger is Baseline
  dispatch wording.

## Open Questions

None.
