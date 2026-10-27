---
spec: 0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote
status: active
created: 2026-10-26
surfaces: [backend]
---

# A profile change keeps the recorded decisions and names what it wrote

Quaestor changed its adopted repository to the `python` Baseline Profile with
Roundfix v0.123.0 and met five frictions: `baseline update` refused
`--decision profile=python` as an unknown decision without naming the command
that changes a profile; `baseline plan` ignored the eleven decisions the Setup
Manifest records and refused once per round; `baseline apply` wrote a root
backup `AGENTS.<sha256>.md` without naming it; `baseline skills reconcile
--revision <newer>` answered "no changes" without saying what it compares;
and `baseline skills restore` installed new skills under `.agents/skills`
without the `.claude/skills` links the repository keeps for Claude Code. After
this Spec the update's refusal names the plan command, the plan carries every
recorded decision the selected profile still selects and asks for everything
else in one refusal, the apply names each Root Backup it wrote, reconcile
names the Setup Snapshot revision the installed skills follow, and a restore
writes a Claude Skill Link for each skill it installs in a repository that
keeps them.

## Sources

One Backlog Entry, written and adopted at birth by this Spec's planning
change (see `references/_index.md`):

1. [`references/2026-10-26-a-baseline-profile-change-loses-the-recorded-decisions-and-writes-without-saying-so.md`](references/2026-10-26-a-baseline-profile-change-loses-the-recorded-decisions-and-writes-without-saying-so.md),
   `type: fix`, from part 1 of ai-memory message
   `01a17abd-20fe-7213-ab5c-f1a1bdb043fe` of
   `default/gesttione-solutions-quaestor` (2026-10-26 21:22 BRT). Part 2 of
   the same message (executable scripts, settle exit 128, archive with
   `qa: declined`) belongs to Spec 0327 and is out of scope here.

No open Backlog Entry or unresolved Finding shares this context: at the
authoring base `docs/backlog/` and `docs/findings/` hold no entry.

## Prerequisites

None in `requires`. Authored against `6a2630e5` (v0.124.0 plus the 0323
delivery and the 0326 plan). Spec 0326 moves the bundled Setup Snapshots to
`b56fc5ff` before this Spec delivers; no test of this Spec names a snapshot
revision as a literal, and the reconcile test reads the revision from the
embedded catalog, so the order of the two does not matter. Spec 0322 edits
`docs/user-guide/context-driven-development.md` in other sections; task_01
edits only the sections this Spec names.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the new planned-change action and warning codes are strings in
  existing arrays. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added, and no Verification or QA row reaches the network.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0263 keeps the profile out of
  `baseline update`, so friction (a) becomes a refusal that names the plan
  command: ADR-0263: "The update stays the Managed Refresh: it never changes
  the profile or the preservation mode." ADR-0264 keeps reconcile and restore
  from moving a vendored skill to an unproven revision, so friction (d)
  changes a message and not a behavior: ADR-0264: "because nothing the binary
  holds can prove a lock-matching copy behind". ADR-0241 gives a
  repository-owned profile its skill contracts from the embedded Setup
  Snapshots, so the revision reconcile names comes from the same contracts:
  ADR-0241: "A repository-owned profile names none". ADR-0341, new, records
  the Carried Decisions and the Claude Skill Links: ADR-0341: "writes it
  inside the same recoverable transaction as the skill tree". Recorded as
  not applying: ADR-0221, ADR-0224 and ADR-0322, listed because they cite
  ADR-0264; no snapshot, installed skill, lock entry or profile asset
  changes. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AQ batch 5 decisions of 2026-10-26 quoted in the authorization
  record, and the Governed Paths this Spec declares. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote/_authorization.md`;
  bounded files:
  `internal/cli/baseline_plan_test.go`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`,
  `.agents/skills/setup-context-driven/SKILL.md`,
  `skills/setup-context-driven/SKILL.md`.

## Goals

- A profile change asked through `baseline update` is refused with the plan
  command that performs it.
- `baseline plan` in an adopted repository needs only `--profile`,
  `preservation.mode` and any decision the new profile introduces; every
  missing answer is named in one refusal.
- Every Root Backup an apply writes is named in its text and JSON output.
- A reconcile that changes nothing says what it removes and which Setup
  Snapshot revision the installed skills follow.
- A restore in a repository that keeps Claude Skill Links leaves every
  restored skill visible to Claude Code.

## Core Features

1. The update refusal (friction a). `roundfix baseline update --decision
   profile=<id>` exits `2` with a message that the profile is not a decision
   and a next action that names `roundfix baseline plan --profile <id>` and
   `roundfix baseline apply`.
2. Carried decisions (friction b). `roundfix baseline plan` reads the Setup
   Manifest and carries each recorded decision the selected profile selects
   and does not fix, unless a flag answers it; the text plan lists the
   Carried Decisions; a recorded decision it does not carry is named in a
   plan warning.
