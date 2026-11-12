package speccheck_test

import (
	"testing"

	"roundfix/internal/speccheck"
)

func TestAPortugueseNotApplicableADRDrawsInNoCitingADR(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name     string
		relation string
	}{
		{name: "accented", relation: "não se aplica; a decisão segue igual."},
		{name: "capitalized", relation: "Não se aplica — a decisão segue igual."},
		{name: "upper case", relation: "NÃO SE APLICA: a decisão segue igual."},
		{name: "unaccented", relation: "nao se aplica; a decisão segue igual."},
		{name: "decomposed", relation: "não se aplica; a decisão segue igual."},
		{name: "line end", relation: "não se aplica"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repoRoot := writeADRRelatedSeedFixture(t, test.relation)
			result := checkADRRelatedSeeds(t, repoRoot)
			for _, code := range []string{speccheck.CodeADRRelated, speccheck.CodeADRUnlisted, speccheck.CodeADRAlternativesUnrelated} {
				if findings := findingsWithCode(result, code); len(findings) != 0 {
					t.Errorf("%s findings = %#v, want none", code, findings)
				}
			}
		})
	}
}

func TestAMarkerOutsideTheClosedListStillDrawsIn(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name     string
		relation string
	}{
		{name: "plural", relation: "não se aplicam; a decisão segue igual."},
		{name: "not leading", relation: "mantida; a alternativa rejeitada não se aplica."},
		{name: "abbreviation", relation: "n/a; unchanged."},
		{name: "other wording", relation: "não aplicável; a decisão segue igual."},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repoRoot := writeADRRelatedSeedFixture(t, test.relation)
			result := checkADRRelatedSeeds(t, repoRoot)
			requireADRRelatedSeedFinding(t, result, "0003", "0002")
		})
	}
}
