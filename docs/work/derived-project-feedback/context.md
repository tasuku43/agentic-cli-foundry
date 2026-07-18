# Work Context: Integrate proven derived-project lessons

## Current behavior

- Foundry HEAD `b37ec8369b2c` is clean before this work and stores
  `profile: template` in `.harness/project.json`.
- Repository-local Git comparison proved that the evaluated derived project
  began from an identical Foundry tree. It therefore supplies direct
  derived-project evidence rather than an unrelated design comparison; no
  copied commit history is retained in this public packet.
- The Foundry catalog bounds each root agent-help entry but deliberately has no
  whole-response scoped-help budget. `help.go` emits one workflow record for
  every matching producer-field/consumer-input pair.
- Offline UTF-8 measurements were 8,359 bytes for Foundry `help sample`, 43,145
  bytes for a derived exact-command scope, and 184,018 bytes for a derived
  namespace scope.
  The latter two contained 56 and 228 expanded workflow records respectively.
- Foundry operation enums implement text marshaling but not strict text
  unmarshaling. Its argv usage validator rejects a single allowed value even
  though an exact literal is a useful input contract.
- `CommandOutput.Completeness` currently uses `complete|paged` both to select a
  public cursor protocol and to describe result completeness. It cannot state
  that an invocation completely delivered a provider-capped or differential
  window without implying exhaustive external coverage.
- One derived presentation packet is marked Complete and has all task boxes
  checked, but its goal acceptance boxes remain unchecked. Its full gate still
  passes, so current repository checks do not enforce work-packet claim
  consistency.
- Foundry already contains derived-project-informed tool-local fixed targets,
  provider-neutral non-secret authentication configuration, a human-handoff
  scorecard, agent-help schema 4, exact toolchain preflight, and fail-closed Git
  enumeration. Those mechanisms must not be replaced by older derived copies.
- The evaluated project is an existing schema-1 `ready` repository. A provisional schema-2
  provenance implementation would reject it before any safe migration decision
  could be made.
- The provisional implementation read ambient `HEAD^{tree}`. It did not prove
  that the object was the template root commit, neutralize `GIT_DIR` and
  `GIT_WORK_TREE`, bind apply to the tree reviewed during preview, or exercise a
  fresh bootstrap followed by the public guard.
- A fresh-bootstrap/public-guard simulation exposed a hard-coded template
  module in the copied `tools/internal/projectconfig/config_test.go`. Skipping
  that package during replacement leaves a ready derived repository that its
  own public guard rejects.
- Existing bootstrap backup removal happens after content updates, the ready
  profile, and renames are committed. A cleanup error can therefore return
  failure after durable identity changes; this predates the provisional
  provenance slice and is not changed here.

## Relevant structure

- Entry point: `cmd/agentic-cli-foundry`
- Domain rule: `internal/domain/operation`
- Application use case: unchanged; this work changes template-wide contracts
- Infrastructure boundary: unchanged
- CLI catalog or presentation: `internal/cli/catalog.go`, `help.go`, and tests
- Existing tests and harness checks: `tools/contractlint`, `tools/repoguard`,
  `scripts/check.sh`, catalog/help contract tests

## Constraints

- `cli.Catalog` remains the only public-command source of truth.
- A machine-readable help shape change requires an explicit schema version.
- Root help remains an outcome index; scoped help must retain complete
  invocation, output, failure, mutation, authentication, and recovery facts.
- Fresh derived repositories must not preserve the template's Git history.
- Project configuration remains schema 1 in this packet. Template provenance
  cannot be added until one independent follow-up solves root-commit evidence,
  existing-ready migration, root-bound Git, preview/apply binding, and a
  fresh-bootstrap public E2E as one contract.
- New checks must have actionable errors and run through `scripts/check.sh`.
- No new network, process, credential, runtime side effect, or dependency is
  required.

## External facts

No network-sourced facts. Lineage and comparison evidence came from local Git
repositories; all serialized scale and behavior fixtures are synthetic.

## Resolved unknowns

- [x] The smallest grouped workflow shape is one record per reference kind with
  deduplicated `producers` and `consumers`. Its Cartesian expansion is tested
  for exact edge equivalence rather than serialized.
