---
task: task_02
spec: 0378-the-catalog-skill-bump-refreshes-with-a-binary-built-after-the-digests
status: pending
type: qa
complexity: low
---

# Task 02: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec makes the Catalog
Skill Bump build its binary again after its `make baseline-digests` step, so
the Managed Refresh it runs records the catalog digest of the tree on disk,
and proves that order in the Catalog Skill Bump Makefile contract.

Every row runs in this repository or in temporary directories it creates
and removes. No row creates, removes or changes any file in `/tmp`,
`/private/tmp` or `$HOME` outside a directory made by `mktemp -d`, runs a
real Catalog Skill Bump, reaches the network, reads another repository,
starts an Agent session or an ACP adapter, reads or writes the real
`~/.roundfix`, sends a request to a provider, the judge or a reviewer,
records an environment value, or sends a signal to any process. Every
`go test` carries an explicit `-timeout` and keeps the session's own `HOME`
and `TMPDIR`. Every row that runs a Roundfix binary sets `HOME` to a fresh
temporary directory, with `GOMODCACHE` and `GOCACHE` kept at their `go env`
values, and removes it afterwards. The QA never writes into the checkout
while a test suite runs. The gate keeps its evidence under names the
repository does not ignore, `.txt` for command output, and links only files
it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST run, with `go test -count=1 -timeout 10m -v ./internal/verifyselect -run '^TestTheSkillsBump'`,
   every Catalog Skill Bump test, recording the pass line of
   `TestTheSkillsBumpRefreshesWithABinaryBuiltAfterTheDigests` and
   `TestTheSkillsBumpPinsTheCatalogMainCommit`; then
   `go test -count=1 -timeout 10m -tags docscontract ./internal/docscontract -run '^TestTheSetupManifestIsCurrent$'`.
3. MUST record the starting behavior in a `git archive` of the Delivery Base
   (the merge base of the audited head with `main`) unpacked in a temporary
   directory and committed there as a fresh Git repository: copy the audited
   head's `internal/verifyselect/skills_bump_makefile_test.go` over, run the
   two named tests with `-timeout`, and record that
   `TestTheSkillsBumpRefreshesWithABinaryBuiltAfterTheDigests` fails with
   `the managed refresh ran a binary that embeds "before digests\n"` and
   `TestTheSkillsBumpPinsTheCatalogMainCommit` fails on its call list, while
   both pass on the audited head (row 2).
4. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive` and committed there as a fresh Git repository, with
   `-timeout` on every `go test`, and record that both
   `TestTheSkillsBumpRefreshesWithABinaryBuiltAfterTheDigests` and
   `TestTheSkillsBumpPinsTheCatalogMainCommit` fail:
   - the line `$(GO) build -o "$$tmp/roundfix" ./cmd/roundfix; \` that
     follows the `baseline-digests` line of the `skills-third-party-bump`
     recipe is removed;
   - the same line is moved to just before the `baseline-digests` line.
5. MUST record Success Metric 2 with the real binary in another `git
   archive` copy of the audited head, committed there as a fresh Git
   repository:
   - in `internal/baseline/assets/modules/backend.json`, change the clause
     text `Keep application use cases independent of HTTP request, response, router, and middleware types.`
     to `Keep application use cases independent of HTTP request, response, router, middleware, and session types.`,
     then run
     `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1 -timeout 10m`;
   - build `go build -buildvcs=false -o <temporary>/before ./cmd/roundfix`,
     run `make baseline-digests`, and record that it rewrites the
     `standard-typescript-monorepo` formatter golden `backend.md` and the
     profile's `goldenDigest`;
   - run `<temporary>/before baseline update --repo . --no-skills --no-history --yes --format text`
     twice, record that the second reports `File changes: 0`, and record
     that `TestTheSetupManifestIsCurrent` then fails with
     `Setup Manifest is stale`;
   - build `<temporary>/after` the same way, run the same refresh with it,
     and record that `TestTheSetupManifestIsCurrent` passes.
6. MUST execute the tests every non-QA Task names in its Verification.
7. MUST record, as the outside-evidence row, sources this Spec did not
   write, with their dates and not blocked: the adopted Backlog Entry
   (`references/`, recorded by the operator from Spec 0370's TechSpec step
   5), whose cause ("a binary it built before `make baseline-digests`
   rewrote the profile's golden digest") matches row 5; and the Go `embed`
   documentation quoted in `_prd.md` → Acceptance evidence
   (https://pkg.go.dev/embed: a `//go:embed` variable is initialized "with
   the contents of files read from the package directory or subdirectories
   at compile time."), compared with row 5's `before` binary. The row reads
   the files in this repository, never the network.
8. MUST perform the glossary check (`docs/agents/domain.md`): record that
   the PRD and TechSpec declare `## Glossary` as `None.` and that the Task
   commit range does not change `CONTEXT.md`.
9. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is listed in `_authorization.md` `paths`;
   - the `Makefile` change is the one recipe line of Requirement 1 of
     task_01, and the `docs/agents/setup-context.json` change touches only
     `declarationDigest` values;
   - the Task commit range changes none of these: any workflow, `go.mod`,
     `go.sum`, `.roundfixrc.yml`, `skills-lock.json`, `CHANGELOG.md`,
     `internal/app/version.go`, `dist/npm/**`, any file under
     `.agents/skills/` or `skills/`, `CONTEXT.md`, any ADR, any Baseline
     module or derived Baseline asset, or the `### QA settlement` section of
     any skill.

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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0378-the-catalog-skill-bump-refreshes-with-a-binary-built-after-the-digests --plan --offline; rm -rf "$h" "$t"`
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
- [ ] Characterize the Delivery Base, run each sabotage and the real-binary row on temporary copies.
- [ ] Record the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The refresh test passes on the audited head and fails on the Delivery Base and under each sabotage.
- [ ] The real binary built before `make baseline-digests` leaves the Setup Manifest stale and the one built after leaves it current.
- [ ] The scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0378-the-catalog-skill-bump-refreshes-with-a-binary-built-after-the-digests/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-2; Core Features 1-2; Success Metric 1; Success Metric 2; Acceptance evidence; Glossary
- `_techspec.md` → API Contract 1; Testing Approach; Build Order 2
