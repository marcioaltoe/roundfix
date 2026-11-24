---
task: task_02
spec: 0419-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed
status: pending
type: qa
complexity: high
---

# Task 02: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. The Spec makes each pre-PR review round that sends a prompt record the model of its Agent Session in `lineage.model`, and makes round 2 of a Reviewer Lineage continue the round-1 session only when the review profile's selection at the recorded index still names that model; otherwise round 2 ends the round-1 session and opens a new one with the round-2 delta prompt. It records the rule in ADR-0451, amends ADR-0197, revises the glossary entry Reviewer Lineage, and rewrites two paragraphs of the review guide and of the Roundfix skill's review reference. It answers the finding reported by fiscus on 2026-11-23 (ai-memory message 01a20a55-66b6-7b31-bd03-ed4f3d6a1fbc). No command, flag, output line, exit code, help text or key of Project Config is added, and the classifier of the reviewer's answer is not changed.

Every row runs on this repository's own files, temporary homes, temporary Git repositories and temporary copies of the audited head or of the Delivery Base. No row:

- starts an Agent harness, an Agent session, an Author Subagent or an ACP adapter, or asks a person a question;
- starts a real Run, a Daemon or a Delivery Owner;
- calls ai-memory, the Indexus or any MCP server, runs the real `gh` or reaches the network;
- reads another repository, the Fiscus's included;
- reads or writes the real `~/.roundfix`, or reads the real home directory except for the toolchain locations below;
- signals a process it did not start, or ends a process by name;
- sends a message or a request to a provider, the judge or a reviewer;
- reads a credential or asserts a wall-clock duration.

Every `go test` it runs carries an explicit `-timeout`.

The gate follows the `qa-gate` section "A fresh home keeps the host toolchains". Every `go test` package or test run, every `make` target and the repository Verification keep the session's own `HOME` and `TMPDIR`. Only a row that runs the built product uses a fresh `HOME`. Every such row runs `bin/roundfix`, the binary `make build` made from the audited head, or the Delivery Base binary named in requirement 6, and never a `roundfix` found on `PATH`, and prepares its environment the same way:

