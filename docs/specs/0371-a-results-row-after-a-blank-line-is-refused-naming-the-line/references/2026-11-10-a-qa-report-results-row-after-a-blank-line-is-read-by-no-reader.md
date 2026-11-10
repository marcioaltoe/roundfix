---
type: fix
status: promoted
created: 2026-11-10
spec: 0371-a-results-row-after-a-blank-line-is-refused-naming-the-line
reason: null
---

# Source for Spec 0371 — a QA Report's Results table broken by a blank line

Recorded by the operator on 2026-11-10 from three inputs: the Pantheon report of
2026-11-09/10 (ai-memory message 01a1c322, consumed in the previous session), the
previous session's own diagnosis (Session Handoff 01a1c38c), and the maintainer's
decision relayed by the Jarvis session (ai-memory message 01a1c4c0).

## What Pantheon reported (reported by pantheon)

Pantheon's Spec 0066 finished its QA with a report whose `## Results` section holds a
blank line in the middle of the table. `roundfix archive` / the eligibility check
answered `rows_blocked_environment is 1; expected 0`, which names the wrong cause.
Pantheon archived with an override. Pantheon first read it as a contradiction between
the qa-gate skill and the binary; it is not.

## Operator diagnosis (measured on Roundfix v0.137.7)

- `qaReportResultsCounts` (`internal/spec/qa.go`, around line 650) reads the first
  table under `## Results`: after the header and separator it reads rows until the
  first line that has no table cells, then looks for another header+separator pair.
  Rows after a blank line have no header, so they are ignored without any diagnostic.
- In the Pantheon report the ignored rows include the pre-PR Pull Request row
  (`blocked (environment: no open Pull Request)`, provenance `Pull Request row`) and
  the declared row (`blocked (declared: …)`). With those two unread, the counts the
  front matter carries no longer match what the reader derives, and eligibility fails
  with the misleading `rows_blocked_environment is 1; expected 0`.
- The table is seeded by the Daemon (`internal/speccheck/report.go`, the block that
  writes `MechanicalRowsHeading`, the `| # | Status | Provenance |` header and the
  carried and blocked rows). The QA Agent then appended its own rows after a blank
  line, in a narrower table without header.

## Shape of the report (reduced; the real file stays in Pantheon)

Read-only reproduction: `git -C ~/dev/pantheon show d2898fb:docs/specs/0066-sonda-do-indexus/qa/qa-report-2026-11-09-01.md`
(lines 44–78). Never copy Pantheon text beyond this reduced shape; fixtures stay synthetic.

```markdown
## Results

| #   | Status                                                    | Provenance   |
| --- | --------------------------------------------------------- | ------------ |
| Q01 | carried (established by: <report path>; head: <40 hex>)   | PRD Goal 1   |
| …   | (20 seeded carried rows, padded columns)                  | …            |
| Q25 | carried (established by: <report path>; head: <40 hex>)   | outside-evidence row |

| Q12 | pass | PRD Core Feature 6 |
| Q19 | pass | Task 02 Acceptance Criterion 1 |
| …   | (rows the QA Agent appended after the blank line) | … |
| Q27 | blocked (environment: no open Pull Request) | Pull Request row |
| Q28 | blocked (declared: <criterion text>) | PRD Unreachable Acceptance: <criterion text> |
| Q29 | pass | PRD Non-Goals / Out of Scope |
| Q30 | pass | PRD Core Feature 6; Task 02 Acceptance Criterion 2 |

## Mechanical skips
```

## Maintainer decision (relayed by jarvis, message 01a1c4c0, 2026-11-10 09:18 UTC)

The operator asked the maintainer through the Jarvis session (message 01a1c379):
"1 Spec 0371 agora (recomendada) / 2 backlog no próximo ciclo (default sem resposta) /
3 só ajustar a skill qa-gate". The Jarvis session answered:

> Resposta à sua msg 01a1c379 ("uma decisão para o Marcio de manhã, defeito do relatório de QA"). Marcio escolheu a opção 1 às 06:20 BRT de 10/11, por AskUserQuestion na sessão do Jarvis: abrir a Spec 0371. O leitor passa a ler a tabela de resultados inteira ou a recusar apontando a linha em branco exata, e a skill/prompt de QA passa a mandar continuar a tabela sem linha em branco. Autorizado escrever, entregar e soltar a v0.137.8, com as mesmas regras do ciclo BA. Mande os blocos SPEC PARA O QUADRO e o resumo final ao default/jarvis.

The maintainer confirmed the authorization directly in the Roundfix session on
2026-11-10 ~15:35 BRT (see the Cycle BB section of the briefing for the exact answer).
