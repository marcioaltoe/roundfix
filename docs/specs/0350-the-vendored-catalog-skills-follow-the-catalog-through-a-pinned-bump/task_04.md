---
task: task_04
spec: 0350-the-vendored-catalog-skills-follow-the-catalog-through-a-pinned-bump
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec adds the Catalog Skill Bump
(`make skills-third-party-bump`), runs it once so the 33 vendored catalog
skills, their lock entries and the five upstream Setup Snapshots name one
Catalog Revision, makes the setup tests read that revision from the lock,
adds a documentation contract that holds every vendored folder to its
snapshot tree, and names the bump in `CONTEXT.md` and the repository rules.

Every row runs on this repository's own files, temporary copies of the
audited head and temporary repositories. No row reads or writes the real
`~/.roundfix`, reads `~/dev/skills` or another repository, reaches the
network, sends a request to a provider, the judge or a reviewer, or reads a
credential. Only Requirement 6 reads the operator's clone
`/Users/marcio/.roundfix-operator/scratch/spec-0350/skills`, read-only, with
`GIT_ALLOW_PROTOCOL=file`. Every `go test` carries an explicit `-timeout`;
every row that runs `bin/roundfix` or the bump sets `HOME` and `TMPDIR` to
fresh temporary directories with `GOMODCACHE` and `GOCACHE` taken from
`go env` first, and removes them afterwards. The QA never writes into the
checkout while a test suite runs. `roundfix doctor` is not run: its
`profiles:` line starts agent sessions; Requirement 5 reads the same skills
comparisons through the managed refresh preview.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1: run
   `go test -timeout 20m -count=1 -v ./internal/verifyselect -run '^TestTheSkillsBump'`
   and record the pass of `TestTheSkillsBumpPinsTheCatalogMainCommit` and of
   each of the six subtests of `TestTheSkillsBumpRefusesWhatItMustNotPin`,
   and that the audited `Makefile` has no `skills-update:` target.
3. MUST verify Success Metric 2 and Invariants 1 and 2 from the audited
   files: every `skills-lock.json` entry has `source: marcioaltoe/skills`
   and `ref` `c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84`, there are 33 of them,
   none is one of the 12 owned skills (`archive-spec`, `brainstorming`,
   `business-analyst`, `evidence-gate`, `implement-task`, `qa-gate`,
   `roundfix`, `setup-context-driven`, `write-idea`, `write-prd`,
   `write-tasks`, `write-techspec`); each snapshot with a `github` source and
   every `github` skill of the five snapshots name that same commit, and the
   composed `go-cli-typescript-bun.json` keeps its `composed` source.
4. MUST verify Success Metric 3 and Invariant 3: in an empty temporary
   repository, fetch the Catalog Revision the audited lock names from the
   operator's clone with `GIT_ALLOW_PROTOCOL=file`, check it out, and record
   that each of the 33 `.agents/skills/<name>` folders equals the catalog's
   `skills/*/<name>` folder (`diff -r` exits `0`); that
   `.agents/skills/systematic-debugging` holds none of `CREATION-LOG.md`,
   `test-academic.md`, `test-pressure-1.md`, `test-pressure-2.md` and
   `test-pressure-3.md`; and that the 12 owned folders under `.agents/skills`
   and `skills/` are byte-identical to the Delivery Base. When the clone is
   unreadable, record the row `blocked (environment: operator clone
   unreadable)` with the reason.
5. MUST verify Success Metric 4 and Invariant 4: run
   `go test -timeout 20m -count=1 ./internal/baseline ./skills ./internal/verifyselect`
   and `go test -timeout 20m -count=1 -tags docscontract ./internal/baseline ./internal/docscontract`,
   recording the passes of
   `TestEveryVendoredCatalogSkillIsTheTreeItsCatalogCommitPublishes`,
   `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot`,
   `TestEveryUpstreamSetupFollowsTheAiMemoryRevision`,
   `TestEveryUpstreamSetupFollowsThePythonPresetRevision`,
   `TestThePythonSetupIsThePublishedUpstreamPreset`,
   `TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree`,
   `TestTheTypeScriptSetupsCarryTheCompletedUiUxProMaxTree`,
   `TestEverySetupCarriesTheKnowledgeWorkspaceThatRetiredTheSecondbrainInbox`,
   `TestCatalogCompatibility`, `TestCheckRepositoryMatchesRealRepository`,
   `TestRepositoryContractTestsDeclareTheirRelevance` and
   `TestTheSetupManifestIsCurrent`. After the suites, in a temporary
   `git clone --no-local` of this checkout detached at the audited head, run
   `make baseline-digests` and record `"changed":false`; run
   `bin/roundfix baseline update --repo <clone> --no-skills --no-history --yes --format text`
   and record `File changes: 0`; run
   `bin/roundfix baseline update --repo <clone> --no-history --format json`
   and record that its `skills` object lists nothing under `missing`,
   `drifted` or `differing`.
