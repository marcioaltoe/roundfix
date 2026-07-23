"""Standard TypeScript Monorepo Profile contract tests.

Suite: strict Standard TypeScript Monorepo Profile assets
Invariant: one declarative 0.0.1 profile preserves the exact stack, topology,
architecture, capabilities, decisions, activations, and Verification contract.
Boundary IN: profile asset loading, HTTP decision normalization, profile plan and
snapshot rendering, and local Repository Capability evaluation.
Boundary OUT: setup audit/apply transaction integration and profile retirement.
"""

import json
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SKILL_ROOT / "scripts"))

from context_assets import (  # noqa: E402
    AssetValidationError,
    build_standard_profile_plan,
    clone_assets_to,
    load_asset_catalog,
    render_standard_profile_plan,
    render_standard_profile_snapshot,
    write_json,
)
from context_capabilities import (  # noqa: E402
    CapabilityStatus,
    evaluate_repository_capabilities,
)


PROFILE_ID = "standard-typescript-monorepo"
EXPECTED_STACK = (
    "TypeScript",
    "Bun",
    "Turborepo",
    "Vite",
    "React",
    "Hono",
    "Drizzle",
    "Zod",
    "Tailwind",
    "shadcn",
    "TanStack Query",
    "TanStack Router",
    "Better Auth",
    "PostgreSQL",
    "LogTape",
    "Oxlint",
    "Oxfmt",
    "Vitest",
)
EXPECTED_ACTIVATION_BUNDLES = (
    "bundle.debugging",
    "bundle.delivery",
    "bundle.frontend-react",
    "bundle.frontend-ui-quality",
    "bundle.hono-endpoint",
    "bundle.hono-endpoint-persistence",
    "bundle.production-code",
    "bundle.qa",
    "bundle.security",
    "bundle.testing",
)
EXPECTED_VERIFICATION = (
    ("verification.format", "Oxfmt", "bun run format"),
    ("verification.lint", "Oxlint", "bun run lint"),
    ("verification.test", "Vitest", "bun run test"),
    ("verification.build", "Turborepo", "bun run build"),
    ("verification.workspace", "Bun", "bun run verify"),
)


