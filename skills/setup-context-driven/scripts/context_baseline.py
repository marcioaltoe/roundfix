"""Load and validate immutable setup-context-driven Source Baselines."""

from __future__ import annotations

import base64
import binascii
import hashlib
import json
import os
import re
import stat
from dataclasses import dataclass
from pathlib import Path, PurePosixPath


SOURCE_BASELINE_VERSION = "0.0.1"
SOURCE_BASELINE_SCHEMA = "setup-context-driven/source-baseline/0.0.1"
SOURCE_BASELINE_MANIFEST_SCHEMA = (
    "setup-context-driven/source-baseline-manifest/0.0.1"
)
SOURCE_BASELINE_INDEX_SCHEMA = "setup-context-driven/source-baseline-index/0.0.1"
DECISION_DOCUMENT_SCHEMA = "setup-context-driven/decisions/0.0.1"
DECISION_DOCUMENT_VERSION = "0.0.1"

_ENTRY_KINDS = {"normative-clause", "recommendation", "operational-contract"}
_NORMATIVE_ENFORCEMENTS = {"mandatory", "prohibited", "stop-and-ask"}
_OPERATIONAL_STRUCTURES = {
    "decision-matrix",
    "lifecycle",
    "ordered-procedure",
    "protocol",
    "template",
}
_IDENTIFIER = re.compile(r"^[a-z0-9](?:[a-z0-9._-]*[a-z0-9])?$")
_SHA256 = re.compile(r"^[0-9a-f]{64}$")
_OPEN_MARKER = re.compile(
    rb"(?m)^<!-- source-baseline-entry: ([a-z0-9][a-z0-9._-]*) -->\r?\n"
)
_MANAGED_BEGIN_MARKER = re.compile(
    rb"<!--\s*setup-context-driven:begin\s+id=([A-Za-z0-9_.-]+)"
    rb"\s+version=([0-9]+)\s*-->"
)
_INSTRUCTION_NAMES = frozenset({"AGENTS.md", "CLAUDE.md"})
_INVENTORY_IGNORED_DIRECTORIES = frozenset(
    {
        ".git",
        ".hg",
        ".svn",
        ".venv",
        "__pycache__",
        "node_modules",
        "vendor",
        "venv",
    }
)
_INVENTORY_IGNORED_PREFIXES = (
    (".agents", "skills"),
    ("skills",),
)
_MANIFEST_CARRIER = Path("docs/agents/setup-context.json")
_REPOSITORY_RULES_PATH = PurePosixPath("docs/agents/repository-rules.md")
_READOPTION_CLASSIFICATIONS = {
    "non-governed",
    "normative-clause",
    "operational-contract",
    "recommendation",
}
_READOPTION_DISPOSITIONS = {
    "managed-entry",
    "rejected",
    "repository-document",
    "repository-rules",
}
_TYPED_DOCUMENTS = {
    "agent-guide",
    "architecture-decision",
    "design-contract",
    "domain-context",
    "http-contract",
}


class SourceBaselineValidationError(ValueError):
    """Raised when Source Baseline documents disagree or are invalid."""

    def __init__(self, diagnostics: list[str]):
        self.diagnostics = tuple(diagnostics)
        super().__init__("\n".join(diagnostics))


@dataclass(frozen=True)
class SourceInventoryDiagnostic:
    code: str
    path: Path
    message: str


class SourceInventoryError(ValueError):
    """Raised when a bounded carrier cannot be inventoried safely."""

    def __init__(self, diagnostics: list[SourceInventoryDiagnostic]):
        self.diagnostics = tuple(diagnostics)
        super().__init__(
            "\n".join(
                f"{item.code}: {item.path.as_posix()}: {item.message}"
                for item in diagnostics
            )
        )


@dataclass(frozen=True)
class DecisionDocumentDiagnostic:
    code: str
    path: Path
    item_id: str
    message: str


class DecisionDocumentError(ValueError):
    """Raised when a structured decision document is malformed or unsafe."""

    def __init__(self, diagnostics: list[DecisionDocumentDiagnostic]):
        self.diagnostics = tuple(diagnostics)
        super().__init__(
            "\n".join(
                f"{item.code}: {item.path.as_posix()}: {item.item_id}: {item.message}"
                for item in diagnostics
            )
        )


@dataclass(frozen=True)
class SourceInventoryLimits:
    max_files: int = 256
    max_file_bytes: int = 2 * 1024 * 1024
    max_total_bytes: int = 8 * 1024 * 1024


@dataclass(frozen=True)
class ReadoptionSourceEntry:
    entry_id: str
    path: Path
    kind: str
    start: int
    end: int
    digest: str
    carrier_digest: str
    source_bytes: bytes
    provenance: tuple[tuple[str, object], ...]

    def to_json(self) -> dict[str, object]:
        return {
            "id": self.entry_id,
            "path": self.path.as_posix(),
            "carrier": self.path.as_posix(),
            "kind": self.kind,
            "start": self.start,
            "end": self.end,
            "digest": self.digest,
            "carrierDigest": self.carrier_digest,
            "sourceBytes": base64.b64encode(self.source_bytes).decode("ascii"),
            "encoding": "base64",
            "structuralProvenance": dict(self.provenance),
        }


@dataclass(frozen=True)
class IncompatibleSourceBaseline:
    baseline_id: str
    declared_identity: str
    digest: str
    carriers: tuple[Path, ...]
    entries: tuple[ReadoptionSourceEntry, ...]
    byte_count: int

    def identity_json(self) -> dict[str, object]:
        return {
            "id": self.baseline_id,
            "declaredIdentity": self.declared_identity,
            "compatibility": "incompatible",
            "digest": self.digest,
            "carrierCount": len(self.carriers),
            "entryCount": len(self.entries),
            "byteCount": self.byte_count,
        }


@dataclass(frozen=True)
class ReadoptionDisposition:
    entry_id: str
    entry_digest: str
    classification: str
    disposition: str
    destination: dict[str, object] | None
    reason: str

    def to_json(self) -> dict[str, object]:
        return {
            "entryId": self.entry_id,
            "entryDigest": self.entry_digest,
            "classification": self.classification,
            "disposition": self.disposition,
            "destination": self.destination,
            "reason": self.reason,
        }


@dataclass(frozen=True)
class ReadoptionDecisions:
    source_baseline_id: str
    source_baseline_digest: str
    dispositions: tuple[ReadoptionDisposition, ...]

    def to_json(self) -> dict[str, object]:
        return {
            "sourceBaseline": {
                "id": self.source_baseline_id,
                "digest": self.source_baseline_digest,
            },
            "dispositions": [item.to_json() for item in self.dispositions],
        }


@dataclass(frozen=True)
class StructuredDecisionDocument:
    source_paths: tuple[Path, ...]
    source_digests: tuple[str, ...]
    decisions: tuple[tuple[str, object], ...]
    readoption: ReadoptionDecisions | None

    def to_json(self) -> dict[str, object]:
        data: dict[str, object] = {
            "schemaVersion": DECISION_DOCUMENT_SCHEMA,
            "version": DECISION_DOCUMENT_VERSION,
            "decisions": [
                {"id": decision_id, "value": value}
                for decision_id, value in self.decisions
            ],
        }
        if self.readoption is not None:
            data["readoption"] = self.readoption.to_json()
        return data


def load_decision_document(path: str | Path) -> StructuredDecisionDocument:
    """Load one strict structured decision document without repository writes."""

    source_path = Path(path).expanduser()
    if not source_path.is_absolute():
        source_path = Path.cwd() / source_path
    source_path = source_path.resolve(strict=False)
    try:
        source_stat = source_path.lstat()
    except OSError as error:
        raise _decision_error(
            "decision-file.read",
            source_path,
            "decision-file",
            f"cannot read decision file: {error}",
        ) from error
    if stat.S_ISLNK(source_stat.st_mode) or not stat.S_ISREG(source_stat.st_mode):
        raise _decision_error(
            "decision-file.type.invalid",
            source_path,
            "decision-file",
            "decision file must be a regular non-symlink file",
        )
    try:
        content = source_path.read_bytes()
    except OSError as error:
        raise _decision_error(
            "decision-file.read",
            source_path,
            "decision-file",
            f"cannot read decision file: {error}",
        ) from error
    try:
        document = json.loads(content, object_pairs_hook=_strict_json_object)
    except _DuplicateJSONKey as error:
        raise _decision_error(
            "decision-file.json.duplicate-key",
            source_path,
            error.key,
            f"duplicate JSON object key {error.key!r}",
        ) from error
    except (UnicodeDecodeError, json.JSONDecodeError) as error:
        raise _decision_error(
            "decision-file.json.invalid",
            source_path,
            "decision-file",
            f"decision file is not valid UTF-8 JSON: {error}",
        ) from error

    decisions, readoption = _parse_decision_document(source_path, document)
    return StructuredDecisionDocument(
        source_paths=(source_path,),
        source_digests=(hashlib.sha256(content).hexdigest(),),
        decisions=decisions,
        readoption=readoption,
    )


