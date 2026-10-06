package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// ProtocolFailureMessageLimit bounds adapter messages in diagnostics and records.
const ProtocolFailureMessageLimit = 512

// ProtocolFailure places the first failure in an acpx prompt process.
type ProtocolFailure struct {
	Step       string
	Code       int
	Message    string
	PromptSent bool
}

// ProtocolDescription describes a prompt-process or session-preparation failure.
type ProtocolDescription struct {
	Step, Message string
	PromptSent    bool
	FromPrompt    bool
}

// DescribeProtocolFailure places only known prompt and preparation failures.
func DescribeProtocolFailure(err error) (ProtocolDescription, bool) {
	if protocol := protocolFailureInChain(err); protocol != nil {
		message := protocol.Message
		if message == "" {
			message = protocolStderrMessage(err)
		}
		return ProtocolDescription{Step: protocol.Step, Message: boundedProtocolMessage(message), PromptSent: protocol.PromptSent, FromPrompt: true}, true
	}
	var rejected *SelectionRejectedError
	var model *ModelNotAdvertisedError
	var access *AccessPolicyError
	var adapter AdapterProbeError
	var adapterPointer *AdapterProbeError
	step := ""
	switch {
	case errors.As(err, &rejected) && rejected != nil:
		switch {
		case strings.HasPrefix(rejected.Operation, "ensure"):
			step = "sessions ensure"
		case strings.HasPrefix(rejected.Operation, "set "):
			step = rejected.Operation
		}
	case errors.As(err, &model) && model != nil:
		step = "sessions ensure"
	case errors.As(err, &access) && access != nil:
		step = "set-mode"
	case errors.As(err, &adapter), errors.As(err, &adapterPointer):
		step = "adapter startup"
	}
	if step == "" {
		return ProtocolDescription{}, false
	}
	return ProtocolDescription{Step: step, Message: boundedProtocolMessage(protocolStderrMessage(err))}, true
}

func protocolFailureInChain(err error) *ProtocolFailure {
	switch failure := err.(type) {
	case *SelectionFailureError:
		if failure != nil && failure.Protocol != nil {
			return failure.Protocol
		}
	case *BatchFailureError:
		if failure != nil && failure.Protocol != nil {
			return failure.Protocol
		}
	}
	switch wrapped := err.(type) {
	case interface{ Unwrap() []error }:
		for _, child := range wrapped.Unwrap() {
			if protocol := protocolFailureInChain(child); protocol != nil {
				return protocol
			}
		}
	case interface{ Unwrap() error }:
		return protocolFailureInChain(wrapped.Unwrap())
	}
	return nil
}

func protocolStderrMessage(err error) string {
	var infrastructure *InfrastructureError
	if errors.As(err, &infrastructure) && infrastructure != nil {
		tail, _ := infrastructureStderrTail(infrastructure.Stderr)
		return strings.ToValidUTF8(tail, "")
	}
	var batch *BatchFailureError
	if errors.As(err, &batch) && batch != nil {
		tail, _ := infrastructureStderrTail(batch.Stderr)
		return strings.ToValidUTF8(tail, "")
	}
	var selection *SelectionFailureError
	if errors.As(err, &selection) && selection != nil && selection.Err != nil {
		tail, _ := infrastructureStderrTail(selection.Err.Error())
		return strings.ToValidUTF8(tail, "")
	}
	return ""
}

func boundedProtocolMessage(message string) string {
	message = strings.Join(strings.Fields(strings.ToValidUTF8(message, "")), " ")
	if len(message) <= ProtocolFailureMessageLimit {
		return message
	}
	end := ProtocolFailureMessageLimit
	for !utf8.RuneStart(message[end]) {
		end--
	}
	return message[:end]
}

func protocolFailureText(protocol *ProtocolFailure) string {
	if protocol == nil {
		return ""
	}
	text := " at " + protocol.Step
	if message := boundedProtocolMessage(protocol.Message); message != "" {
		text += fmt.Sprintf(": %s (JSON-RPC %d)", message, protocol.Code)
	}
	return text
}

// The trace keeps request identifiers and methods only, never request params.
type protocolRequest struct {
	id       string
	method   string
	answered bool
}

type promptProtocolTrace struct {
	requests   []protocolRequest
	promptSent bool
	failure    *ProtocolFailure
}

func (trace *promptProtocolTrace) observe(line []byte) {
	var message acpxJSONRPCMessage
	if json.Unmarshal(line, &message) != nil {
		return
	}
	id := string(bytes.TrimSpace(message.ID))
	hasID := id != "" && id != "null"
	if hasID && trackedProtocolMethod(message.Method) {
		trace.requests = append(trace.requests, protocolRequest{id: id, method: message.Method})
		if message.Method == "session/prompt" {
			trace.promptSent = true
		}
	}
	// Resolve the step before marking the response answered, including null ids.
	if message.Error != nil && trace.failure == nil {
		failure := trace.describe(id)
		failure.Code = message.Error.Code
		failure.Message = boundedProtocolMessage(message.Error.Message)
		trace.failure = &failure
	}
	if hasID && message.Method == "" && (len(message.Result) > 0 || message.Error != nil) {
		for index := len(trace.requests) - 1; index >= 0; index-- {
			if trace.requests[index].id == id {
				trace.requests[index].answered = true
				break
			}
		}
	}
}

func trackedProtocolMethod(method string) bool {
	switch method {
	case "initialize", "authenticate", "session/new", "session/load", "session/set_model", "session/set_mode", "session/set_config_option", "session/prompt":
		return true
	default:
		return false
	}
}

func (trace *promptProtocolTrace) describe(id string) ProtocolFailure {
	failure := ProtocolFailure{PromptSent: trace.promptSent}
	for index := len(trace.requests) - 1; index >= 0; index-- {
		if trace.requests[index].id == id {
			failure.Step = trace.requests[index].method
			return failure
		}
	}
	for index := len(trace.requests) - 1; index >= 0; index-- {
		if !trace.requests[index].answered {
			failure.Step = trace.requests[index].method
			return failure
		}
	}
	switch {
	case len(trace.requests) == 0:
		failure.Step = "adapter startup"
	case trace.promptSent:
		failure.Step = "session/prompt"
	default:
		failure.Step = "session setup"
	}
	return failure
}

func (trace *promptProtocolTrace) result() *ProtocolFailure {
	if trace.failure != nil {
		return trace.failure
	}
	failure := trace.describe("")
	return &failure
}
