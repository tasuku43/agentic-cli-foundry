# Work Plan: Integrate proven derived-project lessons

- Status: Accepted
- Goal: [goal.md](goal.md)
- Context: [context.md](context.md)
- Tasks: [tasks.md](tasks.md)

## Chosen approach

Promote only provider-neutral lessons supported by direct derived-project
evidence. Change governing documents before or with mechanisms, version the
agent-help compatibility break, add narrow shared validators rather than product
types, and use synthetic fixtures for scale and negative paths.

## Alternatives considered

### Copy the derived implementation

Rejected. The derived implementation contains provider vocabulary, a broad
tagged request/result union, a generic task port, provider-specific
authentication policy and limits, and an older Foundry harness/help schema.
Copying it would regress current template contracts and violate task-owned port
guidance.

### Add prose guidance only

Rejected. The observed failures include Cartesian output growth, enum decode
failure, conflicting work-packet state, and reference-kind laundering. Each has
a practical validator or test, so prose alone would leave the claim incomplete.

### Import the full presentation evaluator

Rejected. Its valuable result is the evidence protocol, while the implementation
is large, product-specific, nondeterministic, and has audited scorer/oracle gaps.
The template will adopt a small typed-fixture protocol instead.

## Design

### Public contract

- Agent-help advances from schema 4 to a grouped-workflow schema.
- Output metadata replaces the overloaded completeness value with an explicit
  delivery protocol and collection-coverage declaration. Runtime success JSON
  schemas do not change.
- Supported outcome guidance gains operational closure; this does not silently
  change an intentionally declared raw/export utility into a composed task.
- No public command path, effect, role, authentication requirement, side effect,
  or exit code changes.
- Exact single-literal argv input becomes representable in catalog metadata.

### Layer changes

- Domain: strict text decoding for finite operation enums.
- Application: no production change; guidance gains request/result validation.
- Infrastructure: no production change.
- CLI and catalog: exact-literal validation, global fault signature validation,
  delivery/collection-coverage validation, grouped agent-help workflow
  projection, schema/version and scale tests.

### Data and control flow

Catalog reference declarations remain canonical. Agent help groups their derived
edges by reference kind and lists each unique producer and consumer once. No
runtime command discovers or reconstructs references.

### Error and cancellation behavior

Runtime behavior is unchanged. New catalog validation reports deterministic
contract errors before dispatch. Enum decoding rejects unknown text without
changing the receiver.

### Security and public boundary

Derived-project evidence is summarized from repository-local inspection; all
public serialized data and fixtures are synthetic. Work-packet validation
follows regular-file/no-symlink policy. Template provenance and project schema
2 are withdrawn from this packet; the current public boundary remains schema 1
with no new template-identity exception.

## Implementation slices

1. Work packet, governing thesis/product/architecture/security language, and
   failing focused tests.
2. Domain and catalog contract hardening.
3. Delivery/coverage output metadata plus the versioned grouped agent-help shape
   and derived-scale budget regression.
4. Work-packet lifecycle validation; withdraw the provisional provenance slice
   and record its independent follow-up contract.
5. Skill, readiness protocol, harness claims, and full verification.

## Verification

- Unit and contract tests: operation, catalog, help, repository/work-packet
  validators, focused schema-1 project/bootstrap regressions, and sample
  empty-result/same-label/target-mismatch presentation-boundary tests.
- Negative side-effect tests: no new side effects; invalid catalog contracts
  fail before dispatch.
- Opaque-reference and complete-pagination tests: existing suite plus grouped
  workflow edge-equivalence tests.
- Structured output, hostile-output, and recovery tests: existing help/output
  snapshots and schema tests must pass after the explicit version change.
- Agent-readiness scenario: root-to-scope limits remain two/one invocations; the
  synthetic derived-scale scope remains under its fixed whole-response budget.
- Human-handoff scorecard: not applicable; no setup or authentication change.
- Manual observation: run root and sample scoped agent help and inspect the
  grouped workflow.
- Required profiles: `task check` (includes security, release, and public).
- Generated-diff or artifact checks: repository status and tidy/format diff.

## Rollout and rollback

Agent-help consumers must treat the schema bump as a compatibility change. The
old pair-expanded workflow representation and command-local `next_actions` are
removed rather than maintained as a parallel registry. Consumers expand the
schema 5 producer and consumer sets only when they need individual edges;
fault-local recovery actions remain where the fault is declared. Output
metadata consumers migrate from `completeness` to independent `delivery` and
`collection_coverage` fields. Runtime success envelopes do not change.

The provisional project-schema-2/template-provenance slice is withdrawn. One
independent follow-up must satisfy all of these conditions together before
adding a stored provenance claim:

- verify the template root commit/tree rather than accepting whichever commit
  ambient `HEAD` names;
- define a safe schema-1 `ready` migration, including the explicit outcome when
  recoverable source evidence is unavailable;
- run Git inspection against the selected repository root with ambient
  `GIT_DIR`, `GIT_WORK_TREE`, and equivalent routing state neutralized;
- bind apply to the exact source tree reviewed during preview so a changed tree
  cannot be accepted silently; and
- prove a synthetic fresh bootstrap followed by the public guard, including
  copied internal project-configuration tests and every retained protected-
  identity exception.

That follow-up must also record and decide the existing backup-cleanup behavior:
cleanup currently occurs after durable updates and can report failure after the
profile is already `ready`. This packet does not change that mechanism.

The exact-input value-kind/arity idea is deliberately deferred. Adding metadata
while handlers still parse raw argv would create a second, non-executable
registry. A future packet should first supply representative boolean, integer,
repeatable, and multiword commands, then introduce one catalog-owned parser and
derive handler inputs and help from the same declarations.

The remaining changes are validator and guidance additions. Reverting a public
schema change before a stable release is safe only when its implementation,
documentation, and contract tests are reverted together.

## Documentation promotion

- Operational closure: theses, product contract, add-capability Skill, readiness.
- Semantic-before-presentation and request/result binding: architecture,
  security model, Skill, readiness.
- Grouped workflow and budgets: architecture, product contract, harness.
- Work-packet lifecycle validation: harness and work templates.
- Template provenance/schema 2: explicitly deferred to the independent
  follow-up contract in this plan.
- Capability retirement and lightweight presentation evidence: Skill and work
  templates, without importing provider-specific machinery.
