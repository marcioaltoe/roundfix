---
spec: 0329-the-skills-lock-matches-what-is-installed-and-a-gate-less-run-settles-green
status: active
created: 2026-10-27
surfaces: [backend]
---

# The skills lock matches what is installed, and a Run without a QA gate settles green

Quaestor reported two frictions on Roundfix v0.124.0 and saw both again on
v0.125.0. After `roundfix baseline update` installed the Roundfix-owned
skills the binary carries, the repository's `skills-lock.json` still held the
`computedHash` of the copies it replaced, so the repository's own lock test
failed until someone rewrote five hashes by hand. And two Runs of Specs whose
Task Graph declares `qa: declined` reached Clean while the repository
Verification was red, because the Daemon runs that command at settlement only
in a graph with a QA gate Task. After this Spec the bundle install rewrites
the hash of every lock entry that names an installed owned skill, every Task
settles on the repository Verification whether or not its graph has a QA
gate, and a Run without a QA gate checks its integrated tree before it ends
Clean.

## Sources

Two Backlog Entries, adopted from the operator's unpushed commit `450eb2b9`
(branch `docs/backlog-aq-batch-6`) as copies under `references/` (see
`references/_index.md`):

1. [`references/2026-10-27-baseline-update-leaves-the-skills-lock-hash-of-the-catalog.md`](references/2026-10-27-baseline-update-leaves-the-skills-lock-hash-of-the-catalog.md),
   `type: fix`, Quaestor ai-memory message
   `01a17b1a-38e7-7c82-b268-971c34cfcf4d` (2026-10-26 23:04 BRT), friction 1,
   repeated as item 1 of message `01a17d0a-cd60-7191-a724-2988affea2d0`
   (2026-10-27 08:07 BRT).
2. [`references/2026-10-27-a-run-without-a-qa-gate-ends-clean-with-the-repository-red.md`](references/2026-10-27-a-run-without-a-qa-gate-ends-clean-with-the-repository-red.md),
   `type: fix`, the same first message, friction 3, repeated as item 2 of
   message `01a17d0a`.

The "Outro ponto da 0024" paragraph of message `01a17d0a` (one Unreachable
Acceptance per blocked row, and `settle` reading a stale Run Worktree)
belongs to Spec 0330 and is out of scope here. No open Backlog Entry or
unresolved Finding on `main` shares this context: at the authoring base
`docs/backlog/` holds two entries, both about `roundfix settle`, the subject
of Spec 0330, and `docs/findings/` holds none.

## Prerequisites

None in `requires`. Authored against `5e9475bf` (v0.125.0), with no active
Spec under `docs/specs/`. Specs 0330 and 0331 of the same cycle deliver after
this one; if either changes `internal/daemon/task_engine.go` or
`internal/cli/implement.go`, it rebases on this Spec's `TaskCycle` and
`runTaskWorker` changes.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the new JSON field and the new diagnostics file name are strings
  in existing payloads and directories. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added, and no Verification or QA row reaches the network.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0182 scoped the repository
  Verification at settlement to a gated graph, and ADR-0342 narrows it:
  ADR-0342: "For every non-QA Task, with or without a QA gate Task, the Daemon
  appends the configured repository Verification as the attempt's last
  command". That makes true the Baseline clause every Daemon-assigned turn
  reads in `docs/agents/agent-instructions.md`: "In a Daemon-assigned turn
  the Daemon owns that claim: it runs the Task's declared Verification and
  the repository Verification at settlement." ADR-0189 owns the owned skill
  bundle, and its 2026-10-27 addendum records the lock rewrite: ADR-0189:
  "The minimum version of an owned skill is now the version the binary
  carries for it." ADR-0233 regenerates an owned skill's version raise at
  merge, which task_01's version raise follows, and ADR-0256 holds a release
  until the skills describing a changed surface agree, which task_01
  satisfies for `implement`, `baseline update` and `skills install`.
  Recorded as not applying, each listed because it cites ADR-0189 or one of
  the ADRs above: ADR-0197, ADR-0250, ADR-0272, ADR-0279 and ADR-0286; no
  pre-PR review lineage, Baseline module, coverage record, review record,
  Queue Owner or Correction Review changes. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AR decisions of 2026-10-27 quoted in the authorization record, and
  the Governed Paths this Spec declares. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0329-the-skills-lock-matches-what-is-installed-and-a-gate-less-run-settles-green/_authorization.md`;
  bounded files:
  `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `skills/write-tasks/SKILL.md`.

## Goals

- After `baseline update` or `skills install --target project`, every
  `skills-lock.json` entry that names a Roundfix-owned skill describes the
  installed tree, and no other lock byte changes.
- No Task of any Task Graph settles `completed` while the repository
  Verification it can see is red, when that check is enabled.
- No Run without a QA gate reports Clean while its integrated tree fails the
  repository Verification.

## Core Features

1. The lock follows the bundle (source 1). After installing the owned
   skills, `baseline update` and `skills install --target project` rewrite
   `computedHash` of each existing lock entry that names an owned skill to
   the external skills digest of its installed folder, keep every other
   field and the key order, add no entry, and report the names they
   rewrote.