3. One refusal (friction b). A plan that lacks answers names, in one
   refusal, every missing decision, the decisions a missing decision's
   catalog effect would require, and `preservation.mode`.
4. Named Root Backups (friction c). The text plan labels each backup file
   change, and the apply result carries one warning per Root Backup, naming
   the file and the carrier it copies.
5. Reconcile says what it compared (friction d). A reconcile with no planned
   change prints what it removes and the Setup Snapshot revision the
   installed skills of that source follow.
6. Claude Skill Links (friction e). `roundfix baseline skills restore` plans,
   binds into its Plan Digest and writes, in its transaction, one Claude Skill
   Link per selected skill that lacks one, in a repository whose
   `.claude/skills` keeps them.

## User Experience

An adopter that changes profile runs `roundfix baseline plan --repo .
--profile <id> --decision preservation.mode=preservation --format json`; the
plan reuses the recorded answers and a refusal, if any, names everything
still needed at once. The applied change prints, for each Root Backup,
`Warning: baseline.preservation.backup.written: AGENTS.<digest>.md: …`, which
tells the adopter to commit it. `baseline skills reconcile` that finds
nothing to remove now explains that installed skills move only with a
Roundfix release whose Setup Snapshot names a newer revision. In a repository
whose `.claude/skills` holds per-skill links, `baseline skills restore
--profile <id> --skill <name>` lists `- create-link .claude/skills/<name>
[<name>]` in its preview and writes the link with the skill; Quaestor repairs
its `ai-memory-*` links with one restore per skill.

Declared behavior changes: `baseline plan` in a repository with a Setup
Manifest now succeeds where it refused for decisions the manifest records,
and its Plan Digest reflects those recorded values; a plan refusal can name
more than it did; an apply that writes a Root Backup carries one more
warning; a restore preview can list `create-link` changes and therefore exit
`3` where it reported no changes. `baseline update`'s managed refresh, its
exit codes and its JSON stay as they are.

## Non-Goals / Out of Scope

- Letting `baseline update` change the profile (ADR-0263, ADR-0341).
- Recording Root Backups in the Setup Manifest: the backup is immutable
  evidence the repository owns once written, a later Preservation finds it
  by its content-addressed name, and the manifest records only
  what the Baseline keeps current.
- Moving an installed skill to an arbitrary upstream revision: reconcile only
  removes lock entries and installed skills follow the Setup Snapshot
  (ADR-0264); the snapshots move with Spec 0326 and later releases.
- Doctor or `skills check` reporting a missing Claude Skill Link, links for
  harnesses other than Claude Code, and linking Roundfix-owned skills, which
  `roundfix skills install` already handles.
- Part 2 of the source message, owned by Spec 0327.

## Success Metrics

1. Success Metric: `baseline update --decision profile=python` exits `2`
   and its text and JSON name `roundfix baseline plan` with `--profile
   python`; no file changes.
2. Success Metric: in a temporary repository adopted with `go-cli-tui`,
   `baseline plan --profile python --decision preservation.mode=preservation`
   exits `0` and its plan decisions equal the recorded values of every
   decision `python` selects; the same plan without `preservation.mode`
   exits `3` naming `preservation.mode`, and without a manifest and without
   decisions it exits `3` naming every required decision, `runtime.backend`
   and `runtime.design` with their condition, and `preservation.mode` in one
   message.
3. Success Metric: an apply of a plan that backs up `AGENTS.md` prints one
   `baseline.preservation.backup.written` warning naming the written file in
   text and JSON, and a managed refresh prints none.
4. Success Metric: a reconcile with no planned change prints the Setup
   Snapshot revision the embedded catalog pins for the profile's skills.
5. Success Metric: in a repository whose `.claude/skills` keeps one Claude
   Skill Link, a restore of another skill previews `create-link`, applies the
   link with target `../../.agents/skills/<name>`, rolls it back with the
   tree on a failure, and adds nothing in a repository without such links.
6. Success Metric: after release, Quaestor or another adopter changes a
   profile with `baseline plan` and confirms the recorded decisions carried.

## Acceptance evidence

The outside-evidence row rests on text this Spec did not write. First, the
adopter's report quoted verbatim in the adopted Backlog Entry under
`references/` (ai-memory message `01a17abd`, 2026-10-26). Second, published
documentation of the skills CLI (vercel-labs/skills, "Installation Methods",
retrieved 2026-10-26), whose passage the QA row quotes: "Canonical copy:
Skills are copied to `.agents/skills/`" and "Symlinks: Each agent's skill
directory symlinks to the canonical location", with Claude Code's project
path listed as `.claude/skills/`. The QA compares the restored link layout
with that passage, offline.

## Unreachable Acceptance

