---
type: feat
status: open
created: 2026-10-03
spec: null
reason: null
---

# A repository cannot leave the branch prefix to the commit-type rule

## Symptom

`internal/baseline/assets/modules/core.json` lists `branch.prefix` among `requiredDecisions`, and `templates/guides/agent-instructions.md` interpolates it unconditionally ("The branch-prefix pattern is {{branch.prefix}}"). A repository that wants only the normal rule, `<type>/<description>` from the Conventional Commit type, cannot turn the decision off. The oraculum maintainer asked for this on 2026-10-01; the maintainer's global rule already says branches follow `<type>/<description>`.

## Expected

`branch.prefix` becomes an optional decision. While it is unrecorded, the guide states the Conventional Commit type rule; a recorded value keeps today's sentence. An adopter's next update asks nothing and keeps its recorded value. Clause retention follows ADR-0191 and the optional-decision machinery from Spec 0207.

## Source

Secondbrain `inbox/roundfix/2026-10-01-branch-prefix-e-decisao-obrigatoria-e-nao-deixa-o-repo-usar-a-regra-normal.md` (oraculum).
