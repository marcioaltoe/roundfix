---
task: task_05
spec: 0217-a-cursor-runtime-to-measure-grok-on
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec adds `cursor` as an opt-in ACP
Runtime in `internal/config`, `internal/agent` and `internal/cli`, a login
check, guides and a Roundfix Skill section, and a live measurement. Behavior
rows run the built binary and the named tests against fakes; no row reaches
Cursor, a provider or the network, and the only live evidence is the
measurement record Task 04 committed.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify through the built binary, in a temporary repository and Home
   with a fake `cursor-agent` that is a compiled binary:
   - `implement`, `resolve` and `watch` help list
     `codex, claude, cursor, opencode`, as Surface Transcript 3 shows for
     `implement`, and an unknown runtime is refused with Surface
     Transcript 2;
   - a Project Config naming `cursor` with `grok-4-20[thinking=true]` and an
     empty effort loads, and one with `reasoning_effort: high` is refused
     with the model-value rule;
   - `roundfix doctor` with a fake whose `status` prints `Not logged in`
     prints Surface Transcript 1's `adapter:` line and exits `1`, and with a
     fake that reports a login the Cursor segment passes.
3. MUST verify, by executing `TestCursorCapabilitiesReadWholeModelValues` and
   `TestOtherRuntimesKeepBracketParsing`, that only `cursor` reads bracketed
   values whole, and record that the published fixture still yields
   `malformed_model_value` under `claude`.
4. MUST verify that no production Go file names `CURSOR_API_KEY`,
   `CURSOR_AUTH_TOKEN`, `--api-key`, `--auth-token` or a `login` or
   `logout` argument, and that no Task ran `cursor-agent login`.
5. MUST read the measurement record and the recorded session that Task 04
   committed, run `TestCursorRecordedSessionAdvertisesGrok`, and confirm that
   the record names the advertised Grok values and one row per replay, and
   that no account field is in the recording. When Task 04 ended with
   `cursor_login_required`, the row MUST be recorded as
   `blocked (environment: cursor-agent not logged in)`.
6. MUST record, as evidence this Spec did not author:
   - Cursor's ACP documentation (<https://cursor.com/docs/cli/acp>), for the
     `initialize`, `authenticate`, `session/new` order, the
     `cursor_login` method and the blocking extension methods;
   - the Cursor forum thread on ACP model switching
     (<https://forum.cursor.com/t/bug-agent-acp-model-switching-updates-session-metadata-but-does-not-change-the-inference-backend/157312>),
     for the published model values and the 2026-04-30 fix;
   - the installed acpx agent registry, for `cursor` → `cursor-agent acp`;
   - the live measurement record, as a runtime this Spec did not build.

   When a source is unavailable, it MUST record that row as blocked with its
   reason.
7. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   check whether it introduced, changed or retired a glossary term; record
   whether `CONTEXT.md` needs `Cursor runtime` or the model-value rule.
8. MUST verify from the repository history that each Task's changed files
   stay within its declarations or its `## Recorded paths`, that the only
   Governed Paths changed are the three `_authorization.md` bounds, and that
   the Recommended Profile, the built-in profiles and `.roundfixrc.yml` are
   unchanged. This row reads Task commits, so it declares a `commit_range`
   input.
9. MUST record the non-waivable Pull Request row as
   `blocked (environment: no open Pull Request)`, naming the Pull Request row
   in its provenance. It MUST support that row with the pre-PR equivalent of
   each control:
   - approval: the maintainer's recorded delivery authority;
   - checks and status: the Daemon's repository Verification at the audited
     head;
   - unresolved threads: none, because no Pull Request exists;
   - Merge-Ready acceptance: none yet, because the pre-PR review follows
     archive;
   - review-artifact ancestry: the audited head named as the claimed
     candidate.
10. MUST NOT accept a row whose only evidence is that a file was read.
11. MUST NOT write to the live Run Database under `~/.roundfix`, run a real
    `cursor-agent`, or send a prompt anywhere.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no
      declaration waives.
- [ ] Execute each row against the built tree and record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the
      sources that row covers.
- [ ] The repository Verification result is recorded as a fact, not re-derived.
- [ ] The measurement row is either satisfied by the committed record or
      blocked by the named login blocker.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0217-a-cursor-runtime-to-measure-grok-on/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; User Stories 1-4; Core Features 1-5; Success Metrics 1-5; Acceptance
evidence; `_techspec.md` → Surface Transcripts 1-3; API Contracts 1-3; Testing Approach 1-5;
The measurement protocol; Build Order 5; ADR-0080; ADR-0088; ADR-0091;
ADR-0104; ADR-0155; ADR-0156; ADR-0167; ADR-0217.
