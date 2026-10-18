---
task: task_05
spec: 0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec delivers three things:

- the Correction Review of a fixed round-2 finding, and the retry exemption
  that reaches it;
- `SC-VERIFY-PATH-UNDECLARED`;
- the Setup Manifest contract that runs at every Task settlement.

Every row runs on temporary repositories and homes. No row reads or writes the
real `~/.roundfix`, reaches the network, sends a request to a provider, the
judge or a reviewer, reads a credential, or reads another repository. The
gate never runs `roundfix spec judge`, never runs `roundfix review` against a
real provider, and never runs an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1:
   - Run the tests task_01 names, and record from their output that exactly
     one fake reviewer call reviewed the round-2-head-to-head delta, that its
     record carries `correction: true`, and that a further review made no
     call.
   - Make a temporary copy of the audited head with `git archive`. In it, make
     the applicability check of API Contract 1 always false, and record that
     `TestCorrectionReviewReviewsTheFixOnce` fails (Invariants 1 and 2).
3. MUST verify Success Metric 2:
   - Run the retry tests task_01 names, and record that the fixed round-2 park
     outside the declared paths is accepted and the dismissed case is refused
     with today's text.
   - On a temporary copy of the audited head, remove the exemption from
     `ProveReviewCorrection`, and record that
     `TestRetryAdmitsAFixedRoundTwoParkOutsideDeclaredPaths` fails.
4. MUST verify Success Metric 3 through the built `bin/roundfix` with a
   temporary home:
   - Make a temporary Git repository from `git archive d25aee62`, the 0274
     tree after pull request #559. Record that
     `bin/roundfix spec check 0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd --stage tasks --format json`
     reports `SC-VERIFY-PATH-UNDECLARED` for task_03 naming
     `internal/docscontract/glossary_contract_test.go`.
   - Make a second temporary repository from `git archive f363ad6d`, the
     repaired tree, and record that the same command reports no such finding.
     Both trees are this repository's history, which no Task of this Spec
     authored.
   - Record that the corpus golden counts 0 for the new code and that the
     check tests task_02 names pass (Invariant 3).
5. MUST verify Success Metric 4:
   - Record that `TestTheSetupManifestIsCurrent` passes on the audited head,
     and how long the `internal/docscontract` package takes with it.
   - Make a temporary Git repository from `git archive` of the audited head
     and append one byte to its Makefile. Record that the test fails there,
     naming `docs/agents/setup-context.json` and the regeneration command
     (Invariant 4).
6. MUST execute the tests task_01, task_02 and task_04 name in their
   Verification.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence.
   - Read the repository-history row offline with `git show --numstat` on the
     eight named commits, and record whether each release commit changed only
     the two `declarationDigest` lines.
   - Because this Spec's authorization forbids network access, record each
     web source as `blocked (environment: network prohibited by Spec
     authorization)`, quoting the passage the PRD cites. Such a row never
     decides a qualifying `partial`.
8. MUST perform the glossary check of `docs/agents/domain.md`:
   - record that `CONTEXT.md` defines **Correction Review** and carries the
     revised **Reviewer Lineage**;
   - record that the glossary contract tests pass;
   - record whether any other term the Spec introduced needs an entry.
9. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every Governed Path changed is bounded in `_authorization.md`;
    - these did not change: the Makefile, the CI workflows, `go.mod`,
      `go.sum`, `.roundfixrc.yml`, every Baseline module,
      `docs/agents/setup-context.json`, `internal/cli/spec_judge.go`,
      every file under `internal/judge/` and `internal/delivery/engine.go`;
    - the Roundfix and `write-tasks` skill versions were raised and recorded;
    - no `### QA settlement` section changed;
    - no guide or skill copy still says, without the Correction Review
      exception, that a third review in a lineage never calls the reviewer.

    This row reads Task commits, so it declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
11. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers --plan --offline`
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
- [ ] Replay the 0274 trees and the stale Makefile through the built binary and the contract.
- [ ] Record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A fixed round-2 finding gets exactly one Correction Review, and each removed guard makes its named test fail.
- [ ] The 0274 tree after pull request #559 raises the new code, and the repaired tree does not.
- [ ] A one-byte Makefile change fails the Setup Manifest contract, and the scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; Invariants 1-4; Current behavior; Build Order 4
- ADR-0286; ADR-0287; ADR-0197; ADR-0272; ADR-0104; ADR-0276
