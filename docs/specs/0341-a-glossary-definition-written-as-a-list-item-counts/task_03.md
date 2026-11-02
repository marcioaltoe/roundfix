---
task: task_03
spec: 0341-a-glossary-definition-written-as-a-list-item-counts
status: pending
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the Spec Consistency
Check read a Glossary Definition written as `**Term**:`, `**Term** —`,
`- **Term**:` or `- **Term** —` (list marker `-`, `*` or `+`), names those
forms in `SC-GLOSSARY-MISSING`, prints the line-anchored glossary
Verification in `SC-GLOSSARY-UNPLANNED`, and teaches the same forms and
command in the `write-tasks`, `write-prd` and Roundfix skills.

Every row runs on this repository's own files, temporary repositories and
temporary copies of the audited head. No row starts an Agent session or an
ACP adapter, reads or writes the real `~/.roundfix`, reads another
repository, reaches the network, sends a request to a provider, the judge or
a reviewer, or reads a credential. Every `go test` it runs carries an
explicit `-timeout`. Every row that runs `bin/roundfix` sets `HOME` and
`TMPDIR` to fresh temporary directories, with `GOMODCACHE` and `GOCACHE` kept
at their `go env` values, because the CLI loads `$HOME/.roundfix/config.yml`,
and removes them afterwards. The QA never writes into the checkout while a
test suite runs. No row compares one command's output with another's byte
for byte. The gate keeps its evidence under names the repository does not
ignore, `.txt` for command output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 5 and Invariants 1 to 4 by running
   `go test -count=1 -timeout 20m -v` with
   `TestGlossaryDefinitionFormsAreRead`,
   `TestGlossaryVerificationCommandAgreesWithTheDefinitionReader` and
   `TestGlossaryFindingsNameTheDefinitionForms` over `./internal/speccheck`
   and `TestTheGlossaryGuidesStateTheDefinitionForms` over `./skills`, then
   the whole `./internal/speccheck` and `./skills` packages and
   `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST reproduce the Pantheon case through the built binary in a temporary
   Git repository the row creates and removes: commit
   `.agents/skills/write-prd/references/glossary.md`, then a Spec
   `docs/specs/9001-fixture/` whose PRD declares `- adds: **Acervo Term**`,
   whose `_tasks.md` declares `qa: declined` and one completed Task declaring
   `interface: CONTEXT.md` and naming `**Acervo Term**` in its Verification.
   For each glossary line `- **Acervo Term**: d`, `- **Acervo Term** — d`,
   `**Acervo Term** — d` and `+ **Acervo Term**: d`, run
   `bin/roundfix spec check 9001-fixture` (with the temporary `HOME` and
   `TMPDIR`) and record no `SC-GLOSSARY-MISSING`; for `**Acervo Term**: d`
   record none as well; for `### Acervo Term`, `  - **Acervo Term**: d` and
   `- **Acervo Term** - d` record one `SC-GLOSSARY-MISSING` whose fix names
   `**Term** —`. Then remove the binding Task's Verification term and record
   that `SC-GLOSSARY-UNPLANNED` prints the command
   `grep -F -- '**Acervo Term**' < CONTEXT.md | grep -qE ...`; run that printed
   command in the fixture against each line above and record exit 0 exactly
   for the lines the check accepted. Remove the temporary directories and the
   repository afterwards.
4. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails on the named case:
   - list marker dropped from `glossaryDefinitionPattern`:
     `TestGlossaryDefinitionFormsAreRead` fails on `dash list colon`;
   - leading whitespace allowed in the reader:
     `TestGlossaryDefinitionFormsAreRead` fails on `indented list item`;
   - em dash dropped from the reader: `TestGlossaryDefinitionFormsAreRead`
     fails on `em dash`;
   - the command's pattern also accepts ` - `:
     `TestGlossaryVerificationCommandAgreesWithTheDefinitionReader` fails on
     `ascii hyphen`;
   - `regexp.QuoteMeta` dropped from the command:
     `TestGlossaryVerificationCommandAgreesWithTheDefinitionReader` fails on
     `metacharacter term`;
   - the command reads the glossary as a file operand instead of standard
     input: `TestGlossaryFindingsNameTheDefinitionForms` fails;
   - the `SC-GLOSSARY-MISSING` fix reverted to "Define the declared term in
     the glossary.": `TestGlossaryFindingsNameTheDefinitionForms` fails;
   - the fenced command removed from the canonical and mirror `write-tasks`
     skill: `TestTheGlossaryGuidesStateTheDefinitionForms` fails.
5. MUST execute the tests task_01 and task_02 name in their Verification.
6. MUST search `docs/user-guide/commands/spec.md`, `CONTEXT.md` and every copy
   of the `write-tasks`, `write-prd` and Roundfix skills and record that the
   spec command guide and both copies of the `write-prd` glossary guide and
   the Roundfix `spec` reference name the accepted forms, that both copies of
   `write-tasks` carry the rendered command and no longer the sentence "checks
   the bolded term followed by a colon in the glossary on
   whitespace-normalized text", and that `CONTEXT.md` defines
   **Glossary Definition**.
7. MUST record, as outside-evidence rows, two sources this Spec did not
   write, each with its date and none blocked:
   - the Pantheon report quoted in
     `references/2026-11-02-a-glossary-definition-written-as-a-list-item-is-not-read.md`
     (2026-11-02): compare offline its four measured forms with the results
     of requirement 3;
   - the CommonMark passage `_prd.md` → Acceptance evidence quotes ("A bullet
     list marker is a `-`, `+`, or `*` character."): compare it with the
     `dash`, `star` and `plus` rows of `TestGlossaryDefinitionFormsAreRead`.
8. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that the
   PRD and TechSpec declare `adds: **Glossary Definition**`, that `CONTEXT.md`
   defines it in the `**Term**:` form, and that the Task commit range changes
   no other glossary entry.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`, `docs/references/skill-coverage.json`,
      `internal/docscontract/user_guide_contract_test.go`, anything under
      `internal/baseline/assets/` or `docs/agents/`, the `domain-modeling`
      skill, any skill other than the Roundfix, `write-prd` and `write-tasks`
      files the Tasks declare;
    - no `### QA settlement` section changed.

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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0341-a-glossary-definition-written-as-a-list-item-counts --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the Pantheon case through the binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every accepted form defines the term and every rejected form does not, in the tests and through the binary.
- [ ] The printed glossary Verification passes exactly on the lines the check accepts.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0341-a-glossary-definition-written-as-a-list-item-counts/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-4; Invariants 1-4; Declared breaks; Testing Approach; Build Order 3
