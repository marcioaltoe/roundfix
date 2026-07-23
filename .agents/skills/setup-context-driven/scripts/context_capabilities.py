"""Evaluate Repository Capabilities from bounded, read-only local evidence."""

from __future__ import annotations

import hashlib
import json
import re
import shutil
import stat
from dataclasses import dataclass
from enum import Enum
from pathlib import Path
from types import MappingProxyType
from typing import Callable, Iterable, Mapping


CAPABILITY_SCHEMA = "setup-context-driven/capabilities/0.0.1"
MAX_DECLARED_FILES = 32
MAX_EVIDENCE_FILE_BYTES = 1024 * 1024

_IDENTIFIER = re.compile(r"^[a-z0-9](?:[a-z0-9._-]*[a-z0-9])?$")
_SKILL_NAME = re.compile(r"^[A-Za-z0-9](?:[A-Za-z0-9._-]*[A-Za-z0-9])?$")


class _StringEnum(str, Enum):
    def __str__(self) -> str:
        return self.value


class RequirementStrength(_StringEnum):
    REQUIRED = "required"
    RECOMMENDED = "recommended"
    OPTIONAL = "optional"


class EvidenceKind(_StringEnum):
    DECLARED_FILE = "declared-file"
    EXECUTABLE = "executable"
    INSTALLED_SKILL = "installed-skill"


class EvidenceStrength(_StringEnum):
    NONE = "none"
    DECLARED = "declared"
    DISCOVERED = "discovered"
    VERIFIED = "verified"


class EvidenceStatus(_StringEnum):
    ABSENT = "absent"
    INVALID = "invalid"
    PRESENT = "present"


class CapabilityStatus(_StringEnum):
    INSUFFICIENT = "insufficient"
    MISSING = "missing"
    SATISFIED = "satisfied"


_EVIDENCE_RANK = {
    EvidenceStrength.NONE: 0,
    EvidenceStrength.DECLARED: 1,
    EvidenceStrength.DISCOVERED: 2,
    EvidenceStrength.VERIFIED: 3,
}


def _freeze_mapping(value: Mapping[str, object]) -> Mapping[str, object]:
    frozen = {
        str(key): _freeze_value(item)
        for key, item in sorted(value.items(), key=lambda pair: str(pair[0]))
    }
    return MappingProxyType(frozen)


def _freeze_value(value: object) -> object:
    if isinstance(value, Mapping):
        return _freeze_mapping(value)
    if isinstance(value, (list, tuple)):
        return tuple(_freeze_value(item) for item in value)
    if isinstance(value, set):
        return tuple(sorted((_freeze_value(item) for item in value), key=repr))
    return value


@dataclass(frozen=True)
class RepositoryCapability:
    capability_id: str
    strength: RequirementStrength
    evidence_kind: EvidenceKind
    probe: Mapping[str, object]
    title: str = ""
    minimum_evidence: EvidenceStrength = EvidenceStrength.DECLARED
    explanation: str = ""
    next_action: str = ""

    def __post_init__(self) -> None:
        if not _IDENTIFIER.fullmatch(self.capability_id):
            raise ValueError(f"invalid capability identifier {self.capability_id!r}")
        object.__setattr__(self, "strength", RequirementStrength(self.strength))
        object.__setattr__(self, "evidence_kind", EvidenceKind(self.evidence_kind))
        object.__setattr__(
            self,
            "minimum_evidence",
            EvidenceStrength(self.minimum_evidence),
        )
        if not isinstance(self.probe, Mapping):
            raise TypeError("capability probe must be a mapping")
        object.__setattr__(self, "probe", _freeze_mapping(self.probe))
        if not self.title:
            object.__setattr__(self, "title", self.capability_id)

    def __deepcopy__(self, _memo: dict[int, object]) -> RepositoryCapability:
        """Return this immutable value when a catalog mutation test is copied."""

        return self


