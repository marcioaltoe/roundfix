# Maintainer answers for Spec 0401

The maintainer answered through AskUserQuestion on 2026-11-16 in the Supervisor session that planned Specs 0399 and 0400, and on 2026-11-17 before any Run of this Spec. The operator recorded the answers in Portuguese, and they are copied here verbatim.

- Mensagem do Marcio logo depois: "Esse negocio de criar sessões a cada etapa está atrasando demais."
- Fronteira de fase → "Só por Spec (Recommended)" (~21:00 BRT): sessão encerra só depois da entrega+release da Spec; fases de autoria em agentes; Spec fix ajusta skills e o aviso do roundfix phase.
- 16/11 ~21:45 BRT: Spec 0401 (fronteira só por Spec) autorização → "Mesmas do BE (Recommended)".
- 17/11 ~10:00 BRT: Marcio notou que ao fim de implementação de Spec a sessão para pedindo reinício (prejudica autonomous work) → fronteira "Só quando você pedir (Recommended)": nenhuma parada automática (nem fase nem Spec, presente ou não); handoff só ao encerrar de fato; Phase Record segue medindo; aviso do roundfix phase sai; ajustar ~/.claude/CLAUDE.md. Emendar plano 0401 antes de rodar; skills só depois da 0401.

The operator's brief for this Spec's first plan restated the answer of 2026-11-16 in English: the Supervisor session ends only after the Spec is delivered and released; the Authoring Phases run in agents with their own context; the Phase Record keeps measuring.

The answer of 2026-11-17 supersedes it. The operator's brief for the amendment restates it in English: no automatic stop at all, neither after an Authoring Phase nor after a Spec's delivery and release, with or without a person present. The Supervisor follows the queue — authoring delegated to agents with their own context, delivery, release, next Spec — and writes the Session Handoff and ends only when the maintainer asks or the work truly stops, because nothing is left or it is blocked on the maintainer. The Phase Record keeps recording phases, and the Boundary Notice of `roundfix phase start` goes away. The authorization is unchanged, "Mesmas do BE" of 2026-11-16, which grants implement, commit, push, pull request, merge and release right after the Spec, up to three corrections, and the Makefile, CI, `go.mod` and lint only if a correction requires one. Adjusting the maintainer's own instruction file outside the repository is the operator's work, not this Spec's.
