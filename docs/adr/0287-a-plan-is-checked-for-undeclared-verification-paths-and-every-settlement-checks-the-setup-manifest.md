---
status: accepted
created_at: 2026-10-18T13:00:00Z
updated_at: 2026-10-18T13:00:00Z
deprecated_at: null
superseded_by: null
---

# A plan is checked for undeclared Verification paths, and every settlement checks the Setup Manifest

Cycle AH lost three Runs and a release check to defects that exist before any
code is written:

- **A Verification named a file that no longer existed.** On 2026-10-17, pull
  request #559 amended Spec 0274 to move two contract tests into existing
  files. task_03's Verification still grepped
  `internal/docscontract/glossary_contract_test.go`, a file the Task no longer
  declared, that no Task created and that did not exist. The Run failed twice
  with `grep: … No such file or directory`, and pull request #560 repaired the
  plan.
- **QA rows asked for more than the plan could show.** The QA gates of Spec
  0272 (pull request #554) and Spec 0274 (pull request #561) each failed one
  row. In 0272 the row asked for exactly 243 files, a count that Spec 0273's
  merged record moved to 244. In 0274 the row asked the glossary to name two
  check codes that no Task added.
- **Every Makefile edit left the Setup Manifest stale until release.** The
  Setup Manifest records a `declarationDigest` of the whole Makefile for the
  repository Verification and the incremental Verification. Specs 0247, 0248,
  0259 and 0277 each changed the Makefile. Each time the next release (v0.58.0,
  v0.59.0, v0.72.0 and v0.88.0) found `baseline: plan_ready: 1 file change`
  and rewrote only those two digests.

## Decision

**The Spec Consistency Check reports a Verification path the plan does not
provide.** `SC-VERIFY-PATH-UNDECLARED` is an error from the Tasks stage, with
no horizon. It covers a Verification command of a pending non-QA Task that
names a repository path meeting all of these conditions:

- The Task does not declare the path (`creates:`, `interface:`, `deletes:`,
  `instruction:` or `## Recorded paths`), and no declared path lies under it.
- No Task in the Task's dependency closure creates the path.
- The path does not exist in the checked tree.

A command "names a repository path" only when a shell word matches a path
shape whose first element is an existing top-level directory of the
repository. These words never count:

- the pattern operand of `grep` or `rg`;
- the operands of `printf` and `echo`;
- the operand of a negated existence test;
- a word holding `$`, a glob or `...`;
- a path under the Spec Root.

Measured on 2026-10-18 over the 247 archived Specs with Verification (3,187
commands), the rule raised 34 findings in 11 Specs:

- 19 findings in 7 Specs named a path the Task created without declaring it.
  That is the rule's own subject, and declaring the path answers it.
- 15 findings in 4 Specs were false positives: a build output the command
  itself writes, and a skill test directory that was not in Git. That is
  0.47% of commands and 1.6% of Specs.
- Specs 0131 to 0277 raised none.
- On the 0274 tree after pull request #559, the rule raised exactly task_03's
  `glossary_contract_test.go`.

**Every Task settlement checks the Setup Manifest.** A repository contract
test in an existing always-relevant documentation contract file builds the
same Managed Refresh plan the release check builds. It fails when the plan
holds any file change or history move, and its message names the regeneration
command. It runs in about 0.1 s, so every Task settlement, `make verify-docs`
and the Full Contract Run report a stale manifest before a release does. No
Makefile entry or CI workflow changes.

## Considered Options

- Treating every path-shaped word as a path was rejected. Grep patterns and
  printed messages hold path-shaped text that no command reads.
- Requiring every Verification path to be declared, even when it exists, was
  rejected. Verification reads many repository files a Task never changes.
- Asking the judge whether an implementation Task delivers each QA
  requirement was measured and rejected. The candidate was an advisory Noul
  question per QA requirement, against the requirements of every
  implementation Task. On 2026-10-18 it was asked of 74 QA requirements of
  Specs 0272 to 0277 at their planning commits, through OpenRouter, and every
  answer came from `typesafe/jev-1.13-20260917`, with no error.
  - Both real cases scored above any useful threshold, that is towards
    "delivered": 0272 requirement 2 at 0.88 and 0274 requirement 10 at 0.51. In 0272, task_04 itself stated the
    count of 243 files, so that defect was a wrong count, not an undelivered
    requirement.
  - Below 0.3, 14 of the 72 negatives were raised, every one a requirement
    the gate fulfils by itself: outside evidence, the non-waivable Pull
    Request row, the glossary check, and "MUST NOT accept a row whose only
    evidence …".
  - That is a recall of 0 of 2 and 19% false alarms at 0.3, and no threshold
    separates the real cases from the negatives. The question was dropped
    before the Run, as the maintainer had authorized.
- Keeping the manifest check in the Full Contract Run alone was rejected. No
  existing contract detects the drift, and the closest one, which takes 46 s,
  passes on the stale tree. Running the plan per settlement costs about a
  tenth of a second.
- Adding a Repository Contract Test file was rejected. It needs a Makefile
  entry, which this cycle's authorization excludes, and the nine always-relevant
  documentation contract files are pinned.

## Consequences

- A Task that changes the Makefile also regenerates the Setup Manifest, with
  `roundfix baseline update --repo . --no-skills --yes`, in the same Task.
- The path rule raises an undeclared creation as an error, so an author
  declares the file before the Run, not after it.
- A QA requirement that exceeds its Tasks is still found by the QA gate, not
  before the Run.
