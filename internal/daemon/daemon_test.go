package daemon

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestExecVerifierRunsConfiguredCommand(t *testing.T) {
	var stream bytes.Buffer
	err := ExecVerifier{}.Verify(context.Background(), VerifyRequest{
		WorkDir: t.TempDir(),
		Command: "printf verified",
		Stream:  &stream,
	})

	if err != nil {
		t.Fatalf("verify command: %v", err)
	}
	if stream.String() != "verified" {
		t.Fatalf("expected verification output, got %q", stream.String())
	}
}

func TestExecVerifierReportsFailedCommand(t *testing.T) {
	var stream bytes.Buffer
	err := ExecVerifier{}.Verify(context.Background(), VerifyRequest{
		WorkDir: t.TempDir(),
		Command: "printf broken; exit 7",
		Stream:  &stream,
	})

	if err == nil {
		t.Fatal("expected verification failure")
	}
	if !strings.Contains(err.Error(), "verification command") {
		t.Fatalf("expected verification error context, got %v", err)
	}
	if stream.String() != "broken" {
		t.Fatalf("expected failed verification output, got %q", stream.String())
	}
}

func TestGitCommitterValidatesRequest(t *testing.T) {
	err := GitCommitter{}.Commit(context.Background(), CommitRequest{})

	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "git root is required") {
		t.Fatalf("expected git root validation error, got %v", err)
	}
}

func TestGitPusherValidatesRequest(t *testing.T) {
	err := GitPusher{}.Push(context.Background(), PushRequest{})

	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "git root is required") {
		t.Fatalf("expected git root validation error, got %v", err)
	}
}
