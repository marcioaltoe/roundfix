# Maintainer observation and answers for Spec 0391 (AskUserQuestion via `/roundfix --follow-up`, 2026-11-14 ~06:45–06:58 BRT)

Quote these verbatim in `_authorization.md` and in the PRD `## Decisions`. The source of this Spec is the maintainer himself, in this session; no ai-memory message is adopted.

## The observation

The maintainer wrote it as free text added to his answer to another question (the scope of Spec 0392, whose option text ended with "Cabe em 4 Tasks + QA."):

"E também uma observação sobre a regra de 4 tasks que talvez tenha que virar uma nova redação: O objetivo da spec com limite de 4 tasks era porque tinhamos specs gigantes no inicio com varios contextos na mesma spec e tambem specs com tasks pequenas que consumiam runs com grandes tempos de verificação (rotina normal) e faziam a spec demorar muito mais do que o suficiente."

## The rule decision

Question: "Sobre a sua observação: hoje a regra diz "agrupe fontes numa Spec só enquanto ela couber em quatro Tasks de implementação mais o QA" (Baseline `spec-routing`, skills write-prd e write-techspec, termo Grouping Bound do CONTEXT.md, ADR-0208) e não diz o motivo. Você lembrou os dois motivos: evitar Spec gigante com vários contextos e evitar Task pequena que gasta um ciclo inteiro de verificação. Como a regra deve ficar?"

Options offered: "Manter 4 e dizer o motivo (Recommended)" — "O limite de 4 Tasks + QA continua. A regra passa a dizer os dois motivos e ganha a orientação de não criar Task ou Spec pequena demais para o custo da verificação: mudança pequena do mesmo contexto entra numa Task vizinha. Vira Spec de redação neste ciclo."; "Trocar o número pelo critério" — "Sai o limite fixo. A Spec passa a ser medida por um contexto só e por Tasks que justifiquem cada uma o seu ciclo de verificação; o 4 vira referência, não teto. Muda o ADR-0208 e o termo do glossário. Vira Spec neste ciclo."; "Backlog Entry" — "Registro a sua observação como Backlog Entry, com os dois motivos, e a regra fica como está até uma Spec futura."

Answer: "Trocar o número pelo critério".

## Authorization

Question: "Se a nova redação da regra das 4 Tasks virar Spec neste ciclo (primeira pergunta), qual autorização vale para ela? Ela muda um módulo do Baseline e duas skills de autoria, que os adotantes recebem no próximo `baseline update`."

Answer: "Mesmas do BE (Recommended)" — option text: "Push, PR e merge; release patch depois dela; até 3 correções; Makefile, CI, go.mod e lint só quando uma correção exigir, descrito no relatório."
