---
spec: 0400-planning-references-and-release-preparation-need-no-manual-classification
status: active
created: 2026-11-16
surfaces: [backend, cli, docs]
type: fix
---

# Planning references and release preparation need no manual classification

The Release Plan Command asks for a classification by hand for two kinds of commit that ship nothing. The first is a Spec's planning commit that adds outside evidence under `docs/references/`, as `docs/agents/spec-routing.md` requires. The second is the release-preparation commit that writes the changelog and the version files. A Release Ceiling never answers a classification by hand, so in adopters that record ceilings every release stops and asks again. The Oraculum met it on 2026-11-16, and Roundfix met it at `v0.143.4`. The release guide also tells every adopter to preserve a tag-triggered workflow, while the Quaestor publishes its release when the version change merges. After this Spec, the plan classifies both commits itself, and the release guide describes both publication triggers.

## Sources

Two Backlog Entries, adopted by this Spec's planning change (see `references/_index.md`):

1. [`references/2026-11-16-a-planning-change-with-references-makes-every-release-a-manual-classification.md`](references/2026-11-16-a-planning-change-with-references-makes-every-release-a-manual-classification.md), the Oraculum's finding (ai-memory message `01a1e695-ed77-7d50-8d78-93be226b26dd`, no reply expected). Whole.
2. [`references/2026-11-16-the-release-guide-assumes-tag-triggered-releases.md`](references/2026-11-16-the-release-guide-assumes-tag-triggered-releases.md), the release-guide half of the Quaestor's message (ai-memory message `01a1e522-a608-7563-8568-289c6138d377`, information only). The author split it from `docs/backlog/2026-11-16-the-delivery-queue-and-the-release-guide-assume-loose-branch-rules-and-tag-releases.md`, which keeps the Delivery Queue half open and names this Spec for the other.

Neither message is a Notice Request. The maintainer's answers are in [`references/2026-11-16-maintainer-answers.md`](references/2026-11-16-maintainer-answers.md). The measurement over this repository's releases is in [`references/2026-11-16-release-plan-classification-measurement.md`](references/2026-11-16-release-plan-classification-measurement.md).

No other open Backlog Entry and no unresolved Finding shares this context:

- The Delivery Queue half changes the queue's merge step, not the release plan, and the maintainer kept it in the backlog.
- The phase-boundary and large-context entries change authoring sessions.
- The docs-folders entry is parked by the maintainer.
- The Conexus status entry is about the status view.

## Prerequisites

None in `requires`. Authored against `c2d1f77c` (`docs/plan-0399`, based on v0.147.0). Spec 0399 is active and shares no file with this one: it changes the Baseline, `baseline update` and `doctor`, and this Spec changes the release plan and the release guide. Both raise the version of the Roundfix skill, so whichever Spec delivers second records its skill version on the merge base, as the version record command does.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted. The JSON field is named in the TechSpec. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the plan stays local and read-only; it reads committed Git objects only. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0412 (new) records the decision; it keeps ADR-0349 (the ceiling never answers a classification by hand), ADR-0143 (the plan stays read-only) and ADR-0256 (a release waits for the skills that describe a changed surface); the TechSpec relates each one. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-16, quoted in the authorization record. It covers the Roundfix skill's release reference and its `SKILL.md` pair. It grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0400-planning-references-and-release-preparation-need-no-manual-classification/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/release.md`, `skills/roundfix/SKILL.md`.

## Goals

- A commit that changes only planning evidence and files under `docs/references/` is maintenance-only, so the plan needs no `--impact` for it.
- A commit with no automatic impact that changes only paths listed in the `release_paths` of a Spec the range delivers with a ceiling is a release-preparation commit. The plan treats it as maintenance-only and says why in its evidence.
- A plan rerun after the release-preparation commit reaches the decision it reached before that commit, so a Release Ceiling answers it as before.
- The Release Ceiling's authority does not grow: it still never answers a breaking plan or a classification by hand.
- The release guide tells a session to keep the publication trigger the repository already uses, whether a pushed tag or the merge of the version change.

## Core Features

1. `docs/references/` joins the planning-evidence paths of the release plan's maintenance boundary.
2. The plan reads the delivered Specs' records before it classifies, when a commit may be a release-preparation commit, and recognises that commit by their `release_paths`.
3. The JSON change evidence carries `releasePreparation`, and the text lists such a commit nowhere, as it lists no maintenance-only commit.
4. The Roundfix skill's release reference and `SKILL.md`, the release runbook, the usage page and the glossary describe both rules, and the release reference describes both publication triggers and asks that `release_paths` list every file the preparation commit changes.

