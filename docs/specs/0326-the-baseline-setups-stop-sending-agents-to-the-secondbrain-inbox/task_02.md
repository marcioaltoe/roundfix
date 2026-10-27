---
task: task_02
spec: 0326-the-baseline-setups-stop-sending-agents-to-the-secondbrain-inbox
status: pending
type: qa
complexity: high
---

# Task 02: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec moves the five bundled
Setup Snapshots to the published upstream commit
`b56fc5ffcad1ab366fbc03839a6ce8c662401178` through the asset sync, restores
this repository's `knowledge-workspace` and its lock entry through the
confirmed restore command, regenerates the derived Baseline files and the
Setup Manifest, and moves the shared revision test constant with a new test
that pins the corrected `knowledge-workspace` tree.

Every row runs on this repository's own files, temporary copies of the
audited head and temporary repositories. No row reads or writes the real
`~/.roundfix`, reads `~/dev/skills`, the operator's scratch clone or another
repository, reaches the network, sends a request to a provider, the judge or
a reviewer, or reads a credential. Every `go test` it runs carries an
explicit `-timeout`, and every row that runs `bin/roundfix` sets `HOME` and
`TMPDIR` to fresh temporary directories, because the CLI loads
`$HOME/.roundfix/config.yml`. The QA never writes into the checkout while a
test suite runs. No row compares one command's output with another's byte
for byte.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 and Invariants 1, 3 and 4 from the audited
   `internal/baseline/assets/setups/*.json` and the same five files at the
   parent of task_01's first commit (the start of the Task commit range):
   record that each audited file with a `github` source names
   `b56fc5ffcad1ab366fbc03839a6ce8c662401178` as its `source.ref`; that the
   composed `go-cli-typescript-bun.json`, which has no `source.ref` by
   design, keeps the `composed` source with `setups: [go, typescript]` it had
   at that parent; that all five name that commit as the `ref` of every
   `github` skill; that each file's ordered skill names equal the parent's;
   and that the only `treeDigest` that differs from the parent is
   `knowledge-workspace`'s, in all five files, now
   `490af37df4faba388e560380df0863b88881716d96e6185858918681182f2a96`.
3. MUST verify Success Metric 2: run
   `go test -timeout 20m -count=1 ./internal/baseline ./skills` and
   `go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract`,
   and record the passes of
   `TestEverySetupCarriesTheKnowledgeWorkspaceThatRetiredTheSecondbrainInbox`
   (all five subtests), `TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree`,
   `TestEveryUpstreamSetupFollowsTheAiMemoryRevision`,
   `TestEveryUpstreamSetupFollowsThePythonPresetRevision`,
   `TestThePythonSetupIsThePublishedUpstreamPreset`,
   `TestCatalogCompatibility`,
   `TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill`,
   `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot`,
   `TestCheckRepositoryMatchesRealRepository` and
   `TestTheSetupManifestIsCurrent`.
4. MUST verify Success Metric 3 and Invariant 2 after the suites finish: in a
   temporary `git clone --no-local` of this checkout detached at the audited
   head, run `make baseline-digests` and record `"changed":false`; then, with
   the binary built from the audited head, run
   `bin/roundfix baseline update --repo <clone> --no-skills --no-history --yes --format text`
   and record `File changes: 0`. In the same clone, run
   `bin/roundfix baseline update --repo <clone> --no-history --format json`
   (a preview without `--yes`) and record that its `skills` object lists
   nothing under `missing`, `drifted` or `differing`.
5. MUST verify Success Metric 4 and Invariant 5:
   - no `internal/**/*_test.go` file of the audited head names
     `8d71fade68c54693bedbd917371f0514311af157`, and `upstreamSetupRevision`
     and `knowledgeWorkspaceTreeDigest` are each declared once, in
     `internal/baseline/upstream_setup_revision_test.go`;
   - the audited `.agents/skills/knowledge-workspace/SKILL.md` reads
     `version: 0.6.2`, and its `skills-lock.json` entry names ref
     `b56fc5ffcad1ab366fbc03839a6ce8c662401178`;
   - a scan of the audited tree for `dev/secondbrain/inbox/` and
     `create-only carve-out` (fixed strings, hidden directories included,
     `.git` excluded) finds nothing outside `docs/history/` and this Spec's
     own folder, whose PRD, Tasks and adopted reference quote the old text on
     purpose.
