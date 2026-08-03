// Suite: Source Baseline rule segmentation
// Invariant: admitted ranges cover every source byte exactly once and materialize only locally derived Source Baseline Entries.
// Boundary IN: canonical segmentation snapshots, strict range proposals, local materialization, and manual fallback.
// Boundary OUT: ACP supervision, semantic classification, repository mutation, and final human approval.

package baseline

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestSegmentationSnapshotIsCanonicalAndCheckoutFree(t *testing.T) {
	t.Parallel()

	source := segmentationTestSource()
	first, err := NewRuleSegmentationSnapshot(source)
	if err != nil {
		t.Fatalf("build first Segmentation Snapshot: %v", err)
	}
	second, err := NewRuleSegmentationSnapshot(source)
	if err != nil {
		t.Fatalf("build second Segmentation Snapshot: %v", err)
	}
	firstBytes, err := first.CanonicalBytes()
	if err != nil {
		t.Fatalf("marshal first Segmentation Snapshot: %v", err)
	}
	secondBytes, err := second.CanonicalBytes()
	if err != nil {
		t.Fatalf("marshal second Segmentation Snapshot: %v", err)
	}
	if !bytes.Equal(firstBytes, secondBytes) {
		t.Fatal("equivalent Segmentation Snapshots produced different canonical bytes")
	}
	if first.SnapshotDigest == "" || first.SnapshotDigest != second.SnapshotDigest {
		t.Fatalf("snapshot digests differ: %q != %q", first.SnapshotDigest, second.SnapshotDigest)
	}
	for _, forbidden := range []string{
		"/repository/checkout",
		`"repository"`,
		`"workDir"`,
		`"tools"`,
		`"writeCapability"`,
	} {
		if bytes.Contains(firstBytes, []byte(forbidden)) {
			t.Fatalf("sealed Segmentation Snapshot exposes %q: %s", forbidden, firstBytes)
		}
	}

	tampered := first
	tampered.SourceBaseline = cloneReadoptionSourceBaseline(first.SourceBaseline)
	tampered.SourceBaseline.Entries[0].SourceBytes[0] ^= 0xff
	if _, err := tampered.CanonicalBytes(); err == nil {
		t.Fatal("mutated Segmentation Snapshot retained its sealed identity")
	}
}

func TestSegmentationSnapshotMakesExactTextAndBoundariesAvailableToSealedAnalysis(t *testing.T) {
	t.Parallel()

	snapshot := segmentationTestSnapshot(t)
	canonical, err := snapshot.CanonicalBytes()
	if err != nil {
		t.Fatalf("marshal Segmentation Snapshot: %v", err)
	}
	entry := snapshot.SourceBaseline.Entries[0]
	for _, required := range []string{
		`"semanticEntries"`,
		`"entryId":"` + entry.ID + `"`,
		`"lines":[{"start":0,"end":11,"text":"first rule\n"},{"start":11,"end":12,"text":"\n"},{"start":12,"end":24,"text":"second rule\n"}]`,
		`Split each readable entry into the smallest coherent instruction clauses`,
	} {
		if !bytes.Contains(canonical, []byte(required)) {
			t.Fatalf("Segmentation Snapshot does not expose %q:\n%s", required, canonical)
		}
	}

	proposal := map[string]any{
		"schemaVersion":  RuleSegmentationProposalSchemaVersion,
		"snapshotDigest": snapshot.SnapshotDigest,
		"segments": []map[string]any{
			{"entryId": entry.ID, "start": 0, "end": 11},
			{"entryId": entry.ID, "start": 11, "end": len(entry.SourceBytes)},
		},
	}
	payload, err := json.Marshal(proposal)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseRuleSegmentationProposal(payload, snapshot)
	if err != nil {
		t.Fatalf("range-only Segmentation Proposal was rejected: %v", err)
	}
	if parsed.SourceBaseline.ID != snapshot.SourceBaseline.ID ||
		parsed.SourceBaseline.Digest != snapshot.SourceBaseline.Digest {
		t.Fatalf("Source Baseline identity was not derived locally: %+v", parsed.SourceBaseline)
	}
	for index, segment := range parsed.Segments {
		if segment.Digest == "" {
			t.Fatalf("segment %d digest was not derived locally: %+v", index, segment)
		}
	}
}