@dataclass(frozen=True)
class CapabilityEvidence:
    capability_id: str
    status: EvidenceStatus
    version: str | None
    source_path: Path | None
    source_digest: str | None
    evidence_kind: EvidenceKind = EvidenceKind.DECLARED_FILE
    strength: EvidenceStrength = EvidenceStrength.NONE
    detail: str = ""

    def __post_init__(self) -> None:
        object.__setattr__(self, "status", EvidenceStatus(self.status))
        object.__setattr__(self, "evidence_kind", EvidenceKind(self.evidence_kind))
        object.__setattr__(self, "strength", EvidenceStrength(self.strength))
        if self.source_path is not None:
            object.__setattr__(self, "source_path", Path(self.source_path))


@dataclass(frozen=True)
class CapabilityDiagnostic:
    code: str
    message: str
    next_action: str


@dataclass(frozen=True)
class CapabilityOutcome:
    capability: RepositoryCapability
    status: CapabilityStatus
    blocking: bool
    evidence: tuple[CapabilityEvidence, ...]
    diagnostic: CapabilityDiagnostic
    explanation: str


@dataclass(frozen=True)
class CapabilityEvaluation:
    outcomes: tuple[CapabilityOutcome, ...]
    ready: bool
    guidance: str


LOCAL_RESEARCH_GUIDANCE = (
    "Search local repository code with `rg` first. Context7, Exa, and Firecrawl "
    "support external research but cannot replace local code search."
)


UNIVERSAL_CAPABILITIES = (
    RepositoryCapability(
        capability_id="capability.context7",
        title="Context7",
        strength=RequirementStrength.REQUIRED,
        evidence_kind=EvidenceKind.INSTALLED_SKILL,
        minimum_evidence=EvidenceStrength.VERIFIED,
        probe={"skill": "context7"},
        explanation="Context7 provides current authoritative library and API documentation.",
        next_action=(
            "Add the context7 skill to the installed Repository Skill Set, then rerun "
            "capability evaluation."
        ),
    ),
    RepositoryCapability(
        capability_id="capability.exa",
        title="Exa",
        strength=RequirementStrength.REQUIRED,
        evidence_kind=EvidenceKind.INSTALLED_SKILL,
        minimum_evidence=EvidenceStrength.VERIFIED,
        probe={"skill": "exa-web-search"},
        explanation="Exa provides varied broad-source searches when local and authoritative sources are insufficient.",
        next_action=(
            "Add the exa-web-search skill to the installed Repository Skill Set, then "
            "rerun capability evaluation."
        ),
    ),
    RepositoryCapability(
        capability_id="capability.firecrawl",
        title="Firecrawl",
        strength=RequirementStrength.RECOMMENDED,
        evidence_kind=EvidenceKind.INSTALLED_SKILL,
        minimum_evidence=EvidenceStrength.VERIFIED,
        probe={"skill": "firecrawl"},
        explanation="Firecrawl provides structured web-content extraction for external research.",
        next_action=(
            "Add the firecrawl skill to the installed Repository Skill Set if structured "
            "web extraction is useful."
        ),
    ),
    RepositoryCapability(
        capability_id="capability.rg",
        title="rg",
        strength=RequirementStrength.RECOMMENDED,
        evidence_kind=EvidenceKind.EXECUTABLE,
        minimum_evidence=EvidenceStrength.DISCOVERED,
        probe={"executable": "rg"},
        explanation="rg provides fast bounded local repository search.",
        next_action="Install rg and expose it on PATH if faster local search is useful.",
    ),
    RepositoryCapability(
        capability_id="capability.rtk",
        title="rtk",
        strength=RequirementStrength.RECOMMENDED,
        evidence_kind=EvidenceKind.EXECUTABLE,
        minimum_evidence=EvidenceStrength.DISCOVERED,
        probe={"executable": "rtk"},
        explanation="rtk keeps command evidence compact without changing command behavior.",
        next_action="Install rtk and expose it on PATH if compact command output is useful.",
    ),
)


def evaluate_repository_capabilities(
    repository: str | Path,
    capabilities: Iterable[RepositoryCapability] = UNIVERSAL_CAPABILITIES,
    *,
    executable_lookup: Callable[[str], str | None] = shutil.which,
) -> CapabilityEvaluation:
    """Collect bounded local evidence and evaluate the declared capabilities."""

    repository_root = Path(repository)
    ordered_capabilities = _ordered_capabilities(capabilities)
    evidence = tuple(
        _collect_evidence(repository_root, capability, executable_lookup)
        for capability in ordered_capabilities
    )
    return evaluate_capabilities(ordered_capabilities, evidence)


