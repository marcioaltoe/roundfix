package codex

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestCodexHygieneInspector(t *testing.T) {
	tests := []struct {
		name              string
		goos              string
		configuredPath    string
		pathEntry         string
		quarantined       bool
		accepted          bool
		wantStatus        Status
		wantPath          string
		wantQuarantined   bool
		wantAccepted      bool
		wantNextAction    string
		wantDetail        string
		wantLookPathCalls int
		wantProbeCalls    int
	}{
		{
			name:              "quarantined configured codex fails with curl reinstall",
			goos:              "darwin",
			configuredPath:    "/configured/codex",
			pathEntry:         "/path/codex",
			quarantined:       true,
			accepted:          true,
			wantStatus:        StatusFailed,
			wantPath:          "/configured/codex",
			wantQuarantined:   true,
			wantAccepted:      true,
			wantNextAction:    ReinstallNextAction,
			wantDetail:        "/configured/codex is quarantined with com.apple.quarantine",
			wantLookPathCalls: 0,
			wantProbeCalls:    1,
		},
		{
			name:              "not accepted path codex fails with curl reinstall",
			goos:              "darwin",
			pathEntry:         "/path/codex",
			accepted:          false,
			wantStatus:        StatusFailed,
			wantPath:          "/path/codex",
			wantNextAction:    ReinstallNextAction,
			wantDetail:        "/path/codex is not accepted by Gatekeeper",
			wantLookPathCalls: 1,
			wantProbeCalls:    1,
		},
		{
			name:              "clean path codex passes",
			goos:              "darwin",
			pathEntry:         "/path/codex",
			accepted:          true,
			wantStatus:        StatusOK,
			wantPath:          "/path/codex",
			wantAccepted:      true,
			wantDetail:        "/path/codex is accepted by Gatekeeper and has no com.apple.quarantine attribute",
			wantLookPathCalls: 1,
			wantProbeCalls:    1,
		},
		{
			name:              "non darwin is not applicable and does not inspect",
			goos:              "linux",
			pathEntry:         "/path/codex",
			quarantined:       true,
			accepted:          false,
			wantStatus:        StatusSkipped,
			wantDetail:        "not-applicable on linux",
			wantLookPathCalls: 0,
			wantProbeCalls:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lookPathCalls := 0
			quarantineCalls := 0
			acceptanceCalls := 0
			inspector := Inspector{
				ConfiguredPath: tt.configuredPath,
				GOOS:           tt.goos,
				LookPath: func(command string) (string, error) {
					lookPathCalls++
					if command != BinaryName {
						t.Fatalf("expected look path command %q, got %q", BinaryName, command)
					}
					if tt.pathEntry == "" {
						return "", errors.New("missing codex")
					}
					return tt.pathEntry, nil
				},
				Quarantine: quarantineProbeFunc(func(_ context.Context, path string) (bool, error) {
					quarantineCalls++
					if path != tt.wantPath {
						t.Fatalf("expected quarantine path %q, got %q", tt.wantPath, path)
					}
					return tt.quarantined, nil
				}),
				Acceptance: acceptanceProbeFunc(func(_ context.Context, path string) (bool, error) {
					acceptanceCalls++
					if path != tt.wantPath {
						t.Fatalf("expected acceptance path %q, got %q", tt.wantPath, path)
					}
					return tt.accepted, nil
				}),
			}

			got := inspector.Inspect(context.Background())

			if got.Status != tt.wantStatus {
				t.Fatalf("status = %q, want %q", got.Status, tt.wantStatus)
			}
			if got.Hygiene.Path != tt.wantPath {
				t.Fatalf("path = %q, want %q", got.Hygiene.Path, tt.wantPath)
			}
			if got.Hygiene.Quarantined != tt.wantQuarantined {
				t.Fatalf("quarantined = %v, want %v", got.Hygiene.Quarantined, tt.wantQuarantined)
			}
			if got.Hygiene.Accepted != tt.wantAccepted {
				t.Fatalf("accepted = %v, want %v", got.Hygiene.Accepted, tt.wantAccepted)
			}
			if got.NextAction != tt.wantNextAction {
				t.Fatalf("next action = %q, want %q", got.NextAction, tt.wantNextAction)
			}
			if got.Detail != tt.wantDetail {
				t.Fatalf("detail = %q, want %q", got.Detail, tt.wantDetail)
			}
			if lookPathCalls != tt.wantLookPathCalls {
				t.Fatalf("look path calls = %d, want %d", lookPathCalls, tt.wantLookPathCalls)
			}
			if quarantineCalls != tt.wantProbeCalls || acceptanceCalls != tt.wantProbeCalls {
				t.Fatalf("probe calls quarantine=%d acceptance=%d, want %d each", quarantineCalls, acceptanceCalls, tt.wantProbeCalls)
			}
		})
	}
}

func TestCodexHygieneInspectorReportsResolutionFailure(t *testing.T) {
	inspector := Inspector{
		GOOS: "darwin",
		LookPath: func(string) (string, error) {
			return "", errors.New("not found")
		},
		Quarantine: quarantineProbeFunc(func(context.Context, string) (bool, error) {
			t.Fatal("quarantine probe must not run after resolution failure")
			return false, nil
		}),
		Acceptance: acceptanceProbeFunc(func(context.Context, string) (bool, error) {
			t.Fatal("acceptance probe must not run after resolution failure")
			return false, nil
		}),
	}

	got := inspector.Inspect(context.Background())

	if got.Status != StatusFailed {
		t.Fatalf("status = %q, want %q", got.Status, StatusFailed)
	}
	if got.NextAction != ReinstallNextAction {
		t.Fatalf("next action = %q, want %q", got.NextAction, ReinstallNextAction)
	}
	if !strings.Contains(got.Detail, "resolve codex on PATH") {
		t.Fatalf("expected resolution detail, got %q", got.Detail)
	}
}

type quarantineProbeFunc func(context.Context, string) (bool, error)

func (fn quarantineProbeFunc) Quarantined(ctx context.Context, path string) (bool, error) {
	return fn(ctx, path)
}

type acceptanceProbeFunc func(context.Context, string) (bool, error)

func (fn acceptanceProbeFunc) Accepted(ctx context.Context, path string) (bool, error) {
	return fn(ctx, path)
}
