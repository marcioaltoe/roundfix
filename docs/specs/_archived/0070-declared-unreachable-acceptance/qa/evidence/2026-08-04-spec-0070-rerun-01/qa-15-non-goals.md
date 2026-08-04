# QA-15 — Non-Goals and scope creep

Status: pass

- The QA Gate Skill still says the gate must never declare unreachability and
  may match only a pre-run PRD declaration.
- The explicit override path remains available and was exercised in QA-10;
  the Archive Command gained no override flag.
- The gate still derives a complete matrix from stories, criteria, Non-Goals,
  and surface sweeps. The change adds declaration matching and implicated-row
  scoping; it does not remove matrix derivation.
- `rtk git diff --quiet d8c0403..HEAD -- Makefile` exited 0. No cache-warming
  or cheap-detector ordering change was introduced.
- QA-01's commit-by-commit path audit contains only the parser, count, archive,
  authorized Skills, generated fallout, replay, documentation, test, Task, and
  QA evidence surfaces this Spec names.