- [x] Work-packet lifecycle validation belongs in `repoguard`: the existing tool
  already owns regular-file, symlink, hygiene, and public repository claims and
  runs through every canonical profile.
- [x] The provisional template-provenance/schema-2 slice is withdrawn. An
  independent follow-up must simultaneously validate the template root commit,
  define schema-1 `ready` migration (including unavailable evidence), bind Git
  resolution to the selected repository root, bind apply to the previewed tree,
  and pass a fresh-bootstrap-to-public-guard E2E. It must also resolve or
  explicitly separate the existing post-commit backup-cleanup ambiguity.
- [x] Output contracts need two independent axes: `delivery` describes one
  invocation's complete-or-paged transport, while `collection_coverage`
  describes the declared task scope as not applicable, exhaustive, bounded, or
  differential.

## Implementation evidence

- Agent-help schema 5 groups 324 implicit edges into 24,493 encoded bytes in
  the derived-scale synthetic fixture; the equivalent pair-expanded shape is
  177,759 bytes and exceeds the 65,536-byte response budget.
- Manual root and `sample` scoped help probes encode 1,517 and 8,222 bytes and
  retain invocation, workflow, output, fault, and recovery facts.
- Focused operation, catalog, help, project-configuration, bootstrap, and
  repository-guard tests pass; the project-metadata command package builds (it
  has no package-local tests).
- Sample CLI contracts prove a successful empty JSON collection, preserve two
  distinct opaque IDs behind the same display label, pass only the selected ID
  to read, and suppress presentation when the repository returns another ID.
- Work-packet guard fixtures reject fenced/comment-only metadata and completion,
  inline-code or escaped comment lookalikes, non-ASCII blank/closing/path
  whitespace, tab-separated tasks, tab-relative list fences, raw Markdown-link
  successors, template targets, and successor cycles.
- A default-shell fast check correctly rejected Go 1.26.3. The canonical fast
  gate passes with the repository-required Go 1.26.5 toolchain and isolated
  writable Go caches.

## Thesis evidence

- Repeated design decision or point of agent confusion: deterministic joins and
  semantic interpretation moved from prospective agent post-processing into the
  derived CLI more than once.
- User outcome or friction observed in the minimal slice: scoped help grew to
  184,018 bytes for six related commands because workflow edges were expanded.
- Code workaround or exception being considered: deleting workflow/recovery
  facts to meet a size target would weaken the product contract.
- Current thesis that resolves it, or proposed thesis revision: Thesis 2 should
  define operational closure and retain detail through a compact representation.
- Downstream impact: product, architecture, security, Skill, catalog/help schema,
  agent-readiness tests, and harness checks. Project provenance metadata is
  deferred to the independent follow-up above.

## Reproduction or observation

```sh
go run ./cmd/agentic-cli-foundry help sample --format agent
go test ./internal/cli ./internal/domain/operation ./tools/repoguard
```

The exact derived-scale runtime measurements are retained above as design
evidence. The implementation uses a synthetic catalog corpus rather than a
dependency on the evaluated repository.

## Security and public-boundary notes

- Assets and side effects involved: public machine contracts and repository
  work metadata only; no runtime external side effect.
- Credentials or confidential data involved: none.
- New dependencies, destinations, files, processes, or generated content: none
  expected beyond a repository-local lint and fixtures/tests.
- External schema provenance, publication rights, and drift evidence: not
  applicable; derived evidence is summarized without copied provider data.
- Pagination, timeout, retry, idempotency, and cancellation facts: unchanged.
- Publication and licensing concerns: all new fixtures must be synthetic.

## Glossary

- **Operational closure:** a supported outcome owns the deterministic
  selection, join, and interpretation needed for routine success; users do not
  need an undeclared parser or exploratory call.
- **Contextual reference validation:** checking that a valid opaque value appears
  in a field whose declared reference kind matches its meaning.
- **Grouped workflow:** a reference-kind adjacency representation with producer
  and consumer sets, rather than one record per pair.
- **Deferred template provenance:** a future source-contract claim that is not
  part of the current schema-1 implementation and must satisfy the combined
  follow-up requirements above before adoption.
