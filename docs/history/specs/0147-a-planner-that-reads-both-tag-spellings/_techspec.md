---
spec: 0147-a-planner-that-reads-both-tag-spellings
prd: _prd.md
created: 2026-09-18
---

# A planner that reads both tag spellings — Technical Spec

## Executive Summary

One parser learns a second spelling, one selection spans both, and one refusal
appears where the two meet. The parser accepts a stable version with or without
the `v` prefix and keeps every rejection it has; selection compares versions
across spellings; and a highest version reachable under more than one ref
refuses in preflight rather than choosing.

The trade-off this design accepts is that a repository carrying both spellings
of its highest version loses a plan it used to get. It was getting that plan by
having one of its refs ignored, which is the failure this Spec exists to stop.

## Project Constraints

- Identifier strategy: applicable — a tag's exact spelling is part of its identity, so the planner preserves the selected tag's text rather than normalizing it, and coins no identifier. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable to range planning, which reads local Git refs only. Reset planning already reads GitHub Releases through the existing authenticated path, and this Spec adds no transport and changes no credential handling. Source: `docs/agents/agent-instructions.md` and `docs/agents/cli.md`.
- Active ADR obligations: applicable — the command's read-only contract and its public refusal vocabulary are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, so this Spec's terminal gate covers what its Requirements name.
  ADR-0156 applies: this Spec declares its Success Metrics and API Contracts as numbered units and names each in a Task.
  ADR-0104 applies: acceptance rests on evidence this Spec did not author, which here is this repository's own tag history and the published guidance the command already carries.
- Tooling authority: applicable — the change is public CLI behavior, and the repository's hard rule ships the Roundfix skill update with it. Express maintainer authorization: standing grant of 2026-09-18 for the Roundfix skill and its mirror across this queue's remaining slices, consumed and bounded in [_authorization.md](_authorization.md); bounded files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned regeneration: `make skills-sync`, and `make baseline-digests` for any derived pin it moves. The planner, the command surface and the user guide are ordinary source that no authorization has bounded. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Version parser | `internal/releaseplan` | Accept both spellings, reject everything it rejects today, and carry the spelling it read. |
| Selection | `internal/releaseplan` | Compare across spellings and detect a highest version reachable under more than one ref. |
| Preflight refusal | `internal/cli` release plan command | Refuse an ambiguous highest version with both refs and the resolving selector named. |
| Shipped skill | `.agents/skills/roundfix/SKILL.md` and its mirror | State the accepted spellings and the refusal. |

No new package, state, exit code or configuration is proposed.

## Implementation Design

### Both spellings, one version

The parser's prefix requirement becomes optional: a leading `v` is consumed when
present, and the remaining text must be exactly three canonical numeric
identifiers, as it must today. Prerelease and build-metadata detection keep
reading the same shapes, so `1.2.3-rc.1` and `v1.2.3+build` stay rejected
exactly as they are. The parsed value carries the spelling it was read from, so
nothing downstream has to guess.

### Selecting across spellings

Selection orders by semantic version, not by text, so `1.3.0` outranks `v1.2.3`.
When the highest version is reachable under more than one ref — the two
spellings of the same numbers — selection reports that instead of returning one
of them.

### Refusing the ambiguity

The command turns that report into the preflight refusal it already uses for a
target it cannot resolve. The refusal names both refs and the `--from` selector
that resolves them, proposes no version, approves nothing, and keeps the exit
code preflight refusals already use. Ambiguity is therefore not a new state in
the plan's vocabulary.

### Keeping the spelling

The proposed version, the approval question and reset planning render the
spelling the selection resolved to. Reset planning's inventory keeps each
entry's spelling, its ref and its commit, matches Releases by exact tag name,
and never collapses two refs into one entry.

### Interfaces

```go
// A parsed stable version remembers how it was written.
type Version struct {
    Major, Minor, Patch int
    Prefixed            bool // the tag carried a leading "v"
}
```

### Data Models

No stored record changes. The reset digest keeps its fields and gains no new
one; what changes is that two spellings of one version occupy two entries
instead of colliding.

