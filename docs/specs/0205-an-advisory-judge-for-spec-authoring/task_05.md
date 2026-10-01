---
task: task_05
spec: 0205-an-advisory-judge-for-spec-authoring
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec adds `roundfix spec judge`, an
advisory command that asks TypeSafe's Jev, through OpenRouter's System One
API or else TypeSafe directly, two measured questions about a Spec's
citations and Coverage Map, logs every call, stops at a monthly ceiling and
fails open, and has the `write-prd` and `write-techspec` skills
run it. Every behavior row is exercised through the built binary or by
executing the named tests against the built tree, in a disposable repository
with a disposable Roundfix Home. No command reaches TypeSafe or any other
network service, OpenRouter and TypeSafe included, except the reads of
published pages the outside-evidence row names.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify that `internal/judge/questions.json` is byte-identical to the block under `_techspec.md` → Questions and thresholds, that it pins `jev-1.13`, that its `accepted_model_pattern` admits exactly the shapes of Invariant 1, that it lists `openrouter` (`ROUNDFIX_OPENROUTER_API_KEY`, requesting `jev-1.13`) before `typesafe` (`ROUNDFIX_TYPESAFE_API_KEY`, requesting `jev-1.13.0`), and that its two thresholds equal the operating points of the adopted finding (`docs/history/findings/2026-09-30-jev-judgments-measured-against-the-spec-archive.md`): confidence 0.8 for citation support and probability 0.3 for goal to mechanism.
3. MUST verify, by executing task_01's tests against the built tree, every rule of Citation claims, Goal pairs, The language gate and The only readers, including the refusal of a symbolic-link PRD.
4. MUST verify, by executing task_02's tests against the built tree, the transport selection of Invariant 2, including that the generic `OPENROUTER_API_KEY` alone sends nothing, that each transport sends its own model ID, the model pin of Invariant 1 on both reported shapes and on a mismatching one, the reported and computed cost of Invariant 3, each threshold on both sides of its boundary, that an answer from another model is never raised or cleared, each stop and skip reason of Asking, that no request follows a stop, a missing Jev key, an unreadable Judge Log or a reached ceiling, one Judge Log line per request, and `TestRequestsCarryOnlySpecArtifactText`.
5. MUST verify, by executing task_03's tests against the built tree, Surface Transcripts 1 and 5 and the JSON document of API Contract 2.
6. MUST verify through the built binary, in a disposable repository holding an English fixture Spec and a Portuguese one, with `HOME` set to a disposable directory, `HTTPS_PROXY` and `HTTP_PROXY` pointed at a closed local port, and `ROUNDFIX_OPENROUTER_API_KEY` and `ROUNDFIX_TYPESAFE_API_KEY` removed, that `roundfix spec judge` reproduces Surface Transcript 2 with only a dummy `OPENROUTER_API_KEY` set; with a dummy `ROUNDFIX_TYPESAFE_API_KEY` and a seeded Judge Log of US$5.0003 for the current month, Surface Transcript 3; with a dummy `ROUNDFIX_OPENROUTER_API_KEY`, Surface Transcript 4 for the Portuguese Spec; and Surface Transcript 6 for an unknown Spec. After each run it MUST record that the disposable Judge Log gained no line.
7. MUST verify through the built binary that `roundfix spec --help` and `roundfix --help` name the command as API Contract 5 states, and that `roundfix spec judge --help` names `ROUNDFIX_OPENROUTER_API_KEY`, `ROUNDFIX_TYPESAFE_API_KEY`, the monthly ceiling and the Judge Log.
8. MUST verify, through a read-only probe that writes nothing to the repository (for example a `go test -overlay` test), that the language gate of the built tree classifies as English every `_prd.md` and `_techspec.md` under `docs/history/specs/` and every accepted ADR under `docs/adr/`, and record the lowest English share it saw.
9. MUST verify that no non-test Go file outside `internal/judge`, `internal/cli/spec_judge.go` and `internal/cli/cli.go` imports `net/http`, so the command is Roundfix's only network client, and that `internal/speccheck` does not import `internal/judge`.
10. MUST record, as evidence this Spec did not author, the sources quoted in `_prd.md` → Acceptance evidence: OpenRouter's TypeSafe SDK guide, Jev hub, Jev tutorial, Jev 1.13 model page and errors reference, the TypeSafe API reference, models page, confidence page and citation-check cookbook, the Agent Skills specification and the Claude Code subagent reference. For each page it reaches, it MUST record what it read, and whether the endpoints, the request and response fields (including OpenRouter's `id`, `provider` and `usage.cost`), the model ID mapping, the price and the version-pinning advice still match `_techspec.md`; for each it cannot reach, it MUST record the row as blocked with the reason. It MUST record the measurement figures of the adopted finding that this Spec relies on, and the 2026-10-01 endpoint probes of `_techspec.md` → Measured outside evidence as evidence this Spec did not produce.
11. MUST verify that `docs/user-guide/commands/spec.md` carries `### spec judge`, that `.agents/skills/roundfix/references/spec.md` carries `### Advisory judge`, that `write-prd` and `write-techspec` carry `## Advisory judgment` with the command of their own stage, that each mirror equals its canonical copy, that the `### QA settlement` section of every skill is unchanged, and that `make skills-sync-check` exits `0`.
12. MUST verify that this Spec's own artifacts satisfy the promise rule.
13. MUST verify from the repository history that each Task's changed files stay within its declarations or its `## Recorded paths`, that every Governed Path is bounded in `_authorization.md`, and that `internal/cli/cli_test.go`, `.roundfixrc.yml` and `go.mod` did not change.
14. MUST perform the glossary check of `docs/agents/domain.md` for the terms this Spec introduces, Judge Log and advisory judgment, and record its outcome.
15. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance, and support it with the pre-PR equivalent of each control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows archive;
    - review-artifact ancestry: the audited head named as the claimed candidate.
16. MUST NOT accept a row whose only evidence is that a file was read.
17. MUST NOT send any request to OpenRouter or TypeSafe, read the real `ROUNDFIX_OPENROUTER_API_KEY`, `ROUNDFIX_TYPESAFE_API_KEY` or `OPENROUTER_API_KEY`, or write under the real `~/.roundfix`.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no declaration waives.
- [ ] Execute each row against the built tree and record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the sources that row covers.
- [ ] The repository Verification result is recorded as a fact, not re-derived.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0205-an-advisory-judge-for-spec-authoring/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-5; User Stories 1-5; Core Features 1-11; Success Metrics
1-6; Acceptance evidence; `_techspec.md` → Surface Transcripts 1-6; Testing
Approach 1-6; API Contracts 1-5; ADR-0035; ADR-0080; ADR-0081; ADR-0088;
ADR-0089; ADR-0091; ADR-0093; ADR-0094; ADR-0096; ADR-0104; ADR-0116;
ADR-0117; ADR-0149; ADR-0155; ADR-0156; ADR-0166; ADR-0167; ADR-0176;
ADR-0178; ADR-0182; ADR-0183; ADR-0184; ADR-0187; ADR-0189; ADR-0193;
ADR-0200; ADR-0201.