def evaluate_capabilities(
    capabilities: Iterable[RepositoryCapability],
    evidence: Iterable[CapabilityEvidence],
) -> CapabilityEvaluation:
    """Evaluate pre-collected evidence without probing or mutating external state."""

    ordered_capabilities = _ordered_capabilities(capabilities)
    ordered_evidence = tuple(sorted(tuple(evidence), key=_evidence_sort_key))
    outcomes = tuple(
        _evaluate_capability(
            capability,
            tuple(
                item
                for item in ordered_evidence
                if item.capability_id == capability.capability_id
                and item.evidence_kind is capability.evidence_kind
            ),
        )
        for capability in ordered_capabilities
    )
    return CapabilityEvaluation(
        outcomes=outcomes,
        ready=not any(outcome.blocking for outcome in outcomes),
        guidance=LOCAL_RESEARCH_GUIDANCE,
    )


def render_capability_json(evaluation: CapabilityEvaluation) -> bytes:
    """Render the stable machine-readable capability result."""

    document = {
        "schemaVersion": CAPABILITY_SCHEMA,
        "ready": evaluation.ready,
        "guidance": evaluation.guidance,
        "capabilities": [_outcome_document(outcome) for outcome in evaluation.outcomes],
    }
    return (
        json.dumps(document, sort_keys=True, separators=(",", ":"), ensure_ascii=True)
        + "\n"
    ).encode("utf-8")


def render_capability_guidance(evaluation: CapabilityEvaluation) -> str:
    """Render deterministic human-readable guidance and diagnostics."""

    lines = [evaluation.guidance, ""]
    for outcome in evaluation.outcomes:
        label = outcome.capability.strength.value.capitalize()
        summary = (
            f"- [{label}] {outcome.capability.title}: {outcome.status.value} — "
            f"{outcome.explanation}"
        )
        lines.append(summary)
        if outcome.status is not CapabilityStatus.SATISFIED:
            lines.append(f"  Next action: {outcome.diagnostic.next_action}")
    return "\n".join(lines).rstrip() + "\n"


def _evaluate_capability(
    capability: RepositoryCapability,
    evidence: tuple[CapabilityEvidence, ...],
) -> CapabilityOutcome:
    present = tuple(item for item in evidence if item.status is EvidenceStatus.PRESENT)
    sufficient = tuple(
        item
        for item in present
        if _EVIDENCE_RANK[item.strength]
        >= _EVIDENCE_RANK[capability.minimum_evidence]
    )
    if sufficient:
        status = CapabilityStatus.SATISFIED
        diagnostic = CapabilityDiagnostic(
            code="capability.satisfied",
            message=f"{capability.title} has sufficient local evidence.",
            next_action="",
        )
    elif present:
        status = CapabilityStatus.INSUFFICIENT
        diagnostic = CapabilityDiagnostic(
            code="capability.evidence.insufficient",
            message=(
                f"{capability.title} evidence is weaker than the required "
                f"{capability.minimum_evidence.value} strength."
            ),
            next_action=capability.next_action,
        )
    elif any(item.status is EvidenceStatus.INVALID for item in evidence):
        status = CapabilityStatus.INSUFFICIENT
        diagnostic = CapabilityDiagnostic(
            code="capability.evidence.invalid",
            message=f"{capability.title} local evidence is invalid or outside inspection bounds.",
            next_action=capability.next_action,
        )
    else:
        status = CapabilityStatus.MISSING
        code = (
            "capability.required.missing"
            if capability.strength is RequirementStrength.REQUIRED
            else "capability.recommended.missing"
            if capability.strength is RequirementStrength.RECOMMENDED
            else "capability.optional.missing"
        )
        diagnostic = CapabilityDiagnostic(
            code=code,
            message=f"{capability.title} has no compatible local evidence.",
            next_action=capability.next_action,
        )
    return CapabilityOutcome(
        capability=capability,
        status=status,
        blocking=(
            capability.strength is RequirementStrength.REQUIRED
            and status is not CapabilityStatus.SATISFIED
        ),
        evidence=evidence,
        diagnostic=diagnostic,
        explanation=capability.explanation,
    )


