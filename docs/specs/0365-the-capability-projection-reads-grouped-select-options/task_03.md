---
task: task_03
spec: 0365-the-capability-projection-reads-grouped-select-options
status: pending
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec makes the
capability projection read a grouped ACP select option as its groups'
values in advertised order, under the bounds of a flat select, pins the two
codex-acp payloads recorded by the operator as fixtures, and proves the
Exact Agent Selection Proof over both through the fake `acpx` harness.

Every row runs on fixtures, the fake `acpx` harness and temporary
directories. No row starts an Agent session or an ACP adapter, runs `npx`,
`acpx` or `codex`, reads or writes the real `~/.roundfix` or `~/.codex`,
reaches the network, sends a request to a provider, the judge or a reviewer,
reads a credential, or sends a signal to any process. Every `go test`
carries an explicit `-timeout` and keeps the session's own `HOME` and
`TMPDIR`. Every row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to
fresh temporary directories, with `GOMODCACHE` and `GOCACHE` kept at their
`go env` values, and removes them afterwards. The QA never writes into the
checkout while a test suite runs. The gate keeps its evidence under names
the repository does not ignore, `.txt` for command output, and links only
files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST run, with `go test -count=1 -timeout 20m -v`, the tests tasks 01
   and 02 name over `./internal/agent`, then the whole `./internal/agent`
   package, recording each pass, and record that
   `TestSelectionCapabilitiesProjectTheRecordedCodexACPPayloads/2.6.0`
   reports the `mode` values `read-only`, `workspace-write`, `agent`,
   `agent-full-access`, `permission-profile:localdev` in that order.
3. MUST record, with `cmp`, that both files under `internal/agent/testdata/`
   equal the recorded payloads in this Spec's `references/`, and record
   their SHA-256 digests.
4. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails on the named case:
   - `flattenSelectValues` bypassed, so `projectSelectCapability` reads
     `*raw.Options` as values again:
     `TestSelectionCapabilitiesProjectTheRecordedCodexACPPayloads/2.6.0`
     and `TestProofOverTheRecordedCodexACPPayloads/2.6.0/full-access` fail
     with `invalid_select_option`;
   - the duplicate check reset for each group:
     `TestGroupedSelectOptionsKeepTheEvidenceBounds/value_repeated_across_groups`
     fails;
   - a list mixing values and groups accepted:
     `TestGroupedSelectOptionsKeepTheEvidenceBounds/values_and_groups_mixed`
     fails;
   - the ceiling checked per group instead of on the flattened count:
     `TestGroupedSelectOptionsKeepTheEvidenceBounds/flattened_count_above_the_ceiling`
     fails;
   - an empty group refused:
     `TestGroupedSelectOptionsKeepTheEvidenceBounds/empty_group_beside_a_values_group`
     fails;
   - the groups appended in reverse order:
     `TestSelectionCapabilitiesProjectTheRecordedCodexACPPayloads/2.6.0`
     fails on the `Values` order.
5. MUST execute the tests every non-QA Task names in its Verification.
6. MUST record, as an outside-evidence row, sources this Spec did not
   write, with their date and not blocked: the two payloads recorded by the
   operator on 2026-11-09 (row 3), and the Agent Client Protocol passage
   `_prd.md` → Acceptance evidence quotes ("This is either a flat array of
   values or an array of groups. Values and groups are not mixed in the same
   array."), compared offline with the `values and groups mixed` case and
   the 2.6.0 projection.
7. MUST record the PRD's one Unreachable Acceptance entry, the live
   `roundfix profiles validate` with codex-acp 2.6.0 after the release, as
   exactly one row `blocked (environment: released binary and real codex-acp
   2.6.0 adapter required)` with `declared: ` followed by the entry's
   criterion text, counted once in `rows_blocked_declared`.
8. MUST perform the glossary check (`docs/agents/domain.md`): record that
   the PRD and TechSpec declare `None.` under `## Glossary`, that the Task
   commit range does not change `CONTEXT.md`, and that
   `go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract -run '^TestGlossaryPhrasesAreDefinedTerms$'`
   passes.
9. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - `_authorization.md` declares `paths: []` and the Task commit range
     changes no Governed Path;
   - the Task commit range changes none of these: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
     `CHANGELOG.md`, `internal/app/version.go`,
     `dist/npm/roundfix/package.json`, `CONTEXT.md`, any file under
     `docs/adr/`, `docs/agents/`, `docs/user-guide/`, `.agents/skills/`,
     `skills/`, `internal/cli/`, `internal/spec/` or `internal/speccheck/`;
   - task_02 changed no file outside `internal/agent/selection_proof_grouped_mode_test.go`.

   This row reads Task commits, so it declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
11. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0365-the-capability-projection-reads-grouped-select-options --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - start `acpx`, `npx`, `codex` or any ACP adapter;
    - read or write `~/.roundfix` or `~/.codex`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests, the package and the fixture comparison.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the outside sources, the declared unreachable row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The recorded codex-acp 2.6.0 payload projects and proves, with `set-mode agent-full-access` under full access, and every flat payload projects as before.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.
- [ ] The one Unreachable Acceptance entry is one declared row.

## Verification

- `newest="$(find 'docs/specs/0365-the-capability-projection-reads-grouped-select-options/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-4; Testing Approach; Build Order 3