class StandardTypeScriptMonorepoProfileTests(unittest.TestCase):
    def test_profile_inspection_reports_exact_stack_topology_and_optional_modules(self):
        contract = self._contract()

        self.assertEqual(contract.profile_id, PROFILE_ID)
        self.assertEqual(contract.version, "0.0.1")
        self.assertEqual(contract.marker_version, "0.0.1")
        self.assertEqual(contract.stack, EXPECTED_STACK)
        self.assertEqual(
            contract.required_workspaces,
            ("packages/frontend", "packages/backend"),
        )
        self.assertEqual(contract.optional_modules, ("Inngest", "Docker"))

    def test_architecture_contract_uses_systems_and_layers_not_generic_buckets(self):
        architecture = self._contract().architecture

        self.assertEqual(architecture.frontend_organization, "systems")
        self.assertEqual(architecture.frontend_public_boundary, "public system boundary")
        self.assertEqual(architecture.frontend_internal_imports, "direct")
        self.assertEqual(
            architecture.backend_layers,
            ("domain", "application", "infrastructure"),
        )
        self.assertEqual(architecture.backend_http_handlers, "thin")
        self.assertEqual(architecture.backend_use_cases, "HTTP-independent")
        self.assertEqual(architecture.backend_persistence, "Drizzle-owned")
        self.assertEqual(architecture.rejected_normative_buckets, ("modules", "services"))

    def test_undecided_http_contract_requests_exact_repository_owned_modes(self):
        plan = build_standard_profile_plan(load_asset_catalog(SKILL_ROOT))
        document = json.loads(render_standard_profile_plan(plan))

        self.assertEqual(plan.unresolved_decisions, ("http.contract",))
        self.assertIsNone(document["httpContract"])
        self.assertEqual(
            document["unresolvedDecisions"],
            [
                {
                    "id": "http.contract",
                    "prompt": "Choose the repository HTTP Contract Decision.",
                    "values": ["REST", "Post-only"],
                }
            ],
        )

    def test_selected_http_contract_and_typed_exceptions_persist_in_plan_and_snapshot(self):
        decision = {
            "mode": "REST",
            "exceptions": [
                {
                    "scope": "/webhooks/*",
                    "methods": ["POST"],
                    "owner": "payments",
                    "reason": "The provider owns webhook delivery semantics.",
                }
            ],
            "source": {
                "path": "docs/architecture/http-contract.json",
                "digest": "a" * 64,
            },
        }
        plan = build_standard_profile_plan(load_asset_catalog(SKILL_ROOT), decision)
        plan_document = json.loads(render_standard_profile_plan(plan))
        snapshot_document = json.loads(render_standard_profile_snapshot(plan))

        self.assertEqual(plan.unresolved_decisions, ())
        self.assertEqual(plan_document["httpContract"], decision)
        self.assertEqual(snapshot_document["httpContract"], decision)
        self.assertEqual(
            snapshot_document["profile"]["architecture"]["backend"]["layers"],
            ["domain", "application", "infrastructure"],
        )

    def test_invalid_http_mode_or_untyped_exception_is_rejected(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        cases = (
            {"mode": "GraphQL", "exceptions": [], "source": self._source()},
            {
                "mode": "Post-only",
                "exceptions": [{"scope": "/jobs", "methods": ["POST"]}],
                "source": self._source(),
            },
        )

        for decision in cases:
            with self.subTest(decision=decision):
                with self.assertRaises(ValueError):
                    build_standard_profile_plan(catalog, decision)

    def test_required_postgresql_and_logtape_evidence_blocks_but_optional_modules_do_not(self):
        contract = self._contract()
        with tempfile.TemporaryDirectory() as temp_dir:
            evaluation = evaluate_repository_capabilities(
                Path(temp_dir),
                contract.capabilities,
                executable_lookup=lambda _name: None,
            )

        outcomes = {
            outcome.capability.title: outcome for outcome in evaluation.outcomes
        }
        for title in ("PostgreSQL", "LogTape"):
            with self.subTest(title=title):
                self.assertEqual(outcomes[title].status, CapabilityStatus.MISSING)
                self.assertTrue(outcomes[title].blocking)
        for title in ("Inngest", "Docker"):
            with self.subTest(title=title):
                self.assertEqual(outcomes[title].status, CapabilityStatus.MISSING)
                self.assertFalse(outcomes[title].blocking)

    def test_profile_binds_universal_capabilities_exact_bundles_and_verification(self):
        contract = self._contract()

        self.assertEqual(contract.capability_sets, ("universal",))
        self.assertEqual(contract.activation_bundles, EXPECTED_ACTIVATION_BUNDLES)
        self.assertEqual(
            tuple((entry.entry_id, entry.tool, entry.command) for entry in contract.verification),
            EXPECTED_VERIFICATION,
        )

    def test_profile_plan_and_snapshot_are_byte_stable_and_only_use_0_0_1_owned_versions(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        decision = {"mode": "Post-only", "exceptions": [], "source": self._source()}

        first = build_standard_profile_plan(catalog, decision)
        second = build_standard_profile_plan(catalog, decision)
        plan_bytes = render_standard_profile_plan(first)
        snapshot_bytes = render_standard_profile_snapshot(first)

        self.assertEqual(plan_bytes, render_standard_profile_plan(second))
        self.assertEqual(snapshot_bytes, render_standard_profile_snapshot(second))
        self.assertEqual(
            json.loads(plan_bytes)["schemaVersion"],
            "setup-context-driven/profile-plan/0.0.1",
        )
        snapshot = json.loads(snapshot_bytes)
        self.assertEqual(
            snapshot["schemaVersion"],
            "setup-context-driven/profile-snapshot/0.0.1",
        )
        self.assertEqual(snapshot["version"], "0.0.1")
        self.assertEqual(snapshot["markerVersion"], "0.0.1")
        self.assertEqual(snapshot["profile"]["version"], "0.0.1")

    def test_profile_mutations_fail_at_the_asset_boundary(self):
        cases = (
            ("stack", lambda profile: profile["stack"].pop(), "profile.stack.capability.mismatch"),
            (
                "workspace",
                lambda profile: profile["workspaces"].append(
                    {"path": "packages/shared", "strength": "required"}
                ),
                "profile.workspace.capability.mismatch",
            ),
            (
                "http-default",
                lambda profile: profile["httpContract"].__setitem__("default", "REST"),
                "profile.httpContract.default.invalid",
            ),
            (
                "activation",
                lambda profile: profile["activationBundles"].pop(),
                "profile.activationBundles.mismatch",
            ),
            (
                "version",
                lambda profile: profile.__setitem__("version", "0.0.2"),
                "profile.version.invalid",
            ),
        )

        for name, mutate, diagnostic in cases:
            with self.subTest(name=name):
                with tempfile.TemporaryDirectory() as temp_dir:
                    root = Path(temp_dir)
                    clone_assets_to(SKILL_ROOT, root)
                    profile_path = root / "assets" / "profiles" / f"{PROFILE_ID}.json"
                    profile = json.loads(profile_path.read_text(encoding="utf-8"))
                    mutate(profile)
                    write_json(profile_path, profile)

                    with self.assertRaises(AssetValidationError) as raised:
                        load_asset_catalog(root)

                self.assertIn(diagnostic, "\n".join(raised.exception.diagnostics))

    def _contract(self):
        return load_asset_catalog(SKILL_ROOT).standard_profiles[PROFILE_ID]

    def _source(self):
        return {
            "path": "docs/architecture/http-contract.json",
            "digest": "b" * 64,
        }


if __name__ == "__main__":
    unittest.main()
