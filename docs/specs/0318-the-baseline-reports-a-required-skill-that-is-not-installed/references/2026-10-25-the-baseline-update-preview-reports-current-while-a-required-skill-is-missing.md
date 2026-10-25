---
type: fix
status: promoted
created: 2026-10-25
spec: 0318-the-baseline-reports-a-required-skill-that-is-not-installed
reason: null
---

# The Baseline update preview reports current while a required skill is missing

## Problem

Fluxus (ai-memory message `01a16ea1-6700-7061-8a77-0d729a78b5bb` from
`gesttione-solutions-fluxus`, 2026-10-24 12:57 BRT, consumed by the operator
on 2026-10-25) reproduced it with roundfix 0.117.0 and 0.118.0. The
`standard-typescript-monorepo` Baseline still renders
`docs/agents/skill-dispatch.md` lines 128-129 requiring the `handoff` skill
(`trigger.core.handoff`), while the upstream skills repository removed it.
Fluxus removed it with `bunx skills remove context7 feature-systems-pattern
handoff review triage -y`. Afterwards `roundfix baseline update --repo .
--no-history --format json` returned `ok:true`, `state:current`,
`fileChanges:[]`, `warnings:[]` and `skills.differing:[]` at catalog
`sha256:a5c687aa60ad35f3d967cb4bdc3094f404c8ff103f39447c4c5e83c8d4fa22ba`,
although the dispatch still requires the skill. The configured pre-PR review
(Claude) blocked Fluxus correctly with F1: removing `.claude/skills/handoff`
and `.agents/skills/handoff` while keeping the mandatory trigger tells the
Agent to activate a skill that does not exist.

The message's own words, verbatim:

> "Pedido: corrigir na fonte canônica do Roundfix o catálogo/dispatch/skill
> set para não exigir a skill retirada e alinhar o fluxo de continuidade ao
> mecanismo suportado; verificar também que a validação detecte
> obrigatoriedade sem instalação. Cobrir a regressão, publicar a correção pela
> entrega normal do projeto e responder ao inbox do Fluxus com a versão
> liberada e o procedimento público para regenerar a baseline. Não restaurar
> silenciosamente skills obsoletas nem criar exceção local no Fluxus."

(Fix the catalog, dispatch and skill set at Roundfix's canonical source so
they no longer require the retired skill, and align the continuity flow with
the supported mechanism; also make the validation detect a requirement
without an installation. Cover the regression, publish the fix through the
project's normal delivery, and answer Fluxus's inbox with the released
version and the public procedure to regenerate the Baseline. Do not silently
restore obsolete skills or create a local exception in Fluxus.)

Spec 0317 owns the first half (the catalog stops requiring `handoff` and
moves continuity to ai-memory). This entry owns the second half: "verificar
também que a validação detecte obrigatoriedade sem instalação".

## Shape

When an adopted repository's modules require and dispatch a skill that is
not installed, the `baseline update` preview must not report a clean
`current`: it names each such skill and the command that restores it, as it
already does for an owned skill older than the binary's. `roundfix doctor`
and `roundfix skills check` report it too. The maintainer chose this on
2026-10-25 at about 09:05 BRT: "Spec 0318 logo após a 0317 (Recomendado)".
