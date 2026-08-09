# A Spec that validates itself — Technical Spec

## Executive Summary

The checker exists and works. `internal/speccheck` carries nineteen rules and
answers for one Spec in 0.04 seconds. Two things are missing: a rule that reads
a cited decision against the claim made about it, and a way to ask the checker
about one authoring stage rather than the whole corpus.

Both land in `internal/speccheck`, which is ordinary source. Wiring the result
into the authoring skills and slimming the QA gate are protected-tooling edits,
covered by the standing grant for loop performance; the `write-tasks` wiring
shipped ahead of this Spec under an earlier narrow grant.

## Project Constraints

- Identifier strategy: not applicable — checks key off the existing Spec slug
  and artifact paths. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — file reads only.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0096 places machine facts before the
  Agent turn, and this Spec carries that to its conclusion. ADR-0091 keeps the
  authored QA gate a Task node of its own type. ADR-0104 requires an
  outside-evidence acceptance row. This Spec adds ADR-0116 and ADR-0117.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — steps 5 and 6 edit authored skills under the
  standing grant at
  `docs/workflow/authorizations/2026-08-09-standing-tooling-authority-for-loop-performance.md`,
  Bounded files: `.agents/skills/write-prd/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`,
  and their generated copies `skills/write-prd/SKILL.md`,
  `skills/write-techspec/SKILL.md` and `skills/qa-gate/SKILL.md`, which `make skills-sync` rewrites; ADR-0081 draws authorization around the cause rather than its computable effects, so they follow the authorized skill edit.
  Source: `docs/agents/agent-instructions.md`.

## Vocabulary Contract

- emits: `internal/speccheck/citations.go`
  pattern: `SC-CITATION-[A-Z]+`
  documented-in: `CONTEXT.md`

## System Architecture

No new package. `internal/speccheck` gains one detector and one scope parameter;
`internal/cli` gains one flag.

The semantic check joins `internal/speccheck/citations.go`, which already holds
the listing checks `SC-ADR-UNLISTED` and `SC-ADR-RELATED`. Keeping them in one
file is deliberate: the three answer the same question at increasing strength —
was the record listed, was it accounted for, does it say what the artifact
claims.

**The citation detector.** An artifact sentence that attributes a subject to a
decision record — "ADR-0083 makes `make verify` the authoritative gate" — is
resolved to that record and matched against its text. A claim the record does
not support is `SC-CITATION-UNSUPPORTED`, reported with the claim and the
record's own subject line so a maintainer settles it by reading two quotes.

**Stage scope.** The package's entry point runs every detector over every
artifact.
It gains a scope so a caller can ask for the rules decidable at one authoring
stage: `prd`, `techspec`, or `tasks`. The default stays the full sweep, so
`make verify` is unchanged.

The detectors do not move between stages arbitrarily. A rule belongs to the
earliest stage at which every artifact it reads exists — Project Constraints at
`prd`, coverage mapping at `tasks`, citation wherever the citing artifact is
written.

## Implementation Design

### Interfaces

```go
// Stage names the authoring moment a caller is validating.
type Stage string

const (
    StageAll      Stage = ""         // full sweep; what make verify runs
    StagePRD      Stage = "prd"
    StageTechSpec Stage = "techspec"
    StageTasks    Stage = "tasks"
)

// Claim is one attribution an artifact makes about a decision record.
type Claim struct {
    Artifact string // repository-relative path making the claim
    Line     int
    Target   string // e.g. "ADR-0083"
    Subject  string // what the artifact says the record establishes
}

func detectUnsupportedCitations(result *Result, repoRoot string, claims []Claim) error
```

### Data Models

`Result` is unchanged. `SC-CITATION-UNSUPPORTED` is an error-level code, so it
blocks like the existing constraint codes rather than reporting as a gap.

### API Contracts

`roundfix spec check [<slug>...] [--stage prd|techspec|tasks] [--format text|json]`