def _collect_evidence(
    repository: Path,
    capability: RepositoryCapability,
    executable_lookup: Callable[[str], str | None],
) -> CapabilityEvidence:
    if capability.evidence_kind is EvidenceKind.DECLARED_FILE:
        return _declared_file_evidence(repository, capability)
    if capability.evidence_kind is EvidenceKind.EXECUTABLE:
        return _executable_evidence(capability, executable_lookup)
    if capability.evidence_kind is EvidenceKind.INSTALLED_SKILL:
        return _installed_skill_evidence(repository, capability)
    raise ValueError(f"unsupported evidence kind {capability.evidence_kind!r}")


def _declared_file_evidence(
    repository: Path, capability: RepositoryCapability
) -> CapabilityEvidence:
    paths = capability.probe.get("paths", ())
    contains = capability.probe.get("contains")
    if (
        not isinstance(paths, tuple)
        or not paths
        or len(paths) > MAX_DECLARED_FILES
        or (contains is not None and not isinstance(contains, str))
    ):
        return _invalid_evidence(capability, "declared file probe is invalid")
    for relative in sorted(paths):
        if not isinstance(relative, str):
            return _invalid_evidence(capability, "declared file path is invalid")
        candidate = _bounded_regular_file(repository, relative)
        if candidate is None:
            continue
        try:
            size = candidate.stat(follow_symlinks=False).st_size
            if size > MAX_EVIDENCE_FILE_BYTES:
                return _invalid_evidence(capability, "declared file exceeds byte limit")
            content = candidate.read_bytes()
        except OSError:
            return _invalid_evidence(capability, "declared file cannot be read")
        if contains is not None and contains.encode("utf-8") not in content:
            continue
        return CapabilityEvidence(
            capability_id=capability.capability_id,
            status=EvidenceStatus.PRESENT,
            version=_probe_version(capability.probe),
            source_path=Path(relative),
            source_digest=hashlib.sha256(content).hexdigest(),
            evidence_kind=capability.evidence_kind,
            strength=EvidenceStrength.DECLARED,
        )
    return _absent_evidence(capability)


def _executable_evidence(
    capability: RepositoryCapability,
    executable_lookup: Callable[[str], str | None],
) -> CapabilityEvidence:
    executable = capability.probe.get("executable")
    if not isinstance(executable, str) or not _SKILL_NAME.fullmatch(executable):
        return _invalid_evidence(capability, "executable probe is invalid")
    try:
        discovered = executable_lookup(executable)
    except OSError:
        return _invalid_evidence(capability, "executable lookup failed")
    if not discovered:
        return _absent_evidence(capability)
    return CapabilityEvidence(
        capability_id=capability.capability_id,
        status=EvidenceStatus.PRESENT,
        version=_probe_version(capability.probe),
        source_path=Path(discovered),
        source_digest=None,
        evidence_kind=capability.evidence_kind,
        strength=EvidenceStrength.DISCOVERED,
    )


def _installed_skill_evidence(
    repository: Path, capability: RepositoryCapability
) -> CapabilityEvidence:
    skill = capability.probe.get("skill")
    roots = capability.probe.get("roots", (".agents/skills",))
    if (
        not isinstance(skill, str)
        or not _SKILL_NAME.fullmatch(skill)
        or not isinstance(roots, tuple)
        or not roots
        or len(roots) > MAX_DECLARED_FILES
    ):
        return _invalid_evidence(capability, "installed skill probe is invalid")
    for root in sorted(roots):
        if not isinstance(root, str) or _unsafe_relative_path(root):
            return _invalid_evidence(capability, "installed skill root is invalid")
        relative = (Path(root) / skill / "SKILL.md").as_posix()
        candidate = _bounded_regular_file(repository, relative)
        if candidate is None:
            continue
        try:
            content = candidate.read_bytes()
        except OSError:
            return _invalid_evidence(capability, "installed skill cannot be read")
        if len(content) > MAX_EVIDENCE_FILE_BYTES:
            return _invalid_evidence(capability, "installed skill exceeds byte limit")
        return CapabilityEvidence(
            capability_id=capability.capability_id,
            status=EvidenceStatus.PRESENT,
            version=_probe_version(capability.probe),
            source_path=Path(relative),
            source_digest=hashlib.sha256(content).hexdigest(),
            evidence_kind=capability.evidence_kind,
            strength=EvidenceStrength.VERIFIED,
        )
    return _absent_evidence(capability)


