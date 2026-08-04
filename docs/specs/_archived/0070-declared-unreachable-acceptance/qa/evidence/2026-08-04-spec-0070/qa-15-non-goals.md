# QA-15 — Non-Goals and scope creep

Status: pass

- The QA skill still says the gate must never declare unreachability and may
  match only a pre-run PRD declaration.
- The explicit override operational path remains available and was exercised
  in QA-10; the Archive Command gained no override flag.
- The gate still derives the complete matrix from PRD stories, criteria,
  Non-Goals, and surface sweeps. The change adds declaration matching and
  implicated-row scoping; it does not remove matrix derivation.
- `rtk git diff --quiet d8c0403..HEAD -- Makefile` exited 0. No cache-warming
  or cheap-detector ordering change was smuggled into the implementation.
- Exact changed-path provenance in QA-01 contains only the specified parser,
  report-count, archive, QA-skill, generated-digest, replay, fixture, test,
  and Task files.
