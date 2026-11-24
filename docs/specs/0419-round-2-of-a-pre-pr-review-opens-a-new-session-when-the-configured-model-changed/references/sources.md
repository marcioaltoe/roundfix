# Sources of Spec 0419

Sections 1 to 4 are copied verbatim from the operator's brief and its second-round addendum, which relay them as this project's ai-memory inbox received them. Sections 5 to 7 are the author's measurements and the published sources read for this Spec.

## 1. Fiscus, message 01a20a55-66b6-7b31-bd03-ed4f3d6a1fbc (2026-11-23), nature: measured

Fiscus, 23/11/2026, Spec 0108 na Delivery Queue. Depois de corrigir o achado F1 e rodar `roundfix deliver retry`, a review pré-PR (provider codex, sessão continuada: lineage.sessionOpen=true) devolveu esta resposta:

"Warning: This session was recorded with model `gpt-6-luna` but is resuming with `gpt-6.1-sol`. Consider switching back to `gpt-6-luna` as it may affect Codex performance.

No findings."

O Roundfix parou o item com "review-blocked: unclassifiable agent output: no-findings verdict does not account for other content". O veredito real é "sem achados". O texto a mais é um aviso do próprio Codex ao retomar uma sessão gravada com outro modelo, e não é conteúdo do revisor.

Sugestões: (1) remover ou ignorar a linha de aviso de troca de modelo do Codex antes de classificar a resposta; ou (2) não retomar a sessão da review quando o modelo configurado mudou, e abrir uma sessão nova. Fiz um retry para ver se repete. Se repetir, a fila do Fiscus fica parada na 0108 até a correção.

Evidence: ~/.roundfix/artifacts/2f7556b822176047/pre-pr-review/0a80d011ad299efa/pre-pr-review-answer.txt e pre-pr-review.json (outcome blocked)

The two evidence files are on the Fiscus machine and were not read for this Spec. The answer the reviewer returned is quoted whole in the message above, and that quotation is the evidence this Spec uses.

## 2. DECISÃO DE MARCIO, message 01a20a6c-7dce-7bb3-9a50-b09b27461492 from default/praetor (2026-11-23 18:48 BRT)

```text
DECISÃO DE MARCIO
pedido: 01a20a5a-ebf0-7ba2-8ac0-95a330f45eeb
pergunta: Roundfix: a fila do Fiscus está parada na Spec 0108 porque, na segunda rodada da review antes do PR, o Roundfix retoma a sessão antiga do revisor (Codex) mesmo com o modelo configurado diferente (gpt-6-luna → gpt-6.1-sol); o Codex avisa a troca e o classificador recusa a resposta, mesmo "No findings". Abrir uma Spec curta de correção? Opções: antes da 0418; depois da 0418; não, o Fiscus contorna.
resposta: "Sim, antes da 0418 (Recomendado)"
alcance: repositório roundfix, Spec nova 0419: PRD, TechSpec, Tasks, implementação, PR, merge e release patch v0.153.1 (a 0418 passa a v0.153.2). Código em internal/cli/review*.go e testes; sem Makefile, CI, go.mod, lint nem skill protegida — se precisar de caminho protegido, novo PEDIDO DE DECISÃO. Avise o Fiscus (default/gesttione-solutions-fiscus) quando o release sair, para ele dar o retry da 0108.
quando: 23/11/2026 18:48 BRT
```

## 3. DECISÃO DE MARCIO, message 01a20ac2-3453-71c1-8227-b02e7cb2c905 from default/praetor (2026-11-23 20:33 BRT)

```text
DECISÃO DE MARCIO
pedido: 01a20a8f-f044-7402-ab98-56f7e59b3420
pergunta: Roundfix, Spec 0419 (a revisão abre sessão nova quando o modelo muda): a skill descreve a retomada da sessão e a regra do repositório exige mudá-la no mesmo PR. Três arquivos são protegidos (references/review.md e as duas cópias do SKILL.md). Pode mudar? Opções: sim, os 3 arquivos com a regeneração; não, só o código.
resposta: "Sim, os 3 arquivos (Recomendado)"
alcance: repositório roundfix, Spec 0419, até o release patch v0.153.1. Caminhos: .agents/skills/roundfix/references/review.md, .agents/skills/roundfix/SKILL.md, skills/roundfix/SKILL.md. Regeneração: make skills-sync; go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions; go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage. Sem Makefile, CI, go.mod ou lint. Registre no _authorization.md citando o id desta mensagem.
quando: 23/11/2026 20:33 BRT
```

