---
task: task_05
spec: 0204-a-run-that-reports-the-tokens-and-spend-it-used
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec reads the usage ACP adapters
report, records it in the Run Database, prints it in `runs show`, the
Implement Run summary, the Run Event Stream and `deliver status`, and adds a
token ceiling to the Delivery Queue. Every behavior row is exercised through
the built binary or by executing the named tests against the built tree, in a
disposable repository with a disposable Roundfix Home. No command reaches a
real ACP adapter, a provider, a metering gateway or the network, except the
reads of published pages the outside-evidence row names.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing task_01's tests against the built tree, that the recorded Codex fixture yields basis `request-sum` with the sum of its readings and no split, that the recorded Claude fixture yields basis `turn` with its split and reported cost, that the runtime guard counts a grown report as `turn`, and that a malformed payload never fails a prompt.
3. MUST verify, by executing task_02's tests against the built tree, that each prompt leaves one row and one `daemon.token_usage` Run Event including a failed and a stopped prompt, that a failed write leaves the Task and Run outcome unchanged, that a database at the previous schema migrates, and that every Run ID recorded on a queue item is linked and counted.
4. MUST verify through the built binary that `roundfix runs show` reproduces Surface Transcript 1 for a seeded Run, Surface Transcript 2 for an unknown Run, and schema `roundfix/runs-show/v1` with `--json`.
5. MUST verify, by executing task_03's tests against the built tree, that the Implement Run summary ends with the line of Surface Transcript 3, and through the built binary that `deliver start` prints Surface Transcript 5.
6. MUST verify through the built binary, on a seeded Delivery Queue, that `deliver status` reproduces Surface Transcript 4, and that `deliver retry` at the ceiling reproduces Surface Transcript 6.
7. MUST verify, by executing task_04's tests against the built tree, that a queued item parks at the ceiling without a worktree, that an item in stage `running` continues, that a queue without a ceiling never parks for tokens, and that `--max-tokens 0` is refused and records nothing.
8. MUST verify through the built binary that `roundfix events <run-id> --filter usage` reproduces Surface Transcript 7 for a seeded Run, and that the default stream includes the `usage` record.
9. MUST verify through the built binary, for a seeded Run whose prompts are all unreported, that `runs show`, `deliver status` and the output task_03's test captures for the Implement Run summary each say `none reported` and none prints `0 tokens`.
10. MUST record, as evidence this Spec did not author, the sources quoted in `_prd.md` → Acceptance evidence: the three Agent Client Protocol pages, the published `@agentclientprotocol/codex-acp` 2.0.1 and `@agentclientprotocol/claude-agent-acp` 0.84.0 packages, the Node.js 26.0.0 release note and the undici pages. For each page it reaches, it MUST record what it read; for each it cannot reach, it MUST record the row as blocked with the reason. When the acpx session streams recorded on the QA host are readable, it MUST recount, read-only, how many Codex prompt responses report a total equal to their last `usage_update` reading, and record the count; when they are not, it MUST record that sub-row as blocked with its reason.
11. MUST verify that the user guide states each adapter's basis, the unreported rule, the metering gateway as an operator option with its measured failure and the `allowH2: false` caveat, and the ceiling; that the four skill reference files and their mirrors carry `### Token usage`; that the `### QA settlement` section of every skill is unchanged; and that `make skills-sync-check` exits `0`.
12. MUST verify that this Spec's own artifacts satisfy the promise rule.
13. MUST verify from the repository history that each Task's changed files stay within its declarations or its `## Recorded paths`, that every Governed Path is bounded in `_authorization.md`, and that `internal/cli/cli_test.go`, `docs/references/coverage-record.json` and `.roundfixrc.yml` did not change.
14. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance, and support it with the pre-PR equivalent of each control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows archive;
    - review-artifact ancestry: the audited head named as the claimed candidate.
15. MUST NOT accept a row whose only evidence is that a file was read.
16. MUST NOT write to the live Run Database or the metering gateway files under `~/.roundfix`, MUST NOT start a session on a real ACP adapter, and MUST NOT start a metering gateway.

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

- `newest="$(find 'docs/specs/0204-a-run-that-reports-the-tokens-and-spend-it-used/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; User Stories 1-5; Core Features 1-11; Success Metrics
1-6; Acceptance evidence; `_techspec.md` → Surface Transcripts 1-7; Testing
Approach 1-6; API Contracts 1-7; ADR-0008; ADR-0017; ADR-0020; ADR-0033;
ADR-0051; ADR-0080; ADR-0088; ADR-0091; ADR-0093; ADR-0094; ADR-0096;
ADR-0098; ADR-0104; ADR-0117; ADR-0155; ADR-0156; ADR-0158; ADR-0164;
ADR-0166; ADR-0167; ADR-0176; ADR-0182; ADR-0184; ADR-0198; ADR-0199.