## User Experience

Today, on this repository's history:

```text
$ roundfix release plan --from v0.143.3 --to v0.143.4
Decision: manual_classification_required
...
Blocking commits:
- 576ae51d9d4a docs: plan 0393 a feature is judged for the idea stage before its prd (#871) [impact=none, type=docs, manual]
- 4f8a5b7db838 chore: release v0.143.4 (#873) [impact=none, type=chore, manual]
```

With this Spec, the same command prints `Decision: ready`, `Proposed version: v0.143.4` and `Determining commits:` naming only the `fix:` delivery of Spec 0391.

## Non-Goals / Out of Scope

- Letting a Release Ceiling answer a classification by hand (ADR-0349, ADR-0412).
- Other root files outside every planning path, such as `.ai-memory.toml`, which stopped `v0.142.1`.
- A release-preparation commit in a range that delivers no Spec with a ceiling.
- Changing the Delivery Queue, which keeps the other half of the Quaestor's entry.
- Changing Roundfix's own tag-triggered workflow or its runbook's publication steps.
- Changing the decision states, the exit codes, or the Release Ceiling's five conditions.

## Success Metrics

1. Success Metric: `MaintenanceOnly` returns true for a path under `docs/references/` and keeps every earlier answer, including false for `docs/user-guide/`, `README.md` and `.github/workflows/release.yml`.
2. Success Metric: a committed fixture range with a planning commit that adds a `docs/references/` file, then a `feat:` delivery whose record grants a `minor` ceiling, ends `ready` with `Release authorization: answered`. The same range without the reference file reaches the same decision and proposed version.
3. Success Metric: adding a release-preparation commit that changes only the record's `release_paths` leaves the decision, the proposed version and the authorization answer equal to the range without it. JSON marks that change `releasePreparation: true` and `crossesMaintenanceOnlyBoundary: false`. A preparation commit that also changes one unlisted path, a range whose record has no ceiling, and a `fix:` commit that changes only release paths each keep today's classification.
4. Success Metric: the release reference keeps the trigger the repository uses and names both a pushed tag and the merge of the version change; it no longer says "Preserve the existing tag-triggered workflow".
5. Success Metric, outside evidence: run with the built binary on this repository's tags, the plan at the tag reaches the same decision as at the commit before the tag in each of the 46 ranges the measurement marks `yes`. The range `v0.143.3..v0.143.4` is `ready` with `v0.143.4` at both revisions. The five ranges marked `no` stay `manual_classification_required` at the tag.

## Acceptance evidence

The outside-evidence row rests on this repository's release history from `v0.125.0` to `v0.147.0`, which this Spec did not build. The measurement of 2026-11-16 records each range's decision with the installed Roundfix 0.147.0 and whether the delivered records list the preparation commit's paths. The QA gate reruns the same ranges with the built binary, read-only.

## Unreachable Acceptance

None. The tags are local to every clone of this repository. The Oraculum's own range lives in another repository, so this Spec does not rest on it.

## Validation

The decisions were checked on 2026-11-16 against these sources:

- The measurement in `references/`.
- The code of `internal/releaseplan/classifier.go`, `build.go` and `release_grant.go`, and of `internal/cli/releaseplan_command.go` and `releaseplan_git_source.go`.
- ADR-0349 and ADR-0143.
- `docs/agents/docs-layout.md`, which gives `docs/references/` its job, and `docs/agents/spec-routing.md`, which sends outside evidence there.

## Glossary

- adds: **Release-Preparation Commit**

## Decisions

- 2026-11-16: What happens to the Oraculum's finding? — "Spec fix depois da 0399 (Recommended)": a release-plan Spec that adopts the Oraculum finding and part 2 of the Quaestor entry (release.md assumes tag releases); part 1 (the queue updating a branch behind `main`) stays in the backlog.
- 2026-11-16: Which authorization applies to that Spec? — "Mesmas do BE (Recommended)": implement, commit, push, pull request, merge and release right after the Spec, up to three corrections. Makefile, CI, `go.mod` and lint only if a correction requires one.

Autonomous route: no person present during authoring; the remaining decisions were taken from the operator's brief for the TechSpec of Spec 0400 of 2026-11-16 and the measurement.

- Of the Oraculum's three suggestions, the first (as a path rule) and the third are taken; the second, which would let a ceiling answer a classification by hand, is not (ADR-0412).
- `type: fix` with a `minor` Release Ceiling, because the JSON change evidence gains a field.

## Open Questions

None.
