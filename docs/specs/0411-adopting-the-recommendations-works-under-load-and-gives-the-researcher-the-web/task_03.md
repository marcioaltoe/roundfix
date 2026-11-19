---
task: task_03
spec: 0411-adopting-the-recommendations-works-under-load-and-gives-the-researcher-the-web
status: completed
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. The Spec makes `profiles configure` and `profiles check --apply` prove under the configured `defaults.profile_proof_setup_timeout`, read by a narrow reader in `internal/config`, and makes the Recommended Authoring Declaration list `tools: [WebSearch, WebFetch]` for `research`, with Doctor and the Update Check naming those tools when a declared research entry has no `tools` key; the guides, the Roundfix skill, the glossary entry **Recommended Authoring Declaration**, ADR-0415 and ADR-0441 state both rules.

Every row runs on this repository's own files, temporary homes, temporary repositories and temporary copies of the audited head. No row:

- starts an Agent harness, an Agent session or an ACP adapter;
- starts a Run, a Daemon or a Delivery Owner;
- calls ai-memory, runs the real `gh` or reaches the network;
- reads or writes the real `~/.roundfix`, or reads the real home directory except for the toolchain locations below;
- sends a request to a provider, the judge or a reviewer;
- asserts a wall-clock duration.

Every `go test` it runs carries an explicit `-timeout`.

Every `go test` package or test run, and the repository Verification, keeps the session's own `HOME` and `TMPDIR`. When the gate isolates `HOME` for any row that runs `./internal/cli` or the repository Verification, it keeps the host Rust toolchain readable, with `RUSTUP_HOME` and `CARGO_HOME` pointing read-only at the host's and `PATH` unchanged. Every row that runs the product runs `bin/roundfix`, the binary `make build` made from the audited head, and never a `roundfix` found on `PATH`. Each such row prepares its environment the same way:

