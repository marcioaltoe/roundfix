package speccheck

import (
	"errors"
	"os/exec"
	"strings"
)

// ConcreteContractGuidePath is the guide whose adding commit starts the
// receipt and transcript declaration gaps.
const ConcreteContractGuidePath = ".agents/skills/write-techspec/references/concrete-contracts.md"

type contractHorizon struct {
	held    bool
	missing string
}

func newContractHorizon(repoRoot, prdPath string) contractHorizon {
	if _, ok, err := receiptFilePath(repoRoot, ConcreteContractGuidePath); err != nil {
		return contractHorizon{held: true}
	} else if !ok {
		return contractHorizon{missing: ConcreteContractGuidePath}
	}
	prdCommit, readable := prdAddingCommit(repoRoot, prdPath)
	if !readable {
		return contractHorizon{held: true}
	}
	output, err := adrHorizonGitOutput(repoRoot, "log", "--diff-filter=A", "--format=%H", "--", ConcreteContractGuidePath)
	if err != nil {
		return contractHorizon{held: true}
	}
	commits := strings.Fields(string(output))
	if len(commits) == 0 {
		return contractHorizon{missing: "a committed " + ConcreteContractGuidePath}
	}
	if prdCommit == "" {
		return contractHorizon{held: true}
	}
	_, err = adrHorizonGitOutput(repoRoot, "merge-base", "--is-ancestor", commits[len(commits)-1], prdCommit)
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return contractHorizon{missing: "a PRD committed at or after " + ConcreteContractGuidePath}
	}
	return contractHorizon{held: true}
}
