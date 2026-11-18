---
task: task_05
spec: 0405-authoring-model-fallback
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

Execute the `qa-gate` skill as this Spec's authored terminal gate. It derives
the QA matrix from the PRD and the Task evidence, exercises the key and the
two subcommands through the built binary, checks the limit reader against the
outside evidence the planning change committed, and writes the dated QA
report.

## Requirements

1. MUST build the binary from the worktree and run every check with a fresh temporary `HOME` and temporary fixture repositories it creates and removes. MUST NEVER use the real `~/.roundfix`, `~/.codex` or `~/.claude`. A fresh `HOME` that runs the whole `./internal/cli` package MUST keep the host Rust toolchain: `PATH`, `RUSTUP_HOME` and `CARGO_HOME` pointing at the host's toolchain, read-only.
2. MUST reproduce Surface Transcript 1, Surface Transcript 2, Surface Transcript 3, Surface Transcript 4, Surface Transcript 5, Surface Transcript 6, Surface Transcript 7, Surface Transcript 8 of `_techspec.md` with the built binary, in the fixture the Surface Transcripts section describes, with `CLAUDE_CODE_SESSION_ID` set. Transcripts 3 to 5 run after transcript 2 in the same home.
3. MUST verify Success Metric 2: after transcript 2 the home holds exactly one fallback entry with the keys of `_techspec.md` → Data Models, in order. Each refusal of API Contract 1 MUST leave the day files byte-identical.
4. MUST verify Success Metric 3, through transcripts 4 and 5 and through `roundfix phase model idea --spec 0413-another`.
5. MUST verify Success Metric 4:
   - in a fixture whose Project Config declares no fallback, `roundfix skills install --target project` writes the same skills and subagents as the binary built from the Delivery Base;
   - `roundfix phase start` prints the same lines and exit codes as that binary.
6. MUST verify Success Metric 1 through the four authoring skills' text: each, in both copies, runs `roundfix phase model` before starting its Author Subagent and `roundfix phase fallback` on an error, for its own phase, and restarts on the printed model and effort without waiting.
7. MUST verify the outside-evidence row without the network. Each Claude Code usage-limit message quoted in `references/2026-11-17-harness-and-test-measurements.md` → Usage limit messages is classified as that documentation says, through `roundfix phase fallback` with the built binary:
   - the session and weekly limits as shared;
   - the Opus and Sonnet limits as a model's own;
   - the spend-limit lines as not a model's limit;
   - the observed `You've hit your Fable limit · resets Nov 22 at 8pm` as the Fable limit.
8. MUST verify that a Project Config with a fallback on `research`, a fallback without a model, or a fallback that names its entry's model is refused with the messages of `_techspec.md` → Data Models.
9. MUST run the glossary check of `docs/agents/domain.md` for **Authoring Fallback**, **Authoring Model Declaration**, **Author Subagent** and **Phase Record**. MUST confirm that ADR-0401 carries `Amended by ADR-0427`, ADR-0407 carries `Amended by ADR-0428` and ADR-0415 carries `Amended by ADR-0430`.
10. MUST verify Success Metric 5. Reproduce Surface Transcript 9 with the built binary in a fixture repository where `roundfix skills install --target project` installed the skills and `.roundfixrc.yml` declares no `skills`. Check that `roundfix doctor` prints the same block and that the configuration guide's block equals it.
11. MUST verify that this repository's `.roundfixrc.yml` is byte-identical to the Delivery Base: its fallbacks are the operator's step after the release.
12. MUST verify Success Metric 6 with the built binary, in a fresh temporary Git repository with a fresh temporary `HOME`:
    - adopt it with `roundfix baseline plan --profile go-cli-tui` and every required decision, then `roundfix baseline apply` with the printed Plan Digest; the profile's required skills may be copied into the fixture from this worktree first, and the fixture's `CLAUDE.md` is a symbolic link to `AGENTS.md`, as the measured adopters keep it;
    - confirm that `AGENTS.md`, and `CLAUDE.md` read through its link, each hold the first sentence of the bullet of `_techspec.md` → API Contract 10, which opens with "In Claude Code,", and its last sentence, "This bullet does not apply to Codex or OpenCode, whose models come from the Agent Selection Profiles.", each exactly once, normalized to single spaces, inside the `root.spec-workflow` markers, and that `CLAUDE.md` is still a link;
    - confirm that a second `roundfix baseline update --no-skills --no-history` on the fixture reports `current`;
    - confirm that this repository's `AGENTS.md`, through `CLAUDE.md`, holds the same two sentences once, with the scoping words;
    - confirm that no file outside the fixture and the temporary `HOME` changed, and that no user-level instruction file was read or written.
13. MUST NOT start an Agent session, an Author Subagent or a harness, and MUST NOT reach a provider or the network.

## Subtasks

- [ ] Build the binary and prepare the fixtures.
- [ ] Reproduce the eight transcripts and check the record.
- [ ] Check the outside-evidence row, the refusals and the case without a fallback.
- [ ] Check the skill text, the glossary, the ADR notices and this repository's Project Config.
- [ ] Adopt a fresh fixture and check the root statement in both root files.
- [ ] Write the dated QA report.

## Acceptance Criteria

- [ ] Every Success Metric of the PRD has a passing row with its evidence.
- [ ] The outside-evidence row cites `references/2026-11-17-harness-and-test-measurements.md` and the Claude Code documentation pages it records.
- [ ] The QA report is written under this Spec's folder.

## Context

- instruction: `docs/specs/0405-authoring-model-fallback/references/2026-11-17-harness-and-test-measurements.md`
- instruction: `docs/specs/0405-authoring-model-fallback/references/2026-11-17-maintainer-answers.md`
- instruction: `docs/specs/0405-authoring-model-fallback/references/2026-11-18-root-statement-decision-and-measurement.md`

## Verification

- `newest="$(find 'docs/specs/0405-authoring-model-fallback/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; User Stories 1 to 8; Glossary
- `_techspec.md` → API Contracts 1 to 8; Surface Transcript 1, Surface Transcript 2, Surface Transcript 3, Surface Transcript 4, Surface Transcript 5, Surface Transcript 6, Surface Transcript 7, Surface Transcript 8, Surface Transcript 9; API Contract 9; API Contract 10; Interfaces (invariants 1 to 8); Data Models
- [`references/2026-11-17-harness-and-test-measurements.md`](references/2026-11-17-harness-and-test-measurements.md)
- [`references/2026-11-18-root-statement-decision-and-measurement.md`](references/2026-11-18-root-statement-decision-and-measurement.md)
- ADR-0427, ADR-0428, ADR-0429, ADR-0430, ADR-0431, ADR-0415, ADR-0401, ADR-0407