- before `HOME` changes it takes `m="$(go env GOMODCACHE)"`, `c="$(go env GOCACHE)"`, `r="${RUSTUP_HOME:-$HOME/.rustup}"` and `k="${CARGO_HOME:-$HOME/.cargo}"`;
- it sets `HOME` and `TMPDIR` to fresh temporary directories outside the checkout, with `GOMODCACHE="$m"`, `GOCACHE="$c"`, `RUSTUP_HOME="$r"` and `CARGO_HOME="$k"`, which the command only reads, and `PATH` otherwise unchanged;
- it puts first on `PATH` a temporary directory holding an executable `gh` that writes nothing and exits 1;
- it unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`, `CLAUDE_CODE_SESSION_ID` and every Judge Key Variable for that command;
- it reads standard input from `/dev/null` and removes the temporary directories afterwards, after `chmod -R u+w` on them.

A command that fails only because a fresh `HOME` hid a toolchain is a gate error, never an environment block: the gate runs it again with the toolchain locations and records that result. A temporary copy of the audited head or of the Delivery Base is made with `git archive` of that commit and committed in a fresh Git repository outside the checkout, because the suite guard refuses a tree outside Git. The QA never writes into the checkout while a test suite runs. It keeps its evidence under names the repository does not ignore, `.txt` for command output, and links only files it kept. Evidence records environment variable names, never their values.

## Requirements

1. MUST build the binary with `make build` and run the repository Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 and 6 and Core Features 2 and 3 by running, with `-v` over `./internal/cli`, `TestReviewRoundTwoOpensANewSessionWhenTheConfiguredModelChanged`, recording the test and its subtest `the_classifier_still_refuses_the_reported_answer` by name; and record from the test that the round-1 session is ended before the new name is prepared, that round 2 sends the round-2 delta prompt, stays round 2 with `continued: false`, prints nothing on standard error and records `lineage.model` of `gpt-6.1-sol` (invariants 1 to 3).
3. MUST verify Success Metric 6, the outside-evidence row, with its date and not blocked: compare, looking nothing up, the answer the Fiscus measured on Roundfix 0.153.0 in `references/sources.md`, section 1, with the notice the fake runner `modelNoticeReviewRunner` builds for `gpt-6-luna` and `gpt-6.1-sol` followed by `No findings.`, and record that the two hold the same warning sentence, the same two model names and the same verdict, and that the subtest of requirement 2 refuses that answer with the reason the Fiscus quoted, `unclassifiable agent output: no-findings verdict does not account for other content` (invariant 5, ADR-0174). Quote the passage of the source the comparison rests on.
4. MUST verify Success Metrics 2 and 4 and Core Feature 1 by running, with `-v` over `./internal/cli`, `TestReviewRoundTwoOpensANewSessionWhenTheRecordHoldsNoModel` and `TestReviewRecordsTheModelOfTheSessionOfEachRound`, recording each test and the subtests `preferred_then_continued` and `fallback` by name; then `TestACeilingRecordNamesNoSelection`, recording that the `ceiling-closed` lineage keeps its bytes and holds no model (invariant 4).
5. MUST verify Success Metric 3 by running, with `-v` over `./internal/cli`, the existing controls `TestReviewRoundTwoContinuesTheRecordedSession`, `TestReviewRoundTwoContinuesTheRecordedFallbackSelection`, `TestReviewRoundTwoFallsBackToAFreshSessionWhenPreparationFails` and `TestReviewLineageChangeEndsTheOpenSession`, recording each by name, and record from the Delivery Base with `git diff` that none of these tests, and no other existing test, golden or fixture of `internal/cli` or `internal/delivery`, changed; then run `./internal/delivery` whole.
6. MUST verify invariants 6 and 8 through the built binaries and the Delivery Base source, in a fresh home prepared as above: build a Delivery Base binary from a `git archive` of the Delivery Base, the merge base of the audited head with the default branch, in a temporary copy outside the checkout; for `review` and `review dispose`, record that `bin/roundfix <command> --help` and the Delivery Base binary's print the same bytes and exit code; and record from the Delivery Base source that the function that reads the Pre-PR Review Record decodes it without refusing unknown fields, so a Roundfix released before this Spec reads a record that carries `lineage.model`.
7. MUST verify Success Metric 5 and Core Feature 5 by running, with `-tags docscontract` and `-v` over `./internal/docscontract`, `TestTheReviewModelDecisionAmendsTheLineageRule`, `TestTheReviewerLineageEntryStatesTheModelRule`, `TestTheReviewGuideAndReferenceDescribeTheModelRule` with its subtests `guide`, `canonical_reference` and `mirror_reference`, `TestTheSkillCoverageMapIsCurrent` and `TestGlossaryLinksNameAcceptedADRs`, recording each test and subtest by name, then `TestEveryOwnedSkillVersionIsRecorded`, `TestAuthorialSkillSync` and `TestSettlementGuidanceIsOneTable` with `-v` over `./skills`; and record from the files: ADR-0451's header with `status: accepted` and `knowledge: true` and its title; ADR-0197's `amended_by` ending with ADR-0451, its `updated_at` equal to ADR-0451's `created_at`, and its dated notice after the notice of ADR-0391; the path of ADR-0451 once in the Standard Export List, in ascending order; that ADR-0451 cites no Spec folder path; and that both copies of the Roundfix skill and of its review reference are byte-identical with `cmp`.
8. MUST verify that the regenerated files are the commands' output by running, on a temporary copy of the audited head, `make skills-sync`, `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` and `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`, and recording that they change no file.
9. MUST run each sabotage below on its own temporary copy of the audited head, the three of `_techspec.md` → Testing Approach. Each sabotage keeps the packages compiling and regenerates nothing, so the named test fails on its assertion; record that it fails on the named test.
   1. The round-2 branch no longer compares the recorded model with the model of the selection at the recorded index: `TestReviewRoundTwoOpensANewSessionWhenTheConfiguredModelChanged` fails.
   2. The block that records the session of a round whose prompt was sent no longer records the model: `TestReviewRoundTwoContinuesTheRecordedSession` fails, because round 2 can no longer prove the model.
   3. The added condition of the round-2 branch is always false: `TestReviewRoundTwoContinuesTheRecordedSession` fails.
10. MUST execute the tests task_01 names in its Verification.
11. MUST perform the glossary check (`docs/agents/domain.md`): record that `CONTEXT.md` still defines **Reviewer Lineage** as a line-anchored Glossary Definition, that its first sentence ends with the clause of `_techspec.md` → API Contract 4 and its list of decision records links ADR-0451 after ADR-0391, that no other entry of `CONTEXT.md` differs from the Delivery Base, and that no term the Spec uses needs a definition the glossary lacks.
12. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`, and the Governed Paths changed are only those three;
    - in `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` only the version fields differ from the Delivery Base;
    - the Task commit range changes none of these: the Makefile, the CI workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`, `AGENTS.md`, `CLAUDE.md`, `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/roundfix/package.json`, `docs/references/skill-coverage.json`, `docs/references/coverage-record.json`, any file under `docs/agents/`, any Vendored Skill, any owned skill other than `roundfix`, and any reference of the Roundfix skill other than `review`;
    - the Go source files changed are `internal/cli/review.go` and `internal/cli/review_lineage.go`, besides the two new test files, and the function that classifies the reviewer's answer is byte-identical to the Delivery Base's;
    - no accepted ADR other than ADR-0197 changed, it stays `accepted`, and ADR-0451 is `accepted`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
13. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance. It supports the row with the pre-PR equivalent of each control:
    - approval from the recorded delivery authority, the maintainer's three decisions of 2026-11-23 quoted in `_authorization.md`;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
14. MUST record, as a limitation and not as a failed row, that the real Codex adapter is not exercised (`_techspec.md` → Risks & Considerations): the fix is confirmed on a real reviewer by the Fiscus's Delivery Retry of its Spec 0108 on the released version, which the operator asks for after the release and no row of this gate performs.
15. MUST, for the Outcome for the Archive Record, run only `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0419-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed --plan --offline </dev/null; rm -rf "$h" "$t"` with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It writes the report's Outcome as one paragraph that does not end with a colon. It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix`, or read the real home directory except for the toolchain locations above;
    - run the real `gh`, call ai-memory, the Indexus or any MCP server, or reach the network;
    - start an Agent harness, an Agent session, an Author Subagent, an ACP adapter, a real Run, a Daemon or a Delivery Owner;
    - send any message, or any request to a provider, the judge or a reviewer;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the new review tests, the controls and the Delivery Engine package, and record the outside-evidence row.