def merge_decision_documents(
    documents: tuple[StructuredDecisionDocument, ...],
) -> StructuredDecisionDocument | None:
    """Merge repeated decision files while rejecting conflicting public inputs."""

    if not documents:
        return None
    merged_decisions: dict[str, object] = {}
    readoption: ReadoptionDecisions | None = None
    diagnostics: list[DecisionDocumentDiagnostic] = []
    for document in documents:
        source_path = document.source_paths[0]
        for decision_id, value in document.decisions:
            if decision_id in merged_decisions and merged_decisions[decision_id] != value:
                diagnostics.append(
                    DecisionDocumentDiagnostic(
                        "decision-file.decision.conflict",
                        source_path,
                        decision_id,
                        "repeated decision files contain conflicting values",
                    )
                )
                continue
            merged_decisions[decision_id] = value
        if document.readoption is None:
            continue
        if readoption is not None:
            diagnostics.append(
                DecisionDocumentDiagnostic(
                    "decision-file.readoption.duplicate",
                    source_path,
                    "readoption",
                    "only one repeated decision file may contain Readoption dispositions",
                )
            )
            continue
        readoption = document.readoption
    if diagnostics:
        raise DecisionDocumentError(diagnostics)
    return StructuredDecisionDocument(
        source_paths=tuple(path for item in documents for path in item.source_paths),
        source_digests=tuple(
            digest for item in documents for digest in item.source_digests
        ),
        decisions=tuple(sorted(merged_decisions.items())),
        readoption=readoption,
    )


def validate_readoption_decisions(
    repo: str | Path,
    inventory: IncompatibleSourceBaseline,
    decisions: ReadoptionDecisions,
) -> tuple[tuple[ReadoptionDisposition, ...], tuple[DecisionDocumentDiagnostic, ...]]:
    """Validate individual coverage, stale evidence, and typed destinations."""

    root = Path(repo)
    document_path = Path("decision-file")
    diagnostics: list[DecisionDocumentDiagnostic] = []
    if (
        decisions.source_baseline_id != inventory.baseline_id
        or decisions.source_baseline_digest != inventory.digest
    ):
        diagnostics.append(
            DecisionDocumentDiagnostic(
                "readoption.source.stale",
                document_path,
                decisions.source_baseline_id,
                "decision file Source Baseline identity does not match current bytes",
            )
        )

    expected_by_id = {entry.entry_id: entry for entry in inventory.entries}
    by_id: dict[str, list[ReadoptionDisposition]] = {}
    for disposition in decisions.dispositions:
        by_id.setdefault(disposition.entry_id, []).append(disposition)

    for entry_id, candidates in sorted(by_id.items()):
        if len(candidates) > 1:
            diagnostics.append(
                DecisionDocumentDiagnostic(
                    "readoption.disposition.duplicate",
                    document_path,
                    entry_id,
                    "Source Baseline Entry has more than one disposition",
                )
            )
        entry = expected_by_id.get(entry_id)
        if entry is None:
            diagnostics.append(
                DecisionDocumentDiagnostic(
                    "readoption.disposition.unknown",
                    document_path,
                    entry_id,
                    "disposition names an entry absent from the current inventory",
                )
            )
            continue
        candidate = candidates[0]
        if candidate.entry_digest != entry.digest:
            diagnostics.append(
                DecisionDocumentDiagnostic(
                    "readoption.disposition.stale",
                    document_path,
                    entry_id,
                    "disposition entryDigest does not match the current source bytes",
                )
            )
        diagnostics.extend(_validate_readoption_destination(root, candidate))

    for entry in inventory.entries:
        if entry.entry_id not in by_id:
            diagnostics.append(
                DecisionDocumentDiagnostic(
                    "readoption.disposition.missing",
                    entry.path,
                    entry.entry_id,
                    "Source Baseline Entry has no disposition",
                )
            )

    inventory_order = {entry.entry_id: index for index, entry in enumerate(inventory.entries)}
    ordered = tuple(
        sorted(
            decisions.dispositions,
            key=lambda item: (inventory_order.get(item.entry_id, len(inventory_order)), item.entry_id),
        )
    )
    return ordered, tuple(diagnostics)


def repository_rules_proposed_bytes(
    dispositions: tuple[ReadoptionDisposition, ...],
) -> bytes:
    """Return only exact decision-supplied Repository-Specific Normative Rules bytes."""

    chunks: list[bytes] = []
    for item in dispositions:
        if item.disposition != "repository-rules" or item.destination is None:
            continue
        chunks.append(
            base64.b64decode(str(item.destination["proposedBytes"]), validate=True)
        )
    return b"".join(chunks)


class _DuplicateJSONKey(ValueError):
    def __init__(self, key: str):
        self.key = key
        super().__init__(key)


def _strict_json_object(pairs: list[tuple[str, object]]) -> dict[str, object]:
    result: dict[str, object] = {}
    for key, value in pairs:
        if key in result:
            raise _DuplicateJSONKey(key)
        result[key] = value
    return result


def _decision_error(
    code: str, path: Path, item_id: str, message: str
) -> DecisionDocumentError:
    return DecisionDocumentError(
        [DecisionDocumentDiagnostic(code, path, item_id, message)]
    )


def _parse_decision_document(
    source_path: Path, document: object
) -> tuple[tuple[tuple[str, object], ...], ReadoptionDecisions | None]:
    if not isinstance(document, dict):
        raise _decision_error(
            "decision-file.schema.invalid",
            source_path,
            "decision-file",
            "decision document must be a JSON object",
        )
    allowed_fields = {"schemaVersion", "version", "decisions", "readoption"}
    if (
        not {"schemaVersion", "version", "decisions"}.issubset(document)
        or not set(document).issubset(allowed_fields)
        or document.get("schemaVersion") != DECISION_DOCUMENT_SCHEMA
        or document.get("version") != DECISION_DOCUMENT_VERSION
    ):
        raise _decision_error(
            "decision-file.schema.invalid",
            source_path,
            "decision-file",
            (
                "decision document requires schemaVersion, version, and decisions "
                "with only the optional readoption field"
            ),
        )
    raw_decisions = document.get("decisions")
    if not isinstance(raw_decisions, list):
        raise _decision_error(
            "decision-file.decisions.invalid",
            source_path,
            "decisions",
            "decisions must be an ordered array",
        )
    decisions: dict[str, object] = {}
    diagnostics: list[DecisionDocumentDiagnostic] = []
    for index, record in enumerate(raw_decisions):
        item_id = f"decisions[{index}]"
        if not isinstance(record, dict) or set(record) != {"id", "value"}:
            diagnostics.append(
                DecisionDocumentDiagnostic(
                    "decision-file.decision.invalid",
                    source_path,
                    item_id,
                    "decision record requires exactly id and value",
                )
            )
            continue
        decision_id = record.get("id")
        if not isinstance(decision_id, str) or _IDENTIFIER.fullmatch(decision_id) is None:
            diagnostics.append(
                DecisionDocumentDiagnostic(
                    "decision-file.decision.invalid",
                    source_path,
                    item_id,
                    "decision id is invalid",
                )
            )
            continue
        if decision_id in decisions:
            diagnostics.append(
                DecisionDocumentDiagnostic(
                    "decision-file.decision.duplicate",
                    source_path,
                    decision_id,
                    "decision id appears more than once",
                )
            )
            continue
        decisions[decision_id] = record.get("value")
    readoption = None
    if "readoption" in document:
        try:
            readoption = _parse_readoption_decisions(source_path, document["readoption"])
        except DecisionDocumentError as error:
            diagnostics.extend(error.diagnostics)
    if diagnostics:
        raise DecisionDocumentError(diagnostics)
    return tuple(sorted(decisions.items())), readoption


def _parse_readoption_decisions(
    source_path: Path, raw: object
) -> ReadoptionDecisions:
    if not isinstance(raw, dict) or set(raw) != {"sourceBaseline", "dispositions"}:
        raise _decision_error(
            "decision-file.readoption.invalid",
            source_path,
            "readoption",
            "readoption requires exactly sourceBaseline and dispositions",
        )
    source = raw.get("sourceBaseline")
    if (
        not isinstance(source, dict)
        or set(source) != {"id", "digest"}
        or not isinstance(source.get("id"), str)
        or not isinstance(source.get("digest"), str)
        or _SHA256.fullmatch(source["digest"]) is None
    ):
        raise _decision_error(
            "decision-file.readoption.source.invalid",
            source_path,
            "sourceBaseline",
            "sourceBaseline requires an id and lowercase SHA-256 digest",
        )
    raw_dispositions = raw.get("dispositions")
    if not isinstance(raw_dispositions, list):
        raise _decision_error(
            "decision-file.readoption.dispositions.invalid",
            source_path,
            "dispositions",
            "dispositions must be an ordered array",
        )
    dispositions: list[ReadoptionDisposition] = []
    diagnostics: list[DecisionDocumentDiagnostic] = []
    for index, item in enumerate(raw_dispositions):
        try:
            dispositions.append(_parse_readoption_disposition(source_path, index, item))
        except DecisionDocumentError as error:
            diagnostics.extend(error.diagnostics)
    if diagnostics:
        raise DecisionDocumentError(diagnostics)
    return ReadoptionDecisions(
        source_baseline_id=source["id"],
        source_baseline_digest=source["digest"],
        dispositions=tuple(dispositions),
    )