2. Every Task settles on the repository Verification (source 2). The
   configured repository Verification runs as the last command of every
   non-QA Task's settlement attempt, in a graph with or without a QA gate
   Task, when `verification.repository_at_settlement` is enabled.
3. The Integrated Repository Verification (source 2). A Run without a QA
   gate whose Task Capacity is above one and that completed two or more
   Tasks runs the repository Verification once on its integrated Run
   Worktree, and ends `Unresolved` with a named reason when it fails.

## User Experience

An adopter upgrades Roundfix and runs `roundfix baseline update --yes`; the
text result adds `Skills lock updated: 5` and one `- lock updated <name>`
line per rewritten entry after the installed skills, the JSON result lists
the same names under
`skills.lockUpdated`, and the repository's lock test passes without a hand
edit. `roundfix skills install --target project` prints one
`Updated skills-lock.json computedHash: <name>` line per rewritten entry. A
repository whose lock lists no owned skill sees no new line and no changed
byte.

An adopter whose Spec declares `qa: declined` runs `roundfix implement`; a
Task whose change breaks `make verify` gets the failure as Verification
Feedback in its one repair turn and, if still red, settles `failed`, so the
Run ends `Unresolved` instead of Clean. When two Tasks of one Wave each pass
and integrate red, the Run ends `Unresolved` with
`Integrated Repository Verification failed: make verify; diagnostics: <path>`
on stderr and keeps its Run Worktree.

Declared behavior changes: a Task of a graph without a QA gate now pays one
repository Verification per settlement and can fail on it; such a Run can pay
one more at its end and can end `Unresolved` where it ended Clean;
`baseline update` and `skills install --target project` can rewrite
`skills-lock.json` and print one more line. Gated graphs, the two in-process
Settlement Checks, the QA gate and every other `skills-lock.json` writer keep
their behavior.

## Non-Goals / Out of Scope

- Changing the `source`, `ref`, `sourceType` or `skillPath` of a lock entry,
  or adding lock entries for owned skills the lock does not list (ADR-0189
  addendum).
- Doctor or `skills check` reporting an owned-skill lock entry that differs
  from its tree; both keep reading owned skills by version.
- Running the Spec Consistency and authorization Settlement Checks in a
  graph without a QA gate (ADR-0342).
- Review Runs (`fetch`, `resolve`, `watch`), whose batch settlement is not a
  Task Graph.
- The "Outro ponto da 0024" paragraph of message `01a17d0a`, owned by Spec
  0330.

## Success Metrics

1. Success Metric: in a temporary repository whose `skills-lock.json` lists
   `qa-gate` with a stale hash and an external skill, `skills install
   --target project` leaves the `qa-gate` entry's `computedHash` equal to the
   digest of `.agents/skills/qa-gate`, its other fields and the external
   entry byte-identical, and a second run changes no byte.
2. Success Metric: the `baseline update` skills stage reports the rewritten
   names in text (`Skills lock updated:`) and JSON (`skills.lockUpdated`),
   and omits both when nothing changed.
3. Success Metric: a lock that is a symbolic link or malformed makes both
   commands fail with a message naming `skills-lock.json`, and the lock's
   bytes stay unchanged; a repository without a lock gets no lock.
4. Success Metric: in a graph without a QA gate, a Task's settlement runs its
   declared commands then the repository Verification; a repository failure
   returns as Verification Feedback and a final failure settles the Task
   `failed`; with `verification.repository_at_settlement: false` only the
   declared commands run; the in-process checks do not run.
5. Success Metric: a Run without a QA gate at Task Capacity 2 with two
   completed Tasks ends `Unresolved` when the integrated check fails, naming
   the command and its diagnostics file, and Clean when it passes; at Task
   Capacity 1, with one completed Task, or with the switch off, the check
   does not run.
6. Success Metric: a gated graph's settlement commands, in-process checks and
   QA precondition run exactly as before.
7. Success Metric: after release, Quaestor runs `baseline update` and its
   lock test passes without a hand edit, and a `qa: declined` Run whose
   repository Verification is red does not report Clean.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not write. First,
the adopter's reports quoted verbatim in the two adopted Backlog Entries
under `references/` (ai-memory messages `01a17b1a`, 2026-10-26, and
`01a17d0a`, 2026-10-27). Second, the published source of the skills CLI that
writes `skills-lock.json` (vercel-labs/skills, `src/local-lock.ts`,
retrieved 2026-10-27), whose passage the QA row quotes: "SHA-256 hash
computed from all files in the skill folder. Unlike the global lock which
uses GitHub tree SHA, the local lock computes the hash from actual file
contents on disk." and, for `sourceType`, "The provider/source type (e.g.,
"github", "node_modules", "local")". Third, Graydon Hoare, "not rocket
science (the story of monotone and bors)", 2014, whose passages the QA row
quotes: "automatically maintain a repository of code that always passes all
the tests" and "tested in isolation and then integrated on the basis of that
test (with no guarantee that the integrated combination works)". The QA
compares offline the rewritten hash with the first passage and the
integrated check with the second.