- [ ] Compare the built binaries' help with the Delivery Base binary's and read the record decoder of the Delivery Base.
- [ ] Run the documentation, glossary and skills contracts, and prove the regenerated records are the commands' output on a temporary copy.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the glossary check, the scope audit, the Pull Request row and the limitation, and write the QA report.

## Acceptance Criteria

- [ ] A round 2 after a changed model, or over a record without a model, ends the round-1 session and reaches the reviewer's verdict on a new session with the round-2 delta prompt; a round 2 with the model unchanged continues the round-1 session as before; the classifier still refuses the answer the Fiscus measured.
- [ ] Every round that sent a prompt records its model, a `ceiling-closed` lineage keeps its bytes, the help of `review` and `review dispose` keeps its bytes, and a Roundfix released before this Spec reads the new record.
- [ ] ADR-0451, its amendment of ADR-0197, the glossary entry, the guide page and both copies of the skill reference agree with the behavior; every sabotage makes its named test fail; the scope audit finds no change outside the declared files and the three granted Governed Paths.

## Verification

- `newest="$(find 'docs/specs/0419-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Non-Goals / Out of Scope; Glossary
- `_techspec.md` → Interfaces; Invariants; Data Models; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Coverage Map; Testing Approach; Risks & Considerations; Build Order 2
- `_authorization.md` → What the grant covers; Sanctioned regeneration; Limits
- [`references/sources.md`](references/sources.md), sections 1, 5, 6 and 7 (reported by fiscus, 2026-11-23)
- ADR-0197, ADR-0174, ADR-0256, ADR-0451