def _parse_readoption_disposition(
    source_path: Path, index: int, raw: object
) -> ReadoptionDisposition:
    item_id = f"dispositions[{index}]"
    expected_fields = {
        "entryId",
        "entryDigest",
        "classification",
        "disposition",
        "destination",
        "reason",
    }
    if not isinstance(raw, dict) or set(raw) != expected_fields:
        raise _decision_error(
            "readoption.disposition.invalid",
            source_path,
            str(raw.get("entryId", item_id)) if isinstance(raw, dict) else item_id,
            "disposition record fields are invalid",
        )
    entry_id = raw.get("entryId")
    entry_digest = raw.get("entryDigest")
    classification = raw.get("classification")
    disposition = raw.get("disposition")
    reason = raw.get("reason")
    managed_id = str(entry_id) if isinstance(entry_id, str) else item_id
    if (
        not isinstance(entry_id, str)
        or re.fullmatch(r"source-entry\.[0-9a-f]{64}", entry_id) is None
        or not isinstance(entry_digest, str)
        or _SHA256.fullmatch(entry_digest) is None
        or classification not in _READOPTION_CLASSIFICATIONS
        or disposition not in _READOPTION_DISPOSITIONS
        or not isinstance(reason, str)
    ):
        raise _decision_error(
            "readoption.disposition.invalid",
            source_path,
            managed_id,
            "entry evidence, classification, disposition, or reason is invalid",
        )
    destination = _parse_readoption_destination(
        source_path, managed_id, disposition, raw.get("destination")
    )
    if classification == "non-governed" and disposition != "rejected":
        raise _decision_error(
            "readoption.disposition.invalid",
            source_path,
            managed_id,
            "non-governed evidence must use the rejected disposition",
        )
    if (classification == "non-governed" or disposition == "rejected") and not reason.strip():
        raise _decision_error(
            "readoption.disposition.reason.required",
            source_path,
            managed_id,
            "non-governed or rejected evidence requires an individual reason",
        )
    return ReadoptionDisposition(
        entry_id=entry_id,
        entry_digest=entry_digest,
        classification=classification,
        disposition=disposition,
        destination=destination,
        reason=reason,
    )


def _parse_readoption_destination(
    source_path: Path, entry_id: str, disposition: str, destination: object
) -> dict[str, object] | None:
    if disposition == "rejected":
        if destination is not None:
            raise _decision_error(
                "readoption.disposition.invalid",
                source_path,
                entry_id,
                "rejected disposition must have a null destination",
            )
        return None
    if not isinstance(destination, dict):
        raise _decision_error(
            "readoption.disposition.invalid",
            source_path,
            entry_id,
            "non-rejected disposition requires a typed destination object",
        )
    if disposition == "managed-entry":
        if set(destination) != {"managedId"} or not isinstance(
            destination.get("managedId"), str
        ):
            raise _decision_error(
                "readoption.disposition.invalid",
                source_path,
                entry_id,
                "managed-entry destination requires exactly managedId",
            )
        return {"managedId": destination["managedId"]}
    if disposition == "repository-document":
        if set(destination) != {"documentType", "path", "digest"}:
            raise _decision_error(
                "readoption.disposition.invalid",
                source_path,
                entry_id,
                "repository-document destination fields are invalid",
            )
        if destination.get("documentType") not in _TYPED_DOCUMENTS:
            raise _decision_error(
                "readoption.destination.document-type.invalid",
                source_path,
                entry_id,
                "repository documentType is not supported",
            )
        if not isinstance(destination.get("path"), str) or not isinstance(
            destination.get("digest"), str
        ) or _SHA256.fullmatch(destination["digest"]) is None:
            raise _decision_error(
                "readoption.disposition.invalid",
                source_path,
                entry_id,
                "repository-document destination path or digest is invalid",
            )
        return dict(destination)
    expected = {"documentType", "path", "proposedBytes", "digest"}
    if set(destination) != expected:
        raise _decision_error(
            "readoption.disposition.invalid",
            source_path,
            entry_id,
            "repository-rules destination fields are invalid",
        )
    if (
        destination.get("documentType") != "repository-rules"
        or destination.get("path") != _REPOSITORY_RULES_PATH.as_posix()
        or not isinstance(destination.get("proposedBytes"), str)
        or not isinstance(destination.get("digest"), str)
        or _SHA256.fullmatch(destination["digest"]) is None
    ):
        raise _decision_error(
            "readoption.destination.repository-rules.invalid",
            source_path,
            entry_id,
            "Repository-Specific Normative Rules require the default typed path and exact bytes",
        )
    try:
        proposed = base64.b64decode(destination["proposedBytes"], validate=True)
    except (ValueError, binascii.Error) as error:
        raise _decision_error(
            "readoption.destination.proposed-bytes.invalid",
            source_path,
            entry_id,
            "proposedBytes must be canonical base64",
        ) from error
    if not proposed or hashlib.sha256(proposed).hexdigest() != destination["digest"]:
        raise _decision_error(
            "readoption.destination.proposed-bytes.stale",
            source_path,
            entry_id,
            "proposed bytes are empty or do not match the declared digest",
        )
    if _MANAGED_BEGIN_MARKER.search(proposed) is not None or b"<!-- setup-context-driven:" in proposed:
        raise _decision_error(
            "readoption.destination.proposed-bytes.managed-marker",
            source_path,
            entry_id,
            "Repository-Specific Normative Rules proposed bytes must remain unmarked",
        )
    return dict(destination)


def _validate_readoption_destination(
    root: Path, disposition: ReadoptionDisposition
) -> list[DecisionDocumentDiagnostic]:
    destination = disposition.destination
    if destination is None or disposition.disposition in {"rejected", "managed-entry"}:
        return []
    entry_id = disposition.entry_id
    raw_path = destination.get("path")
    if not isinstance(raw_path, str) or not _decision_destination_path_is_safe(raw_path):
        return [
            DecisionDocumentDiagnostic(
                "readoption.destination.path.unsafe",
                Path("decision-file"),
                entry_id,
                "destination path must be a safe repository-relative path",
            )
        ]
    relative = PurePosixPath(raw_path)
    if disposition.disposition == "repository-rules":
        return []
    document_type = str(destination["documentType"])
    if not _typed_document_path_matches(document_type, relative):
        return [
            DecisionDocumentDiagnostic(
                "readoption.destination.document-type.invalid",
                Path(raw_path),
                entry_id,
                "destination path does not match its declared documentType",
            )
        ]
    target = root.joinpath(*relative.parts)
    try:
        target_stat = target.lstat()
    except OSError:
        return [
            DecisionDocumentDiagnostic(
                "readoption.destination.missing",
                Path(raw_path),
                entry_id,
                "typed repository document does not exist",
            )
        ]
    if stat.S_ISLNK(target_stat.st_mode) or not stat.S_ISREG(target_stat.st_mode):
        return [
            DecisionDocumentDiagnostic(
                "readoption.destination.type.invalid",
                Path(raw_path),
                entry_id,
                "typed repository document must be a regular non-symlink file",
            )
        ]
    try:
        observed = hashlib.sha256(target.read_bytes()).hexdigest()
    except OSError:
        return [
            DecisionDocumentDiagnostic(
                "readoption.destination.read",
                Path(raw_path),
                entry_id,
                "typed repository document cannot be read",
            )
        ]
    if observed != destination["digest"]:
        return [
            DecisionDocumentDiagnostic(
                "readoption.destination.stale",
                Path(raw_path),
                entry_id,
                "typed repository document digest does not match current bytes",
            )
        ]
    return []


def _decision_destination_path_is_safe(value: str) -> bool:
    if not value or "\\" in value:
        return False
    path = PurePosixPath(value)
    return not path.is_absolute() and all(
        part not in {"", ".", ".."} for part in path.parts
    )


def _typed_document_path_matches(document_type: str, path: PurePosixPath) -> bool:
    parts = path.parts
    if document_type == "agent-guide":
        return (
            len(parts) >= 3
            and parts[:2] == ("docs", "agents")
            and path.suffix == ".md"
            and path != _REPOSITORY_RULES_PATH
        )
    if document_type == "architecture-decision":
        return len(parts) == 3 and parts[:2] == ("docs", "adr") and path.suffix == ".md"
    if document_type == "design-contract":
        return path.name == "DESIGN.md"
    if document_type == "domain-context":
        return path.name == "CONTEXT.md"
    if document_type == "http-contract":
        return len(parts) >= 3 and parts[:2] == ("docs", "architecture") and path.suffix == ".json"
    return False


def inventory_incompatible_source_baseline(
    repo: str | Path,
    declared_identity: str,
    *,
    limits: SourceInventoryLimits | None = None,
) -> IncompatibleSourceBaseline:
    """Inventory bounded historical setup carriers without classifying meaning."""

    root = Path(repo)
    active_limits = limits or SourceInventoryLimits()
    paths = _discover_inventory_carriers(root)
    if len(paths) > active_limits.max_files:
        raise SourceInventoryError(
            [
                SourceInventoryDiagnostic(
                    "source-inventory.limit.files",
                    Path("."),
                    f"carrier count {len(paths)} exceeds {active_limits.max_files}",
                )
            ]
        )

    contents: list[tuple[Path, bytes]] = []
    total_bytes = 0
    diagnostics: list[SourceInventoryDiagnostic] = []
    for relative_path in paths:
        try:
            content = _read_inventory_carrier(
                root, relative_path, active_limits.max_file_bytes
            )
        except SourceInventoryError as error:
            diagnostics.extend(error.diagnostics)
            continue
        total_bytes += len(content)
        if total_bytes > active_limits.max_total_bytes:
            diagnostics.append(
                SourceInventoryDiagnostic(
                    "source-inventory.limit.total-bytes",
                    relative_path,
                    f"carrier bytes exceed {active_limits.max_total_bytes}",
                )
            )
            break
        contents.append((relative_path, content))

    if diagnostics:
        raise SourceInventoryError(diagnostics)

    entries: list[ReadoptionSourceEntry] = []
    identity_digest = hashlib.sha256()
    for relative_path, content in contents:
        path_bytes = relative_path.as_posix().encode("utf-8")
        identity_digest.update(len(path_bytes).to_bytes(8, "big"))
        identity_digest.update(path_bytes)
        identity_digest.update(len(content).to_bytes(8, "big"))
        identity_digest.update(content)
        if relative_path == _MANIFEST_CARRIER:
            _validate_inventory_manifest_paths(relative_path, content)
        entries.extend(_partition_inventory_carrier(relative_path, content))

    digest = identity_digest.hexdigest()
    entries.sort(
        key=lambda entry: (
            entry.path.as_posix(),
            entry.start,
            entry.end,
            entry.kind,
        )
    )
    return IncompatibleSourceBaseline(
        baseline_id=f"baseline.readoption.{digest}",
        declared_identity=declared_identity,
        digest=digest,
        carriers=tuple(path for path, _ in contents),
        entries=tuple(entries),
        byte_count=total_bytes,
    )


