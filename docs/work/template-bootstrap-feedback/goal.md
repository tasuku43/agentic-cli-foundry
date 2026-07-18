# Work Goal: Make derived bootstrap evidence accurate and immediately verifiable

- Status: Complete
- Owner: Template maintainers
- Target: Template bootstrap feedback
- Related ADRs: None; product-readiness state remains a follow-up contract decision

## Outcome

A repository derived from Agentic CLI Foundry can retain the template's GitHub owner or license, apply identity renames without staging them, and immediately run the canonical gates. Bootstrap and gate output explain the identity-only state and diagnose local toolchain prerequisites without misleading or late failures.

## Why now

An actual derived-project bootstrap exposed two blocking defects and five contract/documentation gaps. The owner check rejects a valid same-owner repository, and repository guard attempts to inspect tracked rename sources that no longer exist in the working tree. The same exercise also showed self-rewriting harness prose, ambiguous `ready` terminology, hidden `full` profile prerequisites, noisy Go toolchain mismatch failures, and unbounded scoped agent-help size.

## Non-goals

- Changing this template repository's identity or `.harness/project.json` profile.
- Modifying `tools/internal/projectconfig/defaults.go`; it remains the exact bootstrap provenance record.
- Weakening security, release, public, symlink, special-file, or Git failure checks.
- Introducing a product-readiness state migration before its higher-level contract and compatibility impact are reviewed.
- Removing scoped agent-help information before an executable discover/execute/interpret/recover evaluation proves the replacement contract.

## Acceptance criteria

- [x] A derived identity that changes the repository, module, binary, display name, and other project-specific fields may retain the template GitHub owner and license.
- [x] Runnable template identity and individually unchanged required derived fields are rejected with actionable field-level problems.
- [x] Repository guard ignores a tracked path already deleted from the working tree, still inspects its untracked rename destination, and fails closed for links, special files, and Git errors.
- [x] Bootstrap mechanism prose survives exact identity replacement and refers to the protected provenance source rather than a replaceable literal.
- [x] The current `ready` value is described and reported as identity-only; a product-readiness migration recommendation and impact surface are recorded before mechanism work.
- [x] `full` composition and local prerequisites are documented, with missing prerequisites and Go toolchain inconsistency reported before long-running checks.
- [x] Template agent-help sizes are measured and a contract-preserving scoped-budget recommendation is recorded with an evaluation scenario.
- [x] Template identity/profile and the blob of `tools/internal/projectconfig/defaults.go` remain unchanged.
- [x] `task check` and `task public:check` pass.

## Governing documents

- Thesis: Thesis 5 (executable claims), Thesis 6 (public boundary), Thesis 7 (one maintainable path)
- Product contract section: Derived-project completion checklist and machine-readable agent-help contract
- Architecture or security invariant: Fail-closed repository inspection, protected bootstrap provenance, canonical gate environment
- Existing ADR: None

## Completion definition

The blocking fixes have regression tests, safe local improvements for findings 3–7 are promoted to durable documentation and checks, unresolved higher-level changes have explicit recommendations and impact, required profiles pass, and the template identity plus provenance blob are unchanged.
