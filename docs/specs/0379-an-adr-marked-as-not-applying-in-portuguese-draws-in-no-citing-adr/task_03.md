---
task: task_03
spec: 0379-an-adr-marked-as-not-applying-in-portuguese-draws-in-no-citing-adr
status: pending
type: qa
complexity: low
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec makes the
related-ADR seed rule of `roundfix spec check` read a closed list of
not-applicable tokens, `not applicable`, `não se aplica` (precomposed or
decomposed `ã`) and `nao se aplica`, without regard to case; names both
tokens in the `SC-ADR-RELATED` fix text; and aligns ADR-0377, the glossary,
the authoring skills, the `spec` reference and the guides.

Every row runs in this repository or in temporary directories it creates
and removes. No row creates, removes or changes any file in `/tmp`,
`/private/tmp` or `$HOME` outside a directory made by `mktemp -d`, reaches
the network, reads another repository, starts an Agent session or an ACP
adapter, reads or writes the real `~/.roundfix`, sends a request to a
provider, the judge or a reviewer, records an environment value, or sends a
signal to any process. Every `go test` carries an explicit `-timeout` and
keeps the session's own `HOME` and `TMPDIR`. Every row that runs a Roundfix
binary sets `HOME` and `TMPDIR` to fresh temporary directories, with
`GOMODCACHE` and `GOCACHE` kept at their `go env` values, and removes them
afterwards. The QA never writes into the checkout while a test suite runs.
The gate keeps its evidence under names the repository does not ignore,
`.txt` for command output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST run, with `go test -count=1 -timeout 10m -v ./internal/speccheck`,
   the tests `TestAPortugueseNotApplicableADRDrawsInNoCitingADR`,
   `TestAMarkerOutsideTheClosedListStillDrawsIn`,
   `TestAnAppliedADRStillDrawsInItsCitingADR`,
   `TestANotApplicableADRDrawsInNoCitingADR`,
   `TestWithoutATechSpecEveryListedADRSeeds` and
   `TestAnADRCitingAnAppliedADRIsStillRelated`, recording the pass line of
   each subtest; then the whole `./internal/speccheck` package.
3. MUST record the starting behavior in a `git archive` of the Delivery Base
   (the merge base of the audited head with `main`) unpacked in a temporary
   directory and committed there as a fresh Git repository: copy the audited
   head's `internal/speccheck/adr_related_marker_test.go` over and run both
   of its tests with `-timeout`; record that every subtest of
   `TestAPortugueseNotApplicableADRDrawsInNoCitingADR` fails with an
   `SC-ADR-RELATED` finding and that
   `TestAMarkerOutsideTheClosedListStillDrawsIn` passes there, which shows
   the guards hold the starting behavior.
4. MUST record Success Metrics 1 and 2 through the built product: build a
   Delivery Base binary with `go build -buildvcs=false` in the row 3 copy;
   then, for each relation of the second fixture ADR, `não se aplica; a decisão segue igual.`,
   `NÃO SE APLICA: a decisão segue igual.`,
   `nao se aplica; a decisão segue igual.`,
   `não se aplicam; a decisão segue igual.` and
   `not applicable; unchanged.`, create a temporary Git repository with
   `docs/agents/agent-instructions.md`, `docs/agents/cli.md` and
   `docs/agents/domain.md`, three accepted fixture ADRs numbered `0001`, `0002` and `0003`, where
   the second cites the first and the third cites the second, committed
   first, then a Spec `0001-adr-related-seeds` whose PRD lists the first two
   under Active ADR obligations with the four Project Constraints rows and
   whose TechSpec's `## ADR Alternatives` applies the first and gives the
   second the relation, committed; run `spec check 0001-adr-related-seeds --strict` with
   both `bin/roundfix` and the Delivery Base binary, and record the number of
   `SC-ADR-RELATED` lines: the Delivery Base reports one for every relation
   except the English one; the audited head reports none for the three
   accepted Portuguese forms and the English one, and one for the plural,
   whose `fix:` line names `"not applicable" or "não se aplica"`.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive` and committed there as a fresh Git repository, with
   `-timeout` on every `go test`, and record that the named subtest fails:
   - the pattern without the Portuguese alternative,
     `(?i)^not applicable\b`: every subtest of
     `TestAPortugueseNotApplicableADRDrawsInNoCitingADR`;
   - the pattern without `(?i)`:
     `TestAPortugueseNotApplicableADRDrawsInNoCitingADR/capitalized` and
     `/upper_case`;
   - the pattern without `|a\x{0303}`:
     `TestAPortugueseNotApplicableADRDrawsInNoCitingADR/decomposed`;
   - the pattern without the unaccented `|a`:
     `TestAPortugueseNotApplicableADRDrawsInNoCitingADR/unaccented`;
   - the pattern without the closing `\b`:
     `TestAMarkerOutsideTheClosedListStillDrawsIn/plural`;
   - the fix text without ` or \"não se aplica\"`:
     `TestAnAppliedADRStillDrawsInItsCitingADR`.