- before `HOME` changes it takes `m="$(go env GOMODCACHE)"`, `c="$(go env GOCACHE)"`, `r="${RUSTUP_HOME:-$HOME/.rustup}"` and `k="${CARGO_HOME:-$HOME/.cargo}"`;
- it sets `HOME` and `TMPDIR` to fresh temporary directories outside the checkout, with `GOMODCACHE="$m"`, `GOCACHE="$c"`, `RUSTUP_HOME="$r"` and `CARGO_HOME="$k"`, which the command only reads;
- it puts first on `PATH` a temporary directory holding an executable `gh` that writes nothing and exits 1;
- it unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`, `CLAUDE_CODE_SESSION_ID` and every Judge Key Variable for that command;
- it reads standard input from `/dev/null` and removes the temporary directories afterwards, after `chmod -R u+w` on them.

The QA never writes into the checkout while a test suite runs. It keeps its evidence under names the repository does not ignore, `.txt` for command output, and links only files it kept. Evidence records environment variable names, never their values.

## Requirements

1. MUST build the binary with `make build` and run the repository Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 and 2 by running, with `-v`, `TestTheProofSetupDeadlineReaderAgreesWithLoad` in `./internal/config` and `TestProfileAdoptionPassesTheConfiguredSetupDeadline`, `TestProfilesConfigureRefusesAnInvalidSetupDeadlineTranscript`, `TestProfilesConfigureNamesTheFileAndLineOfABrokenConfig`, `TestProfileProofPassesTheConfiguredSetupDeadline`, `TestProfilesConfigureProofScope` and `TestProfilesConfigureChangeSummary` in `./internal/cli`, recording every test and subtest by name, then the whole `./internal/config` and `./internal/cli` packages, recording each pass.
3. MUST reproduce `_techspec.md` → Surface Transcript 1 through `bin/roundfix`, in a fresh home prepared as above, in a temporary Git repository outside the checkout whose `.roundfixrc.yml` holds `defaults:` with `profile_proof_setup_timeout: 0s`, with a fragment file outside the repository holding one `backend` profile, `preferred: {runtime: codex, model: written-model, reasoning_effort: high}` and one fallback `{runtime: claude, model: written-fallback, reasoning_effort: xhigh}`: record that standard output is empty, standard error is exactly the two lines of the transcript, the exit code is 2, and the fresh home holds no `.roundfix/config.yml` afterwards. Because the deadline is refused before any proof, the row starts no Agent session.
   Then reproduce `_techspec.md` → Surface Transcript 2 the same way, in another temporary Git repository whose `.roundfixrc.yml` holds the three lines `defaults:`, `  profile_proof_setup_timeout: 90s` and `broken: a: b`: record that standard output is empty, that standard error is `roundfix: profiles failed: parse config "<path>": yaml: line 3: mapping values are not allowed in this context` with `<path>` the repository's `.roundfixrc.yml`, then `Run 'roundfix profiles --help' for usage.`, that the exit code is 2, and that the fresh home holds no `.roundfix/config.yml` afterwards.
4. MUST verify Success Metrics 3, 4 and 5 by running, with `-v`, `TestTheChecksNameTheRecommendedResearchToolsWhenResearchDeclaresNone` and `TestRecommendedAuthoringDeclarationRendersValidProjectConfig`, `TestTheRecommendedResearchSubagentReadsTheWebAndHoldsNoShellOrWriteTool` and `TestTheResearchSubagentNeverHoldsAShellOrAWriteTool` in `./skills`, and `TestTheRecommendedAuthoringDeclarationMatchesThisRepository` with `-tags docscontract`, recording every test and subtest by name, then the whole `./skills` package, `-tags docscontract ./internal/docscontract` and `make skills-sync-check`, recording each pass.
5. MUST verify Success Metric 4 through `bin/roundfix`, in a fresh home prepared as above, in a temporary Git repository outside the checkout with one empty commit:
   - write a `.roundfixrc.yml` declaring `skills:`, `authoring:`, `claude:` with `prd: {model: opus, reasoning_effort: xhigh}` and `research: {model: sonnet}`, run `bin/roundfix skills install --target project`, then create the directory `.claude/skills` and in it the symbolic link `write-prd` to `../../.agents/skills/write-prd`;
   - run `bin/roundfix update check --format json` and record that `authoringModels` is `{"state":"applied","harness":"claude","declared":["prd","research"],"problems":[],"recommendedResearchTools":["WebSearch","WebFetch"]}`, then `bin/roundfix update check` and record the line `Authoring models: research declares no tools, so the Research Subagent cannot read the web; add tools: [WebSearch, WebFetch] under skills.authoring.claude.research in .roundfixrc.yml, then apply with roundfix skills install --target project.`;
   - change the research entry to `{model: sonnet, tools: []}`, run `bin/roundfix skills install --target project` again and both checks, and record that `authoringModels` lacks `recommendedResearchTools` and no `Authoring models:` line is printed;
   - in a second temporary Git repository with one empty commit and no `.roundfixrc.yml`, run `bin/roundfix skills install --target project` and `bin/roundfix update check --format json`, and record that `authoringModels.state` is `undeclared` and its `recommendation` ends with the line `        tools: [WebSearch, WebFetch]`;
   - record that every one of these commands exits 0.
6. MUST record the outside-evidence row of `_prd.md` → Success Metric 6 and Acceptance evidence, offline: quote the field reports of fiscus and fluxus from `references/2026-11-18-maintainer-answers.md` and the Claude Code passage "`tools` only: the subagent gets only the listed tools." with its URL from `references/2026-11-18-characterization.md` section 5, as published literature, and record that requirements 2 and 4 show the commands now pass the configured deadline and the recommended subagent lists the two web tools.
7. MUST record `_prd.md` → Unreachable Acceptance as one blocked row, `declared: on a loaded machine, an adopter adopts a Recommended Profile with profiles check --apply after raising defaults.profile_proof_setup_timeout`, citing the PRD entry, and count it as one row.
8. MUST run each sabotage below on its own temporary copy of the audited head made with `git archive` and committed in a fresh Git repository outside the checkout, because the suite guard refuses a tree outside Git. Each sabotage keeps the package compiling, so the named test fails on its assertion; record that it fails on the named subtest. No sabotage regenerates or synchronizes anything.
   1. `ProfileProofSetupTimeout` reads Project Config before User Config: `TestTheProofSetupDeadlineReaderAgreesWithLoad/project_wins`.
   2. `ProfileProofSetupTimeout` drops its range check: `TestTheProofSetupDeadlineReaderAgreesWithLoad/zero`.
   3. `writeProfilesConfiguration` passes `0 * ` the reader's value as the deadline: `TestProfileAdoptionPassesTheConfiguredSetupDeadline/configure/project_scope_reads_user`.
   4. `authoringResearchToolsGap` treats an empty list as missing: `TestTheChecksNameTheRecommendedResearchToolsWhenResearchDeclaresNone/empty_tools_list`.
   5. `authoringResearchToolsGap` always returns nil: `TestTheChecksNameTheRecommendedResearchToolsWhenResearchDeclaresNone/no_tools_key`.
   6. The Update Check sets `recommendedResearchTools` beside an inspection problem: `TestTheChecksNameTheRecommendedResearchToolsWhenResearchDeclaresNone/a_problem_keeps_its_report`.
   7. The recommendation lists only `WebSearch`: `TestTheRecommendedResearchSubagentReadsTheWebAndHoldsNoShellOrWriteTool`.
   8. `ProfileProofSetupTimeout` returns the YAML parser's error without the file name: `TestProfilesConfigureNamesTheFileAndLineOfABrokenConfig`.
9. MUST execute the tests task_01 and task_02 name in their Verification.
10. MUST search and record:
    - the sentences of `_techspec.md` → API Contract 2 in `docs/user-guide/configuration.md`, `docs/user-guide/commands/profiles.md` and both copies of the Roundfix skill's `profiles` reference;
    - the changes of `_techspec.md` → API Contract 5 in `docs/user-guide/configuration.md`, `docs/user-guide/commands/doctor.md`, `docs/user-guide/commands/update.md`, both copies of the Roundfix skill's `setup` and `modes` references, the glossary entry and ADR-0415, with the moved `updated_at` of ADR-0415;
    - that `TestTheSkillCoverageMapIsCurrent` passes and `skills/testdata/owned-skill-versions.json` records the version both copies of the Roundfix skill's `SKILL.md` declare;
    - that `ProfileProofSetupTimeout` and `authoringResearchToolsGap` call no function that writes, creates, renames or removes a file;
    - that `ModelRecommendationSnapshotVersion` in `internal/config/recommendations.go` is the same at the audited head as at the Delivery Base.
11. MUST perform the glossary check (`docs/agents/domain.md`). Record that the PRD declares `changes:` for **Recommended Authoring Declaration**, that task_02 wrote that entry and no other, and that `TestGlossaryPhrasesAreDefinedTerms` passes.
12. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI workflows, `go.mod`, `go.sum`, `skills-lock.json`, `.roundfixrc.yml`, `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/roundfix/package.json`, `docs/references/coverage-record.json`, `docs/references/skill-coverage.json`, `skills/authoring.go`, `internal/config/authoring.go`, `internal/agent/fallback.go`, any file under `internal/baselineacp/`, `internal/baseline/` or `docs/agents/`, any owned skill other than the Roundfix skill, any Roundfix skill file other than its `SKILL.md` and its `profiles`, `setup` and `modes` references in both copies, any Vendored Skill under `.agents/skills`, and any ADR other than ADR-0415;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
13. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance. It supports the row with the pre-PR equivalent of each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
14. MUST, for the Outcome for the Archive Record, run only `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0411-adopting-the-recommendations-works-under-load-and-gives-the-researcher-the-web --plan --offline; rm -rf "$h" "$t"` with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It writes the report's Outcome as one paragraph that does not end with a colon. It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix`, or read the real home directory except for the toolchain locations above;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent harness, an Agent session, an ACP adapter, a Run, a Daemon or a Delivery Owner;
    - send any message, or any request to a provider, the judge or a reviewer;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests of both Tasks and the packages.
- [ ] Reproduce Surface Transcripts 1 and 2 and the Update Check rows through the binary.
- [ ] Record the outside evidence and the unreachable row.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the searches, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Both profile writes prove under the configured deadline, and an invalid deadline is refused before any proof through the product.
- [ ] The recommendation and the subagent it produces hold the two web tools and no shell or write tool, and the checks name the tools only for a research entry without a `tools` key.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0411-adopting-the-recommendations-works-under-load-and-gives-the-researcher-the-web/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Interfaces; API Contracts 1-5; Surface Transcript 1; Surface Transcript 2; Testing Approach; Build Order 3
- [`references/2026-11-18-characterization.md`](references/2026-11-18-characterization.md)
- [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md)
- ADR-0441, ADR-0415, ADR-0298, ADR-0402
