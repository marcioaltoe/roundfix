---
task: task_03
spec: 0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration
status: completed
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec makes Doctor and
the Update Check warn when a repository has installed authoring skills it
does not publish and declares no `skills.authoring.claude`, printing the
**Recommended Authoring Declaration** from one table in the binary, and
holds this repository's declaration and the configuration guide equal to
that table. Nothing fails and nothing is written.

Every row runs on this repository's own files, temporary repositories and
temporary homes. No row starts an Agent harness, an Agent session or an ACP
adapter, starts a Run, a Daemon or a Delivery Owner, calls ai-memory, reads
or writes the real `~/.roundfix` or the real home directory, runs the real
`gh`, reaches the network, sends a request to a provider, the judge or a
reviewer, reads another repository, or asserts a wall-clock duration. Every
`go test` it runs carries an explicit `-timeout`.

Every row that runs the product runs a binary the gate builds from the
audited head, never a `roundfix` found on `PATH`. Requirement 1 builds
`bin/roundfix` with `make build`; the product rows of requirements 3, 4 and
6 run `"$b"`, built with
`go build -buildvcs=false -ldflags "-X roundfix/internal/app.Version=0.0.0-dev" -o "$b" ./cmd/roundfix`
into a temporary directory outside the checkout, so the Update Check makes
no release lookup. Every such row takes `m="$(go env GOMODCACHE)"` and
`c="$(go env GOCACHE)"` first, then sets `HOME` and `TMPDIR` to fresh
temporary directories outside the checkout with `GOMODCACHE="$m"` and
`GOCACHE="$c"`, unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`,
`CLAUDE_CODE_SESSION_ID` and every Judge Key Variable for that command,
reads standard input from `/dev/null`, and removes the temporary
directories afterwards. The QA never writes into the checkout while a test
suite runs. The gate keeps its evidence under names the repository does not
ignore, `.txt` for command output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1, Success Metric 2 and Success Metric 3 by
   running, with `-v`, the tests task_01 names in its Verification,
   recording every test and subtest by name, and the whole
   `./internal/config` and `./internal/cli` packages. For the whole
   `./internal/cli` package only, the fresh home keeps the host's Rust
   toolchain: take `RUSTUP_HOME` and `CARGO_HOME` before `HOME` changes
   (their values, or `$HOME/.rustup` and `$HOME/.cargo` when unset) and pass
   them to the command, read-only, so `TestGuidanceCompositionJourney/rust-cli`
   can select its toolchain.
3. MUST reproduce Surface Transcript 1 of `_techspec.md` through `"$b"
   update check` in a temporary Git repository with one commit, where
   `"$b" skills install --target project` installed the skills and
   `.roundfixrc.yml` holds no `skills` key, recording standard output,
   standard error and the exit code, compared line by line with the
   transcript, where each `...` matches zero or more lines.
4. MUST, in the repository of requirement 3, run `"$b" update check
   --format json` and record the whole `authoringModels` member and the
   report's `next`; run `"$b" doctor` and record its `authoring-models:`
   line, the block lines after it, and its exit code, which a failed
   environment check such as `node` may make non-zero; and record a
   listing of every file of that repository with its size and SHA-256
   before and after both commands, unchanged. Then publish
   `skills/<name>/SKILL.md` for all four authoring skills, rerun both
   commands, and record that neither prints an `authoring-models:` line
   or an `authoringModels` member.
5. MUST verify Success Metric 4 by running, with `-v`, the test task_02
   names in its first Verification command, recording both subtests by
   name, the whole `./skills` package and
   `-tags docscontract ./internal/docscontract`; and by searching and
   recording that both copies of the setup reference name the Recommended
   Authoring Declaration, both copies of `modes.md` name `undeclared`,
   each copy under `skills/roundfix/` equals its `.agents/skills/roundfix/`
   source, and the configuration guide's block equals the output of
   `RecommendedAuthoringYAML()`.
6. MUST record the outside-evidence row of `_prd.md` → Success Metric 5
   and Acceptance evidence. Build a temporary Git repository with the
   state
   [`references/2026-11-17-installed-authoring-skills-without-a-declaration-measurement.md`](references/2026-11-17-installed-authoring-skills-without-a-declaration-measurement.md)
   records for Argus: the four authoring skills installed by
   `"$b" skills install --target project`, a `.roundfixrc.yml` with no
   `skills:` key, no `.claude/agents` directory, and no `model:` or
   `effort:` line in the four installed skills. The row passes when
   `"$b" update check --format json` there reports `authoringModels.state`
   `undeclared` with `installed` `["write-idea","write-prd","write-techspec","write-tasks"]`,
   where the measurement recorded `authoringModels` absent with Roundfix
   0.147.1. Record the source as the measurement, not as this Spec's own
   artifacts.
7. MUST run each sabotage below on its own temporary copy of the audited
   head made with `git archive` and committed in a fresh Git repository
   (the suite guard refuses a tree outside Git), with `-timeout` on every
   `go test`, and record that the named test fails on the named subtest. No
   sabotage regenerates or synchronizes anything:
   1. The report ignores `PublishesAuthoringSkill`:
      `TestDoctorRecommendsAnAuthoringDeclarationWhenNoneIsDeclared/published`.
   2. The report runs even when the inspection found a problem:
      `TestUpdateCheckRecommendsAnAuthoringDeclarationWhenNoneIsDeclared/leftover_generated`.
   3. `printDoctorResult` prints continuation lines without their two
      leading spaces:
      `TestDoctorRecommendsAnAuthoringDeclarationWhenNoneIsDeclared/undeclared`.
   4. The rendering gains `tools: [WebSearch]` under `research`:
      `TestRecommendedAuthoringDeclarationRendersValidProjectConfig/exact_bytes`.
   5. The `tasks` model of the table changes from `opus` to `fable`:
      `TestTheRecommendedAuthoringDeclarationMatchesThisRepository/repository_config`.
   6. One line of the configuration guide's block changes:
      `TestTheRecommendedAuthoringDeclarationMatchesThisRepository/configuration_guide`.
8. MUST execute the tests tasks 01 and 02 name in their Verification.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that
   the PRD and the TechSpec declare `adds:` for
   **Recommended Authoring Declaration** and `changes:` for
   **Authoring Model Declaration**, that task_02 wrote those two entries
   and no other, and that `TestGlossaryPhrasesAreDefinedTerms` passes.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `skills-lock.json`, `CHANGELOG.md`,
      `internal/app/version.go`, `dist/npm/roundfix/package.json`,
      `internal/cli/cli_test.go` apart from the `assertCheckResult` helper
      the authorization record's amendment of 2026-11-17 grants (task_01
      may change only that helper, to `reflect.DeepEqual`),
      `internal/docscontract/publicdocs_test.go`,
      `internal/config/recommendations.go`, `internal/agent/catalog.go`,
      `docs/references/model-selection.md`,
      `docs/references/coverage-record.json`,
      `docs/references/skill-coverage.json`, any file under `docs/agents/`
      or `internal/baseline/assets/`, any Vendored Skill under
      `.agents/skills`, any ADR other than ADR-0401, and any owned skill
      other than `roundfix`;
    - `.roundfixrc.yml` differs from the Delivery Base only in the `prd`,
      `techspec` and `tasks` entries of `skills.authoring.claude`, and
      `.claude/agents/` differs only in `roundfix-prd-author.md`,
      `roundfix-techspec-author.md` and `roundfix-tasks-author.md`, each
      matching what `bin/roundfix skills install --target project` writes;
    - ADR-0401 differs from the Delivery Base (the merge base of the
      audited head with the default branch) only by the added amendment
      paragraph;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent
    of each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write
    the report's Outcome as one paragraph that does not end with a colon.
    It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix` or the real home directory;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent harness, an Agent session, an ACP adapter, a Run, a
      Daemon or a Delivery Owner;
    - send any message, or any request to a provider, the judge or a
      reviewer;
    - read Argus or any other repository;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binaries and run the repository Verification.
- [ ] Run the named tests of both Tasks and their packages.
- [ ] Reproduce the Surface Transcript and the product rows on temporary repositories.
- [ ] Reproduce the measured Argus state in a fixture.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the searches, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Installed authoring skills without a declaration get the warning and the block from Doctor and the Update Check, with no exit change and no write.
- [ ] Published skills, a declaration and a leftover generated file keep today's output.
- [ ] This repository's declaration and the configuration guide equal the recommendation, and each sabotage makes its named test fail.
- [ ] The measured Argus state, reproduced in a fixture, is reported as `undeclared`.

## Verification

- `newest="$(find 'docs/specs/0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-4; Surface Transcript 1; Interfaces; Testing Approach; Build Order 3
- [`references/2026-11-17-installed-authoring-skills-without-a-declaration-measurement.md`](references/2026-11-17-installed-authoring-skills-without-a-declaration-measurement.md)
- ADR-0415, ADR-0401, ADR-0408, ADR-0402
