---
task: task_02
spec: 0387-the-techspec-template-states-the-adr-alternatives-instruction-once
status: pending
type: qa
complexity: medium
---

# Task 02: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec deletes the repeated "Begin
the relation with" from the `## ADR Alternatives` comment of the
`write-techspec` TechSpec template and its shipped mirror, raises and records
the skill's version so adopters receive the correction through `roundfix
baseline update`, and adds a contract test that refuses a phrase of two or
more words repeated across a line break in any Roundfix-owned skill, while
`SC-ADR-RELATED` keeps reading its token from the Spec's own TechSpec.

Every row runs on this repository's own files, temporary homes, temporary
repositories and temporary copies of the audited head or the Delivery Base.
No row starts an Agent session or an ACP adapter, calls ai-memory, reads or
writes the real `~/.roundfix` or the real home directory, runs the real `gh`,
reaches the network, reads another repository, sends a request to a provider,
the judge or a reviewer, or reads a credential. Every `go test` it runs
carries an explicit `-timeout`. The `go test` package and test runs keep the
session's own `HOME` and `TMPDIR`; only rows that run `bin/roundfix` set
`HOME` and `TMPDIR` to fresh temporary directories, with `GOMODCACHE` and
`GOCACHE` kept at their `go env` values, unset `GH_TOKEN`, `GITHUB_TOKEN`,
`TYPESAFE_API_KEY` and every Judge Key Variable for that command, and remove
the temporary directories afterwards. The QA never writes into the checkout
while a test suite runs, keeps its evidence under names the repository does
not ignore, `.txt` for command output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 2 to 4 and Invariants 1 and 3 by running
   `go test -count=1 -timeout 20m -v ./skills` and recording a pass for each
   of the six cases of
   `TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak`, for
   `TestNoOwnedSkillRepeatsAPhraseAcrossALineBreak/canonical` and
   `/shipped`, `TestEveryOwnedSkillVersionIsRecorded` and
   `TestAuthorialSkillSync/write-techspec`; then
   `TestAPortugueseNotApplicableADRDrawsInNoCitingADR`,
   `TestAMarkerOutsideTheClosedListStillDrawsIn`,
   `TestANotApplicableADRDrawsInNoCitingADR` and
   `TestAnAppliedADRStillDrawsInItsCitingADR` over `./internal/speccheck`,
   and `go test -count=1 -timeout 20m -tags docscontract
   ./internal/docscontract`, recording each pass.
3. MUST verify Success Metric 1 and Invariant 2:
   - record `git diff` of
     `.agents/skills/write-techspec/references/techspec-template.md` from the
     Delivery Base (the merge base of the audited head with the default
     branch) to the audited head: exactly one line removed and one added, the
     added line equal to the comment line of `_techspec.md` → Interfaces;
   - record that `cmp` finds the canonical template and
     `skills/write-techspec/references/techspec-template.md` identical;
   - record that the whitespace-normalized comment reads "Begin the relation
     with" once, followed by `not applicable` and then `não se aplica`.
4. MUST verify Success Metric 3 from the history: record the `version` and
   `metadata.version` of both `write-techspec` `SKILL.md` copies at the
   Delivery Base and at the audited head, that the head's value is higher and
   the same in all four fields, and that
   `skills/testdata/owned-skill-versions.json` records it.
5. MUST exercise the shipped skill through the built product: with `HOME`
   and `TMPDIR` at fresh temporary directories, in a fresh temporary Git
   repository, run `bin/roundfix skills install --target project` and record
   that the installed `write-techspec` `SKILL.md` carries the head's version
   and that its `references/techspec-template.md` holds the corrected comment
   and not "Begin the relation with" twice in a row.
6. MUST reproduce the characterization on the Delivery Base: copy
   `skills/repeated_phrase_test.go` from the audited head into a temporary
   copy of the Delivery Base made with `git archive`, run
   `TestNoOwnedSkillRepeatsAPhraseAcrossALineBreak` there with `-timeout`,
   and record that it fails with exactly two reports, the canonical and the
   shipped `write-techspec/references/techspec-template.md` at line 87, each
   naming "Begin the relation with".
7. MUST run each sabotage of `_techspec.md` → Sabotages on its own temporary
   copy of the audited head made with `git archive`, with `-timeout` on every
   `go test`, and record that the named test fails on the named case:
   sabotages 1 to 7 as listed there.
8. MUST execute the tests task_01 names in its Verification.
9. MUST record, as outside-evidence rows, two sources this Spec did not
   write, through the passages `_prd.md` → Acceptance evidence quotes:
   - Argus's report, read from
     `references/2026-11-12-argus-the-techspec-template-repeats-begin-the-relation-with.md`
     in the audited checkout: compare "A linha 87 termina com ela e a linha
     88 começa de novo com ela" with requirement 6's two reports at line 87;
   - Matt Might, "3 shell scripts: Kill weasel words, avoid the passive,
     eliminate duplicates" (read 2026-11-12): compare "a second instance of
     the word "the" / when it starts a new line" and "cut and paste
     accidents" with the guard, which compares the end of each line with the
     start of the next without regard to case.
10. MUST record the PRD's one Unreachable Acceptance entry, an adopter's
    managed refresh after the release, as one row counted in
    `rows_blocked_declared`, with its `satisfied-by`.
11. MUST perform the glossary and documentation checks: record that the
    Spec declares no glossary term, that `CONTEXT.md` has no change in the
    Task commit range, and that `TestGlossaryPhrasesAreDefinedTerms` and
    `TestTheSkillCoverageMapIsCurrent` passed in requirement 2.
12. MUST verify from the repository history:
    - the Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`, `CONTEXT.md`,
      `docs/references/coverage-record.json`,
      `docs/references/behavior-surfaces.json`, any file under `internal/`,
      `docs/agents/`, `docs/user-guide/` or `docs/adr/`, any skill other
      than `write-techspec`, any production `.go` file;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
13. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
14. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0387-the-techspec-template-states-the-adr-alternatives-instruction-once --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real home directory;
    - run the real `gh`, call ai-memory or reach the network;
    - send any message;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests, the history checks and the install row.
- [ ] Reproduce the characterization on the Delivery Base and run each sabotage on a temporary copy.
- [ ] Record the outside sources, the unreachable row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Both template copies, and the copy `skills install` writes, state the instruction once with both tokens, at a raised and recorded version.
- [ ] The guard reports the Argus line on the Delivery Base, nothing on the audited head, and each sabotage makes its named test fail.
- [ ] The scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0387-the-techspec-template-states-the-adr-alternatives-instruction-once/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-4; Invariants 1-3; Sabotages; Testing Approach; Build Order 2
