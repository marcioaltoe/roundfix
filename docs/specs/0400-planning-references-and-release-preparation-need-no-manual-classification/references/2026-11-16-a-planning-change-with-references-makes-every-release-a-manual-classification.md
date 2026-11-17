---
type: fix
status: promoted
created: 2026-11-16
spec: 0400-planning-references-and-release-preparation-need-no-manual-classification
reason: null
---

# A planning change with `docs/references/` makes every release a manual classification, so the Release Ceiling never answers

## Symptom

The Oraculum reported on 2026-11-16 (ai-memory message `01a1e695-ed77-7d50-8d78-93be226b26dd`, type finding, measured, no reply expected) that every release of a Spec with `release_ceiling` ends in `manual_classification_required`. The Spec's own planning Pull Request (`docs:`) commits a file under `docs/references/`, as `docs/agents/spec-routing.md` requires for outside evidence, and the release reference counts only `docs/specs`, `docs/adr`, `docs/findings`, `docs/handoffs` and `docs/backlog` as maintenance-only planning evidence. With `--impact minor --reason …` the plan says `Release authorization: not answered — the plan was classified by hand`, while it lists the delivered Spec's `minor` ceiling. The release-preparation commit (`CHANGELOG.md`, manifests, lockfile, user guide, generated catalog) is also listed as a blocking `[impact=none, type=docs, manual]` commit.

Evidence, in the sender's words: "gesttione-solutions-oraculum@760e8673, Roundfix 0.145.1, 16/11/2026: `roundfix release plan` (base v0.37.1, alvo ad5c4e80) → exit 3, `Decision: manual_classification_required`, `Blocking commits: ce6bc35f9200 docs: planejar a Spec 0116 ... (#428) [impact=none, type=docs, manual]`". Roundfix met the same at v0.143.4 (planning PR #871 brought `docs/references/`, released with `--impact patch`).

The sender's suggestions, which leave the decision to Roundfix: treat `docs/references/` as planning evidence when the `docs:` commit touches no code or the file is cited by a Spec of the same commit; or let the Release Ceiling answer when manual classification covers only impact-none `docs:` commits and every impact-setting commit is a delivery with a ceiling; or recognise the release-preparation commit (only `release_paths` and declared derived files) as impact-free.

## Decision

Marcio, 2026-11-16 ~20:55 BRT: "Spec fix depois da 0399 (Recommended)" — one release-plan Spec that adopts this entry and the release-guide half of `2026-11-16-the-delivery-queue-and-the-release-guide-assume-loose-branch-rules-and-tag-releases.md`; the Delivery Queue half stays in the backlog.

## Where

`roundfix release plan` classification, `references/release.md`, `docs/agents/spec-routing.md` (outside evidence under `docs/references/`).
