---
type: feat
status: promoted
created: 2026-10-24
spec: 0316-an-agent-may-start-the-setup-skill
reason: null
---

# An Agent may start the setup skill without a person invoking it

## Opportunity

The `setup-context-driven` skill that Roundfix ships carries
`disable-model-invocation: true`, so an Agent never sees it as an available
skill and cannot start it; the Baseline's skill guide then tells the Agent to
ask the person to run it. Every adopter that wants an Agent to keep its
Baseline current, for example during a fleet sweep or a follow-up session,
needs the person to type the skill's command by hand first. The maintainer
asked, through Quaestor (ai-memory message `01a16dbe`, 2026-10-24, consumed
by the operator at 11:25 BRT), to remove that lock. The operator recorded the
request as: remove `disable-model-invocation: true` from the
setup-context-driven skill and the matching block in skill dispatch, with the
maintainer's own words kept verbatim:

> "preservando as confirmações do próprio comando para planos e operações
> quando aplicáveis"

(keeping the command's own confirmations for plans and operations where they
apply). The maintainer confirmed the direction on 2026-10-24 at about 11:15
BRT: "Tirar a trava da skill (Recomendado)".

## Value

An Agent can run `roundfix baseline update` when the skill's trigger matches,
present the Plan and its Plan Digest, and wait for the person's approval of
that digest, so the person approves the change itself instead of first
starting a skill whose only side effects already wait for that approval.

## Shape

Drop the metadata line from the canonical skill and its shipped mirror, raise
the skill's version, and say in the skill and the guide that starting it
writes nothing: the Baseline Command still writes only after the person
approves the exact Plan Digest (`--confirm-plan <digest>`, `--yes`, or the
interactive confirmation). The skill guide's generic clause about skills only
a person can start stays for the skills that keep the lock.
