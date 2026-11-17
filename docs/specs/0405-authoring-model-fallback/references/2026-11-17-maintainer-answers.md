# Maintainer request and answers for Spec 0405 — 2026-11-17

The maintainer answered through the Supervisor's structured question tool. The
options he chose are quoted as their labels read; free answers are quoted in
his words.

## Request

The maintainer, 2026-11-17: "Se a cota do Fable acabar enquanto estiver
trabalhando em alguma skill, automaticamente deve ir para o Opus High ou X-High
para não perder qualidade. Talvez a gente tenha que definir o fallback para
isso. Não sei se é possível fazer isso no frontmatter da skill ou se vai ter que
ser alguma regra no baseline. Estude isso e já prepare uma spec."

Observed the same day: the TechSpec Author Subagent of Spec 0404, declared
`fable` through `skills.authoring.claude.techspec`, ended with "You've hit your
Fable limit · resets Nov 22 at 8pm" (HTTP 429 rate_limit) in the middle of the
phase, and the Supervisor restarted the work on an Opus agent after asking the
maintainer.

Idea-stage judgment: indicated (shape_open yes 0.90, unverified_assumption no
0.24, market_question no 0.09; judge typesafe/jev-1.13-20260917).

## Idea questions, about 17:50 BRT

1. Shape — "fallbacks na declaração (Recommended)".
2. This repository's fallback — "Opus xhigh / Tasks Opus high (Recommended)".
3. Ask or act — "Troca e avisa (Recommended)".
4. Trigger — "Só limite do modelo (Recommended)".
5. Persistence — "Fica no reserva na spec e volta ao modelo declarado na
   próxima spec".
6. Record — "Sim, no Phase Record (Recommended)".
7. Delivery authorization — "Mesmas do BE (Recommended)".

## Idea draft, about 18:10 BRT

- Draft — "Aprovar e seguir (Recommended)".
- Tasks, whose declared model is Opus and whose Opus fallback would escape no
  Opus limit — "Vamos manter opus high e sonnet high só para fallback ter
  algo."

## Plan review, about 18:40 BRT

- Which later phases stay on the fallback — "Só as do modelo esgotado (Recommended)".
- Whether the Recommended Authoring Declaration that the Doctor Command and the Update Check print carries the fallbacks — "Incluir os reservas (Recommended)": the Idea, the PRD and the TechSpec fall back to Opus xhigh, the Tasks (Opus high) to Sonnet high, and research declares none.

## Authoring model

Every phase of this Spec was authored on Opus 5.5 (`claude-opus-5-5`), not on
the declared Fable: Fable's quota for the account ran out on 2026-11-17 and
resets on 2026-11-22 at 20:00 BRT.
