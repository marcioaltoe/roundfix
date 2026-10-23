---
task: task_01
spec: 0300-reconcile-releases-a-stale-run-owner-record
status: completed
type: docs
complexity: low
---

# Task 01: The reconcile guide, the skill reference and the glossary describe the Run owner record

## Overview

The `reconcile` user guide and the Roundfix skill's `reconcile` reference
describe `process`, `runBranch` and `staging` candidates but no candidate for
a terminal Run whose owner PID now belongs to another process, and the
glossary has no term for that record. This Task describes the
`runOwnerRecord` candidate and its signal-free release in both documents,
raises the Roundfix skill's version, re-records the Behavior Surface Record,
adds **Stale Run Owner Record** to the glossary and revises **Reconcile
Command** (`_prd.md` → Core Feature 4).

## Requirements

1. MUST add to `docs/user-guide/commands/reconcile.md`, right after the
   paragraph that ends "`--apply` releases stale staging, and
   `--carry-forward` releases it before creating that mode's own staging
   worktree.", this paragraph:
   "Reconcile also proves the recorded owner of each terminal Run. When the
   live process at the owner PID presents a start identity that differs from
   the one the Run recorded, the system gave that PID to another process and
   the recorded owner ended: the Run's owner PID and identity are a Stale Run
   Owner Record. The report lists it once in `processCandidates` with
   `kind` `runOwnerRecord`, the owner PID, an empty `processIds` list and
   that proof, and the text report prints `Run owner record candidate: Run
   <run-id> owner PID <pid>`. `--apply` releases it: it clears the Run's
   recorded owner PID and identity, journals what it cleared, and sends no
   signal; it never inspects the processes of the PID's new holder.
   `debrisSummary.runOwnerRecordsReleased`, `run-owner-records-released` in
   the text summary, counts the releases. An owner identity that cannot be
   read or compared, and a Run recorded without one, stay in
   `preservedCandidates` with the refusal reason, and the Stop Command still
   never signals such a PID (ADR-0313)." The ADR reference takes the guide's
   link form, a relative link to `../../adr/0313-a-proven-owner-identity-mismatch-releases-a-terminal-runs-owner-record-without-a-signal.md`,
   as the guide's other ADR references do.
2. MUST add to the canonical Roundfix reference
   `.agents/skills/roundfix/references/reconcile.md`, right after the
   `staging` candidate bullet that ends "an unlocked legacy staging worktree
   is preserved.", this bullet:
   "- A `runOwnerRecord` candidate is a Stale Run Owner Record: the live
   process at a terminal Run's recorded owner PID presents a start identity
   that differs from the recorded one, which proves the recorded owner ended.
   Its report names the Run outcome, owner PID, an empty process list and that
   proof. `--apply` clears the Run's recorded owner PID and identity, journals
   them, sends no signal, and counts the release in
   `debrisSummary.runOwnerRecordsReleased`; an identity that cannot be read
   or compared stays in `preservedCandidates`."
   In the same file, MUST replace "or proven process, Run Branch, and staging
   candidates" with "or proven process, Run owner record, Run Branch, and
   staging candidates". No text inside any `### QA settlement` section of any
   skill may change.
3. MUST add **Stale Run Owner Record** to `CONTEXT.md` through
   `domain-modeling`, right after the **Process Residue** entry, as
   "**Stale Run Owner Record**:" followed on the next line by "The owner PID
   and start identity a terminal Run still records after the live process at
   that PID has presented a different start identity, which proves the
   recorded owner ended and the PID now names another process. The Reconcile
   Command reports it, and `--apply` releases it by clearing the Run's owner
   PID and identity and journaling them, without sending any signal
   (ADR-0313)." and the line "_Avoid_: Orphan PID, dead owner lock, staging
   owner record".
4. MUST revise **Reconcile Command** in `CONTEXT.md` through
   `domain-modeling` by appending to its definition, after "and
   `--carry-forward` performs Task Carry-Forward.", the sentence "`--apply`
   also releases each Stale Run Owner Record without sending a signal
   (ADR-0313)." Its `_Avoid_` line keeps its bytes. In both entries the ADR
   reference takes the glossary's link form, to
   `docs/adr/0313-a-proven-owner-identity-mismatch-releases-a-terminal-runs-owner-record-without-a-signal.md`.
5. MUST run these commands, because this Task edits an owned skill and a
   guide. Before them, raise both version fields (`metadata.version` and the
   top-level `version`) of the canonical `SKILL.md` of the Roundfix skill by
   one patch over the value the tree holds when this Task starts, because the
   record command refuses changed content under an unchanged version. The
   Agent runs them itself before handing back, because the Daemon never
   regenerates a record; no mirror or record is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     which re-records `docs/references/behavior-surfaces.json` because the
     `reconcile` guide's and reference's fingerprints change.