def _discover_inventory_carriers(root: Path) -> tuple[Path, ...]:
    diagnostics: list[SourceInventoryDiagnostic] = []
    try:
        root_stat = root.lstat()
    except OSError as error:
        raise SourceInventoryError(
            [
                SourceInventoryDiagnostic(
                    "source-inventory.root.invalid", Path("."), str(error)
                )
            ]
        ) from error
    if stat.S_ISLNK(root_stat.st_mode) or not stat.S_ISDIR(root_stat.st_mode):
        raise SourceInventoryError(
            [
                SourceInventoryDiagnostic(
                    "source-inventory.root.invalid",
                    Path("."),
                    "repository root must be a real directory",
                )
            ]
        )

    carriers: set[Path] = set()

    def walk(directory: Path, relative_directory: Path) -> None:
        try:
            children = sorted(os.scandir(directory), key=lambda child: child.name)
        except OSError as error:
            diagnostics.append(
                SourceInventoryDiagnostic(
                    "source-inventory.directory.read",
                    relative_directory,
                    str(error),
                )
            )
            return
        for child in children:
            relative = relative_directory / child.name
            if _inventory_path_is_ignored(relative):
                continue
            if relative.parts[:2] == ("docs", "agents"):
                continue
            try:
                child_stat = child.stat(follow_symlinks=False)
            except OSError as error:
                diagnostics.append(
                    SourceInventoryDiagnostic(
                        "source-inventory.carrier.read", relative, str(error)
                    )
                )
                continue
            if stat.S_ISLNK(child_stat.st_mode):
                if child.name in _INSTRUCTION_NAMES:
                    diagnostics.append(
                        SourceInventoryDiagnostic(
                            "source-inventory.carrier.symlink",
                            relative,
                            "instruction carrier must not be a symlink",
                        )
                    )
                continue
            if stat.S_ISDIR(child_stat.st_mode):
                walk(Path(child.path), relative)
                continue
            if child.name not in _INSTRUCTION_NAMES:
                continue
            if stat.S_ISREG(child_stat.st_mode):
                carriers.add(relative)
            else:
                diagnostics.append(
                    SourceInventoryDiagnostic(
                        "source-inventory.carrier.special",
                        relative,
                        "instruction carrier must be a regular file",
                    )
                )

    walk(root, Path())
    _discover_docs_agent_carriers(root, carriers, diagnostics)
    if diagnostics:
        raise SourceInventoryError(
            sorted(diagnostics, key=lambda item: (item.path.as_posix(), item.code))
        )
    return tuple(sorted(carriers, key=lambda path: path.as_posix()))


def _discover_docs_agent_carriers(
    root: Path,
    carriers: set[Path],
    diagnostics: list[SourceInventoryDiagnostic],
) -> None:
    docs = root / "docs"
    agents = docs / "agents"
    if not docs.exists() and not docs.is_symlink():
        return
    for path, relative in ((docs, Path("docs")), (agents, Path("docs/agents"))):
        try:
            path_stat = path.lstat()
        except FileNotFoundError:
            return
        except OSError as error:
            diagnostics.append(
                SourceInventoryDiagnostic(
                    "source-inventory.directory.read", relative, str(error)
                )
            )
            return
        if stat.S_ISLNK(path_stat.st_mode):
            diagnostics.append(
                SourceInventoryDiagnostic(
                    "source-inventory.carrier.symlink",
                    relative,
                    "bounded carrier directory must not be a symlink",
                )
            )
            return
        if not stat.S_ISDIR(path_stat.st_mode):
            diagnostics.append(
                SourceInventoryDiagnostic(
                    "source-inventory.carrier.special",
                    relative,
                    "bounded carrier directory must be a real directory",
                )
            )
            return

    def walk(directory: Path, relative_directory: Path) -> None:
        try:
            children = sorted(os.scandir(directory), key=lambda child: child.name)
        except OSError as error:
            diagnostics.append(
                SourceInventoryDiagnostic(
                    "source-inventory.directory.read",
                    relative_directory,
                    str(error),
                )
            )
            return
        for child in children:
            relative = relative_directory / child.name
            try:
                child_stat = child.stat(follow_symlinks=False)
            except OSError as error:
                diagnostics.append(
                    SourceInventoryDiagnostic(
                        "source-inventory.carrier.read", relative, str(error)
                    )
                )
                continue
            if stat.S_ISLNK(child_stat.st_mode):
                diagnostics.append(
                    SourceInventoryDiagnostic(
                        "source-inventory.carrier.symlink",
                        relative,
                        "bounded carrier must not be a symlink",
                    )
                )
            elif stat.S_ISDIR(child_stat.st_mode):
                walk(Path(child.path), relative)
            elif stat.S_ISREG(child_stat.st_mode):
                carriers.add(relative)
            else:
                diagnostics.append(
                    SourceInventoryDiagnostic(
                        "source-inventory.carrier.special",
                        relative,
                        "bounded carrier must be a regular file",
                    )
                )

    walk(agents, Path("docs/agents"))


def _inventory_path_is_ignored(relative: Path) -> bool:
    if any(part in _INVENTORY_IGNORED_DIRECTORIES for part in relative.parts):
        return True
    return any(
        relative.parts[: len(prefix)] == prefix
        for prefix in _INVENTORY_IGNORED_PREFIXES
    )


def _read_inventory_carrier(root: Path, relative: Path, max_bytes: int) -> bytes:
    if (
        relative.is_absolute()
        or not relative.parts
        or any(part in {"", ".", ".."} for part in relative.parts)
    ):
        raise SourceInventoryError(
            [
                SourceInventoryDiagnostic(
                    "source-inventory.path.unsafe",
                    relative,
                    "carrier path must stay inside the repository",
                )
            ]
        )
    target = root / relative
    flags = os.O_RDONLY
    if hasattr(os, "O_NOFOLLOW"):
        flags |= os.O_NOFOLLOW
    try:
        descriptor = os.open(target, flags)
    except OSError as error:
        raise SourceInventoryError(
            [
                SourceInventoryDiagnostic(
                    "source-inventory.carrier.read", relative, str(error)
                )
            ]
        ) from error
    try:
        target_stat = os.fstat(descriptor)
        if not stat.S_ISREG(target_stat.st_mode):
            raise SourceInventoryError(
                [
                    SourceInventoryDiagnostic(
                        "source-inventory.carrier.special",
                        relative,
                        "bounded carrier must be a regular file",
                    )
                ]
            )
        if target_stat.st_size > max_bytes:
            raise SourceInventoryError(
                [
                    SourceInventoryDiagnostic(
                        "source-inventory.limit.file-bytes",
                        relative,
                        f"carrier size {target_stat.st_size} exceeds {max_bytes}",
                    )
                ]
            )
        content = bytearray()
        while True:
            chunk = os.read(descriptor, min(64 * 1024, max_bytes + 1 - len(content)))
            if not chunk:
                break
            content.extend(chunk)
            if len(content) > max_bytes:
                raise SourceInventoryError(
                    [
                        SourceInventoryDiagnostic(
                            "source-inventory.limit.file-bytes",
                            relative,
                            f"carrier bytes exceed {max_bytes}",
                        )
                    ]
                )
        return bytes(content)
    finally:
        os.close(descriptor)


def _partition_inventory_carrier(
    relative: Path, content: bytes
) -> list[ReadoptionSourceEntry]:
    carrier_digest = hashlib.sha256(content).hexdigest()
    if relative == _MANIFEST_CARRIER:
        spans = _manifest_record_spans(content)
        if spans:
            return _entries_from_structural_spans(
                relative, content, carrier_digest, spans
            )
    if relative.suffix.casefold() == ".md":
        spans = _managed_block_spans(content)
        if spans:
            return _entries_from_structural_spans(
                relative, content, carrier_digest, spans
            )
        return [
            _inventory_entry(
                relative,
                "unmarked-span",
                0,
                len(content),
                content,
                carrier_digest,
                (("markerState", "unmarked"),),
            )
        ]
    return [
        _inventory_entry(
            relative,
            "file",
            0,
            len(content),
            content,
            carrier_digest,
            (("boundary", "whole-file"),),
        )
    ]


