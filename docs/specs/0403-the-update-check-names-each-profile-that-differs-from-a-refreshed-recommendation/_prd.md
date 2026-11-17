---
spec: 0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation
status: active
created: 2026-11-17
surfaces: [backend, cli, docs]
type: fix
---

# The Update Check names each profile that differs from a refreshed recommendation

The Recommended Profile Roundfix ships is dated 2026-09-30. Since then Codex
added `gpt-6.1-terra` and lists `gpt-6-luna` over the `gpt-5.6` models it marks
older, the Claude adapter advertises `fable` for Fable 5.2, and `haiku` became
Haiku 5.5, while the shipped snapshot still leads `review`, `docs` and `chore`
with `gpt-5.6-luna` and the Model Catalog still offers `claude-fable-5-1` and
describes Haiku 4.5. A repository whose Project Config copied an older snapshot
learns of it only from `roundfix profiles check`, the Doctor or an upgrade
notice, none of which a session runs every day: on 2026-11-17 the maintainer
found one running an old configuration. After this Spec the snapshot is dated
2026-11-17 with the table the maintainer approved, the catalog follows the
adapters, and the Update Check every session runs names each configured
category that differs, with the command that fixes it.

## Sources

The maintainer's request and answers of 2026-11-17, recorded in
[`references/2026-11-17-maintainer-answers.md`](references/2026-11-17-maintainer-answers.md),
and the model evidence read the same day, in
[`references/2026-11-17-model-evidence.md`](references/2026-11-17-model-evidence.md).
No Backlog Entry or Finding is adopted.

No open Backlog Entry or unresolved Finding shares this context. Spec 0402,
planned the same day, changes the Update Check and the Doctor for the Authoring
Model Declaration, a different section of Project Config; it shares files with
this Spec, so this Spec requires it.

## Prerequisites

`requires`: `0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration`.
Both Specs change `internal/cli/update_check.go`, the Update Check's guide, the
Roundfix skill's `modes.md` and `setup.md`, `CONTEXT.md` and `.roundfixrc.yml`,
and the maintainer ordered this Spec after it. Authored against `63a9f493`.

## Project Constraints

