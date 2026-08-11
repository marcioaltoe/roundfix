# R07 — vocabulary sweep

- Strict authoring check exited 0 but explicitly skipped
  `SC-VOCABULARY-UNDOCUMENTED` because `_techspec.md` has no
  `## Vocabulary Contract`.
- Production emits four stable diagnostics at
  `internal/speccheck/mechanical.go:28-31`:
  `QA-AUTH-PATHS`, `QA-CONSEQUENT-ORDER`, `QA-REPORT-SHAPE`, and
  `QA-EVIDENCE-PATH`.
- A fresh recursive search of `CONTEXT.md`, `docs/agents/`,
  `docs/user-guide/`, `docs/references/`, and `docs/adr/` for those four tokens
  exited 1 with no match.
- The codes appear only in implementation/tests and Spec-local artifacts.
  Spec artifacts are not durable vocabulary owners under the documentation
  layout contract.
- The typed evidence-input tokens are documented in ADR-0097 and the QA skill;
  the two-tier target is documented in the adopted agent guides. The gap is
  limited to the four emitted diagnostic codes and the absent declaration
  that would have checked them.