def _managed_block_spans(
    content: bytes,
) -> list[tuple[int, int, str, tuple[tuple[str, object], ...]]]:
    spans: list[tuple[int, int, str, tuple[tuple[str, object], ...]]] = []
    cursor = 0
    while True:
        opening = _MANAGED_BEGIN_MARKER.search(content, cursor)
        if opening is None:
            break
        managed_id = opening.group(1)
        closing_pattern = re.compile(
            rb"<!--\s*setup-context-driven:end\s+id="
            + re.escape(managed_id)
            + rb"\s*-->"
        )
        closing = closing_pattern.search(content, opening.end())
        if closing is None:
            cursor = opening.end()
            continue
        end = closing.end()
        if content.startswith(b"\r\n", end):
            end += 2
        elif content.startswith(b"\n", end):
            end += 1
        spans.append(
            (
                opening.start(),
                end,
                "managed-block",
                (
                    ("managedId", managed_id.decode("ascii")),
                    ("markerVersion", int(opening.group(2))),
                ),
            )
        )
        cursor = end
    return spans


def _entries_from_structural_spans(
    relative: Path,
    content: bytes,
    carrier_digest: str,
    spans: list[tuple[int, int, str, tuple[tuple[str, object], ...]]],
) -> list[ReadoptionSourceEntry]:
    entries: list[ReadoptionSourceEntry] = []
    cursor = 0
    for start, end, kind, provenance in sorted(spans):
        if start < cursor or end < start or end > len(content):
            continue
        if start > cursor:
            entries.append(
                _inventory_entry(
                    relative,
                    "unmarked-span",
                    cursor,
                    start,
                    content[cursor:start],
                    carrier_digest,
                    (("markerState", "unmarked"),),
                )
            )
        entries.append(
            _inventory_entry(
                relative,
                kind,
                start,
                end,
                content[start:end],
                carrier_digest,
                provenance,
            )
        )
        cursor = end
    if cursor < len(content):
        entries.append(
            _inventory_entry(
                relative,
                "unmarked-span",
                cursor,
                len(content),
                content[cursor:],
                carrier_digest,
                (("markerState", "unmarked"),),
            )
        )
    return entries


def _inventory_entry(
    relative: Path,
    kind: str,
    start: int,
    end: int,
    source_bytes: bytes,
    carrier_digest: str,
    provenance: tuple[tuple[str, object], ...],
) -> ReadoptionSourceEntry:
    digest = hashlib.sha256(source_bytes).hexdigest()
    identity = hashlib.sha256()
    for value in (
        relative.as_posix(),
        kind,
        str(start),
        str(end),
        digest,
        json.dumps(dict(provenance), sort_keys=True, separators=(",", ":")),
    ):
        encoded = value.encode("utf-8")
        identity.update(len(encoded).to_bytes(8, "big"))
        identity.update(encoded)
    return ReadoptionSourceEntry(
        entry_id=f"source-entry.{identity.hexdigest()}",
        path=relative,
        kind=kind,
        start=start,
        end=end,
        digest=digest,
        carrier_digest=carrier_digest,
        source_bytes=source_bytes,
        provenance=provenance,
    )


def _validate_inventory_manifest_paths(relative: Path, content: bytes) -> None:
    try:
        document = json.loads(content)
    except (UnicodeDecodeError, json.JSONDecodeError):
        return
    if not isinstance(document, dict):
        return
    diagnostics: list[SourceInventoryDiagnostic] = []
    for collection_name in ("managedArtifacts", "repositoryExtensions"):
        records = document.get(collection_name, [])
        if not isinstance(records, list):
            continue
        for index, record in enumerate(records):
            if not isinstance(record, dict) or "path" not in record:
                continue
            value = record["path"]
            if not _inventory_manifest_path_is_safe(value):
                diagnostics.append(
                    SourceInventoryDiagnostic(
                        "source-inventory.manifest.path.unsafe",
                        relative,
                        f"{collection_name}[{index}].path is unsafe: {value!r}",
                    )
                )
    if diagnostics:
        raise SourceInventoryError(diagnostics)


def _inventory_manifest_path_is_safe(value: object) -> bool:
    if not isinstance(value, str) or not value or "\\" in value:
        return False
    path = PurePosixPath(value)
    return not path.is_absolute() and all(
        part not in {"", ".", ".."} for part in path.parts
    )


def _manifest_record_spans(
    content: bytes,
) -> list[tuple[int, int, str, tuple[tuple[str, object], ...]]]:
    try:
        members = _json_object_members(content, 0)
    except ValueError:
        return []
    spans: list[tuple[int, int, str, tuple[tuple[str, object], ...]]] = []
    for key, key_start, value_start, value_end in members:
        if key == "managedArtifacts":
            try:
                elements = _json_array_elements(content, value_start)
            except ValueError:
                elements = []
            if elements:
                for index, (start, end) in enumerate(elements):
                    provenance: list[tuple[str, object]] = [
                        ("recordKey", key),
                        ("recordIndex", index),
                    ]
                    try:
                        record = json.loads(content[start:end])
                    except (UnicodeDecodeError, json.JSONDecodeError):
                        record = None
                    if isinstance(record, dict) and isinstance(record.get("id"), str):
                        provenance.append(("managedId", record["id"]))
                    spans.append(
                        (start, end, "manifest-record", tuple(provenance))
                    )
                continue
        if key == "decisions":
            try:
                decisions = _json_object_members(content, value_start)
            except ValueError:
                decisions = []
            if decisions:
                for decision_id, start, _, end in decisions:
                    spans.append(
                        (
                            start,
                            end,
                            "manifest-record",
                            (("recordKey", key), ("decisionId", decision_id)),
                        )
                    )
                continue
        spans.append(
            (
                key_start,
                value_end,
                "manifest-record",
                (("recordKey", key),),
            )
        )
    return spans


def _json_object_members(
    content: bytes, start: int
) -> list[tuple[str, int, int, int]]:
    cursor = _json_skip_whitespace(content, start)
    if cursor >= len(content) or content[cursor] != ord("{"):
        raise ValueError("expected JSON object")
    cursor += 1
    members: list[tuple[str, int, int, int]] = []
    while True:
        cursor = _json_skip_whitespace(content, cursor)
        if cursor >= len(content):
            raise ValueError("unterminated JSON object")
        if content[cursor] == ord("}"):
            return members
        key_start = cursor
        key_end = _json_string_end(content, cursor)
        try:
            key = json.loads(content[key_start:key_end])
        except (UnicodeDecodeError, json.JSONDecodeError) as error:
            raise ValueError("invalid JSON object key") from error
        if not isinstance(key, str):
            raise ValueError("JSON object key must be a string")
        cursor = _json_skip_whitespace(content, key_end)
        if cursor >= len(content) or content[cursor] != ord(":"):
            raise ValueError("missing JSON member colon")
        value_start = _json_skip_whitespace(content, cursor + 1)
        value_end = _json_value_end(content, value_start)
        members.append((key, key_start, value_start, value_end))
        cursor = _json_skip_whitespace(content, value_end)
        if cursor >= len(content):
            raise ValueError("unterminated JSON object")
        if content[cursor] == ord(","):
            cursor += 1
            continue
        if content[cursor] == ord("}"):
            return members
        raise ValueError("invalid JSON object separator")


def _json_array_elements(content: bytes, start: int) -> list[tuple[int, int]]:
    cursor = _json_skip_whitespace(content, start)
    if cursor >= len(content) or content[cursor] != ord("["):
        raise ValueError("expected JSON array")
    cursor += 1
    elements: list[tuple[int, int]] = []
    while True:
        cursor = _json_skip_whitespace(content, cursor)
        if cursor >= len(content):
            raise ValueError("unterminated JSON array")
        if content[cursor] == ord("]"):
            return elements
        end = _json_value_end(content, cursor)
        elements.append((cursor, end))
        cursor = _json_skip_whitespace(content, end)
        if cursor >= len(content):
            raise ValueError("unterminated JSON array")
        if content[cursor] == ord(","):
            cursor += 1
            continue
        if content[cursor] == ord("]"):
            return elements
        raise ValueError("invalid JSON array separator")


def _json_value_end(content: bytes, start: int) -> int:
    cursor = _json_skip_whitespace(content, start)
    if cursor >= len(content):
        raise ValueError("missing JSON value")
    token = content[cursor]
    if token == ord('"'):
        return _json_string_end(content, cursor)
    if token in {ord("{"), ord("[")}:
        opening = token
        closing = ord("}") if token == ord("{") else ord("]")
        depth = 0
        while cursor < len(content):
            token = content[cursor]
            if token == ord('"'):
                cursor = _json_string_end(content, cursor)
                continue
            if token == opening:
                depth += 1
            elif token == closing:
                depth -= 1
                if depth == 0:
                    return cursor + 1
            elif token in {ord("{"), ord("[")} and token != opening:
                cursor = _json_value_end(content, cursor)
                continue
            cursor += 1
        raise ValueError("unterminated JSON container")
    end = cursor
    while end < len(content) and content[end] not in b",]} \t\r\n":
        end += 1
    if end == cursor:
        raise ValueError("empty JSON primitive")
    return end


def _json_string_end(content: bytes, start: int) -> int:
    if start >= len(content) or content[start] != ord('"'):
        raise ValueError("expected JSON string")
    cursor = start + 1
    while cursor < len(content):
        if content[cursor] == ord("\\"):
            cursor += 2
            continue
        if content[cursor] == ord('"'):
            return cursor + 1
        cursor += 1
    raise ValueError("unterminated JSON string")


def _json_skip_whitespace(content: bytes, cursor: int) -> int:
    while cursor < len(content) and content[cursor] in b" \t\r\n":
        cursor += 1
    return cursor