6. MUST verify Success Metric 5 and Invariant 5: in another temporary
   `git clone --no-local` detached at the audited head, run
   `GIT_ALLOW_PROTOCOL=file make skills-third-party-bump SKILLS_CATALOG=/Users/marcio/.roundfix-operator/scratch/spec-0350/skills SKILLS_REF=<the Catalog Revision the audited lock names>`
   and record exit `0`, its final line and that `git status --porcelain` of
   the clone is empty afterwards. When the clone is unreadable, record the
   row `blocked (environment: operator clone unreadable)` with the reason.
7. MUST verify Success Metric 6 and Invariant 6: no `*_test.go` file under
   `internal/` or `skills/` names `86df8608726979064a8fe46003dfaf403ffe0087`
   or `c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84`, or any `treeDigest` value of
   the audited `internal/baseline/assets/setups/*.json`; and
   `upstreamSetupRevision` is declared once, as a function.
8. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive` and recorded as a Git repository, with `-timeout` on
   every `go test`, and record that the named test fails, with the input that
   makes the mutation observable:
   - in the `Makefile`, replace the owned-skill refusal with `true`: subtest
     `owned skill in the lock` (its lock names `roundfix`) fails;
   - replace the post-restore `diff -r` check with `true`: subtest `tree the
     catalog does not publish` (its `beta` folder differs) fails;
   - delete the second `$(GO) build` line: `TestTheSkillsBumpPinsTheCatalogMainCommit`
     fails on the build count;
   - replace the stale-ref check with `true`: subtest `entry left at another
     commit` (its `beta` entry names `main`) fails;
   - append a line to `.agents/skills/ai-memory-retrieval/SKILL.md` and
     re-hash its lock entry with `SkillFolderHash`: `./skills` and
     `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` pass, and
     `TestEveryVendoredCatalogSkillIsTheTreeItsCatalogCommitPublishes` fails
     naming `ai-memory-retrieval` (the skill `go-cli-tui` does not require);
   - restore the five `systematic-debugging` files from the Delivery Base:
     the new contract fails naming `systematic-debugging`, and the managed
     refresh preview lists it under `missing` with the reason that it no
     longer matches its lock entry (the restored files break its
     `computedHash`, so the preview reports the hash drift, not a Setup
     Snapshot Difference);
   - set one catalog entry's `ref` to `main`: the six setup revision and tree
     tests fail on the lock read;
   - restore the five setup snapshots from the Delivery Base: the six setup
     tests fail naming the revision the lock does not hold;
   - set the `ui-ux-pro-max` `treeDigest` of `typescript.json` to the
     `incompleteUiUxProMaxTree` constant:
     `TestTheTypeScriptSetupsCarryTheCompletedUiUxProMaxTree` fails.
9. MUST execute the tests task_01 and task_02 name in their Verification.
10. MUST record, as the outside-evidence row, the message that `_prd.md` →
    Acceptance evidence names: the ai-memory message `01a1a608` of
    marcioaltoe-skills reproduced in the adopted Backlog Entry under
    `references/`, a source this Spec did not write. The row is checked
    offline, never blocked: record, against the audited head, each item of
    its RESULTADO ESPERADO and O QUE FAZER: 33 entries following the
    catalog's `main` with a coherent `computedHash` (`TestCheckRepositoryMatchesRealRepository`
    and the new contract), none of the 12 named owned skills in the lock,
    `systematic-debugging` without the five named files, the 12 owned skills
    intact, a repeatable update (`make skills-third-party-bump`), and the
    no-local-edit rule (the repository rule and the new contract). Name the
    message id, its sender and its date.
11. MUST perform the glossary and documentation checks: **Catalog Revision**
    and **Catalog Skill Bump** are Glossary Definitions in `CONTEXT.md`, each
    linking ADR-0360, and the audited `docs/agents/specific-repository.md`
    contains, on whitespace-normalized text, `change only through` followed
    by `make skills-third-party-bump`, in code formatting or not
    (`docs/agents/domain.md`).
12. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes no production `.go` file, `go.mod`,
      `go.sum`, `.github/workflows`, `.roundfixrc.yml`, Baseline module,
      profile or template, owned skill, user guide page or ADR, and under
      `.agents/skills/` only deletes the five `systematic-debugging` files;
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
    `bin/roundfix archive 0350-the-vendored-catalog-skills-follow-the-catalog-through-a-pinned-bump --plan --offline`
    with fresh temporary `HOME` and `TMPDIR` and with every Judge Key Variable
    and `TYPESAFE_API_KEY` unset, and write the report's `## Outcome` as one
    paragraph that does not end with a colon. It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests, the catalog comparison, the regeneration rows, the second bump and each sabotage on temporary copies.
- [ ] Record the outside source, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The lock and the five snapshots name one Catalog Revision; no owned skill is in the lock.
- [ ] Every vendored folder equals the catalog's tree; `systematic-debugging` lost the five files.
- [ ] The derived files regenerate with no change, a second bump changes nothing, and each sabotage makes its named test fail.
- [ ] The scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0350-the-vendored-catalog-skills-follow-the-catalog-through-a-pinned-bump/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; User Experience; Success Metrics 1-6; Acceptance evidence; Glossary
- `_techspec.md` → API Contracts 1-5; Invariants 1-6; Testing Approach; Build Order 4
