package agent

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
)

// ProviderSystemErrorStatus is the adapter's failed thread status.
const ProviderSystemErrorStatus = "systemError"

// ProviderSystemErrorReason names a turn the provider failed.
const ProviderSystemErrorReason = "provider system error"

// ProviderSystemError describes a marked provider failure in a prompt turn.
type ProviderSystemError struct {
	Runtime string
	Detail  string
}

// DescribeProviderSystemError describes only marked selection or batch failures.
func DescribeProviderSystemError(err error) (ProviderSystemError, bool) {
	var stop *StopError
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || errors.As(err, &stop) || IsStopError(err) {
		return ProviderSystemError{}, false
	}
	protocol := providerSystemErrorInChain(err)
	if protocol == nil {
		return ProviderSystemError{}, false
	}
	description := ProviderSystemError{Detail: boundedProtocolMessage(protocol.Detail)}
	var selection *SelectionFailureError
	if errors.As(err, &selection) && selection != nil {
		description.Runtime = selection.Runtime
	}
	return description, true
}

func providerSystemErrorInChain(err error) *ProtocolFailure {
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
	if protocol != nil && protocol.ProviderSystemError {
		return protocol
	}

	switch wrapped := err.(type) {
	case interface{ Unwrap() []error }:
		for _, child := range wrapped.Unwrap() {
			if protocol := providerSystemErrorInChain(child); protocol != nil {
				return protocol
			}
		}
	case interface{ Unwrap() error }:
		return providerSystemErrorInChain(wrapped.Unwrap())
	}
	return nil
}

// promptThreadStatus retains the last status and subsequent text up to the
// prompt result. Session metadata is not Agent output.
type promptThreadStatus struct {
	status   string
	text     strings.Builder
	finished bool
}

func (status *promptThreadStatus) observe(message acpxJSONRPCMessage) {
	if status.finished {
		return
	}
	if message.Method == acpMethodSessionUpdate {
		var params struct {
			Update struct {
				SessionUpdate string                 `json:"sessionUpdate"`
				Content       acpContentBlockPayload `json:"content"`
				Meta          struct {
					Codex struct {
						ThreadStatus struct {
							Type string `json:"type"`
						} `json:"threadStatus"`
					} `json:"codex"`
				} `json:"_meta"`
			} `json:"update"`
		}
		if json.Unmarshal(message.Params, &params) != nil {
			return
		}
		update := params.Update
		switch update.SessionUpdate {
		case "session_info_update":
			if next := update.Meta.Codex.ThreadStatus.Type; next != "" {
				status.status = next
				status.text.Reset()
			}
		case "agent_message_chunk":
			if status.status == ProviderSystemErrorStatus && update.Content.Type == "text" {
				status.text.WriteString(update.Content.Text)
			}
		}
		return
	}
	if len(message.Result) > 0 {
		var result struct {
			StopReason string `json:"stopReason"`
		}
		if json.Unmarshal(message.Result, &result) == nil && result.StopReason != "" {
			status.finished = true
		}
	}
}

func (status *promptThreadStatus) failure() *ProtocolFailure {
	if status.status != ProviderSystemErrorStatus {
		return nil
	}
	detail := boundedProtocolMessage(status.text.String())
	if detail == "" {
		detail = "threadStatus " + ProviderSystemErrorStatus
	}
	return &ProtocolFailure{Step: "session/prompt", PromptSent: true, ProviderSystemError: true, Detail: detail}
}