## Unreachable Acceptance

- criterion: Success Metric 7, after release, Quaestor runs `baseline
  update` and its lock test passes without a hand edit, and a `qa: declined`
  Run whose repository Verification is red does not report Clean.
  reason: it needs the released binary in another repository, which no
  Verification or QA row may read.
  satisfied-by: the release notice to `default/gesttione-solutions-quaestor`
  and that project's next upgrade, reported to this project through
  ai-memory.

## Validation

Characterized on 2026-10-27 at `5e9475bf` in the operator's scratch
directory; no test touched the real `~/.roundfix`.

- Lock: a scratch test in package `skills` wrote a `skills-lock.json` whose
  `qa-gate` entry carried Quaestor's pre-v0.125.0 hash `7434ca0f…`, called
  `Install` for the project target and compared the entry with the
  installed folder's digest. The entry still read `7434ca0f…`; the folder's
  digest was `c7fc87c2…`, the exact value Quaestor wrote by hand in
  gesttione-solutions/quaestor PR #61 (commit `7dae8fea`, read through
  `gh api` without cloning). `Install` (`skills/skills.go`) writes only skill
  files; `runBaselineUpdateSkillsStageWith` (`internal/cli/baseline_update.go`)
  and `runSkillsInstall` (`internal/cli/cli.go`) are its only callers, and
  neither reads the lock. Lock entries for owned skills come from the
  upstream tool: the embedded Setup Snapshots mark owned skills
  `type: repo`, `resolveExternalSkillRequirement` removes them from the
  external set, and this repository's own lock lists none.
- Settlement: `TaskCycle` sets `settlementChecks = qaTask != nil`, and
  `runTaskWorker` appends the repository Verification only when
  `enteredOnRedRepository || (plan.settlementChecks &&
  plan.RepositoryVerificationAtSettlement)`. The existing test
  `TestUngatedGraphIgnoresEnabledRepositorySettlementVerification` pins that
  a gate-less graph runs only its declared commands. In a scratch worktree,
  removing `plan.settlementChecks &&` from that condition failed exactly that
  test in `./internal/daemon`, and 21 implement tests of `./internal/cli`
  whose gate-less fixtures inherit the built-in `make verify` without a
  Makefile, one of them by hanging; the TechSpec names the cause and
  task_03 migrates them.
- Secondbrain, read only: `wiki/index.md`; `qmd query "skills-lock
  computedHash roundfix owned skills baseline update" --all --files
  --min-score 0.3` returned this repository's mirrored decision record on vendored skill
  replacement, the `baseline` and `skills` guides and the skills repository's mirrored
  `roundfix` reference, which confirms the upstream skills repository
  publishes copies of the owned skills; `qmd query "roundfix settlement
  repository verification qa declined clean"` returned the mirrored
  `roundfix` and `write-tasks` skills, the configuration guide and the
  changelog, none of which documents a gate-less settlement.
- Exa: the vercel-labs/skills source and Graydon Hoare's essay quoted under
  Acceptance evidence; issue vercel-labs/skills#806 records that the
  upstream tool's own hash can disagree with an installed copy, which is
  why the rewrite uses the digest Roundfix already computes.

## Glossary

- adds: **Integrated Repository Verification**
- changes: **Settlement Check**
- changes: **Clean**

## Decisions

Maintainer answers to the cycle AR questions (AskUserQuestion, 2026-10-27
about 11:40 BRT), verbatim:

- 2026-10-27: Scope — "A, B e C (Recommended)", option text: "Lock + Clean
  vermelho; família do settle + aceite da QA; profile proof sob carga. Todas
  são health e duas se repetem em adotantes." This Spec is A.
- 2026-10-27: QA acceptance — "Casar por linha (Recommended)", option text:
  "Uma declaração cobre cada criterion que nomeia; o accept confere linha a
  linha, não a contagem. write-prd/qa-gate passam a dizer isso." (applies to
  Spec 0330).
- 2026-10-27: Authorization — "Mesmas regras do AQ (Recommended)", option
  text: "Push de branches, PR dos planos, entrega e merge pela Delivery
  Queue, eu escolho o ponto de release e aviso os adotantes por ai-memory."

Autonomous route: no person present for the technical choices; decisions
taken from the operator's cycle AR brief for Spec 0329 of 2026-10-27 and the
characterization above.

- The lock keeps the source the upstream tool recorded and rewrites only the
  hash: the format's `sourceType` names a provider that can fetch the skill
  again and Roundfix is not one, and the adopter's own test requires
  `source` and `ref` unchanged (ADR-0189 addendum).
- The repository Verification at settlement runs in every graph, and a
  gate-less Run checks its integrated tree once, rather than refusing Clean
  without running anything (ADR-0342).
- Both sources fit four implementation Tasks plus QA; nothing moves to a
  Backlog Entry.

## Open Questions

None.
