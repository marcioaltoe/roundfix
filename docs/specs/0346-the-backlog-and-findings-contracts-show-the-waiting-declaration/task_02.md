---
task: task_02
spec: 0346-the-backlog-and-findings-contracts-show-the-waiting-declaration
status: pending
type: qa
complexity: medium
---

# Task 02: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec adds the two commented
Waiting Declaration keys and one sentence to the Baseline's copyable Backlog
and Findings Operational Contracts, records the `context-workflow` module's
next version, regenerates the derived Baseline files and this repository's
managed docs layout guide and Setup Manifest, and keeps the user guide's
Findings template, the glossary, the Skill Coverage Map and the Behavior
Surface Record in step.

Every row runs on this repository's own files, its embedded catalog and
temporary copies of the audited head or the Delivery Base. No row starts an
Agent session or an ACP adapter, reads or writes the real `~/.roundfix`,
reads another repository, reaches the network, sends a request to a
provider, the judge or a reviewer, or reads a credential. Every `go test` it
runs carries an explicit `-timeout`. Every row that runs `bin/roundfix` sets
`HOME` and `TMPDIR` to fresh temporary directories, with `GOMODCACHE` and
`GOCACHE` kept at their `go env` values, because the CLI loads
`$HOME/.roundfix/config.yml`, and removes them afterwards. The QA never
writes into the checkout while a test suite runs. No row compares one
command's output with another's byte for byte. The gate keeps its evidence
under names the repository does not ignore, `.txt` for command output, and
links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 4 by running
   `go test -count=1 -timeout 20m -v` with
   `TestTheOperationalContractsShowTheWaitingDeclarationCommented`,
   `TestAnAdopterKeepsTheWaitingContractClauses`,
   `TestTheRetirementClausesCarryTheirText`,
   `TestLifecycleClausesCarryTheScopedWording`,
   `TestFindingsOperationalContract` and
   `TestBaselineClauseForceIsCharacterized` over `./internal/baseline`, then
   the whole `./internal/baseline` package and
   `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST verify Success Metric 5 through the built binary on two temporary
   clones made with `git clone --shared` of this checkout into fresh
   temporary directories, removed afterwards:
   - a clone checked out at the Delivery Base (the merge base of the audited
     head with `main`): `bin/roundfix baseline update --repo <clone> --no-skills --no-history --format text`
     records the lines `File changes: 2`,
     `- update docs/agents/docs-layout.md (2 managed entries)`,
     `- update docs/agents/setup-context.json (1 managed entries)` and
     `Retention evidence: 0`;
   - a clone checked out at the audited head: the same command records
     `File changes: 0`.
   Only those lines are asserted; digests, paths and warnings vary.
4. MUST verify the byte invariant of Goal 3 against the Delivery Base: in
   `internal/baseline/assets/modules/context-workflow.json`, every rule and
   every clause other than `clause.context.backlog-01-operational-contract`
   and `clause.context.findings-01-frontmatter` is equal, as decoded JSON, to
   its value at the Delivery Base, except the top-level `version` and the
   `version` of `rule.context.docs-layout`, which each rise by exactly one;
   and each of the two clauses, with the two commented lines and its
   sentence paragraph removed, equals its guidance at the Delivery Base.
5. MUST verify Success Metric 2 offline: extract the copyable frontmatter of
   both contracts from this repository's `docs/agents/docs-layout.md`, decode
   each with a YAML 1.2 reader available in the repository's Go module
   (`gopkg.in/yaml.v3`, through a throwaway `go run` program kept under a
   temporary directory outside the checkout, or through the named test of
   requirement 2), and record that neither declares `waiting_on` or `since`.
6. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails on the named case:
   - `# waiting_on:` uncommented to `waiting_on:` in the Backlog clause of
     the module: `TestTheOperationalContractsShowTheWaitingDeclarationCommented`
     fails in subtest `clause.context.backlog-01-operational-contract` with
     "the key must stay commented";
   - the words "`pending` or `partial`" removed from the Findings sentence of
     the module: the same test fails in subtest
     `clause.context.findings-01-frontmatter` with "lacks its waiting
     sentence";
   - `docs/agents/docs-layout.md` replaced by its Delivery Base bytes: the
     same test fails in subtest `repository_guide`;
   - the `id` of `clause.context.findings-01-frontmatter` renamed in the
     module: `TestAnAdopterKeepsTheWaitingContractClauses` fails;
   - the two commented lines removed from the user guide's Findings
     template: `TestGuidanceCompositionDocumentation` fails in subtest
     `Findings_template`.
7. MUST execute the tests task_01 names in its Verification.
8. MUST record, as an outside-evidence row, a source this Spec did not
   write, with its date and not blocked: the YAML 1.2.2 passage `_prd.md` →
   Acceptance evidence quotes ("Comments are a presentation detail and must
   not have any effect on the serialization tree or representation graph."),
   compared offline with the result of requirement 5.
9. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare `changes: **Waiting Declaration**`, that
    `CONTEXT.md` defines it in the `**Term**:` form and names the Backlog
    and Findings Operational Contracts, and that the Task commit range
    changes no other glossary entry.
11. MUST verify from the repository history:
    - task_01's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`, anything under
      `internal/baseline/assets/source-baselines/` or
      `internal/baseline/assets/retention/`, any other module under
      `internal/baseline/assets/modules/`, any guide under `docs/agents/`
      other than `docs-layout.md` and `setup-context.json`, any ADR, any
      skill under `.agents/skills/` or `skills/`, or any Go file outside
      `internal/baseline/waiting_declaration_clause_test.go`;
    - in `docs/references/skill-coverage.json` only the `review` of
      `guide docs/user-guide/context-driven-development.md` changed.

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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0346-the-backlog-and-findings-contracts-show-the-waiting-declaration --plan --offline; rm -rf "$h" "$t"`
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
- [ ] Run the named tests and the two update previews through the binary.
- [ ] Check the byte invariant and the YAML effect against the Delivery Base.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the outside source, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Both contracts show the two commented keys and their sentence in every rendered copy, and a copied frontmatter declares no wait.
- [ ] Every other clause keeps its bytes, both clauses are retained, and the update preview lists only the docs layout guide and the Setup Manifest.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0346-the-backlog-and-findings-contracts-show-the-waiting-declaration/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary; Skills
- `_techspec.md` → API Contracts 1-3; Testing Approach; Build Order 2