6. MUST execute the tests every non-QA Task names in its Verification.
7. MUST record, as the outside-evidence row, sources this Spec did not
   write, with their dates and not blocked: the adopted Backlog Entry
   (`references/`, recorded by the operator from Fiscus's message of
   2026-11-10 about a repository this Spec did not build: "the 76 ADR lines
   behind the cascade Spec 0367 fixed use \"não se aplica\""), compared with
   row 4's accented form; the Unicode Standard Annex #15 passage quoted in
   `_prd.md` → Acceptance evidence (canonical equivalence of a combining
   sequence such as `Ç ↔ C+◌̧`), compared with the `decomposed` subtest; and
   the Go `regexp/syntax` passage quoted there (`\b` "at ASCII word
   boundary"), compared with the `plural` subtest. The row reads the files in
   this repository, never the network.
8. MUST repeat the wording search of task_01: every line under
   `docs/user-guide/`, `.agents/skills/write-techspec/`,
   `.agents/skills/write-tasks/`, `.agents/skills/roundfix/references/spec.md`
   and `CONTEXT.md` that tells an author to begin an ADR Alternatives line
   with `not applicable`, or says such a line keeps the ADR listed, also
   names `não se aplica` in the same sentence; record each match.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that
   the PRD declares `changes: **Spec Consistency Check**`, that `CONTEXT.md`
   keeps its definition line and its sentence names the Portuguese token and
   ADR-0389, and that no other term changed in the Task commit range.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every Governed Path changed is listed in `_authorization.md` `paths`;
    - ADR-0377 carries the `**Amended by [ADR-0389]` notice and its
      `updated_at` moved, and ADR-0093, ADR-0168 and ADR-0389 are unchanged
      by the Task commits;
    - the Task commit range changes none of these: the Makefile, any
      workflow, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/**`,
      `internal/speccheck/constraints.go`, any Baseline module or derived
      Baseline file, `docs/references/skill-coverage.json`, or the
      `### QA settlement` section of any skill.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0379-an-adr-marked-as-not-applying-in-portuguese-draws-in-no-citing-adr --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read another repository, or run `git fetch`, `git push` or `gh`;
    - read or write `~/.roundfix`;
    - create a file in `/tmp`, `/private/tmp` or `$HOME` outside a
      `mktemp -d` directory it removes;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - send a signal to any process;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification and the named suites.
- [ ] Characterize the Delivery Base, run the product rows with both binaries and each sabotage on temporary copies.
- [ ] Record the outside sources, the wording search, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every accepted Portuguese form leaves the citing ADR undrawn through `bin/roundfix spec check --strict`, while the Delivery Base binary reports it.
- [ ] Forms outside the list still report the gap, and each sabotage fails its named subtest.
- [ ] The decision record, glossary, skills and guides name the token, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0379-an-adr-marked-as-not-applying-in-portuguese-draws-in-no-citing-adr/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Acceptance evidence; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; Testing Approach; Build Order 3