6. MUST search the guides, the agent guides and the skill copies for wording
   that says every `processCandidates` entry names live processes or that a
   mismatched owner identity is only ever preserved, and report what it found.
7. MUST NOT change any Go source or ADR, `docs/references/skill-coverage.json`
   or any other guide page.

## Subtasks

- [ ] Add the Run owner record paragraph to the `reconcile` guide and the bullet and sentence to the Roundfix `reconcile` reference.
- [ ] Add **Stale Run Owner Record** and revise **Reconcile Command** in the glossary.
- [ ] Raise the skill version, sync the mirrors, record the version and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The guide and both copies of the reference describe the `runOwnerRecord` candidate, its signal-free release and its count.
- [ ] `CONTEXT.md` defines **Stale Run Owner Record**, revises **Reconcile Command**, and links ADR-0313 from both.
- [ ] The skill version is raised and recorded, the mirrors are in sync, and the Behavior Surface Record is current.

## Consumers

None. This Task changes documentation, the skill mirrors, the skill version record and the Behavior Surface Record only; no code consumes them.

## Context

- instruction: `docs/adr/0313-a-proven-owner-identity-mismatch-releases-a-terminal-runs-owner-record-without-a-signal.md`
- instruction: `internal/docscontract/skill_coverage_test.go`
- instruction: `docs/references/skill-coverage.json`
- interface: `docs/user-guide/commands/reconcile.md`
- interface: `.agents/skills/roundfix/references/reconcile.md`
- interface: `skills/roundfix/references/reconcile.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `for f in docs/user-guide/commands/reconcile.md .agents/skills/roundfix/references/reconcile.md skills/roundfix/references/reconcile.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'runOwnerRecord' || { printf 'missing candidate kind in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'debrisSummary.runOwnerRecordsReleased' || { printf 'missing release count in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'sends no signal' || { printf 'missing signal-free release in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/reconcile.md skills/roundfix/references/reconcile.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'proven process, Run owner record, Run Branch, and staging candidates' || { printf 'missing apply sentence in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Stale Run Owner Record**:' || { echo 'missing glossary term' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'also releases each Stale Run Owner Record without sending a signal' || { echo 'missing Reconcile Command revision' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no document names the `runOwnerRecord` candidate and the glossary lacks **Stale Run Owner Record**, so the command fails at its first phrase check. Afterwards the guide and both reference copies describe the candidate, its count and its signal-free release, both reference copies name it among what `--apply` releases, the glossary defines **Stale Run Owner Record** and revises **Reconcile Command**, the mirrors are in sync, the skill version is recorded, and the documentation contracts, the Skill Coverage Map among them, pass.

## References

- `_prd.md` → Core Feature 4; Goals; Glossary
- `_techspec.md` → API Contract 7; Build Order 1
- ADR-0313

## Result

Implemented the documentation slice for the Stale Run Owner Record. The
reconcile guide now documents the `runOwnerRecord` candidate, its proof,
signal-free `--apply` release, preserved refusal cases, text output and
`debrisSummary.runOwnerRecordsReleased` count. The canonical Roundfix
reference has the matching candidate bullet and apply sentence, and
`make skills-sync` produced identical mirror content. `CONTEXT.md` now defines
Stale Run Owner Record and revises Reconcile Command with the ADR-0313 links.

Focused implementation evidence:

- `make skills-sync` — passed.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — passed; recorded Roundfix skill version `0.1.130` in both version fields and `skills/testdata/owned-skill-versions.json`.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1` — passed; re-recorded the reconcile guide fingerprint in `docs/references/behavior-surfaces.json`.
- `cmp` checks — canonical and mirror `SKILL.md` files and reconcile references are identical; `git diff --check` passed.
- Required wording search — no guide, agent-guide, or skill-copy wording was found that says every `processCandidates` entry names a live process. Existing preserved-candidate wording remains for ambiguous or unreadable evidence; ADR-0313 and the Spec materials contain the historical preserved behavior and the new exception.

Acceptance evidence:

- The guide and both reference copies contain `runOwnerRecord`,
  `debrisSummary.runOwnerRecordsReleased`, and signal-free release wording;
  both references include Run owner record in the `--apply` release list.
- `CONTEXT.md` contains the Stale Run Owner Record glossary entry, the
  Reconcile Command revision, and ADR-0313 links in both entries.
- The canonical skill is `0.1.130`, the mirrors match, and the owned-skill
  version record plus Behavior Surface Record were regenerated by their
  required commands.

The declared Task Verification commands remain for the Daemon; they were not
run in this handoff.

## Carry-forward provenance

- Source Run: `run_20261023T014935Z_334578791ca9af22`
- Source commit: `7ccf2b1f669d4cacc2f958dbdcf8c21caa91745a`
