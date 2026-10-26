---
status: approved
granted: 2026-10-26
action: make the Roundfix skill's --follow-up Skill Mode accept the project's ai-memory handoff and hand its next steps and open questions to the follow-up command, settle every message it handed in before the turn ends even when the mode stops early, keep --status and --unattended from accepting or consuming anything, pin the text with a skill contract test, and add the Session Handoff glossary entry
consuming: 0325-the-follow-up-mode-consumes-the-handoff-and-every-message
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/modes.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0325

On 2026-10-26, at about 17:30 BRT, the maintainer asked in chat, verbatim:
"Precisamos corrigir uma ação para o follow-up de roundfix relacionada ao
ai-memory, é necessário consumir o handoff e mensagens e não só ler."

At about 17:40 BRT the maintainer answered the operator's AskUserQuestion
call. The answers this Spec rests on are, verbatim with their option text:

- Messages: "Toda mensagem até o fim (Recomendado)" — option text: "Nenhuma
  mensagem entregue ao follow-up fica pendente quando o turno termina. A
  registrada vira Spec, Backlog Entry ou Finding e é consumida em seguida. A
  informativa é consumida e citada no relatório. A que ainda não foi tratada
  vira Backlog Entry e é consumida. Vale também quando o modo para cedo."
- Handoff: "Aceitar e usar no plano (Recomendado)" — option text: "No
  inventário, lista os handoffs abertos e aceita (consome) o mais recente do
  projeto. Os próximos passos e as perguntas abertas dele entram como itens do
  plano e aparecem no relatório. Se houver uma sessão interrompida mais nova
  que o handoff, ela é lida antes."
- Order: "Spec 0325 antes da 0323 (Recomendado)" — option text: "Afeta todos
  os adotantes que usam /roundfix follow-up. Entra logo depois da 0322 e sai
  na v0.125.0 junto com a 0321 e a 0322."

This Spec belongs to batch 3 of cycle AQ, approved on 2026-10-26 under the
same rules as batches 1 and 2. Their answers of about 15:00 BRT, verbatim,
were "Sim, pode fazer a higiene." and "Vamos seguir a sua sugestão de
sequencia para a spec e os demais itens. Talvez seja importante usar o jev
para validar." The cycle AQ rules carry the cycle AP rules unchanged: the
Governed Paths each Spec declares, under the cycle AJ grant of 2026-10-19, at
about 16:20 BRT, "Mexer em Governed Paths declarados" (option text:
"'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and "Merge squash + release + push", and the
skills authorization "considere autorizado a ajustar todas as skills se
necessário". Still not authorized: this repository's Makefile,
`.github/workflows`, `go.mod`, `go.sum`, and lint and formatter
configuration.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator chooses the release point.

The governed set was measured on the authoring branch `docs/plan-0325`, based
on `feeffa28`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `interface:`, every planning
  file of this Spec, its adopted Backlog Entry copy, ADR-0338 and
  `docs/references/behavior-surfaces.json`.
- `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/modes.md` and
  `skills/roundfix/SKILL.md` are governed and listed in `paths`.

## What is not governed

- `skills/roundfix/references/modes.md`, `skills/roundfix_modes_test.go`,
  `skills/testdata/owned-skill-versions.json`.
- `CONTEXT.md` and
  `docs/adr/0338-the-follow-up-mode-claims-the-handoff-and-leaves-no-handed-in-message-pending.md`.
- `docs/references/behavior-surfaces.json`, probed although no Task changes
  it: no guide under `docs/user-guide/` changes.
- Every planning file of this Spec and its adopted Backlog Entry copy.

## Sanctioned regeneration

- `make skills-sync`, which copies the canonical `roundfix` skill to
  `skills/`.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
  outputs: `skills/testdata/owned-skill-versions.json`,
  `.agents/skills/roundfix/SKILL.md`.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Go file outside
  `skills/roundfix_modes_test.go`, any Baseline module or derived Baseline
  file, `docs/agents/`, `docs/user-guide/`, `.roundfixrc.yml`, any skill other
  than `roundfix`, or the `### QA settlement` section of any skill.
- No change to the `roundfix follow-up` command, its flags, its
  `roundfix/follow-up-messages/v1` schema or its output.
- No test, Verification command or QA row:
  - calls an ai-memory tool or reaches an ai-memory server;
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
