---
task: task_04
spec: 0367-the-strict-check-reads-unreachable-acceptance-and-related-adrs-as-written
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec adds
`SC-UNREACHABLE-UNREAD`, which names text in a PRD's `## Unreachable
Acceptance` section that no fielded declaration reads, through one reader in
`internal/spec` that the QA gate shares; it makes `SC-ADR-RELATED` seed its
candidates only with the listed ADRs whose TechSpec ADR Alternatives line
does not begin with `not applicable` (ADR-0377); and it documents both in
ADR-0093, ADR-0168, `CONTEXT.md`, two user guides and three owned skills.

Every row runs in this repository or in temporary repositories it creates
and removes. No row reaches the network, reads another repository, starts an
Agent session or an ACP adapter, reads or writes the real `~/.roundfix`,
sends a request to a provider, the judge or a reviewer, records an
environment value, or sends a signal to any process. Every `go test` carries
an explicit `-timeout` and keeps the session's own `HOME` and `TMPDIR`. Every
row that runs `bin/roundfix` or a binary built from the Delivery Base sets
`HOME` and `TMPDIR` to fresh temporary directories, with `GOMODCACHE` and
`GOCACHE` kept at their `go env` values, and removes them afterwards. The QA
never writes into the checkout while a test suite runs. The gate keeps its
evidence under names the repository does not ignore, `.txt` for command
output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST run, with `go test -count=1 -timeout 20m -v`, the tests tasks 01
   and 02 name in `./internal/spec` and `./internal/speccheck`, then the
   whole `./internal/spec`, `./internal/speccheck` and `./skills` packages,
   `-tags docscontract ./internal/docscontract` and
   `-run 'SpecCheck' ./internal/cli`, recording each pass and that none of
   the nine new tests reports `SKIP`.
3. MUST build a second binary from the Delivery Base (the merge base of the
   audited head with `main`), from a `git archive` of that commit unpacked
   in a temporary directory, with `go build -buildvcs=false`, to characterize
   the starting behavior in rows 4 and 5.
4. MUST record Success Metric 1, the outside-evidence row on Spec 0364's
   planned PRD: in a temporary Git repository, write
   `git show 4dce106d:docs/specs/0364-a-rerun-finishes-a-partial-npm-publication/_prd.md`
   to `docs/specs/0364-a-rerun-finishes-a-partial-npm-publication/_prd.md`
   and run `spec check 0364-a-rerun-finishes-a-partial-npm-publication --strict`
   with both binaries. The delivered binary reports exactly one
   `SC-UNREACHABLE-UNREAD`, at line 165, whose summary quotes
   `1. A real rerun of a failed release job on GitHub Actions against the npm`
   and whose fix holds `"- criterion: <text the QA row repeats>"`; the
   Delivery Base binary reports no `SC-UNREACHABLE-UNREAD`. Then replace
   lines 165 to 168 of the copy with the same entry in the fielded form, a
   top-level `- criterion:` line followed by indented `reason:` and
   `satisfied-by:` lines made of the entry's own sentences, and record that
   the delivered binary reports no `SC-UNREACHABLE-UNREAD`. Record where the evidence came from:
   a PRD written for another Spec, whose numbered entry cost its first QA
   gate on 2026-11-06.
5. MUST record Success Metric 3 on the built binaries: in a temporary
   repository holding accepted fixture decisions ADR-<A>, ADR-<B> citing ADR-<A> and
   ADR-<C> citing ADR-<B>, and a Spec whose PRD lists ADR-<A> and
   ADR-<B> and whose TechSpec has `- ADR-<B>: not applicable; …`, the
   delivered binary's `spec check <slug> --strict` reports no
   `SC-ADR-RELATED` and the Delivery Base binary reports one for ADR-<C>;
   with the line changed to `- ADR-<B>: applied; …`, and again without the
   TechSpec, both binaries report one for ADR-<C>, and the delivered one's
   fix holds `begin its ADR Alternatives line in the TechSpec with "not applicable"`.
6. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails:
   - `appliedADRNumbers` never removes a number:
     `TestANotApplicableADRDrawsInNoCitingADR` fails;
   - `notApplicableRelationPattern` loses its `^` anchor:
     `TestAnAppliedADRStillDrawsInItsCitingADR` fails on the relation that
     only contains the words later;
   - `notApplicableRelationPattern` loses `(?i)`:
     `TestANotApplicableADRDrawsInNoCitingADR` fails on `Not applicable`;
   - `checkAuthoringStage` no longer calls `detectUnreachableDeclarations`:
     `TestANumberedUnreachableEntryIsReportedAtItsFirstLine` fails at
     `StagePRD`;
   - every unread line is its own finding, with no grouping:
     `TestANumberedUnreachableEntryIsReportedAtItsFirstLine` fails;
   - HTML comment lines count as unread:
     `TestAnUnreachableSectionReadAsWrittenReportsNothing` fails on the
     template comment;
   - a `None.` paragraph is read even beside a declaration:
     `TestNoneBesideADeclarationAndAPlainBulletAreReported` fails;
   - the missing-field finding is dropped:
     `TestProseAfterADeclarationAndAMissingFieldAreReported` fails.
7. MUST execute the tests every non-QA Task names in its Verification.
8. MUST record, as a second outside-evidence row, sources this Spec did not
   write, with their dates and not blocked: the Fiscus report quoted in
   `_prd.md` → Acceptance evidence ("O fecho levou 7 rodadas e 72 ADRs"),
   compared offline with row 5, where an ADR marked `not applicable` draws in
   nothing; and the CommonMark passage that section quotes ("An ordered list
   marker is a sequence of 1–9 arabic digits"), compared with row 4, where a
   `1.` entry is reported.
9. MUST record Success Metric 4: the delivered
   `bin/roundfix spec check 0367-the-strict-check-reads-unreachable-acceptance-and-related-adrs-as-written --strict`
   on the audited head reports no finding.
10. MUST record the PRD's one Unreachable Acceptance entry (Fiscus's Spec
    0081 passing the strict check after the release) as exactly one row,
    `blocked (declared: Fiscus Spec 0081)`, its leading label,
    counted once in `rows_blocked_declared`, supported by rows 5 and 8.
11. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare `- changes:` for **Unreachable Acceptance**
    and **Spec Consistency Check**, that the Task commits change exactly
    those two entries of `CONTEXT.md`, and that
    `go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryStatesEachRuleOnce)$'`
    passes.
12. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - ADR-0093 and ADR-0168 carry the `Amended by [ADR-0377]` paragraph and
      `updated_at: 2026-11-09T12:00:00Z`;
    - the Task commit range changes none of these: the Makefile, any
      workflow, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/**`, any file
      under `internal/baseline/`, `internal/agent/` or
      `internal/suiteguardcontract/`, `internal/spec/qa*.go`,
      `internal/spec/archive.go`, any guide under `docs/agents/`, any skill
      other than `roundfix`, `write-techspec` and `write-tasks`, any ADR other
      than ADR-0093 and ADR-0168, or text inside the `### QA settlement`
      section of any skill.

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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0367-the-strict-check-reads-unreachable-acceptance-and-related-adrs-as-written --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read another repository, or run `git fetch`, `git push` or `gh`;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - send a signal to any process;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build both binaries and run the repository Verification and the named suites.
- [ ] Replay Spec 0364's PRDs and the related-ADR fixture with both binaries.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the outside sources, the unreachable row, the glossary check, the self-check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Spec 0364's planned PRD is reported once at line 165, its fielded rewrite is not, and the Delivery Base reports neither.
- [ ] An ADR marked `not applicable` draws in no citing ADR, while an applied one and a Spec without a TechSpec keep today's report.
- [ ] Each sabotage makes its named test fail, this Spec passes its own strict check, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0367-the-strict-check-reads-unreachable-acceptance-and-related-adrs-as-written/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-2; Testing Approach; Build Order 4
