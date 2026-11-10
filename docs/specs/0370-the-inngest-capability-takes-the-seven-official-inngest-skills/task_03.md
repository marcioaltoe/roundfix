---
task: task_03
spec: 0370-the-inngest-capability-takes-the-seven-official-inngest-skills
status: completed
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec makes the `backend`
module require and dispatch the seven official `inngest-*` skills in place
of `inngest`, makes `inngest` a Retired Skill (ADR-0380), moves the upstream
Setup Snapshots and the vendored catalog skills to catalog commit `e45466d2`
through the Catalog Skill Bump, and documents the change in ADR-0206,
`CONTEXT.md`, the `baseline` command guide and the Roundfix skill.

Every row runs in this repository or in temporary directories it creates and
removes. No row reaches the network, reads another repository, starts an
Agent session or an ACP adapter, reads or writes the real `~/.roundfix`,
sends a request to a provider, the judge or a reviewer, records an
environment value, or sends a signal to any process. The only outside path a
row reads is the operator's clone
`/Users/marcio/.roundfix-operator/scratch/spec-0370/skills`, read-only, in
row 5. Every `go test` carries an explicit `-timeout` and keeps the
session's own `HOME` and `TMPDIR`. Every row that runs `bin/roundfix` or
`make skills-third-party-bump` sets `HOME` and `TMPDIR` to fresh temporary
directories, with `GOMODCACHE` and `GOCACHE` kept at their `go env` values,
and removes them afterwards. The QA never writes into the checkout while a
test suite runs. The gate keeps its evidence under names the repository does
not ignore, `.txt` for command output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST run, with `go test -count=1 -timeout 30m -v`, the tests task_01
   names (`TestTheBackendModuleTakesTheSevenOfficialInngestSkills` in
   `./internal/baseline`; `TestTheBackendModuleRequiresTheSevenInngestSkills`
   and `TestBaselineUpdateApprovedApplyDeletesTheReplacedInngestSkill` with
   its `unedited` and `edited` subtests in `./internal/cli`), then the whole
   `./internal/baseline` and `./skills` packages, `-run 'Baseline'
   ./internal/cli` and `-tags docscontract ./internal/docscontract`,
   recording each pass and that none of the new tests reports `SKIP`.
3. MUST record Success Metric 4 on the Delivery Base (the merge base of the
   audited head with `main`): unpack a `git archive` of that commit in a
   temporary directory, copy the two new test files of the audited head into
   it, run the three new tests there with `-timeout`, and record that
   `TestTheBackendModuleTakesTheSevenOfficialInngestSkills`,
   `TestTheBackendModuleRequiresTheSevenInngestSkills` and
   `TestBaselineUpdateApprovedApplyDeletesTheReplacedInngestSkill/unedited`
   fail there.
4. MUST record Success Metrics 1 and 2 on the audited head, with `jq` and
   `grep` only:
   - `internal/baseline/assets/modules/backend.json` requires and dispatches
     each of the seven skills with the one trigger `_techspec.md` → API
     Contract 1 names, and names `inngest` in neither list;
   - the formatter golden
     `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`
     holds each trigger line of `_prd.md` → User Experience and no
     `trigger.backend.inngest:` line;
   - each of the five snapshots under `internal/baseline/assets/setups/`
     names `e45466d2b9e15480eb286a4c19b513494a867bce` in every `github`
     source and never `c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84`;
     `typescript.json` and `go-cli-typescript-bun.json` pin each of the
     seven, with equal tree digests in both, and no `inngest`;
   - every `marcioaltoe/skills` entry of `skills-lock.json` names
     `e45466d2b9e15480eb286a4c19b513494a867bce`.
5. MUST record that the bump is repeatable: unpack a `git archive` of the
   audited head in a temporary directory, make it a Git repository with one
   commit, run there
   `GIT_ALLOW_PROTOCOL=file make skills-third-party-bump SKILLS_CATALOG=/Users/marcio/.roundfix-operator/scratch/spec-0370/skills SKILLS_REF=e45466d2b9e15480eb286a4c19b513494a867bce`,
   and record exit `0`, the line
   `skills-third-party-bump: vendored catalog skills already match e45466d2b9e15480eb286a4c19b513494a867bce`
   and an empty `git status --porcelain`. When the clone does not hold that
   commit, record the row `blocked (environment: the operator's catalog clone lacks e45466d2)`.
6. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails:
   - the `"inngest": {},` line removed from `retiredSkills` in
     `internal/baseline/retired_skills.go`:
     `TestTheBackendModuleTakesTheSevenOfficialInngestSkills` and
     `TestBaselineUpdateApprovedApplyDeletesTheReplacedInngestSkill` fail;
   - the `"inngest-middleware",` line removed from `requiredSkills` in
     `internal/baseline/assets/modules/backend.json`:
     `TestTheBackendModuleTakesTheSevenOfficialInngestSkills` and
     `TestTheBackendModuleRequiresTheSevenInngestSkills` fail;
   - `its local Dev Server.` changed to `its Dev Server.` in the
     `inngest-cli` trigger of that module:
     `TestTheBackendModuleTakesTheSevenOfficialInngestSkills` fails.
7. MUST execute the Verification commands every non-QA Task declares.
8. MUST record, as an outside-evidence row, sources this Spec did not write,
   with their dates and not blocked: the marcioaltoe-skills report quoted in
   `_prd.md` → Acceptance evidence, compared offline with rows 4 and 6 (each
   of its seven slugs sits at `skills/05-implementation-loop/<slug>` in both
   TypeScript snapshots, and each trigger follows its suggested mapping);
   and the published `inngest/inngest-skills` README passages that section
   quotes, compared with the `inngest-setup`, `inngest-flow-control`,
   `inngest-middleware` and `inngest-cli` trigger texts.
9. MUST record that the delivered
   `bin/roundfix spec check 0370-the-inngest-capability-takes-the-seven-official-inngest-skills --strict`
   on the audited head reports no finding.
10. MUST record the PRD's one Unreachable Acceptance entry (the adopter
    migration after the release) as exactly one row,
    `blocked (declared: Adopter migration)`, its leading label, counted once
    in `rows_blocked_declared`, supported by rows 2, 4 and 5.
11. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare `- changes:` for **Retired Skill**, that
    the Task commits change exactly that entry of `CONTEXT.md`, and that
    `go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryStatesEachRuleOnce)$'`
    passes.
12. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - ADR-0206 carries the `Amended by [ADR-0380]` paragraph and
      `updated_at: 2026-11-09T12:00:00Z`;
    - the Task commit range changes none of these: the Makefile, any
      workflow, `go.mod`, `go.sum`, `.roundfixrc.yml`, `CHANGELOG.md`,
      `internal/app/version.go`, `dist/npm/**`, any Baseline module other
      than `backend`, any production Go file other than
      `internal/baseline/retired_skills.go`, any file under
      `internal/baseline/testdata/parity-corpus/` other than
      `v1/fixtures/asset-sync.json` and `v1/manifest.json`, any file under
      `.agents/skills/` or `skills/` other than
      `.agents/skills/context7-cli/references/setup.md` and the Roundfix
      skill's `SKILL.md` and `references/baseline.md` pairs, any ADR other
      than ADR-0206, or text inside the `### QA settlement` section of any
      skill.

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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0370-the-inngest-capability-takes-the-seven-official-inngest-skills --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read another repository or `/Users/marcio/dev/skills`, write the
      operator's clone, or run `git fetch` against a network remote,
      `git push` or `gh`;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - send a signal to any process;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification and the named suites.
- [ ] Run the new tests on the Delivery Base and record the module, golden, snapshot and lock facts.
- [ ] Repeat the bump on a copy, and run each sabotage on a temporary copy.
- [ ] Record the outside sources, the unreachable row, the glossary check, the self-check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The module, the golden, the snapshots and the lock show the seven skills at `e45466d2` and no `inngest`, and a second bump changes nothing.
- [ ] The new tests fail on the Delivery Base, pass on the audited head, and each sabotage makes its named test fail.
- [ ] This Spec passes its own strict check, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0370-the-inngest-capability-takes-the-seven-official-inngest-skills/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-3; Testing Approach; Build Order 3