`--stage` narrows which detectors run. Without it the command behaves exactly as
today, which keeps `make verify` and every existing caller unchanged.

## Coverage Map

| PRD goal | Component |
| --- | --- |
| 1 — a citation is checked against what it cites | The citation detector |
| 2 — every file-decidable check runs during authoring, under a second | Stage scope plus the existing detector set |
| 3 — the gate stops spending turns on file questions | The `qa-gate` matrix loses its governance rows, in Build Order step 6 |
| 4 — a rule a finding established becomes executable | `SC-CITATION-UNSUPPORTED`, which is the 2026-08-06 finding made mechanical |

## Integration Points

- `internal/speccheck/citations.go` — the semantic check, beside the listing ones.
- `internal/speccheck/coherence.go` — stage scope and detector registration.
- `internal/cli/spec_check.go` — the `--stage` flag and its help.
- `CONTEXT.md` — the coined code.

## Testing Approach

The corpus is captured first and has an unusually good fixture available: Spec
0090's PRD as it stood at the moment the gate failed it, claiming ADR-0083 makes
`make verify` authoritative. That exact text is the negative control — the
detector must reject it — and the corrected text is the positive control.

The observability control matters here more than usual: a citation detector that
silently matches nothing would pass every artifact. A test asserts that the
detector reports how many claims it resolved, and that an artifact making no
claims is distinguishable from one whose claims were not parsed.

## Build Order

1. **Characterization corpus.** Record that Spec 0090's original PRD text passes
   every current check, and that no detector reads an ADR's body. Declares the
   break. Depends on nothing.
2. **The citation detector.** Parse claims, resolve targets, match against the
   record, report `SC-CITATION-UNSUPPORTED` with both texts. Depends on step 1.
3. **Stage scope.** Add `Stage`, assign each detector its earliest stage, keep
   the default sweep identical. Depends on step 1; runs parallel to step 2.
4. **The `--stage` flag.** Surface the scope through the CLI with help text.
   Depends on steps 2 and 3.
5. **Wire PRD and TechSpec authoring.** Each authoring skill ends by running the
   stage-scoped checker and blocks on an error-level finding. Depends on step 4.
6. **The gate reads product.** Remove from `qa-gate` every row whose rule the
   checker now decides, keeping the post-commit rows and every rule without an
   equivalent. Depends on step 5, because a rule may only leave the gate once it
   is demonstrably running earlier.
7. **QA gate.** Depends on step 6.

## Risks & Considerations

**A false accusation costs more than a missed defect.** A citation check that
misreads a legitimate claim trains authors to ignore it, which is worse than not
having it. Hence reporting both texts, and hence keeping the detector to
attributions of a *subject* rather than judging nuance.

**A mention is not a citation, and the existing checks cannot tell.** Authoring
this Spec hit it immediately: naming ADR-0083 as the worked example of a false
attribution was reported as an unlisted citation, and a rehearsal case using a
real record number as sample data was too. The listing checks match a token, not
a stance. The semantic detector must not inherit that — an artifact discussing a
record is not an artifact claiming what it establishes — and the resolved-claim
count exists so this distinction is observable rather than assumed.

**Stage assignment can hide a rule.** Moving a detector to `prd` means the
`tasks` stage no longer runs it. The default sweep is what protects against
that: `make verify` still runs everything, so a rule can only be skipped early,
never dropped.

**Removing a gate row is the one irreversible-feeling step.** A rule that leaves
the gate and turns out not to be running earlier is a hole nobody sees until it
lets a defect through. Step 6 depends on step 5 for exactly that reason, its
Task requires a named checker rule per removed row, and the QA gate's matrix
maps them one by one. A rule with no equivalent stays.

## Decisions

- A citation is read against its target, per ADR-0116.
- A check belongs to the stage that can produce its defect, per ADR-0117.
- The default sweep is unchanged, so early checking is an addition rather than a
  redistribution that could lose a rule.
