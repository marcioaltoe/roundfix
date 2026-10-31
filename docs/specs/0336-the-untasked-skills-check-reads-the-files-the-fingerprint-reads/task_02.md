---
task: task_02
spec: 0336-the-untasked-skills-check-reads-the-files-the-fingerprint-reads
status: completed
type: qa
complexity: medium
---

# Task 02: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the untasked-skills
rule of `roundfix spec check` and the Behavior Surface fingerprint read one
file set through `skillcoverage.FingerprintReads`, so a Go test file never
asks for a skills Task while a source file still does, and writes the rule
into the spec command guide and the Roundfix skill's spec reference.

Every row runs on this repository's own files, temporary repositories and
temporary copies of the audited head. No row starts an Agent session or an
ACP adapter, reads or writes the real `~/.roundfix`, reads another
repository, reaches the network, sends a request to a provider, the judge or
a reviewer, or reads a credential. Every `go test` it runs carries an
explicit `-timeout`. Every row that runs `bin/roundfix` sets `HOME` and
`TMPDIR` to fresh temporary directories, because the CLI loads
`$HOME/.roundfix/config.yml`, and removes them afterwards. The QA never
writes into the checkout while a test suite runs. No row compares one
command's output with another's byte for byte.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1, 2 and 4 and Invariants 1 to 3 by running
   `go test -count=1 -timeout 20m -v` with
   `TestAGoTestFileUnderACoveredSourceNeedsNoSkillsTask` over
   `./internal/speccheck`, `TestAGoTestFileBelongsToNoBehaviorSurface` over
   `./internal/skillcoverage`, and, with `-tags docscontract`,
   `TestBehaviorSurfaceBytesComeFromParsedValuesAndGuideContent` and
   `TestTheSkillCoverageMapIsCurrent` over `./internal/docscontract`,
   recording each pass.
3. MUST verify Success Metric 3 through the built binary in a temporary
   clone of the audited head that the row creates and removes: commit two
   fixture Specs `9001-test-only` and `9002-source`, each with a minimal
   `_prd.md` (`status: active`), a `_tasks.md` with `qa: declined` and
   one pending `type: backend` Task whose only Context line is
   `interface:` `internal/config/config_test.go` for the first and
   `internal/config/config.go` for the second, and a Verification of
   `true`. With `h="$(mktemp -d)"; t="$(mktemp -d)"`, run
   `HOME="$h" TMPDIR="$t" bin/roundfix spec check <slug>` for each and
   record that the first output has no `SC-SKILLS-UNTASKED` line and the
   second has one for `config keys`; then `rm -rf "$h" "$t"` and the clone.
4. MUST verify Invariant 4 from the repository history: the Task commit
   range changes in `docs/references/behavior-surfaces.json` only the
   `guide docs/user-guide/commands/spec.md` entry.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails, with the case whose input makes the mutation
   observable:
   - `FingerprintReads` always returns `true`:
     `TestAGoTestFileUnderACoveredSourceNeedsNoSkillsTask` fails on
     `internal/cli/example_test.go`,
     `TestAGoTestFileBelongsToNoBehaviorSurface` on
     `internal/config/config_test.go`, and
     `TestBehaviorSurfaceBytesComeFromParsedValuesAndGuideContent` on
     parsing its fixture `internal/config/config_test.go`;
   - `FingerprintReads` always returns `false`:
     `TestAGoTestFileUnderACoveredSourceNeedsNoSkillsTask` fails on
     `internal/cli/example_testing.go` and
     `TestBehaviorSurfaceBytesComeFromParsedValuesAndGuideContent` on the
     `config keys` fingerprint;
   - `FingerprintReads` returns `!strings.Contains(filePath, "_test")`:
     `TestAGoTestFileUnderACoveredSourceNeedsNoSkillsTask` fails on
     `internal/cli/example_testing.go` and
     `TestAGoTestFileBelongsToNoBehaviorSurface` on
     `internal/config/config_testing.go`;
   - the guard removed from `SurfacesForPath`:
     `TestAGoTestFileUnderACoveredSourceNeedsNoSkillsTask` fails on
     `internal/cli/example_test.go` and
     `TestAGoTestFileBelongsToNoBehaviorSurface` on its test-file cases.
6. MUST execute the tests task_01 names in its Verification.
7. MUST search `docs/user-guide/` and every copy of the Roundfix skill for
   `SC-SKILLS-UNTASKED` and record that the guide and both copies of the
   `spec` reference state that a Go test file is never a Behavior Surface
   source, and that no passage says a test file needs a skills Task.
8. MUST record, as outside-evidence rows, two sources this Spec did not
   write. First, the report of Spec 0332's author quoted in the adopted
   Backlog Entry under `references/` (2026-10-30): compare offline its shape
   (a Task changing only `internal/config` tests, asked for a Coverage
   Review) with the row of requirement 3. Second, the Go command
   documentation passage `_prd.md` → Acceptance evidence quotes
   (`cmd/go`, "When compiling packages, build ignores files that end in
   '_test.go'."): compare it with API Contract 1. Neither row is blocked;
   name each source with its date.
9. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare no term and that `CONTEXT.md` is
    unchanged in the Task commit range.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `docs/references/skill-coverage.json`,
      `internal/speccheck/skills_declaration.go`, anything under
      `internal/baseline/assets/` or `docs/agents/`, any skill other than
      the Roundfix files the Task declares;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
13. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" bin/roundfix archive 0336-the-untasked-skills-check-reads-the-files-the-fingerprint-reads --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's ` Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the repository case through the binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A Go test file under a covered source gives no `SC-SKILLS-UNTASKED`, and a non-test source still gives one, in the tests and through the binary.
- [ ] The fingerprint and the lookup share `FingerprintReads`, and only the spec guide fingerprint moved.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0336-the-untasked-skills-check-reads-the-files-the-fingerprint-reads/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-4; Invariants 1-4; Declared breaks; Testing Approach; Build Order 2