## 4. DECISÃO DE MARCIO, message 01a20ace-16f5-7c32-b81d-7e73d0340f02 from default/praetor (2026-11-23 20:46 BRT)

```text
DECISÃO DE MARCIO
pedido: 01a20a93-6180-76a2-b59f-20ee9895af31
pergunta: O plano da 0418 está pronto e pode rodar esta noite (6 a 10 h); a 0419 acabou de receber a autorização da skill (01a20ac2-3453) mas ainda precisa ser planejada, e o Fiscus sai com um segundo retry. Inverter a ordem e começar a 0418 hoje?
resposta: "Sim, 0418 hoje (Recomendado)"
alcance: repositório roundfix: merge do PR #946 (plano da 0418) e roundfix deliver start da 0418 hoje; a 0418 sai como v0.153.1 e a 0419 como v0.153.2, entrando na fila logo que planejada. Nada mais muda nas autorizações das duas Specs (0418: 01a20a59-a318; 0419: 01a20a6c-7dce e 01a20ac2-3453).
quando: 23/11/2026 20:46 BRT
```

## 5. Characterization on the unchanged tree (the author, 2026-11-23)

Measured at `fe247d4c` (v0.153.0). `50143e74`, the commit this Spec is authored against, changes no file under `internal`, `skills`, `.agents`, `docs/user-guide`, `docs/references` or `CONTEXT.md` since then.

The probe is a Go test laid over `internal/cli` with `go test -overlay`, from a file outside the repository, so nothing is written to the tree. It uses the review tests' own fixture (`newReviewCommandFixture`, `runLineageReview`, `commitLineageFile`, `disposeLineageFinding`) and a fake runner that emulates what the Fiscus measured of Codex: a session name prepared again with another model puts the warning of section 1 before its answer.

Steps: the review profile's preferred selection names `gpt-6-luna`; round 1 answers with one finding; a commit fixes it and `roundfix review dispose --fixed-by` records the disposition; the profile's preferred selection changes to `gpt-6.1-sol`; `roundfix review` runs again.

Result on the unchanged tree:

```text
round 1: exit=1 outcome=findings sessionOpen=true
round 2: exit=2 outcome=blocked reason="unclassifiable agent output: no-findings verdict does not account for other content"
         continued=true sessionOpen=false
         stderr="roundfix: review blocked: unclassifiable agent output: no-findings verdict does not account for other content"
prepare[0] session=<round-1 name> model=gpt-6-luna
prepare[1] session=<round-1 name> model=gpt-6.1-sol
events: prepare:<round-1 name>, prepare:<round-1 name>, end:<round-1 name>
```

Round 2 prepares the round-1 session name with the new model, which is the state the Fiscus reported. The blocked round 2 ends its session, so the next `roundfix review` at the same head opens a session with a new name and ends `exit=0 outcome=reviewed`: on v0.153.0 the defect costs one blocked round for each model change. The measurement is of the `roundfix review` command with a fake runner; `roundfix deliver retry` and the real acpx adapter were not exercised.

The probe's fake runner is the one listed in section 6, under another name.

## 6. The prototype (the author, 2026-11-23)

Three changed lines, laid over the same tree with `go test -overlay`:

- `reviewLineage` gains `Model string` with the tag `json:"model,omitempty"` in `internal/cli/review_lineage.go`;
- the block of `internal/cli/review.go` that records the session of a round whose prompt was sent also records `specContext.runtime.Model`;
- the condition that lets round 2 prepare the recorded session name also requires `prior.Model != "" && prior.Model == strings.TrimSpace(selections[selection].Model)`.

Result of the same probe:

```text
round 1: exit=1 outcome=findings sessionOpen=true model=gpt-6-luna
round 2: exit=0 outcome=reviewed continued=false sessionOpen=false model=gpt-6.1-sol stderr=""
prepare[0] session=<round-1 name> model=gpt-6-luna
prepare[1] session=<new name> model=gpt-6.1-sol
events: prepare:<round-1 name>, end:<round-1 name>, prepare:<new name>, end:<new name>
```

