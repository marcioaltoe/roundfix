---
status: accepted
created_at: 2026-07-17T16:00:35Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Review Source Evidence determines review outcomes

Watch classifies CodeRabbit signals as head-bound Review Source Evidence: an explicit skip ends Review Skipped, a current-head CodeRabbit approval with no unresolved CodeRabbit threads can prove Merge-Ready, and an exact Daemon-created artifact-only descendant may inherit its verified parent evidence without another Roundfix review request or wait. Missing accepted evidence still ends Clean Unverified under ADR-0043; this refines ADR-0019 and preserves ADR-0036's separate review-artifact commit.
