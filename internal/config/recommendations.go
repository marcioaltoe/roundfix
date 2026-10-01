package config

// The dated Recommended Profile is advisory. Published evidence and sources
// live in docs/references/model-selection.md.
const (
	ModelRecommendationSnapshotVersion = "2026-09-30"
	ModelRecommendationSnapshotDate    = ModelRecommendationSnapshotVersion
)

type RecommendationRole string

const (
	RecommendationPreferred RecommendationRole = "preferred"
	RecommendationFallback  RecommendationRole = "fallback"
)

type ModelRecommendation struct {
	Category   WorkCategory       `json:"category"`
	Rank       int                `json:"rank"`
	Role       RecommendationRole `json:"role"`
	Selection  AgentSelection     `json:"selection"`
	SourceAsOf string             `json:"source_as_of"`
	Rationale  string             `json:"rationale"`
}

type recommendedProfileEntry struct {
	profile    AgentSelectionProfile
	rationales []string
}

var recommendedProfilesByCategory = map[WorkCategory]recommendedProfileEntry{
	CategoryGeneral: implementationRecommendation(),
	CategoryBackend: implementationRecommendation(),
	CategoryData:    implementationRecommendation(),
	CategoryInfra:   implementationRecommendation(),
	CategoryTest:    implementationRecommendation(),
	CategoryQA:      implementationRecommendation(),
	CategoryFrontend: {
		profile: AgentSelectionProfile{
			Preferred: AgentSelection{Runtime: "claude", Model: "opus", ReasoningEffort: "high"},
			Fallbacks: []AgentSelection{{Runtime: "codex", Model: "gpt-6.1-sol", ReasoningEffort: "xhigh"}},
		},
		rationales: []string{
			"Design judgment; opus resolves to Opus 5.5.",
			"Changes runtime; the provider places connected visual work at extra-high effort.",
		},
	},
	CategoryReview: {
		profile: AgentSelectionProfile{
			Preferred: AgentSelection{Runtime: "codex", Model: "gpt-5.6-luna", ReasoningEffort: "max"},
			Fallbacks: []AgentSelection{{Runtime: "codex", Model: "gpt-6.1-sol", ReasoningEffort: "high"}},
		},
		rationales: []string{
			boundedWorkRationale,
			"Stays on Codex, because every review selection must use the pre-PR review provider's runtime.",
		},
	},
	CategoryDocs:  boundedWorkRecommendation(),
	CategoryChore: boundedWorkRecommendation(),
}

const boundedWorkRationale = "Bounded work that blocks no other Task; about five minutes per session at a fraction of the price."

func implementationRecommendation() recommendedProfileEntry {
	return recommendedProfileEntry{
		profile: AgentSelectionProfile{
			Preferred: AgentSelection{Runtime: "codex", Model: "gpt-6.1-sol", ReasoningEffort: "high"},
			Fallbacks: []AgentSelection{{Runtime: "claude", Model: "opus", ReasoningEffort: "high"}},
		},
		rationales: []string{
			"The Codex workhorse since 2026-09-29, in the quota band of the model it replaces; adopted directly by the maintainer.",
			"Changes runtime; opus resolves to Opus 5.5.",
		},
	}
}

func boundedWorkRecommendation() recommendedProfileEntry {
	return recommendedProfileEntry{
		profile: AgentSelectionProfile{
			Preferred: AgentSelection{Runtime: "codex", Model: "gpt-5.6-luna", ReasoningEffort: "max"},
			Fallbacks: []AgentSelection{{Runtime: "claude", Model: "sonnet", ReasoningEffort: "high"}},
		},
		rationales: []string{
			boundedWorkRationale,
			"Changes runtime; Sonnet 5.5 is the efficient Claude model.",
		},
	}
}

// RecommendedProfile returns a copy of the dated profile Roundfix recommends
// for category. Every Agent Work Category has one.
func RecommendedProfile(category WorkCategory) (AgentSelectionProfile, bool) {
	entry, ok := recommendedProfilesByCategory[category]
	if !ok {
		return AgentSelectionProfile{}, false
	}
	return cloneProfile(entry.profile), true
}

// ModelRecommendations lists the Recommended Profile in order: the Preferred
// Selection at rank 1, then the Fallback Chain.
func ModelRecommendations(category WorkCategory) ([]ModelRecommendation, bool) {
	entry, ok := recommendedProfilesByCategory[category]
	if !ok {
		return nil, false
	}
	selections := append([]AgentSelection{entry.profile.Preferred}, entry.profile.Fallbacks...)
	rows := make([]ModelRecommendation, 0, len(selections))
	for i, selection := range selections {
		role := RecommendationFallback
		if i == 0 {
			role = RecommendationPreferred
		}
		rows = append(rows, ModelRecommendation{
			Category: category, Rank: i + 1, Role: role, Selection: selection,
			SourceAsOf: ModelRecommendationSnapshotVersion, Rationale: entry.rationales[i],
		})
	}
	return rows, true
}
