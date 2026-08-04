# QA-03 — Declaration reader and diagnostics

Status: pass

The built binary was `roundfix 0.3.1 (b6ea034-dirty, built 2026-08-04
16:11:06 -0300)`; the dirty suffix is this in-progress report and its evidence.

- Present declaration: the declared-only public journey archived and stamped
  the author-supplied `satisfied-by` action.
- Absent section: the pass fixture had no declaration section, archived with
  exit 0, and gained no `unproven` field.
- Malformed declaration: `roundfix archive qa-case` exited 2 and named the
  fixture `_prd.md`, line 12, and missing `reason`; the active Spec remained.
- The first fixture attempt was rejected earlier because the fixture Tasks had
  no mandatory Verification section. That correct preflight behavior was not
  treated as a product finding; the fixture was repaired to match the real
  Task contract and the same public commands were rerun.