- Identifier strategy: not applicable — no persistent identifier is minted; the new JSON member and its fields are public CLI contract named in the TechSpec. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the comparison is offline and reads the loaded configuration and the snapshot in the binary; the Update Check's release lookup is unchanged. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0416 (new) amends ADR-0181; ADR-0180 (one dated Recommended Profile), ADR-0235 (OpenAI and Anthropic models only through the subscriptions), ADR-0371 (the Update Check) and ADR-0256 (a changed surface ships with its skill) bind the design; ADR-0217, ADR-0231, ADR-0238 and ADR-0349 do not apply. The TechSpec relates each one. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-17, quoted in the authorization record. It covers this Spec's governed files, `.roundfixrc.yml` among them. It grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation/_authorization.md`; bounded files: `.roundfixrc.yml`, `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/profiles.md`, `.agents/skills/roundfix/references/setup.md`, `.agents/skills/roundfix/references/modes.md`, `internal/docscontract/publicdocs_test.go`.

## Goals

- The Recommended Profile is the table the maintainer approved on 2026-11-17, dated 2026-11-17, with no `ultra` effort.
- The Model Catalog offers every model the snapshot recommends and the identifiers the installed adapters advertise, with current descriptions.
- `roundfix update check` names each configured category that differs from the Recommended Profile, with the command that adopts it, and reports a Profile Deviation declared against an older snapshot as pinned to an older snapshot, with the key that re-pins it.
- The `--update` Skill Mode shows those categories after an upgrade and runs a fix only when the user chooses it.
- This repository's own Project Config follows the new table.

## Core Features

1. A new snapshot in `internal/config/recommendations.go` and the reference `docs/references/model-selection.md`.
2. The Model Catalog gains `gpt-6.1-terra`, offers `fable` in place of `claude-fable-5-1` and describes Haiku 5.5.
3. A `recommendations` member and `Profiles:` lines in the Update Check.
4. The guides, the Roundfix skill's references and the glossary describe both; ADR-0181 is marked amended.
5. `.roundfixrc.yml` takes the new `review`, `docs` and `chore` profiles and current comments.

## User Experience

In an adopter whose Project Config holds the 2026-09-30 table and a `review`
pin declared against it, after upgrading:

```text
$ roundfix update check
Roundfix 0.148.0 is current (latest release 0.148.0; update.mode notify).
Baseline: current with the installed Roundfix.
Profiles: docs (project) differs from the Recommended Profile of 2026-11-17: configured codex / gpt-5.6-luna / max, then claude / sonnet / high; recommended codex / gpt-6.1-terra / high, then claude / sonnet / high; adopt with roundfix profiles check --apply --scope project.
Profiles: review (project) is pinned to an older snapshot (2026-09-30: Keep the validated model); recommended codex / gpt-6-luna / max, then codex / gpt-6.1-sol / high as of 2026-11-17; re-pin with profiles.review.deviation.from: 2026-11-17 in Project Config, or adopt with roundfix profiles check --apply --scope project.
```

## Non-Goals / Out of Scope

- Changing `roundfix profiles check`, its JSON schema, its statuses, `--apply`, the Doctor line or the upgrade notice.
- A command that writes a re-pin, or any write by the Update Check.
- Comparing when a Run starts, in Preflight or in the Daemon (ADR-0181).
- Recommending or adding to the picker the `ultra` effort.
- The Baseline semantic analysis selections of ADR-0180, and the light tier.
- Changing any adopter's configuration; adopters decide through the report.
- Any route of an OpenAI or Anthropic model through OpenRouter (ADR-0235).

## Success Metrics

1. Success Metric: `RecommendedProfile` returns, for every Agent Work Category, exactly the approved table, and `ModelRecommendationSnapshotVersion` is `2026-11-17`; no recommended selection has the effort `ultra`.
2. Success Metric: every recommended Codex and Claude model is in the Model Catalog, which lists `gpt-6.1-terra` after `gpt-6-astra` and `fable` in place of `claude-fable-5-1`, and offers no retired model.
3. Success Metric: in a fixture repository whose Project Config holds the 2026-09-30 table with a `review` deviation from `2026-09-30`, the Update Check text prints one `Profiles:` line per differing category in category order, the `review` line names the older snapshot, the re-pin key and the adopt command, and the JSON `recommendations` member carries the same rows; `next` and the exit code are unchanged.
4. Success Metric: a repository with no profile configuration, or with every category current or pinned to the shipped snapshot, gets the same Update Check bytes as before this Spec.
5. Success Metric, outside evidence: the snapshot's new selections are the ones the installed Codex 0.169.0 model listing and the Claude adapter 0.107.3 advertise, with the efforts each accepts, as recorded in `references/2026-11-17-model-evidence.md`.

## Acceptance evidence

The outside-evidence row rests on the Codex model listing and the Claude adapter
source installed on the planning machine and on the vendors' release pages,
read on 2026-11-17 and committed under `references/`. The QA gate checks the
shipped snapshot against that record without the network.

## Unreachable Acceptance

None. A live `roundfix profiles validate` of the new tuples needs the
subscriptions and runs as the operator's release check, outside the gate.

## Validation

The decisions were checked on 2026-11-17 against:

- `internal/config/recommendations.go`, `recommendation_check.go`, `internal/cli/update_check.go`, `profiles_check.go` and `internal/agent/catalog.go`;
- ADR-0180, ADR-0181, ADR-0235 and ADR-0371;
- `~/.codex/models_cache.json` (Codex 0.169.0) and `@agentclientprotocol/claude-agent-acp` 0.107.3;
- OpenAI's and Anthropic's release pages cited in the evidence reference.

## Glossary

- changes: **Update Check**

## Decisions

- 2026-11-17: How should the outdated model configurations be handled? — "Spec 0403 separada (Recommended)": update the Recommended Profile (research; the proposed table goes to the maintainer before it is applied) and make the update check and `--update` show each category that differs with the command to fix it; amend ADR-0181; after Spec 0402.
- 2026-11-17: Which authorization applies to Spec 0403? — "Mesmas do BE (Recommended)": implement, commit, push, pull request, merge and release right after the Spec, up to three corrections; Makefile, CI, `go.mod` and lint only if a correction requires one.
- 2026-11-17: Frontend? — "Manter opus high (Recommended)".
- 2026-11-17: Review and chore? — "gpt-6-luna max (Recommended)".
- 2026-11-17: Docs? — "gpt-6.1-terra high (Recommended)".
- 2026-11-17: Implementation and QA stay `codex / gpt-6.1-sol / high` with `claude / opus / high`; fallbacks `review` `codex / gpt-6.1-sol / high`, `docs` and `chore` `claude / sonnet / high`, `frontend` `codex / gpt-6.1-sol / xhigh`; snapshot 2026-11-17; no `ultra` effort.

Autonomous route: no person present during authoring; the remaining decisions were taken from the operator's brief for Spec 0403 of 2026-11-17.

- A pin declared against an older snapshot is reported as pinned to an older snapshot, derived from the deviation, without a new `profiles check` status (ADR-0416).
- The re-pin is named as the configuration key, not written by a command (ADR-0416).
- `type: fix` with a `minor` Release Ceiling, because the Update Check's JSON gains a member.

## Open Questions

None.
