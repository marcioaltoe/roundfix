# Respostas do Marcio — retomada do ciclo BG, 18/11/2026 ~23:10 BRT (AskUserQuestion), rótulos verbatim

Origem: Marcio colou na sessão capturas das sessões do Fluxus e do Fiscus (18/11 ~22:55–23:06 BRT). O Fiscus registrou: "profiles check --apply e profiles configure ignoram esse prazo [defaults.profile_proof_setup_timeout], embora a referência diga que o configure o respeita"; carga da máquina 13–23, cada prova Codex ~45 s, estoura os 30 s padrão; o Fiscus editou o modelo no .roundfixrc.yml à mão e provou com profiles validate (que respeita). O Fluxus subiu o prazo para 3m temporariamente e a adoção por profiles check --apply estourou mesmo assim.

Medido pelo operador (leitura do código em main d1c43f9e):
- internal/cli/profiles_configure.go:90-96 prova com roundconfig.Config{Profiles: proofProfiles}: Defaults zerado → SetupTimeout 0 → 30 s padrão (profiles_validate.go:264-267). profiles check --apply passa por writeProfilesConfiguration (profiles_check_apply.go).
- internal/agent/fallback.go:58,64 (probeFallbackCandidate) usa a constante acpxPreflightSetupTimeout (30 s), nunca o valor configurado.
- profiles validate (profiles_validate.go:223), doctor, review e preflight passam config.Defaults.ProfileProofSetupTimeout.

Pergunta: "Defeito confirmado ... Como tratamos?" → "Spec fix 0411 após a 0408 (Recommended)" — TechSpec no Opus, sem Ideia, "Mesmas do BE", release patch; entra na fila logo depois da 0408 e antes da 0409.

## Segunda fonte — Conexus (ai-memory 01a1f18f-94b3-7b32-b903-0ab7b9c3745a, 2026-11-19T02:08Z, nature measured), consumida pelo operador
Assunto: "Generated roundfix-researcher subagent has no web tools". Corpo verbatim:
> Roundfix 0.148.0: `roundfix skills install --target project` generates `.claude/agents/roundfix-researcher.md` from `skills.authoring.claude.research` with `tools: Read, Grep, Glob`, but its own body says it reads "the codebase, the knowledge base and the web" and returns URLs. Without WebFetch/WebSearch (or an MCP search tool) it cannot do the web sweeps the authoring phases delegate to it. A repository cannot fix it locally: adding the tools by hand makes `roundfix update check` report "differs from the declaration" and the next install overwrites it, and `skills.authoring` has no key to declare tools. Found by the codex pre-PR review in the conexus repo (useconexus/conexus), which dismissed the finding as a generator defect. Suggested fix: include WebFetch and WebSearch in the generated researcher's tools, or let `skills.authoring.<harness>.research` declare them. No reply needed.
Evidence: conexus PR review 2026-11-18, finding F1 at .claude/agents/roundfix-researcher.md:8; `roundfix update check` → authoringModels.problems: ".claude/agents/roundfix-researcher.md differs from the declaration".

Medido pelo operador: a chave `skills.authoring.claude.research.tools` EXISTE (este repo a usa: .roundfixrc.yml:293-296; skills/authoring.go:163 junta Read, Grep, Glob + entry.Tools). O defeito real: a Recommended Authoring Declaration (internal/config/authoring_recommendation.go:13, comentário :17 "It declares no research tools", RecommendedAuthoringYAML :30 "without a tools key") não traz ferramentas, e o subagente gerado diz que lê a web. A afirmação "no key to declare tools" do Conexus é falsa (a chave não é descobrível pelo texto que o doctor/update check mostram — verificar a documentação).

Pergunta (AskUserQuestion ~23:35 BRT): "... Onde entra?" → "Juntar à Spec 0411 (Recommended)" — dois defeitos que os adotantes acharam ao adotar as recomendações da 0.148.0; uma Spec só, mesmo release patch, logo após a 0408.
Pergunta: "Mando ao Conexus o contorno?" → "Mandar o contorno (Recommended)".

## Decisões do plano 0411 (AskUserQuestion, 19/11/2026 ~00:05 BRT), rótulos verbatim
- Release: "Patch, como autorizado (Recommended)".
- Aviso do doctor/update check a adotantes sem `research.tools`: "Sim, avisar com o comando (Recommended)" — nada reescrito sozinho; `tools: []` cala de propósito.
- Recusa nova de profiles configure com YAML quebrado no outro escopo: "Aceitar, com mensagem clara (Recommended)" — a mensagem nomeia o arquivo e a linha.
- Backlog Entries (multiSelect): "(c) bug do spec check, (a) remover código morto, (b) prazo da análise do baseline" — as três.