Every test of `internal/cli` and of `internal/delivery` was run at `50143e74` with the prototype laid over the tree (`go test -count=1 -p 2 -overlay <file> ./internal/cli ./internal/delivery`). Every test passed: `internal/delivery` printed `ok`, and `internal/cli` printed `PASS` for its tests. The `internal/cli` run then ended `FAIL` on the suite guard alone, which named four created paths, all of them this Spec's own planning files, written by the author while the tests ran. So the prototype moves no existing test of those two packages. The run was not repeated on a quiet tree.

Other facts read from the code:

- `agent.RuntimeFor` copies the selection's model into the runtime without changing it (`internal/agent/agent.go`), and `runtimeForProfileSelection` trims it, so the recorded model and the compared model are the same string.
- `readReviewRecord` decodes with `json.Unmarshal` and no `DisallowUnknownFields`, so a Roundfix released before this Spec reads a record that carries the new field.
- Config validation refuses an Agent Selection without a model, so a record whose lineage holds no model was written by an earlier Roundfix or by a round that sent no prompt.
- `decideReviewLineage` already starts a new lineage at round 1 when the provider changed, and `validateReviewProfileProvider` holds every selection of the review profile to the provider's runtime, so a round 2 never resumes on another runtime.
- The one test that pins the bytes of a lineage, in `internal/cli/review_lineage_selection_test.go`, pins a `ceiling-closed` record, which sends no prompt and so carries no model.

### The planned tests, probed

A second probe wrote the three tests the TechSpec plans for `internal/cli`, under probe names, and ran them twice at `50143e74` with `go test -overlay`, from files outside the repository.

- Over the prototype: the three tests and their three subtests pass, and the package ends `ok`, the suite guard included.
- On the unchanged tree: the three tests fail. The first stops at round 1, whose printed lineage holds no model. The second fails because round 2 prepares the round-1 session name and records `continued: true`. The third fails in both subtests because the printed lineage holds no model. The subtest about the classifier did not run there, because its parent had already failed; the reason it asserts is the one section 5 measured on the unchanged tree for the same answer.

The built-in review profile of the fixture names `gpt-6-luna` as its preferred selection and `gpt-6.1-sol` as its fallback, which the third test's fallback case shows.

The probe reads `lineage.model` from the JSON the command prints, not from a Go field, so the same file compiles on both trees and fails by assertion on the unchanged one. Its source, as it ran; the Task names its tests as the TechSpec does, not as this file does:

```go
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/runevent"
)

// Probe of the tests the TechSpec of Spec 0419 plans. It reads lineage.model
// from the JSON the command prints, so it compiles on the unchanged tree too.

type modelNoticeReviewRunner struct {
	sessionReviewRunner
	opened map[string]string
}

func (r *modelNoticeReviewRunner) PrepareSession(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) error {
	if r.opened == nil {
		r.opened = map[string]string{}
	}
	if _, ok := r.opened[req.Session.Name]; !ok {
		r.opened[req.Session.Name] = req.Runtime.Model
	}
	return r.sessionReviewRunner.PrepareSession(ctx, req, sink)
}

func (r *modelNoticeReviewRunner) RunPrepared(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	result, err := r.sessionReviewRunner.RunPrepared(ctx, req, sink)
	if was := r.opened[req.Session.Name]; was != "" && was != req.Runtime.Model {
		result.Message = probeReportedNotice(was, req.Runtime.Model) + result.Message
	}
	return result, err
}

func probeReportedNotice(was, now string) string {
	return fmt.Sprintf("Warning: This session was recorded with model `%s` but is resuming with `%s`. Consider switching back to `%s` as it may affect Codex performance.\n\n", was, now, was)
}

func newModelNoticeRunner(answers ...string) *modelNoticeReviewRunner {
	r := &modelNoticeReviewRunner{}
	for _, answer := range answers {
		r.results = append(r.results, reviewCommandRunResult{result: agent.ExecuteResult{Message: answer, StopReason: "end_turn", ACPSessionID: "s-1"}})
	}
	return r
}

func writeProbeProfile(t *testing.T, fixture reviewCommandFixture, model string) {
	t.Helper()
	mustWrite(t, filepath.Join(fixture.repository, ".roundfixrc.yml"), fmt.Sprintf(`defaults:
  artifact_dir: %q
pre_pr_review:
  provider: codex
profiles:
  review:
    preferred:
      runtime: codex
      model: %s
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: fallback-model
        reasoning_effort: high
