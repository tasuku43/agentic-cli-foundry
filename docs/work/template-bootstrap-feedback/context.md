# Work Context: Template bootstrap feedback

## Current behavior

- `.harness/project.json` reports `profile: template`; the worktree was clean before this packet was created.
- `projectconfig.ReadyProblems` compares eight fields independently with `projectconfig.Defaults`, including `github_owner`, so a valid same-owner derivation reports `github_owner still uses the runnable template default`.
- `license_spdx` is already excluded from `ReadyProblems`; retaining MIT is valid, while the public documentation still requires a deliberate license decision.
- `repoguard.repositoryPaths` runs `git ls-files -co --exclude-standard -z` and returns every listed tracked/untracked path. A tracked rename source remains in that output after its working-tree deletion, and `validateRepositoryPaths` then fails its `lstat`.
- The same function currently falls back to a filesystem walk after any Git error. That behavior is not fail closed and can silently change the publishable input set.
- `tools/bootstrap.replacements` is derived from protected `projectconfig.Defaults`, and `buildPlan` excludes `tools/internal/projectconfig/` from content replacement. `docs/04_harness.md` instead embeds replaceable literals while explaining that mechanism.
- The profile enum is currently `template|ready`. Bootstrap changes it transactionally to `ready`, while durable documentation says this proves identity replacement only.
- `scripts/check.sh full` calls `run_fast`, vet, race, tidy/diff checks, then `run_security`, `run_release`, and `run_public`. The Harness profile table currently describes only the first group.
- The gate exports `GOTOOLCHAIN=local`. Before this change it has no common Go preflight; exact-version validation occurs inside the late release profile.
- Local observation on 2026-07-18: PATH resolved the mise-managed Go 1.26.3 binary; automatic selection reported Go 1.26.5 with a downloaded 1.26.5 GOROOT/GOTOOLDIR, while `GOTOOLCHAIN=local` reported Go 1.26.3 and the mise 1.26.3 GOROOT/GOTOOLDIR. `go.mod` requires Go 1.26.5. Machine-specific absolute paths are intentionally omitted from this public packet.
- Template agent-help measurements from a locally built binary: root index 1,517 UTF-8 bytes; exact `help sample read --format agent` 5,359 bytes; namespace `help sample --format agent` 8,359 bytes. The reported derived-project measurements were 1,485, 5,295, and 8,263 bytes respectively.

## Relevant structure

- Entry point: `tools/bootstrap/main.go`, `tools/repoguard/main.go`, and `scripts/check.sh`
- Domain rule: `tools/internal/projectconfig.ReadyProblems`
- Application use case: Not applicable; these are repository harness contracts, not a new runtime capability.
- Infrastructure boundary: Git path enumeration and local Go/tool prerequisite discovery
- CLI catalog or presentation: Bootstrap status text and schema-v3 agent help; no catalog command changes
- Existing tests and harness checks: `tools/internal/projectconfig/config_test.go`, `tools/bootstrap/main_test.go`, `tools/repoguard/main_test.go`, `scripts/test-check-environment.sh`, help shape/growth tests

## Constraints

- Preserve the runnable template identity and `profile: template` in this repository.
- Preserve `tools/internal/projectconfig/defaults.go` byte-for-byte as bootstrap provenance.
- Ignore only confirmed working-tree absence; path validation and filesystem-shape checks must continue to reject links and special files.
- A Git command failure is a guard failure, not permission to inspect a different path set.
- `GOTOOLCHAIN=local` remains the reproducibility policy; diagnostics must help the developer select the exact local version rather than enabling auto-download.
- `full` continues to include every existing subprofile; missing tools are reported earlier, never skipped.
- Agent-help optimization must preserve discover, execute, interpret, and recover data and exact opaque-reference workflows.

## External facts

None. All evidence is repository-local or supplied as a public-safe reproduction. No external content or dependency is introduced.

## Unknowns

