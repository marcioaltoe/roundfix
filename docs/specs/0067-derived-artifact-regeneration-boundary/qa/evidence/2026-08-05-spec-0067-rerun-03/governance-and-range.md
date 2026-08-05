# Governance and implementation-range evidence

Build: `9e668cfe4649c323694c55f8124ed73840260910`.

## Authored gate and Project Constraints

- `_tasks.md` names `task_05` as the sole terminal `qa` node. Its direct
  dependency `task_07` is `completed`, and the dependency chain reaches every
  non-QA node.
- The PRD and TechSpec each account for identifier strategy, authentication and
  HTTP, active ADR obligations, and tooling authority. Each row gives an
  applicability decision, reason, and operative source under `docs/agents/`.
- No implementation commit introduces a project-owned Internal Identifier,
  authentication, HTTP, secret, transport, or deployment contract.
- ADR-0080, ADR-0081, ADR-0085, and ADR-0091 are accepted and applicable.
  ADR-0093 is correctly recorded as a relation candidate that does not govern
  product behavior.

## Tooling authorization and chronology

- Authorization commit `2e560cea708006024286881c5948702e1e4599c2`
  is an ancestor of Task 03 commit `c7ad3f62`; the ancestry check exited 0.
- Task 03's `git diff-tree --no-commit-id --name-only -r c7ad3f62` lists
  exactly `Makefile` and its assigned `task_03.md`.
- Corrective Task 06 commit `c6cf8033` is a separate descendant of Task 03.
  Its paths are the authorized `Makefile`, assigned `task_06.md`, the
  ownership test, and the three ownership records the Task corrects.
- Task 07 commit `9e668cfe` is a separate descendant of Task 06 and does not
  touch protected tooling. Its paths are the assigned Task, ownership code and
  test, and the parity ownership record.
- `git log a188c987..HEAD -- Makefile` names only `c7ad3f62` and
  `c6cf8033`. The authorization record names Spec 0067 and bounds its protected
  tooling mutation to `Makefile` for the regeneration step list and derived
  path scan.

## Derived-content range diagnosis

A broad `a188c987^..HEAD` comparison reports nine derived-file changes. That
range crosses the intervening `Merge main` commit and therefore includes
unrelated deterministic digest changes already on `main`; it is not evidence
that a Spec 0067 Task changed those bytes.

Per-Task `git diff-tree` inspection is the authoritative provenance check:
Task 01 adds only ownership metadata under the derived roots, Tasks 02 and 04
change code and tests, Task 03 changes `Makefile`, Task 06 changes ownership
metadata plus code/test/`Makefile`, and Task 07 changes only the parity
ownership record plus code/test. No Spec implementation commit changes a
digest value or derived artifact byte.