## API Contracts

1. `roundfix release plan` accepts stable tags written as `MAJOR.MINOR.PATCH`
   and as `vMAJOR.MINOR.PATCH`, and keeps rejecting malformed, prerelease and
   build-metadata inputs with their current messages.
2. When the highest reachable version exists under more than one ref spelling,
   the command refuses in preflight, names both refs and the `--from` selector,
   proposes no version, and exits with the code preflight refusals already use.
3. The proposed version, the approval question and the reset digest use the
   spelling of the selected tag; every other state, exit code and field is
   unchanged.

## Coverage Map

- Goal 1 → Version parser, Selection.
- Goal 2 → Regression locks; API Contract 3.
- Goal 3 → Preflight refusal.
- Goal 4 → Keeping the spelling.
- User Story 1 → Version parser.
- User Story 2 → API Contract 3.
- User Story 3 → Preflight refusal; API Contract 2.
- User Story 4 → Keeping the spelling.
- Core Feature 1 → Version parser; API Contract 1.
- Core Feature 2 → Selecting across spellings.
- Core Feature 3 → Refusing the ambiguity.
- Core Feature 4 → Keeping the spelling.
- Core Feature 5 → Keeping the spelling (inventory and digest).
- Core Feature 6 → Refusing the ambiguity; API Contract 2.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 3.
- Success Metric 3 → Testing Approach 2.
- API Contracts 1-3 → Version parser, Preflight refusal, Keeping the spelling.

## Integration Points

- **This repository's release pipeline.** Unchanged. Its tags are prefixed and
  its workflow triggers on `v*`; the planner is what learns the second spelling.
- **Reset planning.** It keeps its inventory, its digest and its approval
  boundary, and gains the spelling as part of each entry's identity.
- **The shipped skill.** It states the accepted spellings and the refusal, under
  the standing authorization this Spec consumes.

## Testing Approach

1. **Bare tags, at the planner's fixture seam.** A repository tagged `1.4.2`
   plans to `1.4.3`. The case fails on the tree as it stands today, where the
   tag is reported malformed.
2. **Ambiguity, at the command seam.** A repository carrying `v2.0.0` and
   `2.0.0` refuses in preflight, naming both refs and the selector, with no
   proposed version and the existing exit code.
3. **Prefixed tags are untouched.** The existing planner tests pass unedited,
   and a case plans this repository's own tag shape to the same result as
   before.
4. **Rejections are untouched.** Malformed, prerelease and build-metadata inputs
   keep their current messages, asserted case by case.
5. **Outside evidence.** This repository's real tag history is planned against
   and must produce the result it produces today. Where those tags cannot be
   read, the row records that reason and does not block.
6. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. Parser accepts both spellings and carries the one it read, with unit tests
   (depends on: none).
2. Selection across spellings, including the ambiguous-highest report, with unit
   tests (depends on: 1).
3. Preflight refusal, spelling-preserving output and reset digest entries, with
   command tests (depends on: 2).
4. The shipped skill states the spellings and the refusal (depends on: 1, 2, 3).
5. Terminal QA (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **A repository that loses its plan.** One carrying both spellings of its
  highest version now refuses. That is deliberate, and the refusal names the
  selector that restores a plan in one flag.
- **A prefix consumed too eagerly.** A tag like `version-1.2.3` must stay
  malformed; the parser consumes only a leading `v` followed immediately by the
  numeric core, and the rejection cases assert it.
- **Two refs at the same commit.** They are still two refs; the Spec refuses on
  spelling multiplicity rather than on commit identity, because a maintainer who
  created both should say which one is the release.

## Decisions

- **Refuse, do not resolve.** Choosing silently would make output depend on ref
  creation order.
- **Carry the spelling in the parsed value.** Recomputing it downstream is how a
  proposal ends up in the wrong shape.
- **No new state.** Ambiguity is a preflight refusal, which the command already
  has; adding a state would change a public vocabulary for one case.

## Vocabulary Contract

No token is coined. Release Plan Command, preflight refusal and stable version
are existing terms used with their existing meanings.
