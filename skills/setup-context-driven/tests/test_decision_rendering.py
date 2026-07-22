import json
import sys
import tempfile
import unittest
from hashlib import sha256
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
SCRIPT = SKILL_ROOT / "scripts" / "context_setup.py"
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from context_setup import parse_managed_blocks  # noqa: E402
from test_audit import (  # noqa: E402
    install_profile_skills,
    run_context_setup as run_fixture_context_setup,
    snapshot_files,
)


class DecisionRenderingTests(unittest.TestCase):
    def test_domain_layout_selects_distinct_guidance_and_audits_clean(self):
        cases = [
            ("single-context", "single `CONTEXT.md`", "CONTEXT-MAP.md", False),
            ("multi-context", "CONTEXT-MAP.md", "single `CONTEXT.md`", False),
        ]

        for layout, expected_text, omitted_text, monorepo_expected in cases:
            with self.subTest(layout=layout):
                with tempfile.TemporaryDirectory() as temp_dir:
                    repo = Path(temp_dir)

                    result = run_apply(repo, decisions_for(domain_layout=layout, autonomous=False))
                    install_profile_skills(repo, "rust-cli")

                    self.assertEqual(result.returncode, 0, result.stderr)
                    domain = (repo / "docs" / "agents" / "domain.md").read_text(
                        encoding="utf-8"
                    )
                    self.assertIn(expected_text, domain)
                    self.assertNotIn(omitted_text, domain)
                    self.assertEqual(
                        (repo / "docs" / "agents" / "monorepo.md").exists(),
                        monorepo_expected,
                    )
                    audit = run_audit(repo)
                    self.assertEqual(audit.returncode, 0, audit.stderr)
                    self.assertEqual(json.loads(audit.stdout)["findings"], [])

    def test_runtime_and_verification_values_render_and_drive_manifest_digest(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            backend = "codex ``edge`` xhigh"
            design = "claude opus xhigh"
            verification = "cargo test --all-targets"

            result = run_apply(
                repo,
                decisions_for(
                    autonomous=True,
                    runtime_backend=backend,
                    runtime_design=design,
                    verification_gate=verification,
                ),
            )

            self.assertEqual(result.returncode, 0, result.stderr)
            autonomous_path = repo / "docs" / "agents" / "autonomous-work.md"
            autonomous = autonomous_path.read_text(encoding="utf-8")
            instructions_path = repo / "docs" / "agents" / "agent-instructions.md"
            instructions = instructions_path.read_text(encoding="utf-8")
            self.assertIn(f"```{backend}```", autonomous)
            self.assertIn(f"`{design}`", autonomous)
            self.assertIn(f"`{verification}`", instructions)
            self.assertNotIn("{{runtime.backend}}", autonomous)
            self.assertNotIn("{{runtime.design}}", autonomous)
            self.assertNotIn("{{verification.gate}}", instructions)

            guide_blocks, _ = parse_managed_blocks(
                Path("docs/agents/agent-instructions.md"), instructions
            )
            manifest = read_manifest(repo)
            artifact = self.managed_artifact(manifest, "guide.agent-instructions")
            self.assertEqual(
                artifact["digest"],
                digest_for(guide_blocks["guide.agent-instructions"].body),
            )

    def test_verification_renders_in_universal_guidance_without_autonomous_work(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            verification = "cargo test --workspace --all-targets"

            result = run_apply(
                repo,
                decisions_for(
                    autonomous=False,
                    verification_gate=verification,
                ),
            )

            self.assertEqual(result.returncode, 0, result.stderr)
            instructions = (repo / "docs" / "agents" / "agent-instructions.md").read_text(
                encoding="utf-8"
            )
            self.assertIn(f"`{verification}`", instructions)
            self.assertNotIn("{{verification.gate}}", instructions)
            self.assertFalse((repo / "docs" / "agents" / "autonomous-work.md").exists())

    def test_frontend_guidance_names_repository_owned_design_contract_only(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            design = repo / "DESIGN.md"
            design.write_text("# Repository-authored design\n", encoding="utf-8")
            expected_design = design.read_bytes()

            result = run_apply_for_profile(
                repo,
                "typescript-bun-monorepo",
                decisions_for(autonomous=False),
            )

            self.assertEqual(result.returncode, 0, result.stderr)
            frontend = (repo / "docs" / "agents" / "frontend.md").read_text(
                encoding="utf-8"
            )
            self.assertIn("repository-owned `DESIGN.md`", frontend)
            self.assertEqual(design.read_bytes(), expected_design)
            for invented_policy in ["authentication policy", "database policy", "transport policy"]:
                self.assertNotIn(invented_policy, frontend)

    def test_typescript_bun_profile_renders_complete_portable_hard_rules(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            (repo / "DESIGN.md").write_text("# Design contract\n", encoding="utf-8")

            result = run_apply_for_profile(
                repo,
                "typescript-bun-monorepo",
                decisions_for(autonomous=False),
            )

            self.assertEqual(result.returncode, 0, result.stderr)
            instructions = (repo / "docs" / "agents" / "agent-instructions.md").read_text(
                encoding="utf-8"
            )
            frontend = (repo / "docs" / "agents" / "frontend.md").read_text(
                encoding="utf-8"
            )
            typescript = (repo / "docs" / "agents" / "typescript-bun.md").read_text(
                encoding="utf-8"
            )
            dispatch = (repo / "docs" / "agents" / "skill-dispatch.md").read_text(
                encoding="utf-8"
            )

            self.assertIn("warnings as errors", typescript)
            self.assertIn("explicit authority", instructions)
            self.assertIn("Read the repository-owned `DESIGN.md`", frontend)
            self.assertIn("dependent interfaces", typescript)
            self.assertIn("Never guess a decision", instructions)
            self.assertIn("external web-research fallback", instructions)
            self.assertIn("Do not use external research tools", instructions)
            for skill in (
                "`conventional-commits`",
                "`github-pr-workflow`",
                "`exa-web-search`",
            ):
                self.assertIn(skill, dispatch)

    def test_unsafe_inline_decision_values_block_apply_without_writes(self):
        cases = [
            ("newline", "cargo test\nmake verify"),
            ("control", "cargo test\x01"),
            (
                "marker",
                "<!-- setup-context-driven:begin id=root.core version=1 -->",
            ),
        ]

        for name, unsafe_value in cases:
            with self.subTest(name=name):
                with tempfile.TemporaryDirectory() as temp_dir:
                    repo = Path(temp_dir)
                    before = snapshot_files(repo)

                    result = run_apply(
                        repo,
                        decisions_for(
                            autonomous=True,
                            verification_gate=unsafe_value,
                        ),
                    )

                    self.assertEqual(result.returncode, 2)
                    self.assertEqual(snapshot_files(repo), before)
                    payload = json.loads(result.stdout)
                    findings = [
                        finding
                        for finding in payload["findings"]
                        if finding["managedId"] == "verification.gate"
                    ]
                    self.assertEqual(len(findings), 1, payload)
                    self.assertEqual(findings[0]["code"], "decision.value.unsafe")
                    self.assertEqual(findings[0]["severity"], "error")

    def test_language_policy_accepts_only_english_without_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            before = snapshot_files(repo)

            result = run_apply(repo, decisions_for(autonomous=False, language="Portuguese"))

            self.assertEqual(result.returncode, 2)
            self.assertEqual(snapshot_files(repo), before)
            payload = json.loads(result.stdout)
            findings = [
                finding
                for finding in payload["findings"]
                if finding["managedId"] == "language.generated"
            ]
            self.assertEqual(len(findings), 1, payload)
            self.assertEqual(findings[0]["code"], "decision.value.invalid")
            self.assertEqual(findings[0]["severity"], "error")

    def test_audit_reports_rendered_drift_and_apply_refreshes_same_artifact(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            decisions = decisions_for(autonomous=True, verification_gate="cargo test --all-targets")
            applied = run_apply(repo, decisions)
            install_profile_skills(repo, "rust-cli")
            self.assertEqual(applied.returncode, 0, applied.stderr)
            guide_path = repo / "docs" / "agents" / "agent-instructions.md"
            guide_path.write_text(
                guide_path.read_text(encoding="utf-8").replace(
                    "cargo test --all-targets",
                    "make verify",
                ),
                encoding="utf-8",
            )
            before_audit = snapshot_files(repo)

            audit = run_audit(repo)

            self.assertEqual(snapshot_files(repo), before_audit)
            payload = json.loads(audit.stdout)
            self.assertFinding(payload, "managed.content.modified", "guide.agent-instructions")
            self.assertIn(
                {
                    "action": "refresh managed content",
                    "path": "docs/agents/agent-instructions.md",
                    "managedId": "guide.agent-instructions",
                },
                payload["plannedChanges"],
            )

            refreshed = run_apply(repo, decisions)

            self.assertEqual(refreshed.returncode, 0, refreshed.stderr)
            self.assertIn("cargo test --all-targets", guide_path.read_text(encoding="utf-8"))
            self.assertEqual(run_audit(repo).returncode, 0)

    def test_qa07_alternate_values_render_selected_guidance_and_are_idempotent(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            decisions = decisions_for(
                spec=False,
                domain_layout="multi-context",
                triage=True,
                autonomous=False,
                runtime_backend="codex alternate-backend xhigh",
                runtime_design="claude alternate-design xhigh",
                verification_gate="cargo test --all-targets",
            )

            first = run_apply(repo, decisions)
            install_profile_skills(repo, "rust-cli")
            after_first = snapshot_files(repo)
            second = run_apply(repo, decisions)

            self.assertEqual(first.returncode, 0, first.stderr)
            self.assertEqual(second.returncode, 0, second.stderr)
            self.assertEqual(snapshot_files(repo), after_first)
            root = (repo / "AGENTS.md").read_text(encoding="utf-8")
            domain = (repo / "docs" / "agents" / "domain.md").read_text(encoding="utf-8")
            self.assertNotIn("id=root.spec-workflow", root)
            self.assertNotIn("id=root.autonomous-work", root)
            self.assertIn("id=root.external-triage", root)
            self.assertIn("CONTEXT-MAP.md", domain)
            self.assertFalse((repo / "docs" / "agents" / "autonomous-work.md").exists())
            audit = run_audit(repo)
            self.assertEqual(audit.returncode, 0, audit.stderr)
            self.assertEqual(json.loads(audit.stdout)["findings"], [])

    def assertFinding(self, payload, code, managed_id):
        matches = [
            finding
            for finding in payload["findings"]
            if finding["code"] == code and finding["managedId"] == managed_id
        ]
        self.assertEqual(len(matches), 1, payload)

    def managed_artifact(self, manifest, managed_id):
        for artifact in manifest["managedArtifacts"]:
            if artifact["id"] == managed_id:
                return artifact
        self.fail(f"missing managed artifact {managed_id}")


def decisions_for(
    *,
    spec=True,
    domain_layout="single-context",
    triage=False,
    autonomous=True,
    runtime_backend="codex gpt-5.5 xhigh",
    runtime_design="claude opus xhigh",
    verification_gate="make verify",
    language="English",
    secondbrain=False,
    repository_extension=False,
):
    return [
        f"spec.scaffold={str(spec).lower()}",
        f"domain.layout={domain_layout}",
        f"triage.external={str(triage).lower()}",
        f"autonomous.enabled={str(autonomous).lower()}",
        f"runtime.backend={runtime_backend}",
        f"runtime.design={runtime_design}",
        f"verification.gate={verification_gate}",
        f"language.generated={language}",
        f"secondbrain.enabled={str(secondbrain).lower()}",
        f"repository.extension.enabled={str(repository_extension).lower()}",
    ]


def run_apply(repo, decisions):
    return run_apply_for_profile(repo, "rust-cli", decisions)


def run_apply_for_profile(repo, profile, decisions):
    args = ["apply", "--repo", str(repo), "--format", "json", "--profile", profile]
    for decision in decisions:
        args.extend(["--decision", decision])
    result = run_context_setup(*args)
    payload = json.loads(result.stdout)
    if result.returncode == 3 and any(
        finding["code"] == "plan.confirmation.required" for finding in payload["findings"]
    ):
        args.extend(["--confirm-plan", payload["planDigest"]])
        return run_context_setup(*args)
    return result


def run_audit(repo):
    return run_context_setup("audit", "--repo", str(repo), "--format", "json")


def read_manifest(repo):
    return json.loads(
        (repo / "docs" / "agents" / "setup-context.json").read_text(encoding="utf-8")
    )


def digest_for(content):
    normalized = content.strip() + "\n"
    return sha256(normalized.encode("utf-8")).hexdigest()


def run_context_setup(*args):
    return run_fixture_context_setup(*args)


if __name__ == "__main__":
    unittest.main()
