// Suite: repository verification planning and timing.
// Invariant: selection and reporting preserve the documented boundaries.
// Boundary IN: in-process functions and isolated fixtures.
// Boundary OUT: live delivery, GitHub and the Run Database.
package gatetime

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNumberAddsTheAllowanceAndNeverPassesTheTarget(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		m, a, n int
		ok      bool
	}{{104, 120, 120, true}, {212, 240, 240, true}, {31, 120, 40, true}, {121, 120, 0, false}, {0, 120, 0, true}} {
		n, ok := Number(c.m, c.a)
		if n != c.n || ok != c.ok {
			t.Fatalf("%+v: %d %t", c, n, ok)
		}
	}
}
func TestValidateRefusesAChangedTargetAndAWrongNumber(t *testing.T) {
	t.Parallel()
	doc, e := Load("testdata/targets-with-rule.json")
	if e != nil {
		t.Fatal(e)
	}
	approved := map[string]Limits{}
	for _, g := range doc.Gates {
		approved[g.ID] = g.Approved
	}
	if errs := Validate(doc, approved); len(errs) > 0 {
		t.Fatal(errs)
	}
	doc.Gates[0].Approved.MedianSeconds++
	if len(Validate(doc, approved)) == 0 {
		t.Fatal("changed target accepted")
	}
	doc.Gates[0].Approved = approved[doc.Gates[0].ID]
	doc.Gates[0].Number.P95Seconds = 241
	if len(Validate(doc, approved)) == 0 {
		t.Fatal("wrong number accepted")
	}
}
func transcript(t *testing.T, args []string, want string, code int) {
	t.Helper()
	var out, err bytes.Buffer
	got := Run(args, &out, &err)
	if got != code || out.String() != want || err.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q; want %d %q", got, out.String(), err.String(), code, want)
	}
}
func reportTranscript(t *testing.T, args []string, want string) {
	t.Helper()
	transcript(t, append([]string{"report", "-targets", "../../docs/references/gate-time-targets.json"}, args...), want+"\n", 0)
}
func TestReportRefusesAnUnknownGate(t *testing.T) {
	t.Parallel()
	var out, err bytes.Buffer
	if code := Run([]string{"report", "-targets", "../../docs/references/gate-time-targets.json", "-gate", "unknown", "-seconds", "1"}, &out, &err); code != 2 || out.Len() != 0 || !strings.Contains(err.String(), "unknown gate") {
		t.Fatalf("%d %s %s", code, &out, &err)
	}
}
func TestTheRunDatabaseIsOpenedForReadingOnly(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "runs.db")
	writer, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = writer.Exec("CREATE TABLE proof (value INTEGER)"); err != nil {
		t.Fatal(err)
	}
	if err = writer.Close(); err != nil {
		t.Fatal(err)
	}
	reader, err := openRunDatabase(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	if _, err = reader.Exec("INSERT INTO proof VALUES (1)"); err == nil {
		t.Fatal("read-only connection accepted a write")
	}
	dsn := readOnlyDSN("/tmp/a b.db")
	u, e := url.Parse(dsn)
	if e != nil || u.Query().Get("mode") != "ro" || u.Path != "/tmp/a b.db" {
		t.Fatal(dsn, e)
	}
}
func TestTestsListsTheSlowestTestsOfARun(t *testing.T) {
	t.Parallel()
	var input, out bytes.Buffer
	for _, e := range []testEvent{{Action: "pass", Package: "p", Test: "TestFast", Elapsed: 1}, {Action: "output", Package: "p", Test: "TestSlow", Output: "failure details\n"}, {Action: "fail", Package: "p", Test: "TestSlow", Elapsed: 5}, {Action: "fail", Package: "p", Elapsed: 6}} {
		if err := json.NewEncoder(&input).Encode(e); err != nil {
			t.Fatal(err)
		}
	}
	seconds, err := renderTests(&input, &out, 1)
	want := "failure details\nfail p 6.00s\nslowest tests:\n5.00s p/TestSlow\n"
	if err != nil || seconds != 6 || out.String() != want {
		t.Fatalf("%d %q %v", seconds, out.String(), err)
	}
}
func TestSurfaceTranscriptTheGateTimeReportOfATaskSettlement(t *testing.T) {
	t.Parallel()
	reportTranscript(t, []string{"-gate", "task-settlement", "-seconds", "31"}, "gate time: task settlement 31s (target: median 120s, p95 240s)")
}
func TestSurfaceTranscriptTheGateTimeReportOfTheCompleteSuiteRun(t *testing.T) {
	t.Parallel()
	reportTranscript(t, []string{"-gate", "complete-suite-run", "-seconds", "352"}, "gate time: complete suite run 352s (target: 360s)")
}
func TestSurfaceTranscriptTheGateTimeReportOfAVerificationJob(t *testing.T) {
	t.Parallel()
	reportTranscript(t, []string{"-gate", "pull-request", "-job", "cli-2", "-seconds", "412"}, "gate time: pull request, job cli-2 412s (target: median 480s, p95 600s)")
}
func TestSurfaceTranscriptTheGateTimeReportOnceARuleIsRecorded(t *testing.T) {
	t.Parallel()
	transcript(t, []string{"report", "-targets", "testdata/targets-with-rule.json", "-gate", "task-settlement", "-seconds", "31"}, "gate time: task settlement 31s (number: median 120s, p95 240s; within the measured result)\n", 0)
}
func TestSurfaceTranscriptTheWindowsFromRecordedRunsAndSettlements(t *testing.T) {
	t.Parallel()
	transcript(t, []string{"windows", "-targets", openWindowsTargets(t), "-fixture", "testdata/windows-fixture.json"}, "gate times: task settlement, 30 of 30 passing: median 104s of 120s, p95 212s of 240s, within\ngate times: pull request, 14 of 20 passing: window open\n", 0)
}