func TestSegmentationProposalRejectsInvalidRangesAndStaleIdentity(t *testing.T) {
	t.Parallel()

	snapshot := segmentationTestSnapshot(t)
	valid := segmentationTestProposal(t, snapshot, 11)
	validJSON, err := json.Marshal(valid)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseRuleSegmentationProposal(validJSON, snapshot); err != nil {
		t.Fatalf("valid Segmentation Proposal rejected: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(RuleSegmentationProposal) RuleSegmentationProposal
	}{
		{
			name: "gap",
			mutate: func(proposal RuleSegmentationProposal) RuleSegmentationProposal {
				proposal.Segments[1].Start++
				return proposal
			},
		},
		{
			name: "overlap",
			mutate: func(proposal RuleSegmentationProposal) RuleSegmentationProposal {
				proposal.Segments[1].Start--
				return proposal
			},
		},
		{
			name: "unadvertised boundary",
			mutate: func(proposal RuleSegmentationProposal) RuleSegmentationProposal {
				proposal.Segments[0].End = 10
				proposal.Segments[0].Digest = ""
				proposal.Segments[1].Start = 10
				proposal.Segments[1].Digest = ""
				return proposal
			},
		},
		{
			name: "duplicate",
			mutate: func(proposal RuleSegmentationProposal) RuleSegmentationProposal {
				proposal.Segments = append(
					proposal.Segments[:1],
					append([]RuleSegmentProposal{proposal.Segments[0]}, proposal.Segments[1:]...)...,
				)
				return proposal
			},
		},
		{
			name: "reordered",
			mutate: func(proposal RuleSegmentationProposal) RuleSegmentationProposal {
				proposal.Segments[0], proposal.Segments[1] =
					proposal.Segments[1], proposal.Segments[0]
				return proposal
			},
		},
		{
			name: "empty",
			mutate: func(proposal RuleSegmentationProposal) RuleSegmentationProposal {
				proposal.Segments[0].End = proposal.Segments[0].Start
				proposal.Segments[0].Digest = segmentationTestDigest(nil)
				return proposal
			},
		},
		{
			name: "stale snapshot",
			mutate: func(proposal RuleSegmentationProposal) RuleSegmentationProposal {
				proposal.SnapshotDigest = "sha256:" + strings.Repeat("0", 64)
				return proposal
			},
		},
		{
			name: "stale Source Baseline",
			mutate: func(proposal RuleSegmentationProposal) RuleSegmentationProposal {
				proposal.SourceBaseline.Digest = strings.Repeat("0", 64)
				return proposal
			},
		},
		{
			name: "unknown entry",
			mutate: func(proposal RuleSegmentationProposal) RuleSegmentationProposal {
				proposal.Segments[0].EntryID = "source-entry." + strings.Repeat("f", 64)
				return proposal
			},
		},
		{
			name: "out of bounds",
			mutate: func(proposal RuleSegmentationProposal) RuleSegmentationProposal {
				proposal.Segments[1].End++
				return proposal
			},
		},
		{
			name: "rewritten digest",
			mutate: func(proposal RuleSegmentationProposal) RuleSegmentationProposal {
				proposal.Segments[0].Digest = strings.Repeat("a", 64)
				return proposal
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proposal := test.mutate(cloneSegmentationTestProposal(valid))
			if _, err := normalizeRuleSegmentationProposal(snapshot, proposal); err == nil {
				t.Fatal("invalid Segmentation Proposal was accepted")
			}
		})
	}
}

func TestSegmentationProposalParserRejectsUntrustedJSON(t *testing.T) {
	t.Parallel()

	snapshot := segmentationTestSnapshot(t)
	valid := segmentationTestProposal(t, snapshot, 11)
	payload, err := json.Marshal(valid)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		data []byte
	}{
		{name: "extra prose", data: append([]byte("proposal:\n"), payload...)},
		{
			name: "unknown field",
			data: bytes.Replace(
				payload,
				[]byte(`"schemaVersion":`),
				[]byte(`"unknown":true,"schemaVersion":`),
				1,
			),
		},
		{
			name: "duplicate key",
			data: bytes.Replace(
				payload,
				[]byte(`"snapshotDigest":`),
				[]byte(`"snapshotDigest":"stale","snapshotDigest":`),
				1,
			),
		},
		{
			name: "oversized output",
			data: []byte(strings.Repeat("x", RuleSegmentationProposalMaxBytes+1)),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := ParseRuleSegmentationProposal(test.data, snapshot); err == nil {
				t.Fatal("untrusted Segmentation Proposal JSON was accepted")
			}
		})
	}
}

