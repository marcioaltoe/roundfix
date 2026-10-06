package judge

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/openrouterkey"
)

func TestRunPrefersTheJudgeStageKey(t *testing.T) {
	checkStageKey(t, map[string]string{openrouterkey.Judge: "stage-fake-value", openrouterkey.Shared: "shared-fake-value"}, openrouterkey.Judge)
}
func TestRunFallsBackToTheSharedKeyAndNamesIt(t *testing.T) {
	checkStageKey(t, map[string]string{openrouterkey.Shared: "shared-fake-value"}, openrouterkey.Shared)
}
func checkStageKey(t *testing.T, keys map[string]string, variable string) {
	t.Helper()
	q, req := runFixture(t)
	req.Keys = keys
	calls := 0
	expected := keys[variable]
	echo := keys[openrouterkey.Judge] + " " + keys[openrouterkey.Shared]
	req.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.String() != q.Transports[0].Endpoint || r.Header.Get("Authorization") != "Bearer "+expected {
			t.Fatal("wrong recipient or bearer key")
		}
		delete(req.Keys, variable)
		body := strings.ReplaceAll(answerBody("jev-1.13.0", 0.8, 0.29, ""), "fixture-id", echo)
		return fixtureResponse(200, body), nil
	})
	got := runChecked(t, q, req)
	if calls != 2 || got.KeyVariable == nil || *got.KeyVariable != variable || got.Transport == nil || *got.Transport != "openrouter" {
		t.Fatalf("selection report=%+v calls=%d", got, calls)
	}
	for _, row := range readLog(t, req) {
		if row.KeyVariable != variable {
			t.Fatalf("logged variable=%q", row.KeyVariable)
		}
	}
	assertStageSecretsAbsent(t, got, req)
}
func assertStageSecretsAbsent(t *testing.T, got Report, req Request) {
	t.Helper()
	report, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	log, err := os.ReadFile(newJudgeLog(req.HomeDir, req.Now()).path)
	if err != nil {
		t.Fatal(err)
	}
	for _, secret := range []string{"stage-fake-value", "shared-fake-value"} {
		if strings.Contains(string(report), secret) || strings.Contains(string(log), secret) {
			t.Fatal("key value in report or log")
		}
	}
}
func TestJudgeLogNamesTheKeyVariableAndNoKey(t *testing.T) {
	q, req := runFixture(t)
	want := []string{openrouterkey.Judge, openrouterkey.Shared, "ROUNDFIX_TYPESAFE_API_KEY"}
	if !reflect.DeepEqual(q.KeyVariables(), want) {
		t.Fatalf("variables=%v", q.KeyVariables())
	}
	duplicate := q
	duplicate.Transports = append(append([]Transport(nil), q.Transports...), q.Transports[0])
	if !reflect.DeepEqual(duplicate.KeyVariables(), want) {
		t.Fatal("duplicate variable retained")
	}
	req.Keys = map[string]string{openrouterkey.Judge: "stage-fake-value", openrouterkey.Shared: "shared-fake-value"}
	req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("stage-fake-value shared-fake-value")
	})
	got := runChecked(t, q, req)
	assertStageSecretsAbsent(t, got, req)
	rows := readLog(t, req)
	if len(rows) != 1 || rows[0].KeyVariable != openrouterkey.Judge {
		t.Fatalf("rows=%+v", rows)
	}
	// Historical records without key_variable still contribute to the ceiling.
	q, req = runFixture(t)
	writeFixture(t, req.HomeDir, ".roundfix/judge/2026-10.jsonl", `{"cost_usd":1.25}`+"\n")
	log := newJudgeLog(req.HomeDir, req.Now())
	if cost, err := log.monthCost(); err != nil || cost != 1.25 {
		t.Fatalf("legacy cost=%v err=%v", cost, err)
	}
	req.Keys = nil
	req.Transport = neverRequest(t)
	got = runChecked(t, q, req)
	b, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	if got.KeyVariable != nil || !strings.Contains(string(b), `"key_variable":null`) {
		t.Fatalf("no-key report=%s", b)
	}
}
