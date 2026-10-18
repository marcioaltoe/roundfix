---
status: pending
created_at: 2026-10-18
updated_at: 2026-10-18
---

# QA gate — the gate commits screenshots that the next mechanical audit refuses (2026-10-18)

Triaged from the Secondbrain Inbox Entry
`inbox/roundfix/2026-10-18-o-qa-gate-versiona-screenshot-que-a-auditoria-mecanica-recusa.md`
(origin `conexus`, captured manually on 2026-10-18) and checked against this
repository at `2ef4f080` (v0.89.0). The observation comes from the Conexus
repository (`useconexus/conexus`), Spec
`0001-usuario-verificado-com-segundo-fator`, whose QA gate ran four times on
2026-10-18 with Roundfix 0.87.0 and the Agent `codex gpt-6.1-sol high`.

## 1. A screenshot the gate committed fails the next gate as a built binary

- Symptom / evidence: the first two QA runs saved PNG screenshots under
  `docs/specs/<slug>/qa/evidence/`, and the Daemon committed them with the
  report (`docs: qa report ... (partial)`). The third run's mechanical stage
  raised 67 `QA-EVIDENCE-PATH` findings, one per PNG, each with the detail
  "is a built binary", and the gate closed `fail` although the report approved
  every behavior row ("Runnable behavior is approved... the behavioral matrix
  would otherwise meet declared-partial conditions"). The run that committed
  the PNGs did not fail for them; the refusal cost the next full QA run, about
  11 minutes and 5.7 million tokens, with no behavior defect. The workaround
  in Conexus, commit `d132944` on the item branch, removed the PNGs, turned
  the old reports' links into text paths and added a text-only evidence
  requirement to the QA Task; Conexus pull request #5 added the same
  requirement to the QA Tasks of its Specs 0002 and 0003.
- Root cause: proven in this repository.
  `detectMechanicalEvidenceScratchState` (`internal/speccheck/mechanical.go:2501`)
  raises `QA-EVIDENCE-PATH` for every tracked blob under the report's
  `evidence/` directory that `mechanicalBlobIsBinary` (`:2574`) classifies as
  binary, and that function calls a blob binary when its first 8000 bytes hold
  a NUL byte (`:2599`). Every PNG holds one: the signature is followed by the
  four-byte length of the IHDR chunk, `00 00 00 0D`. The `qa-gate` skill tells
  the gate to "Store screenshots and focused artifacts under
  `qa/evidence/<YYYY-MM-DD>-<run-slug>/`" (`.agents/skills/qa-gate/SKILL.md:468`),
  and the gate's own report said "Required visual evidence is retained as
  screenshots". The QA Report commit (`internal/daemon/task_engine.go:3006-3053`)
  runs only the Secret Scan over the paths it stages, which skips binary
  content, so nothing at that commit looks at a binary.
- Action / suggestion: choose one side and align the other. Either the
  `qa-gate` skill forbids binaries under `qa/evidence/` and teaches text
  evidence (measurements, an accessibility snapshot, DOM state), or the
  detector accepts images (PNG, JPEG, SVG) as reader-openable evidence and
  refuses only build output. In both cases the Daemon should refuse the binary
  at the report commit that introduces it, not at the next gate.

## 2. A text-only requirement in the QA Task does not keep captures out

- Symptom / evidence: the fourth run, with the text-only requirement in its QA
  Task, closed `partial`, and its report says "No binary capture authored".
  The Daemon's report commit `afe23e9` (`docs: qa report ... (partial)`) still
  versioned 26 new PNGs under `qa/evidence/2026-10-18-final/{375,768,1440,sweep}/`.
  The next gate of the same Spec would fail again. The archive removes the
  Spec folder from the tree, so the PNGs never reach the default branch, but
  they stay in the item branch's history.
- Root cause: partly proven. The mechanical stage lists tracked files only
  (`mechanicalTrackedEvidence`, `git ls-files --stage`,
  `internal/speccheck/mechanical.go:2541`) and runs before the QA Agent's
  turn, so files written during the turn are invisible to it. The QA Report
  commit stages every path the gate changed under the QA directory
  (`prepareQAReportCommit`, `internal/daemon/task_engine.go:3715`). Which
  process wrote the PNGs is unknown; the most likely writer is the gate's
  browser script, which saves captures by itself.
- Action / suggestion: the report commit, which stages the files, is the
  stage that can see them; it should refuse what the next mechanical stage
  would refuse.

## What worked — keep

- The archive kept the PNGs out of the default branch.
- The mechanical stage named every offending path, so the workaround was
  mechanical.