@dataclass(frozen=True)
class SourceBaselineEntry:
    entry_id: str
    path: Path
    kind: str
    enforcement: str
    carrier: Path
    structure: str | None
    start: int
    end: int
    digest: str


@dataclass(frozen=True)
class SourceBaselineManifest:
    baseline_id: str
    profile: str
    version: str
    entries: tuple[SourceBaselineEntry, ...]


@dataclass(frozen=True)
class SourceBaselineIdentity:
    baseline_id: str
    profile: str
    version: str
    corpus_path: Path
    manifest_path: Path
    entry_count: int
    corpus_digest: str
    manifest_digest: str
    denied_project_tokens: tuple[str, ...]


@dataclass(frozen=True)
class SourceBaselineIndexEntry:
    baseline_id: str
    path: Path
    profile: str
    entry_ids: tuple[str, ...]
    entry_count: int
    corpus_digest: str
    manifest_digest: str


@dataclass(frozen=True)
class SourceBaselineIndex:
    version: str
    entries: tuple[SourceBaselineIndexEntry, ...]


@dataclass(frozen=True)
class SourceBaseline:
    baseline_id: str
    profile: str
    version: str
    identity: SourceBaselineIdentity
    manifest: SourceBaselineManifest
    index_entry: SourceBaselineIndexEntry
    entries: tuple[SourceBaselineEntry, ...]
    corpus_digest: str
    manifest_digest: str


@dataclass(frozen=True)
class _CorpusEntry:
    entry_id: str
    path: Path
    start: int
    end: int
    digest: str


def load_source_baselines(skill_root: str | Path) -> tuple[SourceBaseline, ...]:
    """Load every Source Baseline declared by the independent baseline index."""

    root = _source_baselines_root(Path(skill_root))
    diagnostics: list[str] = []
    index_doc, _ = _read_json(root / "index.json", "index", diagnostics)
    index = _parse_index(index_doc, diagnostics)

    declared_paths = {entry.path.as_posix() for entry in index.entries}
    if root.is_dir():
        for child in sorted(root.iterdir(), key=lambda path: path.name):
            if child.name == "index.json":
                continue
            if child.is_symlink() or not child.is_dir():
                diagnostics.append(
                    f"source-baseline.index.path.unknown: unexpected entry {child.name!r}"
                )
            elif child.name not in declared_paths:
                diagnostics.append(
                    "source-baseline.index.entry.missing: "
                    f"directory {child.name!r} has no index record"
                )

    baselines: list[SourceBaseline] = []
    for index_entry in index.entries:
        baseline_root = root / index_entry.path
        if not baseline_root.is_dir() or baseline_root.is_symlink():
            diagnostics.append(
                "source-baseline.corpus.missing: "
                f"index entry {index_entry.baseline_id!r} has no safe directory"
            )
            continue
        baseline = _load_indexed_baseline(baseline_root, index_entry, diagnostics)
        if baseline is not None:
            baselines.append(baseline)

    if diagnostics:
        raise SourceBaselineValidationError(diagnostics)
    return tuple(baselines)


def load_source_baseline(
    skill_root: str | Path, baseline_id: str
) -> SourceBaseline:
    """Load one indexed Source Baseline by its stable identifier."""

    baselines = load_source_baselines(skill_root)
    for baseline in baselines:
        if baseline.baseline_id == baseline_id:
            return baseline
    raise SourceBaselineValidationError(
        [f"source-baseline.id.unknown: no indexed baseline {baseline_id!r}"]
    )


def render_source_baseline_carriers(
    skill_root: str | Path, baseline_id: str
) -> dict[Path, bytes]:
    """Render complete carrier bytes from one validated Source Baseline."""

    root = _source_baselines_root(Path(skill_root))
    baseline = load_source_baseline(root, baseline_id)
    baseline_root = root / baseline.index_entry.path
    fragments: dict[Path, list[bytes]] = {}
    for entry in baseline.entries:
        content = (baseline_root / entry.path).read_bytes()[entry.start : entry.end]
        fragments.setdefault(entry.carrier, []).append(content.rstrip())
    return {
        carrier: b"\n\n".join(parts) + b"\n"
        for carrier, parts in fragments.items()
    }


def compute_corpus_digest(corpus_root: str | Path) -> str:
    """Return the deterministic path-and-byte digest used by Source Baselines."""

    root = Path(corpus_root)
    digest = hashlib.sha256()
    for path in _corpus_files(root):
        relative = path.relative_to(root).as_posix().encode("utf-8")
        content = path.read_bytes()
        digest.update(len(relative).to_bytes(8, "big"))
        digest.update(relative)
        digest.update(len(content).to_bytes(8, "big"))
        digest.update(content)
    return digest.hexdigest()


def _source_baselines_root(root: Path) -> Path:
    if (root / "index.json").is_file():
        return root
    return root / "assets" / "source-baselines"


def _load_indexed_baseline(
    baseline_root: Path,
    index_entry: SourceBaselineIndexEntry,
    diagnostics: list[str],
) -> SourceBaseline | None:
    identity_doc, _ = _read_json(baseline_root / "baseline.json", "identity", diagnostics)
    manifest_doc, manifest_bytes = _read_json(
        baseline_root / "manifest.json", "manifest", diagnostics
    )
    identity = _parse_identity(identity_doc, diagnostics)
    manifest = _parse_manifest(manifest_doc, diagnostics)
    if identity is None or manifest is None or manifest_bytes is None:
        return None

    _compare_baseline_headers(baseline_root, index_entry, identity, manifest, diagnostics)

    corpus_root = baseline_root / identity.corpus_path
    corpus_entries = _load_corpus_entries(
        baseline_root, corpus_root, diagnostics
    )
    try:
        corpus_digest = compute_corpus_digest(corpus_root)
    except (OSError, ValueError) as error:
        diagnostics.append(
            f"source-baseline.corpus.invalid: cannot digest {corpus_root}: {error}"
        )
        corpus_digest = ""
    manifest_digest = hashlib.sha256(manifest_bytes).hexdigest()

    _validate_membership(manifest, index_entry, corpus_entries, diagnostics)
    _validate_counts(identity, manifest, index_entry, corpus_entries, diagnostics)
    _validate_digests(
        identity,
        index_entry,
        corpus_digest,
        manifest_digest,
        diagnostics,
    )
    _validate_corpus_policy(
        corpus_root,
        identity.denied_project_tokens,
        diagnostics,
    )

    corpus_by_id = {entry.entry_id: entry for entry in corpus_entries}
    for entry in manifest.entries:
        corpus_entry = corpus_by_id.get(entry.entry_id)
        if corpus_entry is None:
            continue
        if entry.path != corpus_entry.path:
            diagnostics.append(
                "source-baseline.entry.path.mismatch: "
                f"{entry.entry_id!r} manifest={entry.path} corpus={corpus_entry.path}"
            )
        if (entry.start, entry.end) != (corpus_entry.start, corpus_entry.end):
            diagnostics.append(
                "source-baseline.entry.range.mismatch: "
                f"{entry.entry_id!r} manifest={entry.start}:{entry.end} "
                f"corpus={corpus_entry.start}:{corpus_entry.end}"
            )
        if entry.digest != corpus_entry.digest:
            diagnostics.append(
                "source-baseline.entry.digest.mismatch: "
                f"{entry.entry_id!r} does not match its corpus bytes"
            )

    return SourceBaseline(
        baseline_id=identity.baseline_id,
        profile=identity.profile,
        version=identity.version,
        identity=identity,
        manifest=manifest,
        index_entry=index_entry,
        entries=manifest.entries,
        corpus_digest=corpus_digest,
        manifest_digest=manifest_digest,
    )


def _parse_index(document: object, diagnostics: list[str]) -> SourceBaselineIndex:
    if not _document_fields(
        document, {"schemaVersion", "version", "baselines"}, "index", diagnostics
    ):
        return SourceBaselineIndex(SOURCE_BASELINE_VERSION, ())
    assert isinstance(document, dict)
    _validate_generation(document, SOURCE_BASELINE_INDEX_SCHEMA, "index", diagnostics)
    records = document["baselines"]
    if not isinstance(records, list):
        diagnostics.append("source-baseline.index.document.invalid: baselines must be a list")
        return SourceBaselineIndex(SOURCE_BASELINE_VERSION, ())

    entries: list[SourceBaselineIndexEntry] = []
    seen_ids: set[str] = set()
    seen_paths: set[str] = set()
    for position, record in enumerate(records):
        label = f"index.baselines[{position}]"
        fields = {
            "id",
            "path",
            "profile",
            "entryIds",
            "entryCount",
            "corpusDigest",
            "manifestDigest",
        }
        if not _document_fields(record, fields, label, diagnostics):
            continue
        assert isinstance(record, dict)
        baseline_id = _identifier(record["id"], f"{label}.id", diagnostics)
        profile = _identifier(record["profile"], f"{label}.profile", diagnostics)
        path = _safe_path(record["path"], f"{label}.path", diagnostics)
        entry_ids = _identifier_list(record["entryIds"], f"{label}.entryIds", diagnostics)
        entry_count = _count(record["entryCount"], f"{label}.entryCount", diagnostics)
        corpus_digest = _digest(record["corpusDigest"], f"{label}.corpusDigest", diagnostics)
        manifest_digest = _digest(
            record["manifestDigest"], f"{label}.manifestDigest", diagnostics
        )
        if baseline_id in seen_ids:
            diagnostics.append(
                f"source-baseline.id.duplicate: index repeats {baseline_id!r}"
            )
        seen_ids.add(baseline_id)
        if path.as_posix() in seen_paths:
            diagnostics.append(
                f"source-baseline.path.duplicate: index repeats {path.as_posix()!r}"
            )
        seen_paths.add(path.as_posix())
        entries.append(
            SourceBaselineIndexEntry(
                baseline_id,
                path,
                profile,
                entry_ids,
                entry_count,
                corpus_digest,
                manifest_digest,
            )
        )
    return SourceBaselineIndex(SOURCE_BASELINE_VERSION, tuple(entries))


