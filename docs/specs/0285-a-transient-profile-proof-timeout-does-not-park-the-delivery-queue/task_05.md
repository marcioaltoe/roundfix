---
task: task_05
spec: 0285-a-transient-profile-proof-timeout-does-not-park-the-delivery-queue
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec delivers three changes: the
profile proof's setup deadline is the config key
`defaults.profile_proof_setup_timeout`, a timed-out proof gets a third
attempt after waits of 5 and 10 seconds, and the cleanup of a disposable
Agent Session that acpx reports absent writes no warning.

Every row runs on temporary repositories and homes with the CLI's injected
runner and sleeper, the compiled fake acpx of `internal/agent`, or the built
binary on a command that starts no agent. No row reads or writes the real
`~/.roundfix`, starts a real ACP agent, adapter or acpx, reaches the network,
sends a request to a provider, the judge or a reviewer, reads a credential,
or reads another repository. The gate never runs `roundfix spec judge`, never
runs `roundfix review` against a real provider, and never runs an archive
plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1:
   - Run `TestProfileProofPassesOnTheThirdAttempt`,
     `TestProfileProofRetriesATimedOutSetupOnce` and
     `TestOperationalPreflightPassesAfterAFallbackProofTimesOutOnce`, and
     record the attempt counts and the recorded waits (Invariant 3).
   - On a temporary copy of the audited head made with `git archive`, reduce
     the attempts to two and record that
     `TestProfileProofPassesOnTheThirdAttempt` fails; on another, skip the
     sleeper and record that it fails.
3. MUST verify Success Metric 2:
   - Run `TestProfileProofReportsAThirdTimeoutAsTemporary` and
     `TestProfileProofStopsWhenCancelledDuringTheBackoff`, and record the
     classification, the next action and the exact standard error the test
     compares with API Contract 2 of `_techspec.md` (Invariant 4).
   - On a temporary copy of the audited head, make a cancellation during the
     wait start another attempt and record that
     `TestProfileProofStopsWhenCancelledDuringTheBackoff` fails.
4. MUST verify Success Metric 3:
   - Run `TestProfileProofSetupTimeoutConfig`,
     `TestProfileProofPassesTheConfiguredSetupDeadline`,
     `TestExactSelectionProofSetupUsesTheRequestedDeadline` and
     `TestProfilesShowRefusesAnInvalidSetupDeadlineTranscript` (Invariant 2).
   - On a temporary copy of the audited head, make `ProveExactSelection`
     ignore `SetupTimeout` and record that
     `TestExactSelectionProofSetupUsesTheRequestedDeadline` fails; on
     another, stop `proveProfileSelectionsWithOptions` from setting it and
     record that `TestProfileProofPassesTheConfiguredSetupDeadline` fails.
   - Through the built product, reproduce Surface Transcript 1 of
     `_techspec.md`: in a temporary Git repository whose `.roundfixrc.yml`
     holds `defaults:` with `profile_proof_setup_timeout: 0s`, with `HOME`
     set to a temporary directory, run `bin/roundfix profiles show` and
     record that standard output, standard error and the exit code match the
     transcript. Repeat with `11m` and record the same message and exit `2`;
     repeat with `90s` and record exit `0`.
5. MUST verify Success Metric 4:
   - Run `TestTimedOutProofSetupClosesAnAbsentSessionSilently`,
     `TestMissingSessionIsRecognisedFromBothExitShapes` and
     `TestDisposableSessionCloseIsAppendedWhenAnOpenSessionWillNotClose`
     (Invariant 5).
   - On a temporary copy of the audited head, restore a warning for an
     absent Session and record that
     `TestTimedOutProofSetupClosesAnAbsentSessionSilently` fails.
6. MUST execute the tests task_01 to task_04 name in their Verification.
7. MUST record that every configured tuple is still proven (Invariant 1):
   `TestProfileProofRetriesATimedOutSetupOnce` records one proof of the
   fallback tuple beside the retried preferred tuple.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - record from the adopted Inbox Entry the five parks, their hours and that
     each passed on the next `roundfix deliver retry`;
   - because this Spec's authorization forbids network access, record each
     web source as `blocked (environment: network prohibited by Spec
     authorization)`, quoting the passage the PRD cites. Such a row never
     decides a qualifying `partial`.
9. MUST perform the glossary check of `docs/agents/domain.md`:
   - record that `CONTEXT.md` defines **Proof Setup Deadline**, linking
     ADR-0298;
   - record that the glossary contract tests pass;
   - record whether any other term the Spec introduced needs an entry.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every Governed Path changed is bounded in `_authorization.md`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
      under `internal/baseline/assets/`, `docs/agents/setup-context.json` and
      `internal/agent/fallback.go`;
    - the Roundfix skill version was raised and recorded, and the coverage
      record was re-recorded with its command;
    - no `### QA settlement` section changed;
    - no guide, agent guide or skill copy still says a timed-out proof "is
      retried once", and no source holds `timed out twice` or `found no
      Session to close`.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0285-a-transient-profile-proof-timeout-does-not-park-the-delivery-queue --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It MUST NOT
    do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - start an ACP agent, an adapter or acpx.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Reproduce Surface Transcript 1 through the built binary.
- [ ] Record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A setup that times out once or twice passes after the recorded waits, and each removed guard makes its named test fail.
- [ ] A third timeout is `temporary` with API Contract 2's exact text, and the built `roundfix profiles show` matches Surface Transcript 1.
- [ ] An absent Session at cleanup writes nothing, and the scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0285-a-transient-profile-proof-timeout-does-not-park-the-delivery-queue/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; Surface Transcript 1; Invariants 1-5; Current behavior; Build Order 5
- ADR-0298; ADR-0050; ADR-0104; ADR-0240; ADR-0276
