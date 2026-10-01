package agent

import "encoding/json"

// UsageBasis names the scope of an adapter's token report.
type UsageBasis string

const (
	UsageBasisTurn       UsageBasis = "turn"
	UsageBasisRequestSum UsageBasis = "request-sum"
)

// TurnUsage is what one prompt's adapter reported. Its zero value is unreported.
type TurnUsage struct {
	Basis                               UsageBasis
	TotalTokens                         int64
	InputTokens, OutputTokens           *int64
	CachedReadTokens, CachedWriteTokens *int64
	ThoughtTokens                       *int64
	Readings                            int
	Cost                                *ReportedCost
}

// ReportedCost is the adapter's last cumulative cost reading.
type ReportedCost struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type promptUsage struct {
	TotalTokens       int64  `json:"totalTokens"`
	InputTokens       *int64 `json:"inputTokens"`
	OutputTokens      *int64 `json:"outputTokens"`
	CachedReadTokens  *int64 `json:"cachedReadTokens"`
	CachedWriteTokens *int64 `json:"cachedWriteTokens"`
	ThoughtTokens     *int64 `json:"thoughtTokens"`
}

func countTurnUsage(lastRequestOnly bool, reported *promptUsage, readings []int64, cost *ReportedCost) TurnUsage {
	usage := TurnUsage{Readings: len(readings), Cost: cost}
	if lastRequestOnly && len(readings) > 0 && (reported == nil || reported.TotalTokens <= readings[len(readings)-1]) {
		usage.Basis = UsageBasisRequestSum
		for _, reading := range readings {
			usage.TotalTokens += reading
		}
	} else if reported != nil {
		usage.Basis = UsageBasisTurn
		usage.TotalTokens = reported.TotalTokens
		usage.InputTokens, usage.OutputTokens = reported.InputTokens, reported.OutputTokens
		usage.CachedReadTokens, usage.CachedWriteTokens = reported.CachedReadTokens, reported.CachedWriteTokens
		usage.ThoughtTokens = reported.ThoughtTokens
	}
	return usage
}

type promptUsageCollector struct {
	reported *promptUsage
	readings []int64
	cost     *ReportedCost
}

func (collector *promptUsageCollector) observeUpdate(payload json.RawMessage) bool {
	var note acpSessionNotificationPayload
	if json.Unmarshal(payload, &note) != nil {
		return false
	}
	var header acpSessionUpdateHeader
	if json.Unmarshal(note.Update, &header) != nil || header.SessionUpdate != "usage_update" {
		return false
	}
	var update struct {
		Used *int64        `json:"used"`
		Cost *ReportedCost `json:"cost"`
	}
	// Usage is optional metadata: a malformed report must not alter the prompt outcome.
	if json.Unmarshal(note.Update, &update) == nil {
		if update.Used != nil {
			collector.readings = append(collector.readings, *update.Used)
		}
		if update.Cost != nil {
			collector.cost = update.Cost
		}
	}
	return true
}

func (collector *promptUsageCollector) observeResult(payload json.RawMessage) {
	var reported *promptUsage
	if json.Unmarshal(payload, &reported) == nil && reported != nil {
		collector.reported = reported
	}
}
