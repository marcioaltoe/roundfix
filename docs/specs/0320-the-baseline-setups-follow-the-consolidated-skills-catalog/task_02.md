---
task: task_02
spec: 0320-the-baseline-setups-follow-the-consolidated-skills-catalog
status: completed
type: qa
complexity: high
---

# Task 02: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec moves the five bundled
Setup Snapshots from upstream commit `65d5f365` to the published
`8d71fade68c54693bedbd917371f0514311af157` through the asset sync,
regenerates the derived Baseline files and this repository's Setup Manifest,
and makes the tests that pinned the upstream revision read one constant, with
a new test that pins the republished `onioncry` tree.

Every row runs on this repository's own files, temporary copies of the
audited head and temporary repositories. No row reads or writes the real
`~/.roundfix`, reads `~/dev/skills` or another repository, reaches the
network, sends a request to a provider, the judge or a reviewer, or reads a
credential. Every `go test` it runs carries an explicit `-timeout`, and every
row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary
directories, because the CLI loads `$HOME/.roundfix/config.yml`. The QA never
writes into the checkout while a test suite runs.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 and Invariants 1, 3 and 5 from the audited
   `internal/baseline/assets/setups/*.json` and the same five files at the
   Delivery Base (the merge base of the audited head with `main`): record
   that each audited file with a `github` source names
   `8d71fade68c54693bedbd917371f0514311af157` as its `source.ref`, that the
   composed `go-cli-typescript-bun.json` source (`composed`, `setups: [go,
   typescript]`, no ref) equals the Delivery Base's, that all five name that
   commit as the `ref` of every `github` skill, that its
   ordered skill names equal the Delivery Base's, and that the only
   `treeDigest` that differs is `onioncry`'s, in `typescript.json` and
   `go-cli-typescript-bun.json`, now
   `be502b02544be90059c126e9b07bfeb6a4af115aada86455845b1407a5eab8b3`.
3. MUST verify Success Metric 2: run
   `go test -timeout 20m -count=1 ./internal/baseline ./skills` and
   `go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract`,
   and record the passes of `TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree`
   (both subtests), `TestEveryUpstreamSetupFollowsTheAiMemoryRevision`,
   `TestEveryUpstreamSetupFollowsThePythonPresetRevision`,
   `TestThePythonSetupIsThePublishedUpstreamPreset`,
   `TestCatalogCompatibility`,
   `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` and
   `TestTheSetupManifestIsCurrent`.
4. MUST verify Success Metric 3 and Invariant 2 after the suites finish: in a
   temporary `git clone --no-local` of this checkout detached at the audited
   head, run `make baseline-digests` and record `"changed":false`; then, with
   the binary built from the audited head, run
   `bin/roundfix baseline update --repo <clone> --no-skills --no-history --yes --format text`
   and record `File changes: 0`.
5. MUST verify Success Metric 4 and Invariant 4: record that no
   `internal/**/*_test.go` file of the audited head names
   `65d5f365cace27c4b468db1723213723b10c870b`, and that
   `upstreamSetupRevision` is declared once, in
   `internal/baseline/upstream_setup_revision_test.go`, and read by the
   three revision tests.
6. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails, with the case whose input makes the mutation
   observable:
   - restore the five `internal/baseline/assets/setups/*.json` files from
     the Delivery Base:
     `TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree/typescript` and
     `/go-cli-typescript-bun` fail naming the `onioncry` ref `65d5f365…` and
     tree `f918508d…`, and `TestEveryUpstreamSetupFollowsTheAiMemoryRevision`
     fails;
   - set `upstreamSetupRevision` to `65d5f365cace27c4b468db1723213723b10c870b`:
     `TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree`,
     `TestEveryUpstreamSetupFollowsTheAiMemoryRevision`,
     `TestEveryUpstreamSetupFollowsThePythonPresetRevision` and
     `TestThePythonSetupIsThePublishedUpstreamPreset` fail, because the
     audited snapshots name `8d71fade`;
   - restore `docs/agents/setup-context.json` from the Delivery Base:
     `TestTheSetupManifestIsCurrent` (with `-tags docscontract`) fails with
     "Setup Manifest is stale";
   - restore `internal/baseline/testdata/catalog.digest` from the Delivery
     Base: `TestCatalogCompatibility` fails.
7. MUST execute the tests task_01 names in its Verification.
8. MUST record, as the outside-evidence row, the statements of
   `default/marcioaltoe-skills` that `_prd.md` → Acceptance evidence quotes
   from the adopted Backlog Entry under `references/`, read in the audited
   checkout. The row is checked offline, never blocked: record that no
   audited snapshot lists `skill-architect`, `writing-for-agents` or
   `skill-creator` ("no setup preset named either"), and that `onioncry` is
   the one skill whose tree changed ("onioncry updated"). Name the two
   message ids and their dates as the source.
9. MUST record Success Metric 5 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`.
10. MUST perform the glossary and documentation checks: record that the Spec
    declares no glossary term and that `CONTEXT.md` is byte-identical to the
    Delivery Base (`docs/agents/domain.md`), and that no operational file
    names `skill-architect` or `writing-for-agents`: the scan excludes
    `docs/history/` and this Spec's own folder, whose PRD, Tasks and adopted
    reference quote the upstream messages on purpose.
11. MUST verify from the repository history:
    - the Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every production
      `.go` file, every file under `internal/baseline/assets/` other than the
      five setups, every file under `.agents/skills/` and `skills/`,
      `skills-lock.json`, every file under `docs/agents/` other than
      `setup-context.json`, `docs/user-guide/` and `docs/adr/`;
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
    `HOME="$(mktemp -d)" TMPDIR="$(mktemp -d)" bin/roundfix archive 0320-the-baseline-setups-follow-the-consolidated-skills-catalog --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test, the regeneration rows and each sabotage on temporary copies.
- [ ] Record the outside source, the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every upstream snapshot names `8d71fade`, keeps its membership, and changes only the `onioncry` tree.
- [ ] The derived files and the Setup Manifest regenerate with no change, and each sabotage makes its named test fail.
- [ ] The scope audit finds no change outside the declared files.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0320-the-baseline-setups-follow-the-consolidated-skills-catalog/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-3; User Experience; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-4; Invariants 1-5; Testing Approach; Build Order 2