6. MUST record that the legacy Source Baseline renders into no current file:
   its directory
   `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/`
   is unchanged by the Task commit range; no `managedArtifacts` path of the
   audited `docs/agents/setup-context.json` lies under
   `internal/baseline/assets/source-baselines/`; and neither the audited
   `docs/agents/` nor the rendered golden
   `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/`
   contains `MAY create files under` or `only writable Secondbrain namespace`,
   while both `secondbrain.md` files contain `inbox/. included` as a regular
   expression match of the prohibition.
7. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails, with the case whose input makes the mutation
   observable:
   - restore the five `internal/baseline/assets/setups/*.json` files from
     the parent of task_01's first commit: all five subtests of
     `TestEverySetupCarriesTheKnowledgeWorkspaceThatRetiredTheSecondbrainInbox`
     fail naming the old ref and tree, and
     `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` fails;
   - restore only `.agents/skills/knowledge-workspace/SKILL.md` from that
     parent: `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot`
     fails naming `knowledge-workspace` and the restore command;
   - restore only `skills-lock.json` from that parent:
     `TestCheckRepositoryMatchesRealRepository` (`./skills`) fails with
     `knowledge-workspace` under `OutdatedExternal`;
   - set `knowledgeWorkspaceTreeDigest` to the tree the parent's `go.json`
     pins for `knowledge-workspace`: all five subtests of the new test fail;
   - set `upstreamSetupRevision` to the ref the parent's `go.json` names:
     the new test, `TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree`,
     `TestEveryUpstreamSetupFollowsTheAiMemoryRevision`,
     `TestEveryUpstreamSetupFollowsThePythonPresetRevision` and
     `TestThePythonSetupIsThePublishedUpstreamPreset` fail;
   - restore `docs/agents/setup-context.json` from that parent:
     `TestTheSetupManifestIsCurrent` (with `-tags docscontract`) fails with
     "Setup Manifest is stale";
   - restore `internal/baseline/testdata/catalog.digest` from that parent:
     `TestCatalogCompatibility` fails.
8. MUST execute the tests task_01 names in its Verification.
9. MUST record, as the outside-evidence row, the upstream text that
   `_prd.md` → Acceptance evidence names: the new lines of
   `skills/00-setup/knowledge-workspace/SKILL.md` that the adopted Backlog
   Entry under `references/` quotes from the upstream diff, read in the
   audited checkout. The row is checked offline, never blocked: record that
   the audited `.agents/skills/knowledge-workspace/SKILL.md` contains, on
   whitespace-normalized text, the phrases `send research to`,
   `default/marcioaltoe-secondbrain`, `The filesystem inbox was retired`,
   `do not create a filesystem substitute` and
   `there is no outside-session write carve-out`. Name the upstream pull
   request (marcioaltoe/skills#136) and the two message ids and their dates
   as the source.
10. MUST record Success Metric 5 as exactly one row, `blocked (declared
    unreachable)`, against the PRD's single Unreachable Acceptance entry, and
    count it once in `rows_blocked_declared`.
11. MUST perform the glossary and documentation checks: record that the Spec
    declares no glossary term and that the Task commit range does not change
    `CONTEXT.md` (`docs/agents/domain.md`).
12. MUST verify from the repository history:
    - the Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, any production `.go`
      file, any file under `internal/baseline/assets/` other than the five
      setups, any file under `.agents/skills/` or `skills/` other than
      `.agents/skills/knowledge-workspace/SKILL.md`, any file under
      `docs/agents/` other than `setup-context.json`, `docs/user-guide/`,
      `docs/adr/` and `CONTEXT.md`; and in `skills-lock.json` only the
      `knowledge-workspace` entry changes;
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
    `HOME="$(mktemp -d)" TMPDIR="$(mktemp -d)" bin/roundfix archive 0326-the-baseline-setups-stop-sending-agents-to-the-secondbrain-inbox --plan --offline`
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
- [ ] Record the outside source, the Source Baseline row, the declared row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every upstream snapshot names `b56fc5ff`, keeps its membership, and changes only the `knowledge-workspace` tree.
- [ ] This repository's installed `knowledge-workspace` matches its lock and the `go` snapshot and carries the upstream's corrected text.
- [ ] The derived files and the Setup Manifest regenerate with no change, and each sabotage makes its named test fail.
- [ ] The scope audit finds no change outside the declared files.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0326-the-baseline-setups-stop-sending-agents-to-the-secondbrain-inbox/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; User Experience; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-5; Invariants 1-5; Testing Approach; Build Order 2
