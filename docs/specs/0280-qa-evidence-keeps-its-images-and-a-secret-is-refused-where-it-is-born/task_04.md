---
task: task_04
spec: 0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec delivers three things:

- readable QA evidence (text, PNG, JPEG, SVG), and the QA Report commit's
  refusal of any other binary under `qa/evidence/`;
- the Secret Scan of every Task commit, with the fixture guidance;
- the Secret Rewrite that a Delivery Retry of a `secret-found` item accepts,
  and the next action that describes it.

Every row runs on temporary repositories and homes. No row reads or writes the
real `~/.roundfix`, reaches the network, sends a request to a provider, the
judge or a reviewer, reads a credential, or reads another repository. The
gate never runs `roundfix spec judge`, never runs `roundfix review` against a
real provider, and never runs an archive plan without `--offline`. Every
secret-shaped sample a row needs is built at runtime from parts and is never
written into the report or the evidence.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1:
   - Run the speccheck tests task_01 names, and record from their output that
     PNG, JPEG and SVG evidence raised no `QA-EVIDENCE-PATH` while each
     non-image binary, including the ELF header named `.png`, raised one.
   - Make a temporary copy of the audited head with `git archive`. In it, make
     `ReadableEvidence` return false for the PNG signature, and record that
     `TestEvidenceAcceptsTrackedImages` fails (Invariant 1).
3. MUST verify Success Metric 2:
   - Run the Daemon tests task_01 names, and record that the PNG rode in the
     QA Report commit and that the non-image binary settled the QA Task
     `failed` with `QA evidence carries a non-image binary: <path>` and no
     commit.
   - On a temporary copy of the audited head, remove the QA Report commit
     refusal, and record that
     `TestQAReportWithANonImageBinarySettlesFailedWithoutACommit` fails.
4. MUST record, as evidence this Spec did not author, the classifier's answer
   for four files of this repository's history at `513b22eb`:
   - `docs/history/specs/0005-tui-cockpit/design/roundfix-01.png`;
   - `docs/history/specs/0035-agent-selection-profiles/references/artificial-analysis-coding-agent.jpeg`;
   - `docs/history/specs/0090-a-gate-that-could-have-failed/qa/evidence/2026-08-09-spec-0090/positive-home/.roundfix/roundfix.db`;
   - `docs/history/specs/0209-sources-that-share-a-context-share-a-spec/qa/evidence/2026-10-01-gate/zhang-author.pdf`.

   Extract each with `git show` into a temporary directory, record its first
   16 bytes in hexadecimal, and run a throwaway test in a temporary copy of
   the audited head that prints `ReadableEvidence` of each file's first
   `EvidenceProbeBytes` bytes. Record that the PNG and the JPEG are readable
   and the database and the PDF are not.
5. MUST verify Success Metric 3:
   - Run the tests task_02 names, and record that the matching Task settled
     `failed` with `Task commit carries a secret: <path>:<line>:
     credential-assignment`, that no commit was made, and that the edit
     beside an existing match and the letters-only fixture committed.
   - On a temporary copy of the audited head, remove the Task commit scan,
     and record that
     `TestTaskCommitWithASecretMatchSettlesFailedWithoutACommit` fails
     (Invariant 2).
   - Record that no tracked file outside `docs/history/` at the audited head
     holds a Secret Match, using a throwaway test in a temporary copy that
     runs `secretscan.ScanContent` over `git ls-files`.
6. MUST verify Success Metric 4:
   - Run the tests task_03 names, and record that the rewrite of only the
     matched lines returned the item to `reviewing` with the head as its
     newest candidate. Record each refused rewrite with the archived-head
     reason and its rule.
   - On a temporary copy of the audited head, make `ProveSecretRewrite`
     accept every head, and record that
     `TestSecretRewriteRefusesAnyOtherChange` fails (Invariant 3).
   - Record the `secret-found` next action printed by the park tests, and
     that `grep -rF` finds the old action "remove the secret from every
     commit the blocker names" nowhere under `docs/user-guide`, `docs/agents`,
     `.agents/skills`, `skills` or `internal`.
7. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
8. MUST record, as evidence this Spec did not author, each published source
   of `_prd.md` → Acceptance evidence. Because this Spec's authorization
   forbids network access, record each web source as `blocked (environment:
   network prohibited by Spec authorization)`, quoting the passage the PRD
   cites. Such a row never decides a qualifying `partial`.
9. MUST perform the glossary check of `docs/agents/domain.md`:
   - record that `CONTEXT.md` defines **Secret Rewrite** and carries the
     revised **Secret Scan** and **Mechanical Refusal Code**;
   - record that the glossary contract tests pass;
   - record whether any other term the Spec introduced needs an entry.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every Governed Path changed is bounded in `_authorization.md`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module,
      `docs/agents/setup-context.json`, `internal/secretscan/patterns.go`,
      `internal/cli/spec_judge.go` and every file under `internal/judge/`;
    - the Roundfix, `qa-gate` and `implement-task` skill versions were raised
      and recorded;
    - no `### QA settlement` section changed;
    - no guide or skill copy still says that screenshots under `qa/evidence/`
      are refused, or that the Secret Scan covers only the QA Report commit and
      pushes.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It MUST NOT
    do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - start a workflow.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Classify the four history files and scan the tracked tree on temporary copies.
- [ ] Record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Images pass and every other binary is refused, both by the mechanical stage and at the QA Report commit, and each removed guard makes its named test fail.
- [ ] A Task that adds a matching line fails at its own settlement, and an honest fixture commits.
- [ ] A Secret Rewrite returns the archived item to review, any other rewrite is refused with its rule, and the scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-4; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; API Contract 8; API Contract 9; Invariants 1-4; Current behavior; Build Order 4
- ADR-0291; ADR-0292; ADR-0273; ADR-0104; ADR-0276