`, fixture.artifactDir, model))
}

// runProbeReview returns the decoded record and the model the printed lineage holds.
func runProbeReview(t *testing.T, fixture reviewCommandFixture) (int, reviewRecord, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"review", "--base", fixture.baseCommit}, &stdout, &stderr)
	var record reviewRecord
	if err := json.Unmarshal(stdout.Bytes(), &record); err != nil {
		t.Fatalf("decode %q: %v (stderr %q)", stdout.String(), err, stderr.String())
	}
	var raw struct {
		Lineage map[string]any `json:"lineage"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	model, _ := raw.Lineage["model"].(string)
	return code, record, model, stderr.String()
}

func probeRecordPath(fixture reviewCommandFixture) string {
	return filepath.Join(reviewCheckoutDir(fixture.artifactDir, fixture.repository), reviewRecordFileName)
}

func TestProbeRoundTwoOpensANewSessionWhenTheConfiguredModelChanged(t *testing.T) {
	r := newModelNoticeRunner(sessionReviewFinding, "No findings.")
	fixture := newReviewCommandFixture(t, "codex", r)
	writeProbeProfile(t, fixture, "gpt-6-luna")
	code1, first, model1, _ := runProbeReview(t, fixture)
	if code1 != exitRunFailed || !first.Lineage.SessionOpen || model1 != "gpt-6-luna" {
		t.Fatalf("round 1 exit=%d lineage=%+v model=%q", code1, first.Lineage, model1)
	}
	head := commitLineageFile(t, fixture, "delta.txt", "round two correction\n")
	disposeLineageFinding(t, "--fixed-by", head)
	writeProbeProfile(t, fixture, "gpt-6.1-sol")
	code2, second, model2, stderr := runProbeReview(t, fixture)
	if code2 != exitOK || second.Outcome != reviewOutcomeReviewed || second.Lineage.Round != 2 || second.Lineage.Continued || second.Lineage.SessionOpen || stderr != "" || model2 != "gpt-6.1-sol" {
		t.Fatalf("round 2 exit=%d outcome=%s reason=%q lineage=%+v model=%q stderr=%q", code2, second.Outcome, second.Reason, second.Lineage, model2, stderr)
	}
	if len(r.prepares) != 2 || r.prepares[1].Session.Name == r.prepares[0].Session.Name || r.prepares[0].Runtime.Model != "gpt-6-luna" || r.prepares[1].Runtime.Model != "gpt-6.1-sol" {
		t.Fatalf("prepares=%+v", r.prepares)
	}
	old, fresh := r.prepares[0].Session.Name, r.prepares[1].Session.Name
	if want := []string{"prepare:" + old, "end:" + old, "prepare:" + fresh, "end:" + fresh}; !reflect.DeepEqual(r.events, want) {
		t.Fatalf("events=%v want=%v", r.events, want)
	}
	if !strings.Contains(r.runs[1].Prompt, "--- BEGIN ROUND 2 DELTA DIFF ---") || !strings.Contains(r.runs[1].Prompt, "fixed by "+head) || strings.Contains(r.runs[1].Prompt, "--- BEGIN CANDIDATE DIFF ---") {
		t.Fatalf("round 2 prompt=%q", r.runs[1].Prompt)
	}
	t.Run("the_classifier_still_refuses_the_reported_answer", func(t *testing.T) {
		only := newSessionReviewRunner(probeReportedNotice("gpt-6-luna", "gpt-6.1-sol") + "No findings.")
		other := newReviewCommandFixture(t, "codex", only)
		code, record, _ := runLineageReview(t, other)
		if code != exitPreflight || record.Outcome != reviewOutcomeBlocked || record.Reason != "unclassifiable agent output: no-findings verdict does not account for other content" {
			t.Fatalf("exit=%d outcome=%s reason=%q", code, record.Outcome, record.Reason)
		}
	})
}

