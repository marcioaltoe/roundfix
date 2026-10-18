package speccheck

import "bytes"

// EvidenceProbeBytes bounds the content inspected at the beginning of QA evidence.
const EvidenceProbeBytes = 8000

// ReadableEvidence accepts text (including SVG), PNG and JPEG by content alone.
// Callers supply at most EvidenceProbeBytes bytes from the start of the file.
func ReadableEvidence(probe []byte) bool {
	return bytes.IndexByte(probe, 0) < 0 ||
		bytes.HasPrefix(probe, []byte("\x89PNG\r\n\x1a\n")) ||
		bytes.HasPrefix(probe, []byte("\xff\xd8\xff"))
}
