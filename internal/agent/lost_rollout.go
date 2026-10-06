package agent

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
)

// LostRolloutPhrase is the adapter's diagnostic for a missing Codex rollout.
const LostRolloutPhrase = "no rollout found for thread id"

// LostRollout describes a missing rollout reported by a prompt process.
type LostRollout struct {
	Runtime    string
	Step       string
	Detail     string
	PromptSent bool
}

// DescribeLostRollout describes only marked selection or batch failures.
func DescribeLostRollout(err error) (LostRollout, bool) {
	var stop *StopError
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || errors.As(err, &stop) {
		return LostRollout{}, false
	}
	protocol := lostRolloutInChain(err)
	if protocol == nil {
		return LostRollout{}, false
	}
	description := LostRollout{Step: protocol.Step, Detail: boundedProtocolMessage(protocol.Detail), PromptSent: protocol.PromptSent}
	var selection *SelectionFailureError
	if errors.As(err, &selection) && selection != nil {
		description.Runtime = selection.Runtime
	}
	return description, true
}

func lostRolloutInChain(err error) *ProtocolFailure {
	var protocol *ProtocolFailure
	switch failure := err.(type) {
	case *SelectionFailureError:
		if failure != nil {
			protocol = failure.Protocol
		}
	case *BatchFailureError:
		if failure != nil {
			protocol = failure.Protocol
		}
	}
	if protocol != nil && protocol.LostRollout {
		return protocol
	}
	switch wrapped := err.(type) {
	case interface{ Unwrap() []error }:
		for _, child := range wrapped.Unwrap() {
			if protocol := lostRolloutInChain(child); protocol != nil {
				return protocol
			}
		}
	case interface{ Unwrap() error }:
		return lostRolloutInChain(wrapped.Unwrap())
	}
	return nil
}

func lostRolloutDetail(failure *acpxJSONRPCError) (string, bool) {
	if failure.Code != -32603 && failure.Code != -32600 {
		return "", false
	}
	if strings.Contains(failure.Message, LostRolloutPhrase) {
		return boundedProtocolMessage(failure.Message), true
	}
	// Decode only details. Keeping it raw makes non-string values inert.
	var data struct {
		Details json.RawMessage `json:"details"`
	}
	if json.Unmarshal(failure.Data, &data) != nil {
		return "", false
	}
	var detail string
	if json.Unmarshal(data.Details, &detail) != nil || !strings.Contains(detail, LostRolloutPhrase) {
		return "", false
	}
	return boundedProtocolMessage(detail), true
}