- criterion: Success Metric 6, after release, Quaestor or another adopter
  changes a profile with `baseline plan` and confirms the recorded decisions
  carried.
  reason: it needs the released binary in another repository, which no
  Verification or QA row may read.
  satisfied-by: the release notice to `default/gesttione-solutions-quaestor`
  and that project's next profile change, reported to this project through
  ai-memory.

## Validation

Characterized on 2026-10-26 with `bin/roundfix` built from `6a2630e5` in
temporary repositories under the operator's scratch directory, with a fresh
`HOME`.

- Friction a: with this repository's `docs/agents/setup-context.json`
  copied into a temporary repository, `baseline update --decision
  profile=python --no-skills --no-history` exited `2` with "declare Setup
  Manifest decisions: resolve Setup Manifest decisions: normalize Baseline
  decisions: unknown decision "profile"" and next action "correct the
  declared decisions and rerun roundfix baseline update".
- Friction b: in the same repository `baseline plan --profile python`
  refused with the nine base decisions although the manifest records them;
  with those nine it refused again for `runtime.backend, runtime.design`
  (required by `autonomous.enabled` through the `autonomous-work` module);
  with those it refused a third time for the preservation mode. The CLI
  never reads the manifest (`runBaselinePlanCommand` passes only flag
  values to `BuildPlan`).
- Friction c: a repository adopted with `go-cli-tui` (greenfield) and then
  planned with `--profile python --decision preservation.mode=preservation`
  listed `- create AGENTS.2336a90d….md (1 managed entries)` in the text plan;
  the apply printed only `Verified postimages: 11`; the file is byte-identical
  to the committed `AGENTS.md`, and the Setup Manifest does not name it.
- Friction d: `buildSkillsReconcilePlan` classifies lock entries as present,
  moved, obsolete, required-removed or unrelated and plans only removals; it
  never compares an installed tree with the revision. The reported "no
  changes" is correct behavior with an unexplained message.
- Friction e: `inspectRestoreTarget` and the restore transaction write only
  `.agents/skills/<name>` and `skills-lock.json`; no Baseline code writes
  under `.claude/skills` (`grep -rn '.claude/skills' internal/baseline`
  finds only the retired-skill listing). `roundfix skills install --target
  project` already links `.claude/skills/roundfix` when `.claude/skills` is a
  directory (`maybeCreateProjectClaudeSkillSymlink`).
- Secondbrain, read only: `wiki/index.md`; `qmd query "baseline profile
  change setup manifest decisions roundfix" --all --files --min-score 0.3`
  returned mirrors of this repository's ADR-0241, changelog and doctor guide;
  `qmd query "claude skills symlink agents skills alias"` returned a Drive
  setup note recording that a shared Drive keeps `.claude/skills` as plain
  copies because it stores links as ordinary files, which is why ADR-0341
  links only where links already exist.
- Exa: the skills CLI documentation states the canonical-copy-plus-symlink
  layout quoted under Acceptance evidence, and its issue #693 states that
  Claude Code reads `.claude/skills` and not `.agents/skills`.

## Glossary

- adds: **Carried Decision**
- adds: **Root Backup**
- adds: **Claude Skill Link**

## Decisions

Maintainer answers to the cycle AQ batch 5 questions (AskUserQuestion,
2026-10-26 about 22:35 BRT), verbatim:

- 2026-10-26: Grouping — "Duas Specs (Recomendado)": 0327 Run lifecycle
  (executable scripts, settle exit 128, archive with `qa: declined`); 0328
  Baseline profile-change frictions.
- 2026-10-26: Executables — "_authorization.md nomeia os executáveis
  (Recomendado)", option text: "O commit da Task aceita 0755 só nos caminhos
  que a autorização da Spec lista; todo o resto continua recusado. O
  write-tasks e a recusa passam a citar essa regra." (applies to Spec 0327).
- 2026-10-26: Authorization — "Mesmas regras do ciclo AQ (Recomendado)",
  option text: "implement, commit, push, PR e merge nos Governed Paths que
  cada Spec declarar; fora Makefile, .github/workflows, go.mod/go.sum e lint.
  A release continua a meu critério, pelos critérios gravados."

Autonomous route: no person present for the technical choices; decisions taken from the operator's cycle AQ batch 5 brief for Spec 0328 of 2026-10-26 and the characterization above.

Taken from those sources, with no further interview:

- Friction a is a refusal, not a feature, so the fix names the working
  command: ADR-0263: "The update stays the Managed Refresh: it never changes
  the profile or the preservation mode."
- Friction c keeps the backup every adoption-class plan takes of a root
  carrier it rewrites and names it in output;
  the Setup Manifest does not record it.
- Friction d is not a defect: the message changes, the behavior does not.
- Friction e links only where `.claude/skills` already keeps a Claude Skill
  Link, inside the restore transaction (ADR-0341).
- All five frictions fit four implementation Tasks plus QA; nothing moves to
  a Backlog Entry.

## Open Questions

None.