// openWindowsTargets copies the live targets with every window start cleared,
// so the recorded fixture is read whatever start the live document records.
func openWindowsTargets(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile("../../docs/references/gate-time-targets.json")
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err = json.Unmarshal(data, &doc); err != nil {
		t.Fatal(err)
	}
	for _, gate := range doc["gates"].([]any) {
		gate.(map[string]any)["window"].(map[string]any)["since"] = nil
	}
	data, err = json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "targets.json")
	if err = os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunsInWindowKeepsTheOldestRunsSinceTheStart(t *testing.T) {
	t.Parallel()
	start := time.Date(2026, 11, 26, 10, 58, 57, 0, time.UTC)
	gate := &Gate{Window: Window{Size: 2, Since: &start}}
	runs := []ciRun{
		{DatabaseID: 4, CreatedAt: start.Add(3 * time.Hour)},
		{DatabaseID: 3, CreatedAt: start.Add(2 * time.Hour)},
		{DatabaseID: 2, CreatedAt: start.Add(time.Hour)},
		{DatabaseID: 1, CreatedAt: start.Add(-time.Hour)},
	}
	got := runsInWindow(runs, gate)
	if len(got) != 2 || got[0].DatabaseID != 2 || got[1].DatabaseID != 3 {
		t.Fatalf("runsInWindow = %+v, want runs 2 and 3", got)
	}
	if runsInWindow(runs, nil) != nil {
		t.Fatal("no pull request gate must keep no run")
	}
}
func TestWindowsHoldsAReleaseOnlyWithARule(t *testing.T) {
	t.Parallel()
	doc, e := Load("testdata/targets-with-rule.json")
	if e != nil {
		t.Fatal(e)
	}
	source := Source{}
	for i := 0; i < 30; i++ {
		source["task-settlement"] = append(source["task-settlement"], Sample{Seconds: 300})
	}
	_, ok := Windows(doc, source)
	if ok {
		t.Fatal("overrun accepted with rule")
	}
	doc.Rule = nil
	_, ok = Windows(doc, source)
	if !ok {
		t.Fatal("target failed without rule")
	}

}

func TestSettlementSamplesPairTaskAndAttempt(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "events.db")
	writer, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer writer.Close()
	for _, statement := range []string{"CREATE TABLE runs (id TEXT,git_root TEXT,created_at TEXT,repository_root TEXT NOT NULL DEFAULT '')", "CREATE TABLE run_events (run_id TEXT,created_at TEXT,payload TEXT,kind TEXT,cursor INTEGER)", "INSERT INTO runs VALUES ('run','/repo','2026-11-26T00:00:00Z','')",
		// A Delivery Queue Run settles in the Item Worktree; its repository_root names the checkout.
		"INSERT INTO runs VALUES ('queued','/home/.roundfix/worktrees/repo-1/item','2026-11-26T00:01:00Z','/repo')",
		"INSERT INTO runs VALUES ('other','/home/.roundfix/worktrees/other-2/item','2026-11-26T00:01:00Z','/other')"} {
		if _, err := writer.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
	for i, event := range []struct{ at, payload string }{
		{"2026-11-26T00:00:01Z", `{"command":"make verify-changed","phase":"started","task":"task_01","attempt":1}`},
		{"2026-11-26T00:00:02Z", `{"command":"make verify-changed","phase":"command-passed","task":"task_01","attempt":2}`},
		{"2026-11-26T00:00:32Z", `{"command":"make verify-changed","phase":"command-passed","task":"task_01","attempt":1}`},
	} {
		if _, err := writer.Exec("INSERT INTO run_events VALUES ('run',?,?,'daemon.verification',?)", event.at, event.payload, i); err != nil {
			t.Fatal(err)
		}
	}
	for i, run := range []string{"queued", "other"} {
		for j, event := range []struct{ at, phase string }{{"2026-11-26T00:01:10Z", "started"}, {"2026-11-26T00:03:10Z", "command-passed"}} {
			payload := `{"command":"make verify-changed","phase":"` + event.phase + `","task":"task_01","attempt":1}`
			if _, err := writer.Exec("INSERT INTO run_events VALUES (?,?,?,'daemon.verification',?)", run, event.at, payload, 10+2*i+j); err != nil {
				t.Fatal(err)
			}
		}
	}
	reader, err := openRunDatabase(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	samples, err := settlementSamples(reader, "/repo", nil)
	if err != nil || len(samples) != 2 || samples[0].Seconds != 31 || samples[1].Seconds != 120 {
		t.Fatal("the checkout's Run and its Delivery Queue Run, not another repository's", samples, err)
	}
	since := time.Date(2026, 11, 26, 0, 0, 10, 0, time.UTC)
	samples, err = settlementSamples(reader, "/repo", &since)
	if err != nil || len(samples) != 1 || samples[0].Seconds != 120 {
		t.Fatal("Run started before the window", samples, err)
	}
}
