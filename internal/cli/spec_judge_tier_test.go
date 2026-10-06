// Suite: advisory model tier command
// Invariant: suggestions and logs are advisory and never mutate Spec files.
// Boundary IN: public CLI with disposable repository and Home
// Boundary OUT: HTTP through the existing fake transport
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/judge"
)

func specJudgeTierFixture(t *testing.T, key string) (commandEnvironment, *specJudgeTransport) {
	t.Helper()
	env, fake := specJudgeFixture(t, key)
	dir := filepath.Join(env.workDir, "docs/specs/0300-example")
	mustWrite(t, filepath.Join(dir, "_tasks.md"), "---\ngraph:\n  nodes:\n    - file: task_01.md\n    - file: task_02.md\n    - file: task_03.md\n---\n")
	for i, kind := range []string{"backend", "docs", "qa"} {
		mustWrite(t, filepath.Join(dir, fmt.Sprintf("task_%02d.md", i+1)), fmt.Sprintf("---\ntask: task_%02d\nstatus: pending\ncomplexity: medium\ntype: %s\n---\n# Task\n## Overview\n%s\n## Verification\n- `true`\n## Result\nPrivate result.\n", i+1, kind, specJudgeEnglish))
	}
	return env, fake
}

func specJudgeTierSnapshot(t *testing.T, dir string) map[string]string {
	t.Helper()
	result := map[string]string{}
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		result[path] = string(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func TestSpecJudgeSuggestsModelTier(t *testing.T) {
	for _, format := range []string{"text", "json"} {
		t.Run(format, func(t *testing.T) {
			env, fake := specJudgeTierFixture(t, "ROUNDFIX_OPENROUTER_API_KEY")
			dir := filepath.Join(env.workDir, "docs/specs/0300-example")
			before := specJudgeTierSnapshot(t, dir)
			var out, stderr bytes.Buffer
			code := runWithContext(context.Background(), []string{"spec", "judge", "0300-example", "--stage", "tasks", "--format", format}, &out, &stderr, env)
			if code != 0 || stderr.Len() != 0 || fake.calls != 2 {
				t.Fatalf("exit=%d calls=%d stderr=%s", code, fake.calls, &stderr)
			}
			if format == "text" {
				want := "suggested model-tier docs/specs/0300-example/task_01.md: light at confidence 0.93\nsuggested model-tier docs/specs/0300-example/task_02.md: light at confidence 0.93\n" + fmt.Sprintf("Judge: 0 advisory, 2 suggested, 0 clear, 0 skipped; 2 call(s), 1684 input tokens, US$0.0001; month US$0.0001 of US$%.2f; model jev-1.13 via openrouter on ROUNDFIX_OPENROUTER_API_KEY\n", specJudgeCeiling(t))
				if out.String() != want {
					t.Fatalf("stdout=%q want=%q", out.String(), want)
				}
			} else {
				var report judge.Report
				if err := json.Unmarshal(out.Bytes(), &report); err != nil {
					t.Fatal(err)
				}
				if len(report.Judgments) != 2 {
					t.Fatalf("report=%+v", report)
				}
				for _, j := range report.Judgments {
					if j.Kind != "model-tier" || j.Outcome != "suggested" || j.Answer == nil || *j.Answer != "light" {
						t.Fatalf("judgment=%+v", j)
					}
				}
			}
			if !reflect.DeepEqual(before, specJudgeTierSnapshot(t, dir)) {
				t.Fatal("Spec directory changed")
			}
			data, err := os.ReadFile(filepath.Join(env.homeDir, ".roundfix/judge/2026-10.jsonl"))
			if err != nil {
				t.Fatal(err)
			}
			lines := strings.Split(strings.TrimSpace(string(data)), "\n")
			if len(lines) != 2 {
				t.Fatalf("log=%s", data)
			}
			for i, line := range lines {
				var row struct{ Judgment, Outcome, Artifact, Model string }
				// question_id uses the existing snake-case log schema.
				var fields map[string]json.RawMessage
				if err := json.Unmarshal([]byte(line), &row); err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal([]byte(line), &fields); err != nil {
					t.Fatal(err)
				}
				q, err := judge.Load()
				if err != nil {
					t.Fatal(err)
				}
				if row.Judgment != "model-tier" || row.Outcome != "suggested" || row.Artifact != fmt.Sprintf("docs/specs/0300-example/task_%02d.md", i+1) || string(fields["question_id"]) != fmt.Sprintf("%q", q.ModelTier.QuestionID) || row.Model == "" {
					t.Fatalf("row=%s", line)
				}
			}
		})
	}
}

func TestSpecJudgeTasksSkipsWithoutKey(t *testing.T) {
	env, fake := specJudgeTierFixture(t, "")
	specJudgeRun(t, env, []string{"0300-example", "--stage=tasks"}, "Judge: skipped: ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set (nor ROUNDFIX_OPENROUTER_API_KEY, nor ROUNDFIX_TYPESAFE_API_KEY); 2 judgment(s) not asked\n", "", 0)
	if fake.calls != 0 {
		t.Fatal("asked without key")
	}
	if _, err := os.Stat(filepath.Join(env.homeDir, ".roundfix/judge")); !os.IsNotExist(err) {
		t.Fatalf("unexpected Judge Log: %v", err)
	}
}

func TestSpecJudgeTasksRequiresGraph(t *testing.T) {
	env, fake := specJudgeFixture(t, "ROUNDFIX_OPENROUTER_API_KEY")
	var out, stderr bytes.Buffer
	code := runWithContext(context.Background(), []string{"spec", "judge", "0300-example", "--stage=tasks"}, &out, &stderr, env)
	if code != 2 || out.Len() != 0 || fake.calls != 0 || !strings.Contains(stderr.String(), "read Spec Task Graph") {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &stderr)
	}
}

func TestSpecJudgeDefaultIgnoresTasks(t *testing.T) {
	baseline, baseFake := specJudgeFixture(t, "ROUNDFIX_OPENROUTER_API_KEY")
	withTasks, taskFake := specJudgeTierFixture(t, "ROUNDFIX_OPENROUTER_API_KEY")
	invoke := func(env commandEnvironment) string {
		t.Helper()
		var out, stderr bytes.Buffer
		code := runWithContext(context.Background(), []string{"spec", "judge", "0300-example"}, &out, &stderr, env)
		if code != 0 || stderr.Len() != 0 {
			t.Fatalf("exit=%d stderr=%s", code, &stderr)
		}
		return out.String()
	}
	before, after := invoke(baseline), invoke(withTasks)
	if before != after || baseFake.calls != 5 || taskFake.calls != 5 {
		t.Fatalf("default behavior changed: before=%q after=%q", before, after)
	}
}

func TestSpecJudgeTasksHonorsCeiling(t *testing.T) {
	env, fake := specJudgeTierFixture(t, "ROUNDFIX_OPENROUTER_API_KEY")
	dir := filepath.Join(env.homeDir, ".roundfix/judge")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "2026-10.jsonl")
	before := fmt.Sprintf("{\"cost_usd\":%g}\n", specJudgeCeiling(t))
	mustWrite(t, path, before)
	want := fmt.Sprintf("Judge: skipped: monthly ceiling reached (US$%.4f of US$%.2f); 2 judgment(s) not asked\n", specJudgeCeiling(t), specJudgeCeiling(t))
	specJudgeRun(t, env, []string{"0300-example", "--stage=tasks"}, want, "", 0)
	after, err := os.ReadFile(path)
	if err != nil || string(after) != before || fake.calls != 0 {
		t.Fatalf("ceiling: calls=%d log=%s error=%v", fake.calls, after, err)
	}
}
