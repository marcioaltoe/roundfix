---
task: task_04
spec: 0285-a-transient-profile-proof-timeout-does-not-park-the-delivery-queue
status: completed
type: docs
complexity: low
---

# Task 04: Document the setup deadline and the third attempt

## Overview

The inbox entry of 2026-10-20 noted that the Roundfix `profiles` reference
offers no option for the proof's setup deadline (`_prd.md` → Sources 1), and
two documents still say a timed-out proof "is retried once". This Task
documents the attempts and the backoff in the user guides and the Roundfix
references, beside the key task_01 documented, adds **Proof Setup Deadline** to the glossary, raises
the Roundfix skill version and re-records the Behavior Surface Record.

## Requirements

1. MUST answer the documentation remark of the Secondbrain Inbox Entry of
   2026-10-20,
   [`references/2026-10-20-a-prova-de-perfil-do-preflight-estoura-o-tempo-em-serie.md`](references/2026-10-20-a-prova-de-perfil-do-preflight-estoura-o-tempo-em-serie.md)
   (`_prd.md` → Sources 1), with the behavior of API Contracts 1-3 of
   `_techspec.md` (ADR-0298).
2. MUST document, in these files:
   - `docs/user-guide/configuration.md`: in the subsection task_01 added
     under Agent selection profiles, state that a timed-out attempt is
     "retried after 5 seconds and again after 10 seconds";
   - `docs/user-guide/commands/profiles.md`: replace the paragraph that says
     a timed-out proof "is retried once" with the three attempts, written as
     "retried after 5 seconds and again after 10 seconds", the `temporary`
     classification of a proof that timed out on 3 attempts, and the next
     action naming `defaults.profile_proof_setup_timeout`;
   - `.agents/skills/roundfix/references/profiles.md`: in the section task_01
     added, the two waits written as "retried after 5 seconds and again after
     10 seconds" and the `temporary` classification of the third timeout;
   - `.agents/skills/roundfix/references/runtime.md`: replace "is retried
     once. A second timeout is classified `temporary`" with the same three
     attempts, "timed out on 3 attempts", and the key.
   Each of the four files names `defaults.profile_proof_setup_timeout`.
3. MUST search `docs/user-guide`, `.agents/skills/roundfix`, `skills/roundfix`
   and `README.md` for "retried once" and "timed out twice" describing the
   profile proof, and change each occurrence. No text inside any `### QA
   settlement` section may change.
4. MUST add **Proof Setup Deadline** to `CONTEXT.md` through domain-modeling:
   the wall-clock bound on one attempt of an Exact Agent Selection Proof's
   disposable Agent Session setup, set by
   `defaults.profile_proof_setup_timeout` (default 30 s); an attempt that
   exceeds it is retried after 5 seconds and again after 10 seconds, and a
   third timeout is classified temporary. Link
   `docs/adr/0298-a-timed-out-profile-proof-gets-a-third-attempt-and-its-setup-deadline-is-configurable.md`.
5. MUST run these commands, because this Task edits user guide pages and the
   Roundfix references. Before them, raise both version fields of the
   Roundfix skill's canonical `SKILL.md` (`metadata.version` and the
   top-level `version`) by one patch over the value the tree holds when this
   Task starts, because the record command refuses changed content under an
   unchanged version. The Agent runs these commands itself before handing
   back: the Daemon only runs the Verification and never regenerates a mirror
   or a record. No record or mirror is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to record the raised Roundfix skill version;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Document the key and the attempts in the configuration and `profiles` guides.
- [ ] Document them in the Roundfix `profiles` and `runtime` references.
- [ ] Add **Proof Setup Deadline** to the glossary.
- [ ] Raise the skill version, sync and record.

## Acceptance Criteria

- [ ] The four documents name the key, and the guides and references describe three attempts with waits of 5 and 10 seconds.
- [ ] No guide or skill copy still says a timed-out proof is retried once.
- [ ] The glossary defines **Proof Setup Deadline**, linking ADR-0298.
- [ ] The skill mirror, the skill version and the Behavior Surface Record are current.

## Consumers

- the Roundfix `profiles` and `runtime` references and the guide fingerprints — consumers: `./skills`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0298-a-timed-out-profile-proof-gets-a-third-attempt-and-its-setup-deadline-is-configurable.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/profiles.md`
- interface: `.agents/skills/roundfix/references/profiles.md`
- interface: `skills/roundfix/references/profiles.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `for file in docs/user-guide/configuration.md docs/user-guide/commands/profiles.md .agents/skills/roundfix/references/profiles.md .agents/skills/roundfix/references/runtime.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'defaults.profile_proof_setup_timeout' || { printf 'missing the key in %s\n' "$file" >&2; exit 1; }; done; for file in docs/user-guide/configuration.md docs/user-guide/commands/profiles.md .agents/skills/roundfix/references/profiles.md .agents/skills/roundfix/references/runtime.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'retried after 5 seconds and again after 10 seconds' || { printf 'missing the backoff in %s\n' "$file" >&2; exit 1; }; done; for file in docs/user-guide/commands/profiles.md .agents/skills/roundfix/references/runtime.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'timed out on 3 attempts' || { printf 'missing the attempt count in %s\n' "$file" >&2; exit 1; }; done; for file in docs/user-guide/configuration.md .agents/skills/roundfix/references/profiles.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'at most 10m' || { printf 'missing the range in %s\n' "$file" >&2; exit 1; }; done; if grep -rq 'is retried once' docs/user-guide .agents/skills/roundfix skills/roundfix; then printf 'a guide still says the proof is retried once\n' >&2; exit 1; fi; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Proof Setup Deadline**:' || { printf 'missing glossary term\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '0298-a-timed-out-profile-proof-gets-a-third-attempt-and-its-setup-deadline-is-configurable.md' || { printf 'glossary term does not link ADR-0298\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no document names the key, so the command fails at its first check. Afterwards the guides and references name the key, the backoff, the attempts and the range, none still says "is retried once", the glossary defines the term with its ADR, and the skill mirror, the skill version and the documentation contracts are current.

## References

- `_prd.md` → Sources; Core Feature 4; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; Glossary; Build Order 4
- ADR-0298; ADR-0189; ADR-0256; ADR-0280

## Result

Implemented the documentation slice for the configurable Proof Setup Deadline
and the third profile-proof attempt. The configuration guide, `profiles`
command guide, canonical Roundfix `profiles` and `runtime` references, their
synced skill copies, and `CONTEXT.md` now describe the `5s` and `10s` waits,
the `temporary` classification after three timeouts, and
`defaults.profile_proof_setup_timeout`. The glossary entry links ADR-0298.

Implementation commands:

- `rtk proxy make skills-sync` — exit 0; canonical and mirror skill copies
  are byte-identical.
- `rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — exit 0; Roundfix skill `0.1.112` recorded.
- `rtk proxy go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1` — exit 0; Behavior Surface Record regenerated.

Focused checks:

- Repository search found no `is retried once` or `timed out twice` occurrence
  in `docs/user-guide`, `.agents/skills/roundfix`, `skills/roundfix` or
  `README.md`.
- Repository search found the required key, backoff, three-attempt wording and
  glossary entry in the declared surfaces; `git diff --check` exited 0.
- The diff changes no text under the Roundfix skill's `### QA settlement`
  section.

The Daemon must run the authored Verification and decide Task settlement.