func TestRuleSegmentationMaterializesExactBytesAndStableIdentities(t *testing.T) {
	t.Parallel()

	snapshot := segmentationTestSnapshot(t)
	firstProposal := segmentationTestProposal(t, snapshot, 11)
	encoded, err := json.Marshal(firstProposal)
	if err != nil {
		t.Fatal(err)
	}
	secondProposal, err := ParseRuleSegmentationProposal(encoded, snapshot)
	if err != nil {
		t.Fatalf("parse equivalent Segmentation Proposal: %v", err)
	}

	first, err := MaterializeRuleSegments(snapshot, firstProposal)
	if err != nil {
		t.Fatalf("materialize first Segmentation Proposal: %v", err)
	}
	second, err := MaterializeRuleSegments(snapshot, secondProposal)
	if err != nil {
		t.Fatalf("materialize second Segmentation Proposal: %v", err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("equivalent proposals produced different Source Baselines:\nfirst=%+v\nsecond=%+v", first, second)
	}
	if first.ID != snapshot.SourceBaseline.ID ||
		first.Digest != snapshot.SourceBaseline.Digest ||
		first.ByteCount != snapshot.SourceBaseline.ByteCount {
		t.Fatalf("materialization changed Source Baseline identity: %+v", first)
	}
	var rebuilt []byte
	for _, entry := range first.Entries {
		rebuilt = append(rebuilt, entry.SourceBytes...)
	}
	if !bytes.Equal(rebuilt, snapshot.SourceBaseline.Entries[0].SourceBytes) {
		t.Fatalf("materialized bytes = %q, want %q", rebuilt, snapshot.SourceBaseline.Entries[0].SourceBytes)
	}
	if len(first.Entries) != 2 ||
		first.Entries[0].Start != snapshot.SourceBaseline.Entries[0].Start ||
		first.Entries[1].End != snapshot.SourceBaseline.Entries[0].End {
		t.Fatalf("materialized ranges are not canonical Source Baseline Entries: %+v", first.Entries)
	}
	if _, err := NewAnalysisSnapshot(first); err != nil {
		t.Fatalf("classification contract rejected materialized entries: %v", err)
	}
}

func TestRuleSegmentationManualFallbackRetainsOriginalEntry(t *testing.T) {
	t.Parallel()

	snapshot := segmentationTestSnapshot(t)
	proposal, err := ManualRuleSegmentationProposal(snapshot)
	if err != nil {
		t.Fatalf("build manual Segmentation Proposal: %v", err)
	}
	materialized, err := MaterializeRuleSegments(snapshot, proposal)
	if err != nil {
		t.Fatalf("materialize manual Segmentation Proposal: %v", err)
	}
	if !reflect.DeepEqual(materialized, snapshot.SourceBaseline) {
		t.Fatalf("manual fallback mutated original Source Baseline:\n got=%+v\nwant=%+v",
			materialized, snapshot.SourceBaseline)
	}
}

func segmentationTestSource() ReadoptionSourceBaseline {
	content := []byte("first rule\n\nsecond rule\n")
	entry := newReadoptionSourceEntry(
		"AGENTS.md",
		"unmarked-span",
		0,
		len(content),
		content,
		segmentationTestDigest(content),
		map[string]any{"markerState": "unmarked"},
	)
	return ReadoptionSourceBaseline{
		ID:               "baseline.readoption." + strings.Repeat("c", 64),
		DeclaredIdentity: "unconfigured",
		Compatibility:    "incompatible",
		Digest:           strings.Repeat("c", 64),
		CarrierCount:     1,
		EntryCount:       1,
		ByteCount:        len(content),
		Entries:          []ReadoptionSourceEntry{entry},
	}
}

func segmentationTestSnapshot(t *testing.T) RuleSegmentationSnapshot {
	t.Helper()
	snapshot, err := NewRuleSegmentationSnapshot(segmentationTestSource())
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}

func segmentationTestProposal(
	t *testing.T,
	snapshot RuleSegmentationSnapshot,
	split int,
) RuleSegmentationProposal {
	t.Helper()
	entry := snapshot.SourceBaseline.Entries[0]
	proposal := RuleSegmentationProposal{
		SchemaVersion:  RuleSegmentationProposalSchemaVersion,
		SnapshotDigest: snapshot.SnapshotDigest,
		SourceBaseline: ClassificationSource{
			ID:     snapshot.SourceBaseline.ID,
			Digest: snapshot.SourceBaseline.Digest,
		},
		Segments: []RuleSegmentProposal{
			{
				EntryID: entry.ID,
				Start:   0,
				End:     split,
				Digest:  segmentationTestDigest(entry.SourceBytes[:split]),
			},
			{
				EntryID: entry.ID,
				Start:   split,
				End:     len(entry.SourceBytes),
				Digest:  segmentationTestDigest(entry.SourceBytes[split:]),
			},
		},
	}
	normalized, err := normalizeRuleSegmentationProposal(snapshot, proposal)
	if err != nil {
		t.Fatal(err)
	}
	return normalized
}

func cloneSegmentationTestProposal(proposal RuleSegmentationProposal) RuleSegmentationProposal {
	proposal.Segments = append([]RuleSegmentProposal(nil), proposal.Segments...)
	return proposal
}

func segmentationTestDigest(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
