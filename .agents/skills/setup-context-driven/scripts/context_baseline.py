"""Load and validate immutable setup-context-driven Source Baselines."""

from __future__ import annotations

import hashlib
import json
import re
from dataclasses import dataclass
from pathlib import Path, PurePosixPath


SOURCE_BASELINE_VERSION = "0.0.1"
SOURCE_BASELINE_SCHEMA = "setup-context-driven/source-baseline/0.0.1"
SOURCE_BASELINE_MANIFEST_SCHEMA = (
    "setup-context-driven/source-baseline-manifest/0.0.1"
)
SOURCE_BASELINE_INDEX_SCHEMA = "setup-context-driven/source-baseline-index/0.0.1"

_ENTRY_KINDS = {"normative-clause", "recommendation", "operational-contract"}
_IDENTIFIER = re.compile(r"^[a-z0-9](?:[a-z0-9._-]*[a-z0-9])?$")
_SHA256 = re.compile(r"^[0-9a-f]{64}$")
_OPEN_MARKER = re.compile(
    rb"(?m)^<!-- source-baseline-entry: ([a-z0-9][a-z0-9._-]*) -->\r?\n"
)


class SourceBaselineValidationError(ValueError):
    """Raised when Source Baseline documents disagree or are invalid."""

    def __init__(self, diagnostics: list[str]):
        self.diagnostics = tuple(diagnostics)
        super().__init__("\n".join(diagnostics))


@dataclass(frozen=True)
class SourceBaselineEntry:
    entry_id: str
    path: Path
    kind: str
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
        fields = {"id", "path", "kind", "start", "end", "digest"}
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
