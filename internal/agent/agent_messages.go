package agent

import "strings"

type agentMessageLog struct {
	entries           []string
	lastUpdateKind    StreamUpdateKind
	lastMessageID     string
	totalMessageBytes int
}

func (log *agentMessageLog) observe(update StreamUpdate) {
	if update.Kind == StreamUpdateMessage {
		startMessage := len(log.entries) == 0
		if !startMessage {
			switch {
			case log.lastMessageID != "" && update.MessageID != "":
				startMessage = log.lastMessageID != update.MessageID
			default:
				startMessage = updateSplitsAgentMessage(log.lastUpdateKind)
			}
		}

		if startMessage {
			log.entries = append(log.entries, update.Text)
		} else {
			last := len(log.entries) - 1
			log.entries[last] += update.Text
		}
		log.lastMessageID = update.MessageID
		log.totalMessageBytes += len(update.Text)
	}
	log.lastUpdateKind = update.Kind
}

func updateSplitsAgentMessage(kind StreamUpdateKind) bool {
	switch kind {
	case StreamUpdateThought, StreamUpdateToolStarted, StreamUpdateToolUpdated, StreamUpdatePlan:
		return true
	default:
		return false
	}
}

func (log *agentMessageLog) messages() []string {
	return append([]string(nil), log.entries...)
}

func (log *agentMessageLog) final() string {
	for index := len(log.entries) - 1; index >= 0; index-- {
		if strings.TrimSpace(log.entries[index]) != "" {
			return log.entries[index]
		}
	}
	return ""
}

func (log *agentMessageLog) totalBytes() int {
	return log.totalMessageBytes
}