def _parse_identity(
    document: object, diagnostics: list[str]
) -> SourceBaselineIdentity | None:
    fields = {
        "schemaVersion",
        "version",
        "id",
        "profile",
        "corpus",
        "manifest",
        "entryCount",
        "corpusDigest",
        "manifestDigest",
        "deniedProjectTokens",
    }
    if not _document_fields(document, fields, "identity", diagnostics):
        return None
    assert isinstance(document, dict)
    _validate_generation(document, SOURCE_BASELINE_SCHEMA, "identity", diagnostics)
    return SourceBaselineIdentity(
        _identifier(document["id"], "identity.id", diagnostics),
        _identifier(document["profile"], "identity.profile", diagnostics),
        SOURCE_BASELINE_VERSION,
        _safe_path(document["corpus"], "identity.corpus", diagnostics),
        _safe_path(document["manifest"], "identity.manifest", diagnostics),
        _count(document["entryCount"], "identity.entryCount", diagnostics),
        _digest(document["corpusDigest"], "identity.corpusDigest", diagnostics),
        _digest(document["manifestDigest"], "identity.manifestDigest", diagnostics),
        _project_tokens(
            document["deniedProjectTokens"],
            "identity.deniedProjectTokens",
            diagnostics,
        ),
    )


def _parse_manifest(
    document: object, diagnostics: list[str]
) -> SourceBaselineManifest | None:
    fields = {"schemaVersion", "version", "id", "profile", "entries"}
    if not _document_fields(document, fields, "manifest", diagnostics):
        return None
    assert isinstance(document, dict)
    _validate_generation(document, SOURCE_BASELINE_MANIFEST_SCHEMA, "manifest", diagnostics)
    raw_entries = document["entries"]
    if not isinstance(raw_entries, list):
        diagnostics.append("source-baseline.manifest.document.invalid: entries must be a list")
        return None

    entries: list[SourceBaselineEntry] = []
    seen: set[str] = set()
    for position, raw_entry in enumerate(raw_entries):
        label = f"manifest.entries[{position}]"
        fields = {
            "id",
            "path",
            "kind",
            "enforcement",
            "carrier",
            "structure",
            "start",
            "end",
            "digest",
        }
        if not _document_fields(raw_entry, fields, label, diagnostics):
            continue
        assert isinstance(raw_entry, dict)
        entry_id = _identifier(raw_entry["id"], f"{label}.id", diagnostics)
        if entry_id in seen:
            diagnostics.append(
                f"source-baseline.entry.id.duplicate: manifest repeats {entry_id!r}"
            )
        seen.add(entry_id)
        path = _safe_path(raw_entry["path"], f"{label}.path", diagnostics)
        if not path.parts or path.parts[0] != "corpus":
            diagnostics.append(
                f"source-baseline.entry.path.invalid: {entry_id!r} must be under corpus/"
            )
        kind = raw_entry["kind"]
        if not isinstance(kind, str) or kind not in _ENTRY_KINDS:
            diagnostics.append(
                f"source-baseline.entry.kind.invalid: {entry_id!r} has kind {kind!r}"
            )
            kind = ""
        enforcement = _entry_enforcement(
            raw_entry["enforcement"], entry_id, kind, diagnostics
        )
        carrier = _safe_path(
            raw_entry["carrier"], f"{label}.carrier", diagnostics
        )
        structure = _entry_structure(
            raw_entry["structure"], entry_id, kind, diagnostics
        )
        start = _count(raw_entry["start"], f"{label}.start", diagnostics)
        end = _count(raw_entry["end"], f"{label}.end", diagnostics)
        if end <= start:
            diagnostics.append(
                f"source-baseline.entry.range.invalid: {entry_id!r} has {start}:{end}"
            )
        entries.append(
            SourceBaselineEntry(
                entry_id,
                path,
                kind,
                enforcement,
                carrier,
                structure,
                start,
                end,
                _digest(raw_entry["digest"], f"{label}.digest", diagnostics),
            )
        )
    return SourceBaselineManifest(
        _identifier(document["id"], "manifest.id", diagnostics),
        _identifier(document["profile"], "manifest.profile", diagnostics),
        SOURCE_BASELINE_VERSION,
        tuple(entries),
    )


def _load_corpus_entries(
    baseline_root: Path, corpus_root: Path, diagnostics: list[str]
) -> tuple[_CorpusEntry, ...]:
    if not corpus_root.is_dir() or corpus_root.is_symlink():
        diagnostics.append(
            f"source-baseline.corpus.missing: no safe corpus directory at {corpus_root}"
        )
        return ()
    try:
        paths = _corpus_files(corpus_root)
    except ValueError as error:
        diagnostics.append(f"source-baseline.entry.path.invalid: {error}")
        return ()
    if not paths:
        diagnostics.append("source-baseline.corpus.empty: corpus contains no files")
        return ()

    entries: list[_CorpusEntry] = []
    seen: set[str] = set()
    for path in paths:
        relative = path.relative_to(baseline_root)
        data = path.read_bytes()
        cursor = 0
        while True:
            opening = _OPEN_MARKER.search(data, cursor)
            if opening is None:
                if data[cursor:].strip():
                    diagnostics.append(
                        "source-baseline.corpus.structure.invalid: "
                        f"{relative} has bytes outside an entry"
                    )
                break
            if data[cursor : opening.start()].strip():
                diagnostics.append(
                    "source-baseline.corpus.structure.invalid: "
                    f"{relative} has bytes before an entry marker"
                )
            entry_id = opening.group(1).decode("ascii")
            closing = re.compile(
                rb"\r?\n<!-- /source-baseline-entry: "
                + re.escape(opening.group(1))
                + rb" -->"
            ).search(data, opening.end())
            if closing is None:
                diagnostics.append(
                    "source-baseline.corpus.structure.invalid: "
                    f"{relative} entry {entry_id!r} has no matching close marker"
                )
                break
            nested = _OPEN_MARKER.search(data, opening.end(), closing.start())
            if nested is not None:
                diagnostics.append(
                    "source-baseline.corpus.structure.invalid: "
                    f"{relative} entry {entry_id!r} contains a nested marker"
                )
            content = data[opening.end() : closing.start()]
            if not content.strip():
                diagnostics.append(
                    f"source-baseline.entry.content.invalid: {entry_id!r} is empty"
                )
            if entry_id in seen:
                diagnostics.append(
                    f"source-baseline.entry.id.duplicate: corpus repeats {entry_id!r}"
                )
            seen.add(entry_id)
            entries.append(
                _CorpusEntry(
                    entry_id,
                    relative,
                    opening.end(),
                    closing.start(),
                    hashlib.sha256(content).hexdigest(),
                )
            )
            cursor = closing.end()
    return tuple(entries)


def _corpus_files(root: Path) -> tuple[Path, ...]:
    if not root.is_dir() or root.is_symlink():
        raise ValueError(f"unsafe or missing corpus directory {root}")
    paths: list[Path] = []
    for path in root.rglob("*"):
        if path.is_symlink():
            raise ValueError(f"corpus path {path} is a symlink")
        if path.is_file():
            paths.append(path)
    return tuple(sorted(paths, key=lambda path: path.relative_to(root).as_posix()))


def _compare_baseline_headers(
    baseline_root: Path,
    index_entry: SourceBaselineIndexEntry,
    identity: SourceBaselineIdentity,
    manifest: SourceBaselineManifest,
    diagnostics: list[str],
) -> None:
    expected = index_entry.baseline_id
    if baseline_root.name != index_entry.path.as_posix():
        diagnostics.append(
            "source-baseline.path.mismatch: "
            f"index path {index_entry.path} resolves to {baseline_root.name!r}"
        )
    for owner, value in (("identity", identity.baseline_id), ("manifest", manifest.baseline_id)):
        if value != expected:
            diagnostics.append(
                f"source-baseline.id.mismatch: {owner} has {value!r}, expected {expected!r}"
            )
    for owner, value in (("identity", identity.profile), ("manifest", manifest.profile)):
        if value != index_entry.profile:
            diagnostics.append(
                "source-baseline.profile.mismatch: "
                f"{owner} has {value!r}, expected {index_entry.profile!r}"
            )
    if identity.corpus_path != Path("corpus"):
        diagnostics.append(
            "source-baseline.entry.path.invalid: identity corpus must be 'corpus'"
        )
    if identity.manifest_path != Path("manifest.json"):
        diagnostics.append(
            "source-baseline.entry.path.invalid: identity manifest must be 'manifest.json'"
        )


