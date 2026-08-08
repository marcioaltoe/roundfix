package config

// Model identifiers below are the ones each ACP adapter advertises, not display
// names from the benchmark. The claude adapter advertises `opus` (as
// `opus[1m]`), `claude-fable-5`, `sonnet`, `haiku`, and `default`; reasoning
// effort is a separate adapter option. `claude-opus-4-8` appeared in the
// 2026-07-16 snapshot and is advertised by no adapter, so it is gone.
//
// Rationale, sources, and the cost/latency/step trade-offs behind each ordering
// live in docs/references/model-selection.md.
const (
	ModelRecommendationSnapshotVersion = "2026-08-07"
	ModelRecommendationSnapshotDate    = ModelRecommendationSnapshotVersion
	modelRecommendationBenchmark       = "DeepSWE v1.1"
)

type ModelRecommendation struct {
	Category         WorkCategory   `json:"category"`
	Rank             int            `json:"rank"`
	Selection        AgentSelection `json:"selection"`
	Benchmark        string         `json:"benchmark"`
	ResultPercent    float64        `json:"result_percent"`
	AverageCostUSD   float64        `json:"average_cost_usd"`
	SourceAsOf       string         `json:"source_as_of"`
	Rationale        string         `json:"rationale"`
	CategorySpecific bool           `json:"category_specific"`
}

func ModelRecommendations(category WorkCategory) ([]ModelRecommendation, WorkCategory, bool) {
	if _, ok := ParseWorkCategory(string(category)); !ok {
		return nil, "", false
	}
	source := category
	if isOptionalWorkCategory(category) {
		source = CategoryGeneral
	}
	rows, ok := modelRecommendationsByCategory[source]
	if !ok {
		return nil, "", false
	}
	recommendations := append([]ModelRecommendation(nil), rows...)
	for index := range recommendations {
		recommendations[index].Category = category
	}
	return recommendations, source, true
}

var modelRecommendationsByCategory = map[WorkCategory][]ModelRecommendation{
	CategoryGeneral: {
		modelRecommendation(CategoryGeneral, 1, "codex", "gpt-5.6-sol", "high", 69, 3.47, "Built-in quality/cost default; 37 steps is fewer than anything near its result, which keeps Run wall clock down."),
		modelRecommendation(CategoryGeneral, 2, "claude", "opus", "high", 73, 6.08, "Highest result on this snapshot at a moderate setting; within one point of its own max for half the cost."),
		modelRecommendation(CategoryGeneral, 3, "codex", "gpt-5.6-terra", "max", 70, 3.96, "Quality-preserving alternative that costs less than raising the default to xhigh."),
		modelRecommendation(CategoryGeneral, 4, "codex", "gpt-5.6-luna", "max", 67, 0.61, "Cost floor: two points below the default for a sixth of the price, and its 4.1s latency offsets its extra steps."),
		modelRecommendation(CategoryGeneral, 5, "codex", "gpt-5.5", "xhigh", 67, 7.23, "Retained proven-generation option, now dominated by every GPT-5.6 setting above it."),
	},
	CategoryBackend: {
		modelRecommendation(CategoryBackend, 1, "codex", "gpt-5.6-sol", "high", 69, 3.47, "Best default for complex repository changes, debugging, and multi-file implementation."),
		modelRecommendation(CategoryBackend, 2, "claude", "opus", "high", 73, 6.08, "Cross-runtime option with the highest result, for architecture-heavy or broad refactoring work."),
		modelRecommendation(CategoryBackend, 3, "codex", "gpt-5.6-terra", "max", 70, 3.96, "Quality-preserving alternative with a strong long-horizon result."),
		modelRecommendation(CategoryBackend, 4, "codex", "gpt-5.6-luna", "max", 67, 0.61, "Lower-cost option for bounded backend slices where a Verification gate catches what the model misses."),
		modelRecommendation(CategoryBackend, 5, "codex", "gpt-5.5", "xhigh", 67, 7.23, "Stable legacy choice when a repository has already validated it."),
	},
	CategoryFrontend: {
		modelRecommendation(CategoryFrontend, 1, "claude", "opus", "high", 73, 6.08, "Preferred frontend profile on design judgment; task-fit is qualitative because no benchmark on file has a UI slice."),
		modelRecommendation(CategoryFrontend, 2, "codex", "gpt-5.6-sol", "high", 69, 3.47, "Strong cross-runtime fallback for implementation and integration correctness."),
		modelRecommendation(CategoryFrontend, 3, "codex", "gpt-5.6-terra", "max", 70, 3.96, "Quality-focused Codex alternative when visual specialization matters less than repository completion."),
		modelRecommendation(CategoryFrontend, 4, "codex", "gpt-5.6-luna", "max", 67, 0.61, "Cost-conscious option for small, well-specified UI slices."),
		modelRecommendation(CategoryFrontend, 5, "claude", "claude-fable-5", "high", 69, 9.18, "Ranked last despite a competitive result: it answers in about 1.7 minutes, roughly 25x Sol, and that multiplies by turn count in a Run."),
	},
	CategoryQA: {
		modelRecommendation(CategoryQA, 1, "codex", "gpt-5.6-sol", "high", 69, 3.47, "Built-in QA default for broad behavior checks and evidence synthesis."),
		modelRecommendation(CategoryQA, 2, "claude", "opus", "high", 73, 6.08, "Cross-runtime option with the highest result, for exploratory behavior analysis."),
		modelRecommendation(CategoryQA, 3, "codex", "gpt-5.6-luna", "max", 67, 0.61, "Cost-conscious option for deterministic, well-bounded QA matrices."),
		modelRecommendation(CategoryQA, 4, "codex", "gpt-5.6-terra", "max", 70, 3.96, "Quality-preserving fallback for long validation flows."),
		modelRecommendation(CategoryQA, 5, "codex", "gpt-5.5", "xhigh", 67, 7.23, "Legacy option for established QA setups already proven on GPT-5.5."),
	},
	CategoryReview: {
		modelRecommendation(CategoryReview, 1, "codex", "gpt-5.6-luna", "max", 67, 0.61, "Review work is bounded and has no dependency graph, so its extra steps block nobody while its price is a sixth of the alternatives."),
		modelRecommendation(CategoryReview, 2, "codex", "gpt-5.6-sol", "high", 69, 3.47, "Balances review depth with cost when two points of result are worth six times the price."),
		modelRecommendation(CategoryReview, 3, "claude", "opus", "medium", 69, 3.29, "Cross-runtime option at the same result and price as Sol, for provider diversity."),
		modelRecommendation(CategoryReview, 4, "codex", "gpt-5.6-terra", "xhigh", 60, 1.70, "Middle-cost alternative when Sol cannot start and Luna's result is too low for the change under review."),
		modelRecommendation(CategoryReview, 5, "claude", "claude-fable-5", "high", 69, 9.18, "Cross-runtime review option; its 1.7-minute latency matters less on a single bounded review than inside a Run."),
	},
}

func modelRecommendation(category WorkCategory, rank int, runtime string, model string, reasoningEffort string, resultPercent float64, averageCostUSD float64, rationale string) ModelRecommendation {
	return ModelRecommendation{
		Category: category,
		Rank:     rank,
		Selection: AgentSelection{
			Runtime:         runtime,
			Model:           model,
			ReasoningEffort: reasoningEffort,
		},
		Benchmark:        modelRecommendationBenchmark,
		ResultPercent:    resultPercent,
		AverageCostUSD:   averageCostUSD,
		SourceAsOf:       ModelRecommendationSnapshotDate,
		Rationale:        rationale,
		CategorySpecific: false,
	}
}
