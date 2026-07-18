# Work Plan: Template bootstrap feedback

- Status: Complete; higher-level migrations remain proposed follow-ups
- Goal: [goal.md](goal.md)
- Context: [context.md](context.md)

## Chosen approach

Fix identity readiness at the narrow project-config boundary and Git path selection at the repository-guard boundary. Add regression tests before implementation. Then make three local, contract-compatible improvements: explain bootstrap through `projectconfig.Defaults`, label current `ready` output as identity-only without changing the stored enum, and add an early gate preflight that validates the exact local Go toolchain and release/full prerequisites. Record product-readiness and scoped-help schema changes as explicit follow-ups.

## Alternatives considered

### Identity readiness

- **Chosen:** keep field-level checks for project-specific name, binary, module, repository, description, Formula class, and security contact; allow `github_owner` and `license_spdx` to equal template defaults. This states what makes a derivation meaningful and preserves actionable diagnostics.
- Reject only when the entire `Project` equals `Defaults`: this would allow partially renamed repositories with stale runnable identity.
- Infer identity from owner/repository as one tuple: this is better than owner alone but still fails to require module, binary, and display identity updates used throughout generated files.

### Repository path selection

- **Chosen:** require successful Git enumeration, validate returned paths as local, ignore only entries confirmed absent at enumeration time, then run the existing no-link/regular-file validation and full working-tree shape scan. This handles `D old` plus `?? new` without staging.
- Parse porcelain rename status: unstaged renames are commonly represented as deletion plus untracked addition and rename detection is heuristic.
- Keep filesystem fallback: this changes the publication input set on Git errors and is not fail closed.

### Bootstrap explanation

- **Chosen:** document `projectconfig.Defaults` as the protected source of exact replacement values and add a test that the actual Harness section contains no replaceable identity literal.
- Exclude all of `docs/04_harness.md` from replacement: derived-project examples elsewhere in that document may legitimately need their identity updated and a broad exclusion creates drift.

### Bootstrap versus product readiness

- **Immediate compatible improvement:** keep schema-v1 `profile: ready`, but label it `identity-ready` in output and durable documentation.
- Rename to `identity_ready`/`bootstrap_ready`: clear but breaks config validation, Skills, repoguard, docs, tests, and existing derived repositories.
- Split `bootstrap_profile` and `product_profile`: recommended follow-up. A schema-v2 migration can preserve explicit bootstrap state while a new product gate verifies concrete theses/security/capability removal or replacement. It requires a definition of executable product evidence before implementation.
- Add an unversioned product gate while retaining the field: risks a second ambiguous source of truth unless its manifest and completion gate are designed together.

### Full-gate prerequisites and toolchain diagnosis

- **Chosen:** run a shared preflight before every profile, require exact `go.mod` Go version under `GOTOOLCHAIN=local`, compare the selected binary, reported version, GOVERSION, GOROOT/GOTOOLDIR compiler, and report one diagnostic. For `full`/`release`, detect ShellCheck before long checks. Document the complete composition plus local tools and network-or-cache requirement.
- Rely on Go's compiler errors and late release checks: produces many misleading errors and wastes the fast/security/race work.
- Enable automatic toolchain selection: conflicts with the established reproducibility boundary.

### Scoped agent-help budget

- A fixed whole-response UTF-8 budget now would be simple but an arbitrary threshold either rejects useful namespace scopes or permits excessive repetition.
- **Recommended follow-up:** design schema v4 with dictionaries for common global faults/types, command-local code references, and an explicit minimal-execution versus complete-contract request. Fix a tokenizer and derived-catalog corpus, then set both UTF-8 and token regression budgets.
- Evaluation must prove: unknown outcome reaches a scoped contract within two discovery calls; exact inputs and effects remain executable; output completeness and reference bytes remain interpretable; every classified failure retains a structured recovery command. No information is removed in this change.

## Design

### Public contract

No new public CLI command or capability is added. Bootstrap accepts same-owner derivations, reports `identity ready`, and the gate fails earlier with a stable human diagnostic. The stored profile enum and agent-help schema remain unchanged.

### Layer changes

- Domain: no production domain change.
- Application: no application change.
- Infrastructure: repository-tool Git enumeration filters absent paths and fails on Git errors.
- CLI and catalog: no catalog change; bootstrap tool status text only.

### Data and control flow

Project identity is validated before bootstrap planning. Repository guard obtains the canonical Git path set, drops confirmed absent deletions, validates every remaining path without following links, scans the full working-tree shape, and only then reads content. Gate preflight sanitizes the Go environment, validates the local toolchain, then dispatches the unchanged profile.

### Error and cancellation behavior

Missing derived identity remains a list of field-level problems. Git failures abort repository inspection. Only `os.IsNotExist` at enumeration time is ignored; all other filesystem errors fail. Toolchain mismatches produce one diagnostic block and no Go-backed gate work. Cancellation and external retry behavior are unchanged.

### Security and public boundary

No credentials, network destination, or dependency is added. Temporary Git repositories use synthetic paths. Links and special files remain rejected through both selected-path validation and full shape inspection. The provenance file is not edited.

## Implementation slices

1. Work packet, baseline hashes/measurements, and failing regression tests.
2. `ReadyProblems` and Git path-selection fixes.
3. Protected-defaults documentation contract and identity-only status wording.
4. Gate preflight plus full-profile/prerequisite documentation.
5. Scoped-help findings and schema-v4 follow-up recorded in agent-readiness documentation.

## Verification

- Unit and contract tests: projectconfig, bootstrap, repoguard, CLI help, and preflight shell fixtures.
- Negative side-effect tests: Git error, symbolic link, special file, and missing toolchain paths fail closed.
- Opaque-reference and complete-pagination tests: unchanged existing suite.
- Structured output, hostile-output, and recovery tests: unchanged existing suite; scoped contract data is not removed.
- Agent-readiness scenario and discovery-round-trip count: existing two/one invocation bounds; record byte measurements.
- Manual observation: owner-only equality and unstaged rename regression tests.
- Required profiles: `task check`, `task public:check`.
- Generated-diff or artifact checks: confirm `defaults.go` blob and template profile are unchanged.

Final evidence: focused Go tests and mixed-toolchain shell fixtures pass; `task check` and independent `task public:check` pass; profile remains `template`; the provenance blob remains `5beb6c11f2f92ad4ebbdaf5f44e70eb66224ecb0`.

## Rollout and rollback

The two bug fixes are backward-compatible relaxations/corrections. Earlier guards required staging or rejected same-owner derivations; no valid protected path is newly skipped. The diagnostic preflight may reject an environment earlier than before but enforces the already documented exact-toolchain and ShellCheck contracts. Rollback is a normal source revert; there is no persisted state migration.

## Documentation promotion

- Promote meaningful derived identity, protected defaults, identity-only readiness, full composition, and prerequisites to `docs/04_harness.md`, README, CONTRIBUTING, and public/release guidance as applicable.
- Promote scoped-help measurement and the contract-preserving evaluation requirements to `docs/09_agent_readiness_validation.md`.
- A future bootstrap/product profile split requires an ADR or schema migration design before code changes.
