---
status: accepted
created_at: 2026-11-09T12:00:00Z
updated_at: 2026-11-09T12:00:00Z
deprecated_at: null
superseded_by: null
---

# A QA Report relies only on evidence its own Run produced

On 2026-11-06 Quaestor reported a terminal QA gate that failed its Spec on
evidence it never produced. An earlier Run of the same Spec had committed an
evidence directory, `qa/evidence/2026-11-06-gate/`, on the item branch. The
later Run's gate cited it: one file ran the console script inside the earlier
Run's Run Worktree (`…/run_20261106T154834Z_…/.venv/bin/quaestor`), and a
replay script pinned a head older than the Task that fixed the text under
test. At the audited head the test that Task added passed. The item parked
`run-unresolved` and was retried by hand. Nothing between the gate's turn and
the QA Report commit asked whether the cited evidence belonged to the Run that
cited it; on 0.137.1 `roundfix qa-report accept` accepts such a report.

## Decision

**A QA Report relies only on evidence its own Run produced, at its audited
head, with executables inside its own Run Worktree.** Evidence a report
relies on that fails this is **Borrowed Evidence**, named by one of four
reasons:

- `reused`: a local link into the Spec's `qa/evidence/` names a file whose
  bytes the audited head already holds at that path;
- `imported`: in a Run, the same link names a file whose bytes equal the copy
  the QA Report Import placed there;
- `outside the worktree`: an entry of the report's `user_flow_executables`
  frontmatter list lies outside the report's worktree, read lexically, so a
  virtual-environment link to a system interpreter stays inside;
- `names another Run Worktree <path>`: the report, or a file it links into
  `qa/evidence/`, names a path inside a registered worktree of the repository
  whose branch is another Run's Run Branch.

A link inside a Carried Row, its Results line or its `### <row-id>` block,
is exempt: its Evidence Snapshot already proves that row's evidence
unmoved. The audited
head is `HEAD` while the report is uncommitted, and otherwise the first parent
of the commit that added the report. Outside Git nothing is checked.

**Every QA settlement path applies it, after the evidence-link rule of
ADR-0351.** `roundfix qa-report accept` and the Settle Command for a `qa`
Task refuse with `QA Report evidence does not belong to this Run: <subject>
(line <n>, <reason>)`, joined by `; `. The Daemon returns the QA turn once,
naming each item, the Run Worktree, the audited head and the pass's evidence
directory, `qa/evidence/<report name>/`; an item left at the QA Report commit
boundary settles the QA Task `failed` with that text and creates no QA Report
commit. `roundfix archive` is unchanged and never re-judges an archived
report.

**The gate records its executables.** The qa-gate skill writes each pass's
evidence under a directory no earlier pass used, builds and runs every
executable from its own Run Worktree, and lists their paths in
`user_flow_executables`. A report without the key, every report written
before this rule among them, is read as declaring none.

## Considered Options

- Prose in the qa-gate skill only. The Quaestor gate already had the skill's
  rule to record the build it audited and still cited another Run's files;
  only a check at the commit boundary stops the next one.
- Compare a commit recorded inside each evidence file with the audited head.
  Evidence formats belong to adopters (Go test output, Python replay scripts,
  Node logs); no reader can find a commit in all of them without guessing,
  and evidence legitimately names other commits, such as the Delivery Base.
- Require every evidence path to sit under the report's own directory. It
  would refuse every adopter report written before the convention, and a
  copied file would still pass; comparing bytes with the audited head refuses
  exactly the reused file wherever it sits.
- Flag every other worktree, the user checkout included. Evidence names the
  user checkout for honest reasons, such as the QA prompt's own fact; another
  Run's Run Worktree is never this Run's business, so only those are named.
- Resolve executable symbolic links before comparing. A Python virtual
  environment links `bin/python` to an interpreter outside the worktree, so a
  resolved comparison would refuse every honest Python gate.
