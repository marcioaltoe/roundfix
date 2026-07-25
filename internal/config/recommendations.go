package config

const (
	ModelRecommendationSnapshotVersion = "2026-07-16"
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
		modelRecommendation(CategoryGeneral, 1, "codex", "gpt-5.6-sol", "high", 69, 3.47, "Built-in quality/cost default; Sol leads the broader composite at max without requiring max for routine work."),
		modelRecommendation(CategoryGeneral, 2, "codex", "gpt-5.6-terra", "max", 70, 4.95, "Strong quality-preserving fallback with lower average cost than the largest Claude settings."),
		modelRecommendation(CategoryGeneral, 3, "claude", "claude-fable-5", "high", 69, 9.18, "Cross-runtime high-quality option when provider diversity matters."),
		modelRecommendation(CategoryGeneral, 4, "codex", "gpt-5.6-luna", "max", 67, 3.03, "Cost-conscious option that remains competitive at maximum effort."),
		modelRecommendation(CategoryGeneral, 5, "codex", "gpt-5.5", "xhigh", 67, 7.23, "Retained proven-generation option, but less attractive than GPT-5.6 on this snapshot."),
	},
	CategoryBackend: {
		modelRecommendation(CategoryBackend, 1, "codex", "gpt-5.6-sol", "high", 69, 3.47, "Best default for complex repository changes, debugging, and multi-file implementation."),
		modelRecommendation(CategoryBackend, 2, "codex", "gpt-5.6-terra", "max", 70, 4.95, "Quality-preserving alternative with a strong long-horizon result."),
		modelRecommendation(CategoryBackend, 3, "claude", "claude-fable-5", "high", 69, 9.18, "Cross-runtime option for architecture-heavy or broad refactoring work."),
		modelRecommendation(CategoryBackend, 4, "codex", "gpt-5.5", "xhigh", 67, 7.23, "Stable legacy choice when a repository has already validated it."),
		modelRecommendation(CategoryBackend, 5, "codex", "gpt-5.6-luna", "max", 67, 3.03, "Lower-cost option for bounded backend slices."),
	},
	CategoryFrontend: {
		modelRecommendation(CategoryFrontend, 1, "claude", "claude-fable-5", "medium", 65, 6.09, "Dated benchmark-ranked frontend option; task-fit is qualitative because the benchmark has no UI-specific slice."),
		modelRecommendation(CategoryFrontend, 2, "codex", "gpt-5.6-sol", "high", 69, 3.47, "Strong cross-runtime fallback for implementation and integration correctness."),
		modelRecommendation(CategoryFrontend, 3, "claude", "claude-opus-4-8", "high", 52, 4.28, "Design-heavy alternative retained for repositories that have validated Opus; general benchmark result is weaker."),
		modelRecommendation(CategoryFrontend, 4, "codex", "gpt-5.6-terra", "max", 70, 4.95, "Quality-focused Codex alternative when visual specialization is less important than repository completion."),
		modelRecommendation(CategoryFrontend, 5, "codex", "gpt-5.6-luna", "max", 67, 3.03, "Cost-conscious option for small, well-specified UI slices."),
	},
	CategoryQA: {
		modelRecommendation(CategoryQA, 1, "codex", "gpt-5.6-sol", "high", 69, 3.47, "Built-in QA default for broad behavior checks and evidence synthesis."),
		modelRecommendation(CategoryQA, 2, "codex", "gpt-5.6-terra", "max", 70, 4.95, "Quality-preserving fallback for long validation flows."),
		modelRecommendation(CategoryQA, 3, "codex", "gpt-5.6-luna", "max", 67, 3.03, "Cost-conscious option for deterministic, well-bounded QA matrices."),
		modelRecommendation(CategoryQA, 4, "claude", "claude-fable-5", "high", 69, 9.18, "Cross-runtime option for exploratory behavior analysis."),
		modelRecommendation(CategoryQA, 5, "codex", "gpt-5.5", "xhigh", 67, 7.23, "Legacy option for established QA setups already proven on GPT-5.5."),
	},
	CategoryReview: {
		modelRecommendation(CategoryReview, 1, "codex", "gpt-5.6-sol", "high", 69, 3.47, "Matches the OpenClaw autoreview default and balances review depth with cost."),
		modelRecommendation(CategoryReview, 2, "codex", "gpt-5.6-terra", "max", 70, 4.95, "Strong alternative when Sol cannot start or a quality-preserving fallback is required."),
		modelRecommendation(CategoryReview, 3, "claude", "claude-fable-5", "high", 69, 9.18, "Cross-runtime review option with a competitive general coding result."),
		modelRecommendation(CategoryReview, 4, "codex", "gpt-5.5", "xhigh", 67, 7.23, "Established predecessor for repositories that prefer a validated older generation."),
		modelRecommendation(CategoryReview, 5, "claude", "claude-opus-4-8", "high", 52, 4.28, "Design- and architecture-oriented alternative; lower general benchmark result keeps it fifth."),
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
