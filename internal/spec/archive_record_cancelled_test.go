package spec

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func cancelledRecordFixture() ArchiveRecord {
	return ArchiveRecord{Spec: "demo", Title: "Demo", Created: "2026-10-10", Archived: "2026-10-12", Disposition: ArchiveCancelled, Source: "docs/specs/demo", SourceRevision: strings.Repeat("a", 40), CancelReason: "no consumer: scope changed", CompletedTasks: []string{"task_01", "task_02"}, PendingTasks: []string{"task_03"}, ReturnedSources: []string{"docs/backlog/entry.md"}, Outcome: "Original intent."}
}

func TestACancelledArchiveRecordRoundTrips(t *testing.T) {
	t.Parallel()
	for _, empty := range []bool{false, true} {
		t.Run(fmt.Sprint(empty), func(t *testing.T) {
			record := cancelledRecordFixture()
			if empty {
				record.CompletedTasks = nil
				record.PendingTasks = nil
				record.ReturnedSources = nil
			}
			content, err := RenderArchiveRecord(record)
			if err != nil {
				t.Fatal(err)
			}
			parsed, err := ParseArchiveRecord(content)
			if err != nil {
				t.Fatal(err)
			}
			next, err := RenderArchiveRecord(parsed)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(content, next) {
				t.Fatalf("round trip changed bytes:\n%s\n%s", content, next)
			}
			if parsed.CancelReason != record.CancelReason {
				t.Fatalf("reason=%q", parsed.CancelReason)
			}
			for _, key := range []string{"cancel_reason", "completed_tasks", "pending_tasks", "returned_sources"} {
				lines := strings.Split(string(content), "\n")
				var kept []string
				removing := false
				for _, line := range lines {
					if strings.HasPrefix(line, key+":") {
						removing = true
						continue
					}
					if removing && strings.HasPrefix(line, "  ") {
						continue
					}
					removing = false
					kept = append(kept, line)
				}
				if _, err := ParseArchiveRecord([]byte(strings.Join(kept, "\n"))); err == nil || !strings.Contains(err.Error(), "archive record missing "+key) {
					t.Fatalf("missing %s: %v", key, err)
				}
			}
			if empty {
				for _, key := range []string{"completed_tasks", "pending_tasks", "returned_sources"} {
					if !bytes.Contains(content, []byte(key+": []\n")) {
						t.Fatalf("missing empty %s: %s", key, content)
					}
				}
			}
		})
	}
}

func TestOnlyACancelledRecordCarriesCancelFields(t *testing.T) {
	t.Parallel()
	for _, disposition := range []ArchiveDisposition{ArchivePass, ArchivePartial, ArchiveQAOverride, ArchiveFailedQA, ArchiveSuperseded, ArchiveNoQA} {
		for _, key := range []string{"cancel_reason", "completed_tasks", "pending_tasks", "returned_sources"} {
			t.Run(string(disposition)+"/"+key, func(t *testing.T) {
				record := cancelledRecordFixture()
				record.Disposition = disposition
				record.SupersededBy = "successor"
				record.QAOverride = &QAArchiveOverrideRecord{Approval: "maintainer", Reason: "reason"}
				record.QAVerdict, record.QAReport = VerdictFail, "qa-report.md"
				content, err := RenderArchiveRecord(record)
				if err != nil {
					t.Fatal(err)
				}
				for _, forbidden := range []string{"cancel_reason", "completed_tasks", "pending_tasks", "returned_sources"} {
					if bytes.Contains(content, []byte(forbidden+":")) {
						t.Fatalf("renderer added %s to %s", forbidden, disposition)
					}
				}
				value := "[]"
				if key == "cancel_reason" {
					value = "reason"
				}
				content = bytes.Replace(content, []byte("disposition: "+string(disposition)+"\n"), []byte("disposition: "+string(disposition)+"\n"+key+": "+value+"\n"), 1)
				want := fmt.Sprintf("archive record disposition %s cannot carry %s", disposition, key)
				if _, err := ParseArchiveRecord(content); err == nil || err.Error() != want {
					t.Fatalf("got %v, want %s", err, want)
				}
			})
		}
	}
}

func TestArchivedTaskCompletedForACancelledRecord(t *testing.T) {
	t.Parallel()
	record := cancelledRecordFixture()
	record.QATask = "task_02"
	for _, task := range []string{"task_01", "task_02", "task_03", "unknown", ""} {
		want := task == "task_01" || task == "task_02"
		if got := ArchivedTaskCompleted(record, task); got != want {
			t.Fatalf("task %q completed=%v, want %v", task, got, want)
		}
	}
	record.CompletedTasks = nil
	if ArchivedTaskCompleted(record, "task_01") {
		t.Fatal("empty completed list completes a Task")
	}
}
