//go:build linux

package store

import (
	"strings"
	"testing"
)

func TestParseProcStatStartTimeCountsFromLastClosingParenthesis(t *testing.T) {
	stat := []byte("42 (worker ) pool) S 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 424242 20")

	startTime, err := parseProcStatStartTime(stat)
	if err != nil {
		t.Fatalf("parse process stat start time: %v", err)
	}
	if startTime != "424242" {
		t.Fatalf("process start time = %q, want 424242", startTime)
	}
}

func TestParseProcStatStartTimeRejectsMissingField(t *testing.T) {
	_, err := parseProcStatStartTime([]byte("42 (worker) S 1 2"))
	if err == nil || !strings.Contains(err.Error(), "missing start time field") {
		t.Fatalf("parse error = %v, want missing start time field", err)
	}
}
