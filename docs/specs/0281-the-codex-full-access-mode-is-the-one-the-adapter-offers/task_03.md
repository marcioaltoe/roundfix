---
task: task_03
spec: 0281-the-codex-full-access-mode-is-the-one-the-adapter-offers
status: completed
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec delivers two things:

- a full-access request that selects the mode the Agent Session advertises,
  sends the Codex sandbox preset when the session offers it, and reports a
  degraded effective access policy when no full-access mode is advertised;
- the documentation of that selection, the `.roundfixrc.yml` comment and the
  `qa-gate` guidance on the browser driver inside a sandbox.

Every row runs on temporary copies, temporary homes and the fake acpx runner.
No row reads or writes the real `~/.roundfix` or `~/.acpx`, starts acpx, an
adapter or a browser, reaches the network, sends a request to a provider, the
judge or a reviewer, reads a credential, or reads another repository. The gate
never runs `roundfix spec judge`, never runs `roundfix review` against a real
provider, never runs `roundfix doctor` or `roundfix profiles validate` against
the real adapters, and never runs an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1:
   - Run `TestChooseFullAccessModeTakesTheFirstAdvertisedCandidate` and
     `TestProofSelectsTheAdvertisedCodexFullAccessMode`, and record from their
     output that the codex-acp 2.x shape received one `set-mode
     agent-full-access`, no `set sandbox_mode`, and the effective policy
     `full-access`.
   - Make a temporary copy of the audited head with `git archive`. In it, make
     `chooseFullAccessMode` return the first candidate whatever the session
     advertises, and record that
     `TestProofSelectsTheAdvertisedCodexFullAccessMode` fails (Invariant 1).
3. MUST verify Success Metric 2:
   - Run `TestRunAppliesTheSandboxPresetWhenTheSessionAdvertisesIt` and
     `TestRunSelectsBypassPermissionsForClaude`, and record the invocation
     order `set-mode full-access`, `set sandbox_mode danger-full-access` and
     the Claude `set-mode bypassPermissions`.
   - On a temporary copy of the audited head, make `applyFullAccess` skip the
     sandbox preset on every Codex session, and record that
     `TestRunAppliesTheSandboxPresetWhenTheSessionAdvertisesIt` fails.
   - Run the five pinned tests named by task_01 Requirement 7 and record that
     they pass under their unchanged names (Invariant 4).
4. MUST verify Success Metric 3:
   - Run `TestFullAccessWithoutAnAdvertisedModeIsDegraded` and record the
     exact degraded text, the status `full_access_mode_not_advertised`, the
     warning, and that the Run reached its prompt.
   - On a temporary copy of the audited head, make the unadvertised case
     return the policy `full-access`, and record that
     `TestFullAccessWithoutAnAdvertisedModeIsDegraded` fails (Invariant 2).
   - Run `TestProfilesValidateTextNamesADegradedPolicy` and
     `TestDoctorNamesADegradedPolicy` in `./internal/cli`, and record that the
     existing degraded format still prints a degraded policy.
5. MUST execute the commands task_01 and task_02 name in their Verification,
   and record that a default-access Run sends no `set-mode` (Invariant 3)
   from `TestACPXRunAppliesFullAccessSessionSetup`'s OpenCode case.
6. MUST record, as evidence this Spec did not author, each published source
   of `_prd.md` → Acceptance evidence. Because this Spec's authorization
   forbids network access, record each web source as `blocked (environment:
   network prohibited by Spec authorization)`, quoting the passage the PRD
   cites. Such a row never decides a qualifying `partial`.
7. MUST perform the glossary check of `docs/agents/domain.md`:
   - record that `CONTEXT.md` defines **Full-Access Mode** with ADR-0293's
     link;
   - record that the glossary contract tests pass;
   - record whether any other term the Spec introduced needs an entry.
8. MUST verify from the repository history, comparing with the Delivery Base
   (the merge base of the audited head with `main`):
   - each Task's changed files stay within its declarations or its
     `## Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, every Baseline module under
     `internal/baseline/assets/modules/`, `docs/agents/setup-context.json`
     and every file under `internal/judge/`;
   - the values of `PinnedCodexAdapterVersion` and
     `PinnedClaudeAdapterVersion` in `internal/agent/acpx_runner.go` are
     unchanged;
   - in `.roundfixrc.yml` only comment lines changed, and the line
     `agent_full_access: false` is unchanged;
   - the Roundfix and `qa-gate` skill versions were raised and recorded, and
     no `### QA settlement` section changed;
   - `grep -rnF '1.1.9'` finds codex-acp 1.1.9 nowhere under `docs/user-guide`,
     `docs/agents`, `docs/adr`, `.agents/skills`, `skills` or `.roundfixrc.yml`.

   This row reads Task commits, so it declares a `commit_range` input.
9. MUST record the non-waivable Pull Request row as
   `blocked (environment: no open Pull Request)`, naming the Pull Request row
   in its provenance. It supports the row with the pre-PR equivalent of each
   control:
   - approval from the recorded delivery authority;
   - checks from the Daemon's repository Verification at the audited head;
   - no unresolved threads, because no Pull Request exists;
   - no Merge-Ready acceptance yet;
   - the audited head named as the claimed candidate.
10. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0281-the-codex-full-access-mode-is-the-one-the-adapter-offers --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It MUST NOT
    do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or `~/.acpx`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - start a workflow.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] An advertised full-access mode is selected, the sandbox preset follows the session's controls, and each removed guard makes its named test fail.
- [ ] A session without an advertised full-access mode is reported degraded with its modes, and no unadvertised mode is sent.
- [ ] The documentation describes the selection, `.roundfixrc.yml` keeps `agent_full_access: false`, and the scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0281-the-codex-full-access-mode-is-the-one-the-adapter-offers/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Acceptance evidence; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; Invariants 1-4; Current behavior; Build Order 3
- ADR-0293; ADR-0011; ADR-0104