def _validate_membership(
    manifest: SourceBaselineManifest,
    index_entry: SourceBaselineIndexEntry,
    corpus_entries: tuple[_CorpusEntry, ...],
    diagnostics: list[str],
) -> None:
    corpus_ids = tuple(entry.entry_id for entry in corpus_entries)
    manifest_ids = tuple(entry.entry_id for entry in manifest.entries)
    index_ids = index_entry.entry_ids
    for entry_id in sorted(set(corpus_ids) - set(manifest_ids)):
        diagnostics.append(
            f"source-baseline.manifest.entry.missing: corpus entry {entry_id!r} is unlisted"
        )
    for entry_id in sorted(set(manifest_ids) - set(corpus_ids)):
        diagnostics.append(
            f"source-baseline.corpus.entry.missing: manifest entry {entry_id!r} is absent"
        )
    for entry_id in sorted(set(manifest_ids) - set(index_ids)):
        diagnostics.append(
            f"source-baseline.index.entry.missing: manifest entry {entry_id!r} is unpinned"
        )
    for entry_id in sorted(set(index_ids) - set(manifest_ids)):
        diagnostics.append(
            f"source-baseline.manifest.entry.missing: index entry {entry_id!r} is absent"
        )
    if set(corpus_ids) == set(manifest_ids) and corpus_ids != manifest_ids:
        diagnostics.append(
            "source-baseline.entry.order.mismatch: corpus and manifest order differ"
        )
    if set(manifest_ids) == set(index_ids) and manifest_ids != index_ids:
        diagnostics.append(
            "source-baseline.entry.order.mismatch: manifest and index order differ"
        )


def _validate_counts(
    identity: SourceBaselineIdentity,
    manifest: SourceBaselineManifest,
    index_entry: SourceBaselineIndexEntry,
    corpus_entries: tuple[_CorpusEntry, ...],
    diagnostics: list[str],
) -> None:
    actual = len(manifest.entries)
    counts = {
        "identity": identity.entry_count,
        "index": index_entry.entry_count,
        "index membership": len(index_entry.entry_ids),
        "corpus": len(corpus_entries),
    }
    for owner, count in counts.items():
        if count != actual:
            diagnostics.append(
                "source-baseline.entry.count.mismatch: "
                f"{owner} records {count}, manifest contains {actual}"
            )


def _validate_digests(
    identity: SourceBaselineIdentity,
    index_entry: SourceBaselineIndexEntry,
    corpus_digest: str,
    manifest_digest: str,
    diagnostics: list[str],
) -> None:
    for owner, pinned in (
        ("identity", identity.corpus_digest),
        ("index", index_entry.corpus_digest),
    ):
        if pinned != corpus_digest:
            diagnostics.append(
                "source-baseline.corpus.digest.mismatch: "
                f"{owner} pin does not match corpus bytes"
            )
    for owner, pinned in (
        ("identity", identity.manifest_digest),
        ("index", index_entry.manifest_digest),
    ):
        if pinned != manifest_digest:
            diagnostics.append(
                "source-baseline.manifest.digest.mismatch: "
                f"{owner} pin does not match manifest bytes"
            )


def _read_json(
    path: Path, label: str, diagnostics: list[str]
) -> tuple[object, bytes | None]:
    try:
        content = path.read_bytes()
    except OSError as error:
        diagnostics.append(
            f"source-baseline.{label}.document.invalid: cannot read {path}: {error}"
        )
        return {}, None
    try:
        return json.loads(content), content
    except (UnicodeDecodeError, json.JSONDecodeError) as error:
        diagnostics.append(
            f"source-baseline.{label}.document.invalid: cannot parse {path}: {error}"
        )
        return {}, content


def _document_fields(
    value: object,
    expected: set[str],
    label: str,
    diagnostics: list[str],
) -> bool:
    if not isinstance(value, dict):
        diagnostics.append(
            f"source-baseline.{label}.document.invalid: expected an object"
        )
        return False
    actual = set(value)
    if actual != expected:
        missing = sorted(expected - actual)
        unknown = sorted(actual - expected)
        diagnostics.append(
            f"source-baseline.{label}.fields.invalid: missing={missing} unknown={unknown}"
        )
        return False
    return True


def _validate_generation(
    document: dict, expected_schema: str, label: str, diagnostics: list[str]
) -> None:
    schema = document["schemaVersion"]
    if not isinstance(schema, str) or schema != expected_schema:
        diagnostics.append(
            f"source-baseline.schema.invalid: {label} expected {expected_schema!r}, got {schema!r}"
        )
    version = document["version"]
    if not isinstance(version, str):
        diagnostics.append(
            f"source-baseline.version.invalid: {label} version must be string '0.0.1'"
        )
    elif version != SOURCE_BASELINE_VERSION:
        diagnostics.append(
            "source-baseline.version.mismatch: "
            f"{label} expected '0.0.1', got {version!r}"
        )


def _identifier(value: object, label: str, diagnostics: list[str]) -> str:
    if not isinstance(value, str) or _IDENTIFIER.fullmatch(value) is None:
        diagnostics.append(
            f"source-baseline.identifier.invalid: {label} has invalid value {value!r}"
        )
        return ""
    return value


def _identifier_list(
    value: object, label: str, diagnostics: list[str]
) -> tuple[str, ...]:
    if not isinstance(value, list):
        diagnostics.append(f"source-baseline.entry.ids.invalid: {label} must be a list")
        return ()
    identifiers = tuple(
        _identifier(item, f"{label}[{position}]", diagnostics)
        for position, item in enumerate(value)
    )
    if len(set(identifiers)) != len(identifiers):
        diagnostics.append(f"source-baseline.entry.id.duplicate: {label} repeats an id")
    return identifiers


def _safe_path(value: object, label: str, diagnostics: list[str]) -> Path:
    if not isinstance(value, str) or not value or "\\" in value:
        diagnostics.append(
            f"source-baseline.entry.path.invalid: {label} has invalid path {value!r}"
        )
        return Path()
    path = PurePosixPath(value)
    if path.is_absolute() or any(part in {"", ".", ".."} for part in path.parts):
        diagnostics.append(
            f"source-baseline.entry.path.invalid: {label} has unsafe path {value!r}"
        )
        return Path()
    return Path(*path.parts)


def _count(value: object, label: str, diagnostics: list[str]) -> int:
    if type(value) is not int or value < 0:
        diagnostics.append(
            f"source-baseline.entry.count.invalid: {label} must be a non-negative integer"
        )
        return 0
    return value


def _digest(value: object, label: str, diagnostics: list[str]) -> str:
    if not isinstance(value, str) or _SHA256.fullmatch(value) is None:
        diagnostics.append(
            f"source-baseline.digest.invalid: {label} must be a lowercase SHA-256 digest"
        )
        return ""
    return value


def _entry_enforcement(
    value: object, entry_id: str, kind: str, diagnostics: list[str]
) -> str:
    allowed = {"recommended"} if kind == "recommendation" else _NORMATIVE_ENFORCEMENTS
    if not isinstance(value, str) or value not in allowed:
        diagnostics.append(
            "source-baseline.entry.enforcement.invalid: "
            f"{entry_id!r} kind {kind!r} cannot use {value!r}"
        )
        return ""
    return value


def _entry_structure(
    value: object, entry_id: str, kind: str, diagnostics: list[str]
) -> str | None:
    if kind == "operational-contract":
        if not isinstance(value, str) or value not in _OPERATIONAL_STRUCTURES:
            diagnostics.append(
                "source-baseline.entry.structure.invalid: "
                f"{entry_id!r} requires a supported Operational Contract structure"
            )
            return None
        return value
    if value is not None:
        diagnostics.append(
            "source-baseline.entry.structure.invalid: "
            f"{entry_id!r} is not an Operational Contract"
        )
    return None


def _project_tokens(
    value: object, label: str, diagnostics: list[str]
) -> tuple[str, ...]:
    if not isinstance(value, list) or not value:
        diagnostics.append(
            f"source-baseline.project-token.list.invalid: {label} must be a non-empty list"
        )
        return ()
    tokens: list[str] = []
    for position, item in enumerate(value):
        if not isinstance(item, str) or not item.strip() or item != item.casefold():
            diagnostics.append(
                "source-baseline.project-token.invalid: "
                f"{label}[{position}] must be a non-empty case-folded string"
            )
            continue
        tokens.append(item)
    if len(tokens) != len(set(tokens)):
        diagnostics.append(
            f"source-baseline.project-token.duplicate: {label} repeats a token"
        )
    return tuple(tokens)


def _validate_corpus_policy(
    corpus_root: Path,
    denied_project_tokens: tuple[str, ...],
    diagnostics: list[str],
) -> None:
    try:
        paths = _corpus_files(corpus_root)
    except ValueError:
        return
    for path in paths:
        try:
            text = path.read_text(encoding="utf-8")
        except (OSError, UnicodeDecodeError) as error:
            diagnostics.append(
                f"source-baseline.corpus.invalid: cannot inspect {path}: {error}"
            )
            continue
        folded = text.casefold()
        for token in denied_project_tokens:
            if token in folded:
                diagnostics.append(
                    "source-baseline.project-token.denied: "
                    f"{path.relative_to(corpus_root)} contains {token!r}"
                )
        if "<!-- setup-context-driven:" in folded:
            diagnostics.append(
                "source-baseline.generated-artifact.denied: "
                f"{path.relative_to(corpus_root)} contains a generated managed marker"
            )
        if re.search(r"(?i)(?:^|[\s`\"'])/(?:users|home)/", text):
            diagnostics.append(
                "source-baseline.path-token.denied: "
                f"{path.relative_to(corpus_root)} contains a machine-specific path"
            )
        if re.search(r"(?i)(?:^|[\s`\"'])[a-z]:\\users\\", text):
            diagnostics.append(
                "source-baseline.path-token.denied: "
                f"{path.relative_to(corpus_root)} contains a machine-specific path"
            )
