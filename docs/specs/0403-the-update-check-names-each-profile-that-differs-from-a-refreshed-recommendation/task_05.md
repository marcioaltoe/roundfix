---
task: task_05
spec: 0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

Execute the `qa-gate` skill as this Spec's authored terminal gate: derive the
QA matrix from the PRD and the Task evidence, exercise the Update Check, the
profile commands and the shipped snapshot through the built binary, check the
snapshot against the outside evidence the planning change committed, and
write the dated QA report.

## Requirements

1. MUST build the binary from the worktree and run every check with a fresh temporary `HOME` and temporary fixture repositories it creates and removes, never the real `~/.roundfix`, `~/.codex` or `~/.claude`; a fresh `HOME` that runs the `./internal/cli` package keeps the host's Rust toolchain: take `RUSTUP_HOME` and `CARGO_HOME` before `HOME` changes (their values, or `$HOME/.rustup` and `$HOME/.cargo` when unset) and pass them to the command, read-only, so `TestGuidanceCompositionJourney/rust-cli` can select its toolchain.
2. MUST verify Success Metric 1: `roundfix profiles show --json` in a repository with no profile configuration reports, for every category, the approved table of `references/2026-11-17-maintainer-answers.md` with `source_as_of` `2026-11-17`, and no selection with the effort `ultra`.
3. MUST verify Success Metric 2: every recommended Codex and Claude model is offered by the Model Catalog, the Codex catalog lists `gpt-6.1-terra` after `gpt-6-astra`, the Claude catalog lists `fable` and no `claude-fable-5-1`, and no retired model is offered, through the catalog tests and the guides.
4. MUST reproduce Surface Transcript 1 of `_techspec.md` with the built binary in a fixture repository, the `...` line standing for the release and Baseline lines, and the same fixture with `--format json`, checking the `recommendations` member against API Contract 2 and that `next` and the exit code equal the run without profiles (Success Metric 3).
5. MUST verify Success Metric 4: with no profile configuration, with every category current, and with a pin to the shipped snapshot, the Update Check prints no `Profiles:` line and no `recommendations` member; and that `roundfix profiles check` and `roundfix profiles check --json` print the same bytes as the binary built from the Delivery Base for a fixture whose profiles equal the new snapshot, differing only where the snapshot itself differs.
6. MUST verify Success Metric 5, the outside-evidence row, without the network: every selection of the shipped Recommended Profile names a model that `references/2026-11-17-model-evidence.md` lists for its runtime, with an effort that record lists as accepted, and every selection runs on the `codex` or `claude` runtime.
7. MUST verify that this repository's `.roundfixrc.yml` defines every profile as the Recommended Profile and that its keys outside `profiles` are byte-identical to the Delivery Base.
8. MUST check the `--update` Skill Mode text in both copies of the Roundfix skill against `_techspec.md` → API Contract 4, including that no fix runs without the user's choice and that `update.mode: auto` only reports.
9. MUST run the glossary check of `docs/agents/domain.md` for the changed term **Update Check** and confirm ADR-0181 carries its amendment line.
10. MUST NOT run `roundfix profiles validate`, `roundfix profiles check --apply` or any command that starts an Agent session or reaches a provider; a live proof of the new tuples is the operator's release check.

## Subtasks

- [ ] Build the binary and prepare the fixtures.
- [ ] Verify the snapshot, the catalog and the outside evidence.
- [ ] Reproduce the transcript and the JSON member, and the cases that report nothing.
- [ ] Check the skill text, this repository's Project Config and the glossary.
- [ ] Write the dated QA report.

## Acceptance Criteria

- [ ] Every Success Metric of the PRD has a passing row with its evidence.
- [ ] The outside-evidence row cites `references/2026-11-17-model-evidence.md`.
- [ ] The QA report is written under this Spec's folder.

## Context

- instruction: `docs/specs/0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation/references/2026-11-17-model-evidence.md`
- instruction: `docs/specs/0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation/references/2026-11-17-maintainer-answers.md`

## Verification

- `newest="$(find 'docs/specs/0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Glossary
- `_techspec.md` → API Contracts 1 to 5; Surface Transcript 1; Interfaces (invariants 1 to 7)
- [`references/2026-11-17-model-evidence.md`](references/2026-11-17-model-evidence.md)
- ADR-0416, ADR-0181, ADR-0180, ADR-0235, ADR-0371, ADR-0256