def _bounded_regular_file(repository: Path, relative: str) -> Path | None:
    if _unsafe_relative_path(relative):
        return None
    current = repository
    parts = Path(relative).parts
    for index, part in enumerate(parts):
        current = current / part
        try:
            entry_stat = current.lstat()
        except OSError:
            return None
        if stat.S_ISLNK(entry_stat.st_mode):
            return None
        if index < len(parts) - 1 and not stat.S_ISDIR(entry_stat.st_mode):
            return None
    return current if stat.S_ISREG(entry_stat.st_mode) else None


def _unsafe_relative_path(value: str) -> bool:
    candidate = Path(value)
    return (
        not value
        or value != candidate.as_posix()
        or candidate.as_posix() == "."
        or candidate.is_absolute()
        or ".." in candidate.parts
        or "\\" in value
        or re.match(r"^[A-Za-z]:", value) is not None
    )


def _absent_evidence(capability: RepositoryCapability) -> CapabilityEvidence:
    return CapabilityEvidence(
        capability_id=capability.capability_id,
        status=EvidenceStatus.ABSENT,
        version=None,
        source_path=None,
        source_digest=None,
        evidence_kind=capability.evidence_kind,
        strength=EvidenceStrength.NONE,
    )


def _invalid_evidence(
    capability: RepositoryCapability, detail: str
) -> CapabilityEvidence:
    return CapabilityEvidence(
        capability_id=capability.capability_id,
        status=EvidenceStatus.INVALID,
        version=None,
        source_path=None,
        source_digest=None,
        evidence_kind=capability.evidence_kind,
        strength=EvidenceStrength.NONE,
        detail=detail,
    )


def _ordered_capabilities(
    capabilities: Iterable[RepositoryCapability],
) -> tuple[RepositoryCapability, ...]:
    ordered = tuple(sorted(tuple(capabilities), key=lambda item: item.capability_id))
    identifiers = tuple(capability.capability_id for capability in ordered)
    if len(identifiers) != len(set(identifiers)):
        raise ValueError("capability identifiers must be unique")
    return ordered


def _evidence_sort_key(evidence: CapabilityEvidence) -> tuple[str, ...]:
    return (
        evidence.capability_id,
        evidence.evidence_kind.value,
        evidence.strength.value,
        evidence.status.value,
        evidence.source_path.as_posix() if evidence.source_path is not None else "",
        evidence.source_digest or "",
        evidence.version or "",
        evidence.detail,
    )


def _outcome_document(outcome: CapabilityOutcome) -> dict[str, object]:
    return {
        "id": outcome.capability.capability_id,
        "title": outcome.capability.title,
        "requirement": outcome.capability.strength.value,
        "evidenceKind": outcome.capability.evidence_kind.value,
        "minimumEvidence": outcome.capability.minimum_evidence.value,
        "status": outcome.status.value,
        "blocking": outcome.blocking,
        "explanation": outcome.explanation,
        "diagnostic": {
            "code": outcome.diagnostic.code,
            "message": outcome.diagnostic.message,
            "nextAction": outcome.diagnostic.next_action,
        },
        "evidence": [
            {
                "status": evidence.status.value,
                "kind": evidence.evidence_kind.value,
                "strength": evidence.strength.value,
                "version": evidence.version,
                "sourcePath": evidence.source_path.as_posix()
                if evidence.source_path is not None
                else None,
                "sourceDigest": evidence.source_digest,
                "detail": evidence.detail,
            }
            for evidence in outcome.evidence
        ],
    }


def _probe_version(probe: Mapping[str, object]) -> str | None:
    version = probe.get("version")
    return version if isinstance(version, str) else None