func TestProbeRoundTwoOpensANewSessionWhenTheRecordHoldsNoModel(t *testing.T) {
	r := newSessionReviewRunner(sessionReviewFinding, "No findings.")
	fixture := newReviewCommandFixture(t, "codex", r)
	_, first, _, _ := runProbeReview(t, fixture)
	path := probeRecordPath(fixture)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var record map[string]any
	if err := json.Unmarshal(content, &record); err != nil {
		t.Fatal(err)
	}
	delete(record["lineage"].(map[string]any), "model")
	rewritten, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, rewritten, 0o644); err != nil {
		t.Fatal(err)
	}
	commitLineageFile(t, fixture, "delta.txt", "fix\n")
	code, second, model, stderr := runProbeReview(t, fixture)
	if code != exitOK || second.Lineage.Continued || second.Lineage.SessionOpen || model == "" {
		t.Fatalf("exit=%d lineage=%+v model=%q stderr=%q", code, second.Lineage, model, stderr)
	}
	if len(r.prepares) != 2 || r.prepares[1].Session.Name == first.Lineage.Session || len(r.ends) != 2 || r.ends[0].Name != first.Lineage.Session {
		t.Fatalf("prepares=%+v ends=%+v", r.prepares, r.ends)
	}
	if want := []string{"prepare:" + first.Lineage.Session, "end:" + first.Lineage.Session, "prepare:" + r.prepares[1].Session.Name, "end:" + r.prepares[1].Session.Name}; !reflect.DeepEqual(r.events, want) {
		t.Fatalf("events=%v want=%v", r.events, want)
	}
}

func TestProbeRecordsTheModelOfTheSessionOfEachRound(t *testing.T) {
	t.Run("preferred_then_continued", func(t *testing.T) {
		r := newSessionReviewRunner(sessionReviewFinding, "No findings.")
		fixture := newReviewCommandFixture(t, "codex", r)
		_, _, model1, _ := runProbeReview(t, fixture)
		if model1 == "" || model1 != r.prepares[0].Runtime.Model {
			t.Fatalf("round 1 model=%q runtime=%q", model1, r.prepares[0].Runtime.Model)
		}
		content, err := os.ReadFile(probeRecordPath(fixture))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(content), fmt.Sprintf("%q:%q", "model", model1)) && !strings.Contains(string(content), fmt.Sprintf("%q: %q", "model", model1)) {
			t.Fatalf("record file holds no model: %s", content)
		}
		commitLineageFile(t, fixture, "delta.txt", "fix\n")
		_, second, model2, _ := runProbeReview(t, fixture)
		if !second.Lineage.Continued || model2 != model1 {
			t.Fatalf("round 2 lineage=%+v model=%q want %q", second.Lineage, model2, model1)
		}
	})
	t.Run("fallback", func(t *testing.T) {
		r := newSessionReviewRunner(sessionReviewFinding)
		r.prepareErrors = []error{&agent.SelectionFailureError{Runtime: "codex", Reason: "adapter startup"}, &agent.SelectionFailureError{Runtime: "codex", Reason: "adapter startup"}, nil}
		fixture := newReviewCommandFixture(t, "codex", r)
		_, first, model, _ := runProbeReview(t, fixture)
		if first.Lineage.Selection == nil || *first.Lineage.Selection != 1 || model == "" || model != r.prepares[2].Runtime.Model || model == r.prepares[0].Runtime.Model {
			t.Fatalf("lineage=%+v model=%q prepares=%q,%q", first.Lineage, model, r.prepares[0].Runtime.Model, r.prepares[2].Runtime.Model)
		}
	})
}
```

## 7. Published sources (read on 2026-11-23)

- openai/codex, issue 9054, "After changing model_providers, the resume command cannot view the conversation history of the previous provider", https://github.com/openai/codex/issues/9054. A contributor from OpenAI, `etraut-openai`, on 2026-01-12: "Conversations include details like CoT (chain of thought) tokens that do not transfer between models. You should stick to one model per thread. If you want to use a different model, create a new thread."
- openai/codex, commit `fe8b474acd43cb8894d661944c6b5c8db0ef0ad1` of 2026-02-05, "fix(core,app-server) resume with different model (#10719)", https://github.com/openai/codex/commit/fe8b474acd43cb8894d661944c6b5c8db0ef0ad1. Codex treats a resume with another model as a switch of model and tells the model about it; the warning the Fiscus quoted compares only the model.
- Read only by the Research Subagent in the first authoring round, not verified by the author: the Codex ACP adapter forwards a Codex warning as agent text without a message identifier (`agentclientprotocol/codex-acp`, `src/CodexEventHandler.ts`), and the Agent Client Protocol treats the model as a session option that a client may set (https://agentclientprotocol.com/protocol/session-setup).
