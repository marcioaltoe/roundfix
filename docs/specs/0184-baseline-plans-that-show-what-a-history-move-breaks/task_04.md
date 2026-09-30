---
task: task_04
spec: 0184-baseline-plans-that-show-what-a-history-move-breaks
status: failed
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

This Task is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec changes Baseline
planning in `internal/baseline` and the documents that describe it.

- **Where commands run.** Every binary command runs in a disposable repository
  or a disposable clone, with the built `roundfix`. No command reaches a
  provider or the network.
- **External repository.** The one external repository, Fluxus at `~/dev/fluxus`, is read
  only through `git archive` and `git show`, and never written.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify through the built binary that a planned History Relocation's
   citations are reported. It uses a disposable Git repository seeded like the
   package's plan fixture:
   - `.agents/skills/context7/SKILL.md`;
   - `.agents/skills/exa-web-search/SKILL.md`;
   - a `Makefile` with `verify` and `verify-incremental` targets;
   - an ADR at `docs/adr/` whose lifecycle status is `superseded`;
   - a tracked file citing that ADR by repository path;
   - a tracked Markdown file linking it relatively;
   - an untracked file citing it.

   The plan is run as `roundfix baseline plan --repo <repo> --profile go-cli-tui` with
   these decisions:
   - `language.generated=English`
   - `verification.gate=make verify`
   - `verification.incremental=make verify-incremental`
   - `branch.prefix=ma/`
   - `spec.scaffold=true`
   - `domain.layout=single-context`
   - `triage.external=false`
   - `autonomous.enabled=true`
   - `runtime.backend=codex gpt-5.5 xhigh`
   - `runtime.design=claude opus xhigh`
   - `secondbrain.enabled=false`
   - `repository.extension.enabled=true`
   - `preservation.mode=greenfield`

   It MUST check both output formats:
   - `--format=text` prints one `Warning: baseline.history.citation:` line for
     each tracked citing file, and none for the untracked file;
   - `--format=json` carries the same entries in `warnings`.

   The same repository without the superseded ADR MUST print no
   `baseline.history.citation` warning.
3. MUST verify, by executing the scanner, plan-level and CLI tests against the
   built tree, the rest of the contract:
   - co-relocated links and already-broken citations stay unreported;
   - symbolic links are not followed;
   - binary files are skipped and oversized files are summarized;
   - the caps hold;
   - the Plan Digest binds the warnings;
   - apply is unchanged;
   - a digest confirmed before a citing file changed is refused.
4. MUST replay Fluxus commit `986f9dc` as evidence this Spec did not author:
   - export `docs/` at `986f9dc` from `~/dev/fluxus` into a disposable
     directory;
   - for every rename that `git show --name-status -M 986f9dc` reports, move
     the destination back to its source path;
   - add the fixture files from Requirement 2, initialize and commit;
   - run the built `roundfix baseline plan` as in Requirement 2 with
     `--format=json`.

   The gate MUST record that each of the sixteen relocated ADRs yields exactly
   one `baseline.history.citation` warning, citing the ADR that consolidated
   it. It MUST compare every reported citing file with a recount over the
   plan's `historyMoves` that is independent of Roundfix. The gate itself runs
   that recount in a script it writes outside the repository, reporting links
   and repository paths that resolve before the moves and not after. The gate
   MUST record the sixteen Fluxus Spec 0067 measured, from the Secondbrain
   mirror
   `projects/fluxus/mirror/docs/history/specs/0067-os-links-de-documentacao-resolvem/_prd.md`.
   When the Fluxus repository, that commit or the mirror is unavailable, or the
   plan cannot be produced, the gate MUST record the row as blocked with its
   reason.
5. MUST measure the cost on this repository:
   - make a disposable clone of this repository and commit one added legacy
     relocation to it, `docs/specs/_archived/9999-probe/_prd.md`;
   - time `roundfix baseline update --repo <clone> --format json` with the
     built binary, and with a v0.20.0 binary built from tag `v0.20.0` in a
     second disposable clone with `go build -buildvcs=false`;
   - record both times, three runs each.

   The added time MUST be at most two seconds. The gate MUST also record that a
   clone with no added relocation gives the same result from both binaries.
   When the tag or the build is unavailable, the gate MUST record the row as
   blocked with its reason.
6. MUST verify that both Roundfix skill copies, the context-driven development
   guide and `CONTEXT.md` describe the Relocation Citation warnings, their
   three codes and the digest binding, and that `make skills-sync-check` exits
   `0`.
7. MUST verify that this Spec's own artifacts satisfy the promise rule.
8. MUST verify from the repository history that the changed files stay within
   the paths the Tasks declare, and that every governed path is bounded in
   `_authorization.md`.
9. MUST NOT accept a row whose only evidence is that a file was read.
10. MUST NOT write to `~/dev/fluxus`, to the Secondbrain, or to the live Run
    Database under `~/.roundfix`.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no
      declaration waives.
- [ ] Execute each row against the built tree and record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the
      sources that row covers.
- [ ] The repository Verification result is recorded as a fact, not re-derived.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0184-baseline-plans-that-show-what-a-history-move-breaks/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals 1-4; User Stories 1-3; Core Features 1-6; Acceptance evidence
- `_prd.md` → Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4
- `_techspec.md` → Testing Approach 1-4; API Contract 1; API Contract 2
- ADR-0064; ADR-0068; ADR-0070; ADR-0071; ADR-0073; ADR-0080; ADR-0091;
  ADR-0093; ADR-0094; ADR-0096; ADR-0103; ADR-0104; ADR-0117; ADR-0120;
  ADR-0155; ADR-0156; ADR-0173
