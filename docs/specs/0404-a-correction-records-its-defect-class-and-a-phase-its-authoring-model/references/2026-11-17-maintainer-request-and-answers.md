# The maintainer's request and answers for Spec 0404 (2026-11-17)

Recorded by the operator on 2026-11-17 on the Mac mini, in files outside the repository, and copied here by the author of the Idea. The request is copied whole. Of the operator's answers file, only the lines about Spec 0404 are copied, as written; one line about the authoring models is quoted in part, because the rest of it holds an e-mail address.

## The request, as the operator recorded it

The Idea-stage judgment read this text.

> Roundfix records why each delivery correction happened and which model authored each phase, so the maintainer can compare authoring models over the next weeks.
>
> The maintainer asked on 2026-11-17 to evaluate whether Fable is effective for Ideas, PRDs and TechSpecs compared with Opus. A one-off study found 4 Fable-authored and 14 Opus-authored Specs, no measurable difference in plan defects, and Fable costing about 5x more tokens; the defect cause of each correction had to be reconstructed by hand from memory notes and commits, and the Phase Record (since Spec 0397) does not store the model or effort.
>
> Wanted: (1) each correction made during a Spec delivery (a fix commit on the item branch, a review disposition, a QA repair) carries a defect class — plan, implementation, environment or operator — with a one-line reason, recorded where Roundfix can read it (for example a `Defect-Class:` commit trailer validated by Roundfix, or the review disposition); (2) the Phase Record entry stores the model and reasoning effort of the session or Author Subagent that ran the phase; (3) `roundfix runs causes` (or a sibling read-only view) counts corrections per Spec by defect class and by the authoring model of its TechSpec and Tasks, so a `git log`-free report answers "plan-defect rate per authoring model". Decision of the maintainer: "Spec já", authorization "Mesmas do BE".

## The Idea-stage judgment

`roundfix spec judge --stage idea --request`, on 2026-11-17: indicated. `shape_open` yes (0.81), `unverified_assumption` no (0.22), `market_question` no (0.07). Judge: `typesafe/jev-1.13-20260917`.

## The answers, as the operator recorded them

The six answers to the questions of the Idea came through AskUserQuestion at about 14:45 BRT. The approval of the draft and two more answers, on the minimum sample and on this repository, came at about 15:30 BRT, and are the last line. Each quoted label is the option the maintainer chose; the answer to Q4 is his own text.

```text
- Registro da causa de cada correção → "Spec já".
- Spec 0404 (Roundfix grava a causa de cada correção — plan/implementation/environment/operator — e modelo/esforço de cada fase no Phase Record; runs causes conta por causa e por modelo de autoria) → "Mesmas do BE (Recommended)": push, PR, merge, entra na release única se pronta a tempo; até 5 correções; autoria e entrega depois da 0403.
- 0404 Ideia (judge: indicated, shape_open 0.81, unverified_assumption 0.22, market_question 0.07; typesafe/jev-1.13-20260917). Respostas (AskUserQuestion ~14:45 BRT):
  - Q1 formato → "Completo e honesto (Recommended)": três partes, a visão mostra nº de Specs, tipo e tamanho ao lado de cada taxa e 'amostra insuficiente' abaixo de um mínimo.
  - Q3 quem declara → "Quem corrige (Recommended)": quem corrige declara classe e motivo; guarda quem declarou; maintainer pode redeclarar; vale a última.
  - Q6 comparação → "Comparar o que vier": sem alternar modelo; a visão mostra modelo por tipo e pode ficar inconclusiva.
  - Q2 unidade → "Uma por causa (Recommended)".
  - Q4 obrigatória → "Só onde é ligado e por default é desligado, nos projetos da gesttione vai ser configurado para ligar".
  - Q5 artefato → "Sim (Recommended)": defeito de plano nomeia PRD, TechSpec ou Tasks e conta contra o modelo daquela fase.
- 0404 rascunho da Ideia (~15:30 BRT): "Aprovar e seguir (Recommended)"; amostra mínima → "20 Specs (Recommended)"; Roundfix → "Ligar e manter unclear (Recommended)".
```

On the authoring models, the same file records that the maintainer keeps the recommended declaration, Fable on the Idea, the PRD and the TechSpec, for now: "[...] Então podemos manter, vai ser interessante para saber o consumo até lá."