- [ ] Whether a future schema should rename `ready` or split `bootstrap_profile` and `product_profile`; compatibility with existing derived repositories and Skills must be designed first.
- [ ] Which scoped-help representation and UTF-8/token budget best predict agent performance across realistic derived catalogs; a schema-v4 fixture corpus and fixed tokenizer would be needed before removing or dictionary-encoding fields.
- [ ] Whether network availability should be actively probed. The current recommendation is to document network-or-cache requirements because an active probe is itself nondeterministic and provider-specific.

## Thesis evidence

- Repeated design decision or point of agent confusion: one state word (`ready`) is carrying an identity-bootstrap meaning while readers infer product completion.
- User outcome or friction observed in the minimal slice: a valid same-owner project could not start, and an otherwise successful transactional rename could not run its immediate gates.
- Code workaround or exception being considered: staging before checks, changing owner, enabling automatic toolchain download, or deleting scoped help detail would route around the template's stated contracts.
- Current thesis that resolves it, or proposed thesis revision: Theses 5 and 7 require the canonical path to model the real working-tree state and emit actionable failures; no thesis revision is needed for the two defects.
- Downstream impact: bootstrap/harness/public documentation and regression checks change. Product-readiness state and agent-help schema changes remain reviewed follow-ups rather than incidental mechanisms.

## Reproduction or observation

```sh
go run ./tools/projectmeta --field profile
git ls-files -co --exclude-standard -z
go version
go env GOVERSION GOROOT GOTOOLDIR
GOTOOLCHAIN=local go version
GOTOOLCHAIN=local go env GOVERSION GOROOT GOTOOLDIR
```

Observed profile: `template`. The exact failing owner and unstaged-rename behaviors are captured as failing regression tests before their fixes. The help byte counts above were measured from the template executable without network access.

Before implementation, the focused regressions failed with eight `ReadyProblems` including `github_owner`, retained both `new-path.txt` and missing `old-path.txt`, and silently returned a filesystem path set outside a Git repository. After implementation, `go test ./tools/internal/projectconfig ./tools/repoguard ./tools/bootstrap` passes.

The real mixed installation was also reproduced after preflight implementation: the selected Go command reported 1.26.5 while stale `GOROOT`/`GOTOOLDIR` selected a 1.26.3 compiler. The gate emitted one `check preflight: Go toolchain mismatch` block containing required version, binary, reported version, `GOVERSION`, `GOROOT`, `GOTOOLDIR`, compiler version, and the local-toolchain remediation. The synthetic shell fixture proves the same classification deterministically.

Final verification on 2026-07-18 used the exact local Go 1.26.5 installation with stale `GOROOT` removed. `task check` passed its fast, vet, race, security (`No vulnerabilities found`), release, and public subprofiles; `task public:check` passed independently. The first sandboxed full run confirmed the documented network-or-cache prerequisite by failing only when the pinned gosec module could not resolve the public Go proxy; the authorized network run completed.

## Security and public-boundary notes

- Assets and side effects involved: public source path selection and bootstrap working-tree integrity; tests use temporary Git repositories only.
- Credentials or confidential data involved: none.
- New dependencies, destinations, files, processes, or generated content: no dependency or network destination; existing `git`, `go`, `gofmt`, ShellCheck, and release prerequisites are diagnosed.
- External schema provenance, publication rights, and drift evidence: not applicable.
- Pagination, timeout, retry, idempotency, and cancellation facts: not applicable.
- Publication and licensing concerns: retaining the template license is allowed, but derived owners must still review it deliberately; the public guard remains fail closed.

## Glossary

- **Bootstrap/identity readiness:** exact repository identity replacement completed; it does not assert product, security, legal, or release readiness.
- **Product readiness:** project-specific theses, tasks, trust boundaries, and capability decisions are concrete enough to begin capability work.
- **Provenance defaults:** the immutable runnable template values in `projectconfig.Defaults` used as exact replacement sources.
